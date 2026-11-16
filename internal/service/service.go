package service

import (
	"studybuddy/internal/models"
	"studybuddy/internal/store"
)

type GoalService struct{ Store *store.Store }
type PlanService struct{ Store *store.Store }
type SessionService struct{ Store *store.Store }

func NewGoalService(s *store.Store) *GoalService       { return &GoalService{Store: s} }
func NewPlanService(s *store.Store) *PlanService       { return &PlanService{Store: s} }
func NewSessionService(s *store.Store) *SessionService { return &SessionService{Store: s} }

func (g *GoalService) Create(in models.GoalInput) (models.Goal, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return models.Goal{}, err
	}
	return g.Store.CreateGoal(&in)
}

func (g *GoalService) Update(id int64, in models.GoalInput) (models.Goal, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return models.Goal{}, err
	}
	return g.Store.UpdateGoal(id, &in)
}

func (g *GoalService) Get(id int64) (models.Goal, error)         { return g.Store.GetGoal(id) }
func (g *GoalService) List(status string) ([]models.Goal, error) { return g.Store.ListGoals(status) }
func (g *GoalService) Delete(id int64) error                     { return g.Store.DeleteGoal(id) }

func (p *PlanService) List(goalID *int64) ([]models.StudyPlan, error) {
	if goalID != nil {
		if _, err := p.Store.GetGoal(*goalID); err != nil {
			return nil, err
		}
	}
	return p.Store.ListPlans(goalID)
}
func (p *PlanService) Get(id int64) (models.StudyPlan, error) { return p.Store.GetPlan(id) }
func (p *PlanService) Create(in models.PlanInput) (models.StudyPlan, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return models.StudyPlan{}, err
	}
	if _, err := p.Store.GetGoal(in.GoalID); err != nil {
		return models.StudyPlan{}, err
	}
	return p.Store.CreatePlan(&in)
}
func (p *PlanService) Update(id int64, in models.PlanInput) (models.StudyPlan, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return models.StudyPlan{}, err
	}
	if _, err := p.Store.GetGoal(in.GoalID); err != nil {
		return models.StudyPlan{}, err
	}
	return p.Store.UpdatePlan(id, &in)
}
func (p *PlanService) Delete(id int64) error { return p.Store.DeletePlan(id) }
func (p *PlanService) Topics(id int64) ([]models.PlanTopic, error) {
	if _, err := p.Store.GetPlan(id); err != nil {
		return nil, err
	}
	return p.Store.ListTopics(id)
}
func (p *PlanService) CreateTopic(in models.TopicInput) (models.PlanTopic, error) {
	in.Normalize()
	if in.Status == "" {
		in.Status = "pending"
	}
	if err := in.Validate(); err != nil {
		return models.PlanTopic{}, err
	}
	if _, err := p.Store.GetPlan(in.PlanID); err != nil {
		return models.PlanTopic{}, err
	}
	return p.Store.CreateTopic(&in)
}
func (p *PlanService) UpdateTopic(id int64, in models.TopicInput) (models.PlanTopic, error) {
	in.Normalize()
	if in.Status == "" {
		in.Status = "pending"
	}
	if err := in.Validate(); err != nil {
		return models.PlanTopic{}, err
	}
	if _, err := p.Store.GetPlan(in.PlanID); err != nil {
		return models.PlanTopic{}, err
	}
	old, err := p.Store.GetTopic(id)
	if err != nil {
		return models.PlanTopic{}, err
	}
	if old.Status == "completed" && in.Status != "completed" {
		return models.PlanTopic{}, models.ErrInvalidTransition
	}
	return p.Store.UpdateTopic(id, &in)
}
func (p *PlanService) DeleteTopic(id int64) error { return p.Store.DeleteTopic(id) }
func (p *PlanService) Schedule(from, to string) ([]models.PlanTopic, error) {
	return p.Store.ListScheduledTopics(from, to)
}

func (s *SessionService) Create(in models.SessionInput) (models.StudySession, error) {
	if err := in.Validate(); err != nil {
		return models.StudySession{}, err
	}
	// ensure goal exists
	if _, err := s.Store.GetGoal(in.GoalID); err != nil {
		return models.StudySession{}, err
	}
	return s.Store.CreateSession(&in)
}
func (s *SessionService) Get(id int64) (models.StudySession, error) { return s.Store.GetSession(id) }
func (s *SessionService) List(goalID *int64, status string) ([]models.StudySession, error) {
	return s.Store.ListSessions(goalID, status)
}
