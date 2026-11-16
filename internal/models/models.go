package models

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidTransition = errors.New("completed activities cannot be reopened")

type Goal struct {
	ID                 int64     `json:"id"`
	Title              string    `json:"title"`
	Subject            *string   `json:"subject,omitempty"`
	Description        *string   `json:"description,omitempty"`
	Status             string    `json:"status"`
	TargetDate         *string   `json:"target_date,omitempty"`
	CurrentLevel       *string   `json:"current_level,omitempty"`
	WeeklyStudyMinutes *int      `json:"weekly_study_minutes,omitempty"`
	PreferredStudyDays []int     `json:"preferred_study_days"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type GoalInput struct {
	Title              string  `json:"title"`
	Subject            *string `json:"subject"`
	Description        *string `json:"description"`
	Status             string  `json:"status"`
	TargetDate         *string `json:"target_date"`
	CurrentLevel       *string `json:"current_level"`
	WeeklyStudyMinutes *int    `json:"weekly_study_minutes"`
	PreferredStudyDays []int   `json:"preferred_study_days"`
}

func (g *GoalInput) Normalize() {
	g.Title = strings.TrimSpace(g.Title)
	if g.Status == "" {
		g.Status = "active"
	}
}

func (g *GoalInput) Validate() error {
	if g.Title == "" {
		return errors.New("title is required")
	}
	switch g.Status {
	case "active", "paused", "completed", "archived":
	default:
		return errors.New("status must be one of: active, paused, completed, archived")
	}
	if g.TargetDate != nil && *g.TargetDate != "" {
		if _, err := time.Parse("2006-01-02", *g.TargetDate); err != nil {
			return errors.New("target_date must be YYYY-MM-DD")
		}
	}
	if g.WeeklyStudyMinutes != nil && (*g.WeeklyStudyMinutes < 1 || *g.WeeklyStudyMinutes > 10080) {
		return errors.New("weekly_study_minutes must be between 1 and 10080")
	}
	seen := map[int]bool{}
	for _, d := range g.PreferredStudyDays {
		if d < 0 || d > 6 || seen[d] {
			return errors.New("preferred_study_days must contain unique weekdays from 0 (Sunday) to 6 (Saturday)")
		}
		seen[d] = true
	}
	return nil
}

type StudyPlan struct {
	ID          int64     `json:"id"`
	GoalID      int64     `json:"goal_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PlanInput struct {
	GoalID      int64   `json:"goal_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func (p *PlanInput) Normalize() { p.Title = strings.TrimSpace(p.Title) }
func (p *PlanInput) Validate() error {
	if p.GoalID <= 0 {
		return errors.New("goal_id is required")
	}
	if p.Title == "" {
		return errors.New("title is required")
	}
	return nil
}

type PlanTopic struct {
	ID               int64     `json:"id"`
	PlanID           int64     `json:"plan_id"`
	Title            string    `json:"title"`
	Description      *string   `json:"description,omitempty"`
	Position         int       `json:"position"`
	EstimatedMinutes *int      `json:"estimated_minutes,omitempty"`
	ScheduledDate    *string   `json:"scheduled_date,omitempty"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TopicInput struct {
	PlanID           int64   `json:"plan_id"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	Position         int     `json:"position"`
	EstimatedMinutes *int    `json:"estimated_minutes"`
	ScheduledDate    *string `json:"scheduled_date"`
	Status           string  `json:"status"`
}

func (t *TopicInput) Normalize() { t.Title = strings.TrimSpace(t.Title) }
func (t *TopicInput) Validate() error {
	if t.PlanID <= 0 {
		return errors.New("plan_id is required")
	}
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.EstimatedMinutes != nil && (*t.EstimatedMinutes < 1 || *t.EstimatedMinutes > 1440) {
		return errors.New("estimated_minutes must be between 1 and 1440")
	}
	if t.ScheduledDate != nil && *t.ScheduledDate != "" {
		if _, err := time.Parse("2006-01-02", *t.ScheduledDate); err != nil {
			return errors.New("scheduled_date must be YYYY-MM-DD")
		}
	}
	if t.Status == "" {
		t.Status = "pending"
	}
	switch t.Status {
	case "pending", "in_progress", "completed":
	default:
		return errors.New("status must be pending, in_progress, or completed")
	}
	return nil
}

type StudySession struct {
	ID          int64      `json:"id"`
	GoalID      int64      `json:"goal_id"`
	PlanTopicID *int64     `json:"plan_topic_id,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type SessionInput struct {
	GoalID      int64   `json:"goal_id"`
	PlanTopicID *int64  `json:"plan_topic_id"`
	StartedAt   *string `json:"started_at"`
}

func (s *SessionInput) Validate() error {
	if s.GoalID <= 0 {
		return errors.New("goal_id is required")
	}
	if s.StartedAt != nil && *s.StartedAt != "" {
		if _, err := time.Parse(time.RFC3339, *s.StartedAt); err != nil {
			return errors.New("started_at must be RFC3339")
		}
	}
	return nil
}

type SessionEvent struct {
	ID         int64     `json:"id"`
	SessionID  int64     `json:"session_id"`
	EventType  string    `json:"event_type"`
	OccurredAt time.Time `json:"occurred_at"`
	Note       *string   `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type EventInput struct {
	SessionID  int64   `json:"session_id"`
	EventType  string  `json:"event_type"`
	OccurredAt *string `json:"occurred_at"`
	Note       *string `json:"note"`
}

func (e *EventInput) Validate() error {
	if e.SessionID <= 0 {
		return errors.New("session_id is required")
	}
	switch e.EventType {
	case "start", "pause", "resume", "stop":
	default:
		return errors.New("event_type must be one of: start, pause, resume, stop")
	}
	if e.OccurredAt != nil && *e.OccurredAt != "" {
		if _, err := time.Parse(time.RFC3339, *e.OccurredAt); err != nil {
			return errors.New("occurred_at must be RFC3339")
		}
	}
	return nil
}

func NowUTC() time.Time                        { return time.Now().UTC() }
func FormatRFC3339(t time.Time) string         { return t.UTC().Format(time.RFC3339) }
func ParseRFC3339(s string) (time.Time, error) { return time.Parse(time.RFC3339, s) }
