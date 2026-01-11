package worker

import (
	"context"
	"testing"

	"flip/internal/model"
)

func Test_New(t *testing.T) {
	// ctx := context.Background()
	busCh := make(chan model.Event)

	// called := false
	h := func(ctx context.Context, ev model.Event) *HandleResult {
		// called = true
		return &HandleResult{Success: true}
	}

	type fields struct {
		bus     <-chan model.Event
		handler func(context.Context, model.Event) *HandleResult
	}

	type args struct{}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "worker new initializes internal state",
			fields: fields{
				bus:     busCh,
				handler: h,
			},
			args: args{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := New(tt.fields.bus, tt.fields.handler)

			if w.bus != tt.fields.bus {
				t.Errorf("bus not wired")
			}
			if w.handler == nil {
				t.Errorf("handler not set")
			}
			if w.stop == nil {
				t.Errorf("stop chan not created")
			}
			if w.tracker == nil || w.tracker.seen == nil {
				t.Errorf("tracker not initialized")
			}
		})
	}
}

func Test_StartDefault(t *testing.T) {
	// ctx := context.Background()
	busCh := make(chan model.Event)

	// called := false
	h := func(ctx context.Context, ev model.Event) *HandleResult {
		// called = true
		return &HandleResult{Success: true}
	}

	type fields struct {
		bus     <-chan model.Event
		handler func(context.Context, model.Event) *HandleResult
	}

	type args struct{}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "start default wires worker and starts goroutine",
			fields: fields{
				bus:     busCh,
				handler: h,
			},
			args: args{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := StartDefault(tt.fields.bus, tt.fields.handler)
			defer w.Stop()
			if w == nil {
				t.Fatalf("worker nil")
			}
			if w.bus != tt.fields.bus {
				t.Errorf("bus not wired")
			}
			if w.handler == nil {
				t.Errorf("handler nil")
			}

		})
	}
}
