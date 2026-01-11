package api

import (
	"flip/internal/model"
	"flip/internal/repository/storage"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_GetIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockRepository(ctrl)

	h := New(mockStorage, nil, nil)

	type fields struct {
		storage storage.Repository
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    []*gomock.Call
		want    int
		wantLen int
	}{
		{
			name: "returns filtered issues",
			fields: fields{
				storage: mockStorage,
			},
			args: args{
				url: "/transactions/issues?upload_id=u1&status=FAILED&limit=10&offset=0",
			},
			mock: []*gomock.Call{
				mockStorage.EXPECT().IssuesFiltered(gomock.Any(), "u1", []string{"FAILED"}, 0, 10).
					Return([]model.Transaction{}, 1, nil),
			},
			want:    http.StatusOK,
			wantLen: 1,
		},
		{
			name: "missing upload_id returns 400",
			fields: fields{
				storage: mockStorage,
			},
			args: args{
				url: "/transactions/issues?status=FAILED",
			},
			want: http.StatusBadRequest,
		},
		{
			name: "storage error returns 500",
			fields: fields{
				storage: mockStorage,
			},
			args: args{
				url: "/transactions/issues?upload_id=u1",
			},
			mock: []*gomock.Call{
				mockStorage.EXPECT().
					IssuesFiltered(gomock.Any(), "u1", []string{"FAILED", "PENDING"}, 0, 50).
					Return(nil, 0, fmt.Errorf("boom")),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "invalid upload_id returns 400",
			fields: fields{
				storage: mockStorage,
			},
			args: args{
				url: "/transactions/issues?limit=1",
			},
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.mock {
				_ = m
			}
			mux := http.NewServeMux()
			h.Register(mux)
			req := httptest.NewRequest(http.MethodGet, tt.args.url, nil)
			rec := httptest.NewRecorder()

			h.issues(rec, req)

			if rec.Code != tt.want {
				t.Fatalf("status=%v want=%v", rec.Code, tt.want)
			}
		})
	}
}
