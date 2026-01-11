package worker

import (
	"context"
	"testing"
	"time"

	"flip/internal/model"
)

type spyHandler struct {
	calls int
	fn    func(context.Context, model.Event) *HandleResult
}

func (s *spyHandler) handle(ctx context.Context, ev model.Event) *HandleResult {
	s.calls++
	return s.fn(ctx, ev)
}

func Test_WorkerLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type fields struct {
		events  chan model.Event
		handler *spyHandler
	}
	type args struct {
		workers int
		ev1     model.Event
		ev2     *model.Event
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCalls int
	}{
		{
			name: "FAILED handled once without retry",
			fields: func() fields {
				evch := make(chan model.Event, 1)
				spy := &spyHandler{
					fn: func(ctx context.Context, ev model.Event) *HandleResult {
						return &HandleResult{Success: false, Retryable: false}
					},
				}
				return fields{evch, spy}
			}(),
			args: args{
				workers: 1,
				ev1: model.Event{
					UploadID: "w1",
					Tx:       model.Transaction{Timestamp: time.Now(), Status: model.StatusFailed},
				},
			},
			wantCalls: 1,
		},
		{
			name: "PENDING succeeds on first try",
			fields: func() fields {
				evch := make(chan model.Event, 1)
				spy := &spyHandler{
					fn: func(ctx context.Context, ev model.Event) *HandleResult {
						return &HandleResult{Success: true}
					},
				}
				return fields{evch, spy}
			}(),
			args: args{
				workers: 1,
				ev1: model.Event{
					UploadID: "w2",
					Tx:       model.Transaction{Timestamp: time.Now(), Status: model.StatusPending},
				},
			},
			wantCalls: 1,
		},
		{
			name: "PENDING retryable then success",
			fields: func() fields {
				evch := make(chan model.Event, 1)
				call := 0
				spy := &spyHandler{
					fn: func(ctx context.Context, ev model.Event) *HandleResult {
						call++
						if call == 1 {
							return &HandleResult{Success: false, Retryable: true}
						}
						return &HandleResult{Success: true}
					},
				}
				return fields{evch, spy}
			}(),
			args: args{
				workers: 1,
				ev1: model.Event{
					UploadID: "w3",
					Tx:       model.Transaction{Timestamp: time.Now(), Status: model.StatusPending},
				},
			},
			wantCalls: 2,
		},
		{
			name: "Duplicate event ignored by tracker",
			fields: func() fields {
				evch := make(chan model.Event, 2)
				spy := &spyHandler{
					fn: func(ctx context.Context, ev model.Event) *HandleResult {
						return &HandleResult{Success: true}
					},
				}
				return fields{evch, spy}
			}(),
			args: args{
				workers: 1,
				ev1: model.Event{
					UploadID: "w4",
					Tx:       model.Transaction{Timestamp: time.Now(), Status: model.StatusFailed},
				},
				ev2: &model.Event{
					UploadID: "w4",
					Tx:       model.Transaction{Timestamp: time.Now(), Status: model.StatusFailed},
				},
			},
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			events := tt.fields.events
			handler := tt.fields.handler

			w := New(events, handler.handle)

			// override sleep so UT runs fast and predictable
			w.sleep = func(time.Duration) {}

			w.Start(ctx, tt.args.workers)

			events <- tt.args.ev1
			if tt.args.ev2 != nil {
				events <- *tt.args.ev2
			}

			// give scheduler time (no real sleep)
			time.Sleep(10 * time.Millisecond)

			w.Stop()

			if handler.calls != tt.wantCalls {
				t.Errorf("%s calls = %d, want %d",
					tt.name, handler.calls, tt.wantCalls)
			}
		})
	}
}
