package memory

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/whywaita/myshoes/pkg/datastore"
)

// Memory is implement datastore on-memory
type Memory struct {
	mu      *sync.RWMutex
	targets map[uuid.UUID]datastore.Target
	jobs    map[uuid.UUID]datastore.Job
	runners map[uuid.UUID]datastore.Runner
}

// New create map
func New() (*Memory, error) {
	m := &sync.RWMutex{}
	t := map[uuid.UUID]datastore.Target{}
	j := map[uuid.UUID]datastore.Job{}
	r := map[uuid.UUID]datastore.Runner{}

	return &Memory{
		mu:      m,
		targets: t,
		jobs:    j,
		runners: r,
	}, nil
}

// CreateTarget create a target
func (m *Memory) CreateTarget(ctx context.Context, target datastore.Target) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.targets[target.UUID] = target
	return nil
}

// GetTarget get a target
func (m *Memory) GetTarget(ctx context.Context, id uuid.UUID) (*datastore.Target, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.targets[id]
	if !ok {
		return nil, datastore.ErrNotFound
	}
	return &t, nil
}

// GetTargetByScope get a target from scope
func (m *Memory) GetTargetByScope(ctx context.Context, scope string) (*datastore.Target, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.targets {
		if t.Scope == scope {
			// found
			return &t, nil

		}
	}

	return nil, datastore.ErrNotFound
}

// ListTargets get a all targets
func (m *Memory) ListTargets(ctx context.Context) ([]datastore.Target, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var targets []datastore.Target

	for _, t := range m.targets {
		targets = append(targets, t)
	}

	return targets, nil
}

// DeleteTarget delete a target
func (m *Memory) DeleteTarget(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.targets, id)
	return nil
}

// UpdateTargetStatus update status in target
func (m *Memory) UpdateTargetStatus(ctx context.Context, targetID uuid.UUID, newStatus datastore.TargetStatus, description string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.targets[targetID]
	if !ok {
		return fmt.Errorf("not found")
	}

	t.Status = newStatus
	if description != "" {
		t.StatusDescription.Valid = true
	} else {
		t.StatusDescription.Valid = false
	}
	t.StatusDescription.String = description

	m.targets[targetID] = t

	return nil
}

// UpdateToken update token in target
func (m *Memory) UpdateToken(ctx context.Context, targetID uuid.UUID, newToken string, newExpiredAt time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.targets[targetID]
	if !ok {
		return fmt.Errorf("not found")
	}
	t.GitHubToken = newToken
	t.TokenExpiredAt = newExpiredAt

	m.targets[targetID] = t
	return nil
}

// UpdateTargetParam update parameter of target
func (m *Memory) UpdateTargetParam(ctx context.Context, targetID uuid.UUID, newResourceType datastore.ResourceType, newProviderURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.targets[targetID]
	if !ok {
		return fmt.Errorf("not found")
	}
	t.ResourceType = newResourceType
	t.ProviderURL = sql.NullString{
		String: newProviderURL,
		Valid:  true,
	}

	m.targets[targetID] = t
	return nil
}

// EnqueueJob add a job
func (m *Memory) EnqueueJob(ctx context.Context, job datastore.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobs[job.UUID] = job
	return nil
}

// ListJobs get all jobs
func (m *Memory) ListJobs(ctx context.Context) ([]datastore.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var jobs []datastore.Job
	for _, j := range m.jobs {
		jobs = append(jobs, j)
	}

	return jobs, nil
}

// DeleteJob delete a job
func (m *Memory) DeleteJob(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.jobs, id)
	return nil
}

// CreateRunner add a runner
func (m *Memory) CreateRunner(ctx context.Context, runner datastore.Runner) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.runners[runner.UUID] = runner

	return nil
}

// ListRunners get a all runners
func (m *Memory) ListRunners(ctx context.Context) ([]datastore.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var runners []datastore.Runner
	for _, r := range m.runners {
		runners = append(runners, r)
	}

	return runners, nil
}

// ListRunnersByTargetID get a not deleted runners that has target_id
func (m *Memory) ListRunnersByTargetID(ctx context.Context, targetID uuid.UUID) ([]datastore.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var runners []datastore.Runner
	for _, r := range m.runners {
		if uuid.Equal(r.TargetID, targetID) {
			runners = append(runners, r)
		}
	}

	return runners, nil
}

// ListRunnersLogByUntil ListRunnerLog get a runners until time
func (m *Memory) ListRunnersLogByUntil(ctx context.Context, until time.Time) ([]datastore.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var runners []datastore.Runner
	for _, r := range m.runners {
		if r.CreatedAt.Before(until) {
			runners = append(runners, r)
		}
	}

	return runners, nil
}

// GetRunner get a runner
func (m *Memory) GetRunner(ctx context.Context, id uuid.UUID) (*datastore.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.runners[id]
	if !ok {
		return nil, datastore.ErrNotFound
	}

	return &r, nil
}

// DeleteRunner delete a runner
func (m *Memory) DeleteRunner(ctx context.Context, id uuid.UUID, deletedAt time.Time, reason datastore.RunnerStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.runners, id)
	return nil
}

// GetLock get lock
func (m *Memory) GetLock(ctx context.Context) error {
	return nil
}

// IsLocked return status of lock
func (m *Memory) IsLocked(ctx context.Context) (string, error) {
	return datastore.IsNotLocked, nil
}
