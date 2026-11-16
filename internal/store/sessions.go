package store

import (
	"database/sql"
	"time"

	"studybuddy/internal/models"
)

func (s *Store) CreateSession(in *models.SessionInput) (models.StudySession, error) {
	started := models.NowUTC()
	if in.StartedAt != nil && *in.StartedAt != "" {
		t, err := time.Parse(time.RFC3339, *in.StartedAt)
		if err != nil {
			return models.StudySession{}, err
		}
		started = t.UTC()
	}
	nowStr := models.FormatRFC3339(models.NowUTC())
	startStr := models.FormatRFC3339(started)
	var topic any
	if in.PlanTopicID != nil {
		topic = *in.PlanTopicID
	}
	res, err := s.DB.Exec(
		`INSERT INTO study_sessions (goal_id, plan_topic_id, started_at, status, created_at, updated_at) VALUES (?, ?, ?, 'running', ?, ?)`,
		in.GoalID, topic, startStr, nowStr, nowStr,
	)
	if err != nil {
		return models.StudySession{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.StudySession{}, err
	}
	return s.GetSession(id)
}

func (s *Store) GetSession(id int64) (models.StudySession, error) {
	row := s.DB.QueryRow(
		`SELECT id, goal_id, plan_topic_id, started_at, ended_at, status, created_at, updated_at FROM study_sessions WHERE id = ?`, id,
	)
	return scanSession(row)
}

func (s *Store) ListSessions(goalID *int64, status string) ([]models.StudySession, error) {
	q := `SELECT id, goal_id, plan_topic_id, started_at, ended_at, status, created_at, updated_at FROM study_sessions`
	var clauses []string
	var args []any
	if goalID != nil {
		clauses = append(clauses, "goal_id = ?")
		args = append(args, *goalID)
	}
	if status != "" {
		clauses = append(clauses, "status = ?")
		args = append(args, status)
	}
	if len(clauses) > 0 {
		q += " WHERE " + joinClauses(clauses)
	}
	q += ` ORDER BY started_at DESC`
	rows, err := s.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.StudySession
	for rows.Next() {
		ss, err := scanSessionRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ss)
	}
	return out, rows.Err()
}

func (s *Store) AddEvent(e *models.EventInput) (models.SessionEvent, error) {
	occurred := models.NowUTC()
	if e.OccurredAt != nil && *e.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, *e.OccurredAt)
		if err != nil {
			return models.SessionEvent{}, err
		}
		occurred = t.UTC()
	}
	nowStr := models.FormatRFC3339(models.NowUTC())
	res, err := s.DB.Exec(
		`INSERT INTO session_events (session_id, event_type, occurred_at, note, created_at) VALUES (?, ?, ?, ?, ?)`,
		e.SessionID, e.EventType, models.FormatRFC3339(occurred), nullStringPtr(e.Note), nowStr,
	)
	if err != nil {
		return models.SessionEvent{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.SessionEvent{}, err
	}
	// update session status on stop
	if e.EventType == "stop" {
		now := models.FormatRFC3339(models.NowUTC())
		s.DB.Exec(`UPDATE study_sessions SET status='completed', ended_at=?, updated_at=? WHERE id=?`, models.FormatRFC3339(occurred), now, e.SessionID)
	}
	if e.EventType == "pause" {
		now := models.FormatRFC3339(models.NowUTC())
		s.DB.Exec(`UPDATE study_sessions SET status='paused', updated_at=? WHERE id=?`, now, e.SessionID)
	}
	if e.EventType == "resume" || e.EventType == "start" {
		now := models.FormatRFC3339(models.NowUTC())
		s.DB.Exec(`UPDATE study_sessions SET status='running', updated_at=? WHERE id=?`, now, e.SessionID)
	}
	row := s.DB.QueryRow(
		`SELECT id, session_id, event_type, occurred_at, note, created_at FROM session_events WHERE id = ?`, id,
	)
	var ev models.SessionEvent
	var note sql.NullString
	var occ, created string
	if err := row.Scan(&ev.ID, &ev.SessionID, &ev.EventType, &occ, &note, &created); err != nil {
		return models.SessionEvent{}, err
	}
	if note.Valid {
		v := note.String
		ev.Note = &v
	}
	o, _ := time.Parse(time.RFC3339, occ)
	c, _ := time.Parse(time.RFC3339, created)
	ev.OccurredAt = o.UTC()
	ev.CreatedAt = c.UTC()
	return ev, nil
}

func (s *Store) ListEvents(sessionID int64) ([]models.SessionEvent, error) {
	rows, err := s.DB.Query(
		`SELECT id, session_id, event_type, occurred_at, note, created_at FROM session_events WHERE session_id = ? ORDER BY occurred_at`, sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.SessionEvent
	for rows.Next() {
		var ev models.SessionEvent
		var note sql.NullString
		var occ, created string
		if err := rows.Scan(&ev.ID, &ev.SessionID, &ev.EventType, &occ, &note, &created); err != nil {
			return nil, err
		}
		if note.Valid {
			v := note.String
			ev.Note = &v
		}
		o, _ := time.Parse(time.RFC3339, occ)
		c, _ := time.Parse(time.RFC3339, created)
		ev.OccurredAt = o.UTC()
		ev.CreatedAt = c.UTC()
		out = append(out, ev)
	}
	return out, rows.Err()
}

func scanSession(row *sql.Row) (models.StudySession, error) {
	var ss models.StudySession
	var topic sql.NullInt64
	var started, created, updated string
	var ended sql.NullString
	if err := row.Scan(&ss.ID, &ss.GoalID, &topic, &started, &ended, &ss.Status, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return models.StudySession{}, ErrNotFound
		}
		return models.StudySession{}, err
	}
	if topic.Valid {
		v := topic.Int64
		ss.PlanTopicID = &v
	}
	st, _ := time.Parse(time.RFC3339, started)
	ss.StartedAt = st.UTC()
	if ended.Valid {
		et, _ := time.Parse(time.RFC3339, ended.String)
		et = et.UTC()
		ss.EndedAt = &et
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	ss.CreatedAt = ct.UTC()
	ss.UpdatedAt = ut.UTC()
	return ss, nil
}

func scanSessionRows(rows *sql.Rows) (models.StudySession, error) {
	var ss models.StudySession
	var topic sql.NullInt64
	var started, created, updated string
	var ended sql.NullString
	if err := rows.Scan(&ss.ID, &ss.GoalID, &topic, &started, &ended, &ss.Status, &created, &updated); err != nil {
		return models.StudySession{}, err
	}
	if topic.Valid {
		v := topic.Int64
		ss.PlanTopicID = &v
	}
	st, _ := time.Parse(time.RFC3339, started)
	ss.StartedAt = st.UTC()
	if ended.Valid {
		et, _ := time.Parse(time.RFC3339, ended.String)
		et = et.UTC()
		ss.EndedAt = &et
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	ss.CreatedAt = ct.UTC()
	ss.UpdatedAt = ut.UTC()
	return ss, nil
}

func joinClauses(c []string) string {
	// simple join with AND
	out := c[0]
	for i := 1; i < len(c); i++ {
		out += " AND " + c[i]
	}
	return out
}
