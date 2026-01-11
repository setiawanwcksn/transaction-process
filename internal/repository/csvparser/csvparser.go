package csvparser

import (
	"context"
	"io"
)

// mockgen -source=internal/repository/csvparser/csvparser.go -destination=internal/repository/csvparser/mock_csvparser.go -package=csvparser
type Repository interface {
	Process(ctx context.Context, uploadID string, r io.Reader) error
}
