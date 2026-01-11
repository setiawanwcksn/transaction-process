package storage_test

import (
	"testing"

	storage "flip/internal/repository/storage/service"
)

func Test_NewStorage(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "storage init ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewStorage()
			if s == nil {
				t.Fatalf("NewStorage() returned nil")
			}
		})
	}
}
