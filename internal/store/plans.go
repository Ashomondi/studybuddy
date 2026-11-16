package store

import (
	"database/sql"
	"time"

	"studybuddy/internal/models"
)

func (s *Store) CreatePlan(p *models.PlanInput) (models.StudyPlan, error) {
	now := models.NowUTC()
	nowStr := models.FormatRFC3339(now)
	res, err := s.DB.Exec(
		`INSERT INTO study_plans (goal_id, title, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		p.GoalID, p.Title, nullStringPtr(p.Description), nowStr, nowStr,
	)
	if err != nil {
		return models.StudyPlan{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.StudyPlan{}, err
	}
	return s.GetPlan(id)
}

func (s *Store) UpdatePlan(id int64, p *models.PlanInput) (models.StudyPlan, error) {
	res, err := s.DB.Exec(`UPDATE study_plans SET goal_id=?, title=?, description=?, updated_at=? WHERE id=?`, p.GoalID, p.Title, nullStringPtr(p.Description), models.FormatRFC3339(models.NowUTC()), id)
	if err != nil {
		return models.StudyPlan{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.StudyPlan{}, ErrNotFound
	}
	return s.GetPlan(id)
}
func (s *Store) DeletePlan(id int64) error {
	res, err := s.DB.Exec(`DELETE FROM study_plans WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) GetPlan(id int64) (models.StudyPlan, error) {
	row := s.DB.QueryRow(
		`SELECT id, goal_id, title, description, created_at, updated_at FROM study_plans WHERE id = ?`, id,
	)
	var p models.StudyPlan
	var desc sql.NullString
	var created, updated string
	if err := row.Scan(&p.ID, &p.GoalID, &p.Title, &desc, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return models.StudyPlan{}, ErrNotFound
		}
		return models.StudyPlan{}, err
	}
	if desc.Valid {
		v := desc.String
		p.Description = &v
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	p.CreatedAt = ct.UTC()
	p.UpdatedAt = ut.UTC()
	return p, nil
}

func (s *Store) ListPlans(goalID *int64) ([]models.StudyPlan, error) {
	q := `SELECT id, goal_id, title, description, created_at, updated_at FROM study_plans`
	var args []any
	if goalID != nil {
		q += ` WHERE goal_id = ?`
		args = append(args, *goalID)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := s.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.StudyPlan
	for rows.Next() {
		var p models.StudyPlan
		var desc sql.NullString
		var created, updated string
		if err := rows.Scan(&p.ID, &p.GoalID, &p.Title, &desc, &created, &updated); err != nil {
			return nil, err
		}
		if desc.Valid {
			v := desc.String
			p.Description = &v
		}
		ct, _ := time.Parse(time.RFC3339, created)
		ut, _ := time.Parse(time.RFC3339, updated)
		p.CreatedAt = ct.UTC()
		p.UpdatedAt = ut.UTC()
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) ListTopics(planID int64) ([]models.PlanTopic, error) {
	rows, err := s.DB.Query(
		`SELECT id, plan_id, title, description, position, estimated_minutes, scheduled_date, status, created_at, updated_at FROM plan_topics WHERE plan_id = ? ORDER BY position, id`, planID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.PlanTopic
	for rows.Next() {
		var t models.PlanTopic
		var desc, date sql.NullString
		var est sql.NullInt64
		var created, updated string
		if err := rows.Scan(&t.ID, &t.PlanID, &t.Title, &desc, &t.Position, &est, &date, &t.Status, &created, &updated); err != nil {
			return nil, err
		}
		if desc.Valid {
			v := desc.String
			t.Description = &v
		}
		if est.Valid {
			v := int(est.Int64)
			t.EstimatedMinutes = &v
		}
		if date.Valid {
			v := date.String
			t.ScheduledDate = &v
		}
		ct, _ := time.Parse(time.RFC3339, created)
		ut, _ := time.Parse(time.RFC3339, updated)
		t.CreatedAt = ct.UTC()
		t.UpdatedAt = ut.UTC()
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) GetTopic(id int64) (models.PlanTopic, error) {
	rows, err := s.DB.Query(`SELECT id,plan_id,title,description,position,estimated_minutes,scheduled_date,status,created_at,updated_at FROM plan_topics WHERE id=?`, id)
	if err != nil {
		return models.PlanTopic{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return models.PlanTopic{}, ErrNotFound
	}
	var t models.PlanTopic
	var desc, date sql.NullString
	var est sql.NullInt64
	var created, updated string
	if err := rows.Scan(&t.ID, &t.PlanID, &t.Title, &desc, &t.Position, &est, &date, &t.Status, &created, &updated); err != nil {
		return models.PlanTopic{}, err
	}
	if desc.Valid {
		v := desc.String
		t.Description = &v
	}
	if est.Valid {
		v := int(est.Int64)
		t.EstimatedMinutes = &v
	}
	if date.Valid {
		v := date.String
		t.ScheduledDate = &v
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	t.CreatedAt = ct.UTC()
	t.UpdatedAt = ut.UTC()
	return t, nil
}

func (s *Store) ListScheduledTopics(from, to string) ([]models.PlanTopic, error) {
	rows, err := s.DB.Query(`SELECT id,plan_id,title,description,position,estimated_minutes,scheduled_date,status,created_at,updated_at FROM plan_topics WHERE scheduled_date>=? AND scheduled_date<=? ORDER BY scheduled_date,position,id`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.PlanTopic{}
	for rows.Next() {
		var t models.PlanTopic
		var desc, date sql.NullString
		var est sql.NullInt64
		var created, updated string
		if err := rows.Scan(&t.ID, &t.PlanID, &t.Title, &desc, &t.Position, &est, &date, &t.Status, &created, &updated); err != nil {
			return nil, err
		}
		if desc.Valid {
			v := desc.String
			t.Description = &v
		}
		if est.Valid {
			v := int(est.Int64)
			t.EstimatedMinutes = &v
		}
		if date.Valid {
			v := date.String
			t.ScheduledDate = &v
		}
		ct, _ := time.Parse(time.RFC3339, created)
		ut, _ := time.Parse(time.RFC3339, updated)
		t.CreatedAt = ct.UTC()
		t.UpdatedAt = ut.UTC()
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) CreateTopic(in *models.TopicInput) (models.PlanTopic, error) {
	now := models.FormatRFC3339(models.NowUTC())
	var date any
	if in.ScheduledDate != nil && *in.ScheduledDate != "" {
		date = *in.ScheduledDate
	}
	res, err := s.DB.Exec(`INSERT INTO plan_topics(plan_id,title,description,position,estimated_minutes,scheduled_date,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, in.PlanID, in.Title, nullStringPtr(in.Description), in.Position, nullIntPtr(in.EstimatedMinutes), date, in.Status, now, now)
	if err != nil {
		return models.PlanTopic{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.PlanTopic{}, err
	}
	topics, err := s.ListTopics(in.PlanID)
	if err != nil {
		return models.PlanTopic{}, err
	}
	for _, t := range topics {
		if t.ID == id {
			return t, nil
		}
	}
	return models.PlanTopic{}, ErrNotFound
}
func (s *Store) UpdateTopic(id int64, in *models.TopicInput) (models.PlanTopic, error) {
	now := models.FormatRFC3339(models.NowUTC())
	var date any
	if in.ScheduledDate != nil && *in.ScheduledDate != "" {
		date = *in.ScheduledDate
	}
	res, err := s.DB.Exec(`UPDATE plan_topics SET plan_id=?,title=?,description=?,position=?,estimated_minutes=?,scheduled_date=?,status=?,updated_at=? WHERE id=?`, in.PlanID, in.Title, nullStringPtr(in.Description), in.Position, nullIntPtr(in.EstimatedMinutes), date, in.Status, now, id)
	if err != nil {
		return models.PlanTopic{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.PlanTopic{}, ErrNotFound
	}
	topics, err := s.ListTopics(in.PlanID)
	if err != nil {
		return models.PlanTopic{}, err
	}
	for _, t := range topics {
		if t.ID == id {
			return t, nil
		}
	}
	return models.PlanTopic{}, ErrNotFound
}
func (s *Store) DeleteTopic(id int64) error {
	res, err := s.DB.Exec(`DELETE FROM plan_topics WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
