package storage

import (
	"context"
	"flip/internal/model"
)

// mockgen -source=internal/repository/storage/storage.go -destination=internal/repository/storage/mock_storage.go -package=storage
type Repository interface {
	Create(ctx context.Context, id string)
	AddSuccess(ctx context.Context, id string, tx model.Transaction) error
	AddFailed(ctx context.Context, id string, tx model.Transaction) error
	Finish(ctx context.Context, id string)
	Get(ctx context.Context, id string) (*model.Statement, error)
	Balance(ctx context.Context, id string) (int64, error)
	Progress(ctx context.Context, id string) (*model.Statement, error)
	IssuesFiltered(ctx context.Context, id string, statuses []string, offset, limit int) ([]model.Transaction, int, error)
}
