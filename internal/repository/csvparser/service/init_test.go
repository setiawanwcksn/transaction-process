package csvparser

import (
	"flip/internal/repository/bus"
	"flip/internal/repository/storage"
	"flip/internal/repository/worker"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockRepository(ctrl)
	mockBus := bus.NewMockRepository(ctrl)
	mockWorker := worker.NewMockRepository(ctrl)

	type fields struct {
		storage storage.Repository
		bus     bus.Repository
		worker  worker.Repository
	}
	tests := []struct {
		name   string
		fields fields
		mock   []*gomock.Call
		want   Repository
	}{
		{
			name: "success",
			fields: fields{
				storage: mockStorage,
				bus:     mockBus,
				worker:  mockWorker,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.mock {
				_ = m
			}
			svc := NewParser(tt.fields.storage, tt.fields.bus, tt.fields.worker)
			if svc == nil {
				t.Fatalf("NewParser() returned nil")
			}
		})
	}
}
