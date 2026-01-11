package worker

import "context"

// mockgen -source=internal/repository/worker/worker.go -destination=internal/repository/worker/mock_worker.go -package=worker
type Repository interface {
	Start(ctx context.Context, workers int)
	Stop()
}
