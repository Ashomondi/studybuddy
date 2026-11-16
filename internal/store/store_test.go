package store

import (
	"database/sql"
	"path/filepath"
	"testing"

	"studybuddy/internal/database"
	"studybuddy/internal/models"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPlanActivitiesPersistAfterReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "studybuddy.db")
	db, err := database.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	s := New(db)
	g, err := s.CreateGoal(&models.GoalInput{Title: "Learn Go", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := s.CreatePlan(&models.PlanInput{GoalID: g.ID, Title: "Go basics"})
	if err != nil {
		t.Fatal(err)
	}
	minutes := 45
	date := "2026-10-29"
	activity, err := s.CreateTopic(&models.TopicInput{PlanID: p.ID, Title: "Variables", EstimatedMinutes: &minutes, ScheduledDate: &date, Status: "pending"})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = database.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	got, err := New(db).GetTopic(activity.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Variables" || got.ScheduledDate == nil || *got.ScheduledDate != date || got.EstimatedMinutes == nil || *got.EstimatedMinutes != minutes {
		t.Fatalf("activity did not persist: %+v", got)
	}
}

func TestGoalCRUD(t *testing.T) {
	db := newTestDB(t)
	s := New(db)

	g, err := s.CreateGoal(&models.GoalInput{Title: "Learn Go", Status: "active"})
	if err != nil {
		t.Fatal(err)
	}
	if g.ID == 0 || g.Title != "Learn Go" {
		t.Fatalf("unexpected goal: %+v", g)
	}

	got, err := s.GetGoal(g.ID)
	if err != nil || got.Title != "Learn Go" {
		t.Fatalf("get: %v %+v", err, got)
	}

	got2, err := s.UpdateGoal(g.ID, &models.GoalInput{Title: "Learn Go deeply", Status: "active"})
	if err != nil || got2.Title != "Learn Go deeply" {
		t.Fatalf("update: %v %+v", err, got2)
	}

	if err := s.DeleteGoal(g.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetGoal(g.ID); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestGoalValidation(t *testing.T) {
	in := models.GoalInput{Title: "", Status: "active"}
	in.Normalize()
	if err := in.Validate(); err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestSessionCreateAndEvents(t *testing.T) {
	db := newTestDB(t)
	s := New(db)
	g, _ := s.CreateGoal(&models.GoalInput{Title: "G", Status: "active"})
	sess, err := s.CreateSession(&models.SessionInput{GoalID: g.ID})
	if err != nil {
		t.Fatal(err)
	}
	if sess.Status != "running" {
		t.Fatalf("status: %s", sess.Status)
	}
	ev, err := s.AddEvent(&models.EventInput{SessionID: sess.ID, EventType: "stop"})
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventType != "stop" {
		t.Fatal("event")
	}
	got, _ := s.GetSession(sess.ID)
	if got.Status != "completed" || got.EndedAt == nil {
		t.Fatalf("session not closed: %+v", got)
	}
}
