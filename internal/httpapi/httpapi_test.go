package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"studybuddy/internal/database"
	"studybuddy/internal/service"
	"studybuddy/internal/store"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	st := store.New(db)
	deps := RouterDeps{
		Goals:    service.NewGoalService(st),
		Plans:    service.NewPlanService(st),
		Sessions: service.NewSessionService(st),
		DBPing:   func() error { return db.Ping() },
	}
	return NewRouter("http://localhost:5173", deps)
}

func TestHealth(t *testing.T) {
	r := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rr.Code, rr.Body.String())
	}
}

func TestGoalCreateAndGet(t *testing.T) {
	r := newTestRouter(t)
	body := []byte(`{"title":"Learn Go","status":"active"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	json.Unmarshal(rr.Body.Bytes(), &created)
	id := int64(created["id"].(float64))

	req2 := httptest.NewRequest(http.MethodGet, "/api/goals/"+itoa(id), nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("get status %d: %s", rr2.Code, rr2.Body.String())
	}
}

func TestGoalInvalid(t *testing.T) {
	r := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewReader([]byte(`{"title":""}`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGoalCreateReportsUnknownJSONField(t *testing.T) {
	r := newTestRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/api/goals", bytes.NewBufferString(`{"title":"Learn Go","unexpected":true}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest || !bytes.Contains(rr.Body.Bytes(), []byte("unexpected")) {
		t.Fatalf("expected useful unknown-field validation, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestManualPlanActivityCRUDAndSchedule(t *testing.T) {
	r := newTestRouter(t)
	call := func(method, path, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		return rr
	}
	g := call("POST", "/api/goals", `{"title":"Learn Go","current_level":"Beginner","weekly_study_minutes":300,"preferred_study_days":[1,3,5]}`)
	if g.Code != 201 {
		t.Fatalf("goal %d: %s", g.Code, g.Body)
	}
	var goal map[string]any
	_ = json.Unmarshal(g.Body.Bytes(), &goal)
	gid := int(goal["id"].(float64))
	p := call("POST", "/api/plans", `{"goal_id":`+itoa(int64(gid))+`,"title":"Go foundations"}`)
	if p.Code != 201 {
		t.Fatalf("plan %d: %s", p.Code, p.Body)
	}
	var plan map[string]any
	_ = json.Unmarshal(p.Body.Bytes(), &plan)
	pid := int(plan["id"].(float64))
	a := call("POST", "/api/activities", `{"plan_id":`+itoa(int64(pid))+`,"title":"Variables","estimated_minutes":30,"scheduled_date":"2026-10-29"}`)
	if a.Code != 201 {
		t.Fatalf("activity %d: %s", a.Code, a.Body)
	}
	var activity map[string]any
	_ = json.Unmarshal(a.Body.Bytes(), &activity)
	aid := itoa(int64(activity["id"].(float64)))
	if got := call("GET", "/api/plans/"+itoa(int64(pid))+"/topics", ""); got.Code != 200 || !bytes.Contains(got.Body.Bytes(), []byte("Variables")) {
		t.Fatalf("activities not persisted: %d %s", got.Code, got.Body)
	}
	if got := call("GET", "/api/schedule?from=2026-10-29&to=2026-10-29", ""); got.Code != 200 || !bytes.Contains(got.Body.Bytes(), []byte("Variables")) {
		t.Fatalf("schedule missing activity: %d %s", got.Code, got.Body)
	}
	if got := call("PUT", "/api/activities/"+aid, `{"plan_id":`+itoa(int64(pid))+`,"title":"Variables","status":"completed","estimated_minutes":30,"scheduled_date":"2026-10-29"}`); got.Code != 200 {
		t.Fatalf("complete activity: %d %s", got.Code, got.Body)
	}
	if got := call("PUT", "/api/activities/"+aid, `{"plan_id":`+itoa(int64(pid))+`,"title":"Variables","status":"pending","estimated_minutes":30,"scheduled_date":"2026-10-29"}`); got.Code != 400 {
		t.Fatalf("expected invalid transition, got %d %s", got.Code, got.Body)
	}
	if got := call("DELETE", "/api/plans/"+itoa(int64(pid)), ""); got.Code != 204 {
		t.Fatalf("delete plan: %d %s", got.Code, got.Body)
	}
	if got := call("GET", "/api/plans/"+itoa(int64(pid))+"/topics", ""); got.Code != 404 {
		t.Fatalf("deleted plan still visible: %d", got.Code)
	}
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
