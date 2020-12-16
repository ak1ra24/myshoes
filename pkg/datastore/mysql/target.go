package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/whywaita/myshoes/pkg/datastore"
)

// CreateTarget create a target
func (m *MySQL) CreateTarget(ctx context.Context, target datastore.Target) error {
	query := `INSERT INTO targets(uuid, scope, ghe_domain, github_personal_token, resource_type, runner_user) VALUES (?, ?, ?, ?, ?, ?)`
	if _, err := m.Conn.ExecContext(ctx, query, target.UUID, target.Scope, target.GHEDomain, target.GitHubPersonalToken, target.ResourceType, target.RunnerUser); err != nil {
		return fmt.Errorf("failed to execute INSERT query: %w", err)
	}

	return nil
}

// GetTarget get a target
func (m *MySQL) GetTarget(ctx context.Context, id uuid.UUID) (*datastore.Target, error) {
	var t datastore.Target
	query := fmt.Sprintf(`SELECT uuid, scope, ghe_domain, github_personal_token, resource_type, runner_user, created_at, updated_at FROM targets WHERE uuid = ?`)
	if err := m.Conn.GetContext(ctx, &t, query, id.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrNotFound
		}

		return nil, fmt.Errorf("failed to execute SELECT query: %w", err)
	}

	return &t, nil
}

// GetTargetByScope get a target from scope
func (m *MySQL) GetTargetByScope(ctx context.Context, gheDomain, scope string) (*datastore.Target, error) {
	var t datastore.Target
	query := fmt.Sprintf(`SELECT uuid, scope, ghe_domain, github_personal_token, resource_type, runner_user, created_at, updated_at FROM targets WHERE scope = "%s"`, scope)
	if gheDomain != "" {
		query = fmt.Sprintf(`%s AND ghe_domain = "%s"`, query, gheDomain)
	}
	if err := m.Conn.GetContext(ctx, &t, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datastore.ErrNotFound
		}

		return nil, fmt.Errorf("failed to execute SELECT query: %w", err)
	}

	return &t, nil
}

// DeleteTarget delete a target
func (m *MySQL) DeleteTarget(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM targets WHERE uuid = ?`
	if _, err := m.Conn.ExecContext(ctx, query, id.String()); err != nil {
		return fmt.Errorf("failed to execute DELETE query: %w", err)
	}

	return nil
}