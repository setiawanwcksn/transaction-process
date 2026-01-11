package storage

import (
	"context"
	"reflect"
	"testing"

	"flip/internal/model"
)

func TestRepository_AddSuccess(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		store *Repository
	}
	type args struct {
		id string
		tx model.Transaction
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantBal   int64
		wantLines int
	}{
		{
			name: "AddSuccess CREDIT should increase balance",
			fields: fields{
				store: NewStorage(),
			},
			args: args{
				id: "u1",
				tx: model.Transaction{
					Amount: 50000,
					Type:   "CREDIT",
					Status: model.StatusSuccess,
				},
			},
			wantBal:   50000,
			wantLines: 1,
		},
		{
			name: "AddSuccess DEBIT should decrease balance",
			fields: fields{
				store: NewStorage(),
			},
			args: args{
				id: "u2",
				tx: model.Transaction{
					Amount: 20000,
					Type:   "DEBIT",
					Status: model.StatusSuccess,
				},
			},
			wantBal:   -20000,
			wantLines: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.store.Create(ctx, tt.args.id)

			err := tt.fields.store.AddSuccess(ctx, tt.args.id, tt.args.tx)
			if err != nil {
				t.Errorf("AddSuccess() err = %v", err)
				return
			}

			stmt, _ := tt.fields.store.Get(ctx, tt.args.id)

			if stmt.Balance != tt.wantBal {
				t.Errorf("balance = %v, want %v", stmt.Balance, tt.wantBal)
			}
			if stmt.Lines != tt.wantLines {
				t.Errorf("lines = %v, want %v", stmt.Lines, tt.wantLines)
			}
		})
	}
}

func TestRepository_AddFailed(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		store *Repository
	}
	type args struct {
		id string
		tx model.Transaction
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantIssues int
		wantBal    int64
		wantLines  int
	}{
		{
			name: "AddFailed should append to issues and not change balance",
			fields: fields{
				store: NewStorage(),
			},
			args: args{
				id: "u3",
				tx: model.Transaction{
					Amount: 15000,
					Type:   "DEBIT",
					Status: model.StatusFailed,
				},
			},
			wantIssues: 1,
			wantBal:    0,
			wantLines:  1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.store.Create(ctx, tt.args.id)

			err := tt.fields.store.AddFailed(ctx, tt.args.id, tt.args.tx)
			if err != nil {
				t.Errorf("AddFailed() err = %v", err)
				return
			}

			stmt, _ := tt.fields.store.Get(ctx, tt.args.id)

			if !reflect.DeepEqual(len(stmt.Issues), tt.wantIssues) {
				t.Errorf("issues = %v, want %v", len(stmt.Issues), tt.wantIssues)
			}
			if stmt.Balance != tt.wantBal {
				t.Errorf("balance = %v, want %v", stmt.Balance, tt.wantBal)
			}
			if stmt.Lines != tt.wantLines {
				t.Errorf("lines = %v, want %v", stmt.Lines, tt.wantLines)
			}
		})
	}
}

func TestRepository_Finish(t *testing.T) {
	ctx := context.Background()
	store := NewStorage()
	store.Create(ctx, "u4")

	store.Finish(ctx, "u4")

	stmt, _ := store.Get(ctx, "u4")
	if !stmt.Done {
		t.Errorf("finish flag not set, got=%v", stmt.Done)
	}
}

func Test_IssuesFiltered(t *testing.T) {
	type args struct {
		ctx      context.Context
		id       string
		statuses []string
		offset   int
		limit    int
	}
	tests := []struct {
		name     string
		setup    func() *Repository
		args     args
		wantLen  int
		wantTot  int
		wantErr  bool
	}{
		{
			name: "filter FAILED only with pagination",
			setup: func() *Repository {
				s := NewStorage()
				s.Create(context.Background(), "u4")
				for i := 0; i < 5; i++ {
					s.AddFailed(context.Background(), "u4", model.Transaction{Status: "FAILED"})
				}
				return s
			},
			args:    args{context.Background(), "u4", []string{"FAILED"}, 2, 2},
			wantLen: 2,
			wantTot: 5,
		},
		{
			name: "missing id returns error",
			setup: func() *Repository { return NewStorage() },
			args: args{
				ctx:      context.Background(),
				id:       "x",
				statuses: []string{"FAILED"},
				offset:   0,
				limit:    10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			got, tot, err := s.IssuesFiltered(tt.args.ctx, tt.args.id, tt.args.statuses, tt.args.offset, tt.args.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("IssuesFiltered() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if tot != tt.wantTot {
				t.Errorf("total = %d, want %d", tot, tt.wantTot)
			}
			if len(got) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func Test_Balance(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		setup      func() *Repository
		args       args
		want       int64
		wantErr    bool
	}{
		{
			name: "compute correct balance",
			setup: func() *Repository {
				s := NewStorage()
				s.Create(context.Background(), "u1")
				s.AddSuccess(context.Background(), "u1", model.Transaction{Type: "CREDIT", Amount: 200, Status: "SUCCESS"})
				s.AddSuccess(context.Background(), "u1", model.Transaction{Type: "DEBIT", Amount: 50, Status: "SUCCESS"})
				return s
			},
			args:    args{context.Background(), "u1"},
			want:    150,
			wantErr: false,
		},
		{
			name: "missing upload returns error",
			setup: func() *Repository { return NewStorage() },
			args:    args{context.Background(), "x"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup()
			got, err := s.Balance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Balance() err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Balance = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Progress(t *testing.T) {
	type fields struct {
		setup func() *Repository
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantLen  int
		wantDone bool
		wantErr  bool
	}{
		{
			name: "returns progress with correct line count",
			fields: fields{setup: func() *Repository {
				s := NewStorage()
				s.Create(context.Background(), "u1")
				_ = s.AddSuccess(context.Background(), "u1", model.Transaction{Status: "SUCCESS"})
				_ = s.AddFailed(context.Background(), "u1", model.Transaction{Status: "FAILED"})
				return s
			}},
			args:     args{context.Background(), "u1"},
			wantLen:  2,
			wantDone: false,
			wantErr:  false,
		},
		{
			name: "still returns when done",
			fields: fields{setup: func() *Repository {
				s := NewStorage()
				s.Create(context.Background(), "u2")
				_ = s.AddFailed(context.Background(), "u2", model.Transaction{Status: "FAILED"})
				s.Finish(context.Background(), "u2")
				return s
			}},
			args:     args{context.Background(), "u2"},
			wantLen:  1,
			wantDone: true,
			wantErr:  false,
		},
		{
			name: "missing upload returns error",
			fields: fields{setup: func() *Repository {
				return NewStorage()
			}},
			args:    args{context.Background(), "x"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t := t
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.setup()
			stmt, err := s.Progress(tt.args.ctx, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Progress() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if stmt.Lines != tt.wantLen {
				t.Errorf("Lines = %d, want %d", stmt.Lines, tt.wantLen)
			}
			if stmt.Done != tt.wantDone {
				t.Errorf("Done = %v, want %v", stmt.Done, tt.wantDone)
			}
		})
	}
}
