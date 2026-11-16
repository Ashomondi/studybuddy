package store

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"studybuddy/internal/models"
)

func (s *Store) CreateGoal(g *models.GoalInput) (models.Goal, error) {
	now := models.NowUTC()
	nowStr := models.FormatRFC3339(now)
	res, err := s.DB.Exec(
		`INSERT INTO goals (title, subject, description, status, target_date, current_level, weekly_study_minutes, preferred_study_days, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		g.Title, nullStringPtr(g.Subject), nullStringPtr(g.Description), g.Status, nullStringPtr(g.TargetDate), nullStringPtr(g.CurrentLevel), nullIntPtr(g.WeeklyStudyMinutes), encodeDays(g.PreferredStudyDays), nowStr, nowStr,
	)
	if err != nil {
		return models.Goal{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Goal{}, err
	}
	return s.GetGoal(id)
}

func (s *Store) GetGoal(id int64) (models.Goal, error) {
	row := s.DB.QueryRow(
		`SELECT id, title, subject, description, status, target_date, current_level, weekly_study_minutes, preferred_study_days, created_at, updated_at FROM goals WHERE id = ?`, id,
	)
	return scanGoal(row)
}

func (s *Store) ListGoals(status string) ([]models.Goal, error) {
	q := `SELECT id, title, subject, description, status, target_date, current_level, weekly_study_minutes, preferred_study_days, created_at, updated_at FROM goals`
	var args []any
	if status != "" {
		q += ` WHERE status = ?`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := s.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Goal
	for rows.Next() {
		g, err := scanGoalRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (s *Store) UpdateGoal(id int64, g *models.GoalInput) (models.Goal, error) {
	nowStr := models.FormatRFC3339(models.NowUTC())
	res, err := s.DB.Exec(
		`UPDATE goals SET title=?, subject=?, description=?, status=?, target_date=?, current_level=?, weekly_study_minutes=?, preferred_study_days=?, updated_at=? WHERE id=?`,
		g.Title, nullStringPtr(g.Subject), nullStringPtr(g.Description), g.Status, nullStringPtr(g.TargetDate), nullStringPtr(g.CurrentLevel), nullIntPtr(g.WeeklyStudyMinutes), encodeDays(g.PreferredStudyDays), nowStr, id,
	)
	if err != nil {
		return models.Goal{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.Goal{}, ErrNotFound
	}
	return s.GetGoal(id)
}

func (s *Store) DeleteGoal(id int64) error {
	res, err := s.DB.Exec(`DELETE FROM goals WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanGoal(row *sql.Row) (models.Goal, error) {
	var g models.Goal
	var subject, description, targetDate, level, days sql.NullString
	var weekly sql.NullInt64
	var created, updated string
	if err := row.Scan(&g.ID, &g.Title, &subject, &description, &g.Status, &targetDate, &level, &weekly, &days, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return models.Goal{}, ErrNotFound
		}
		return models.Goal{}, err
	}
	if subject.Valid {
		v := subject.String
		g.Subject = &v
	}
	if description.Valid {
		v := description.String
		g.Description = &v
	}
	if targetDate.Valid {
		v := targetDate.String
		g.TargetDate = &v
	}
	if level.Valid {
		v := level.String
		g.CurrentLevel = &v
	}
	if weekly.Valid {
		v := int(weekly.Int64)
		g.WeeklyStudyMinutes = &v
	}
	_ = json.Unmarshal([]byte(days.String), &g.PreferredStudyDays)
	if g.PreferredStudyDays == nil {
		g.PreferredStudyDays = []int{}
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	g.CreatedAt = ct.UTC()
	g.UpdatedAt = ut.UTC()
	return g, nil
}

func scanGoalRows(rows *sql.Rows) (models.Goal, error) {
	var g models.Goal
	var subject, description, targetDate, level, days sql.NullString
	var weekly sql.NullInt64
	var created, updated string
	if err := rows.Scan(&g.ID, &g.Title, &subject, &description, &g.Status, &targetDate, &level, &weekly, &days, &created, &updated); err != nil {
		return models.Goal{}, err
	}
	if subject.Valid {
		v := subject.String
		g.Subject = &v
	}
	if description.Valid {
		v := description.String
		g.Description = &v
	}
	if targetDate.Valid {
		v := targetDate.String
		g.TargetDate = &v
	}
	if level.Valid {
		v := level.String
		g.CurrentLevel = &v
	}
	if weekly.Valid {
		v := int(weekly.Int64)
		g.WeeklyStudyMinutes = &v
	}
	_ = json.Unmarshal([]byte(days.String), &g.PreferredStudyDays)
	if g.PreferredStudyDays == nil {
		g.PreferredStudyDays = []int{}
	}
	ct, _ := time.Parse(time.RFC3339, created)
	ut, _ := time.Parse(time.RFC3339, updated)
	g.CreatedAt = ct.UTC()
	g.UpdatedAt = ut.UTC()
	return g, nil
}

func encodeDays(days []int) string { b, _ := json.Marshal(days); return string(b) }

func nullStringPtr(p *string) any {
	if p == nil || strings.TrimSpace(*p) == "" {
		// still store empty? For subject/description allow null if nil/empty
		if p == nil {
			return nil
		}
		// store empty string? We'll store null for empty
		return nil
	}
	return *p
}

func nullIntPtr(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
