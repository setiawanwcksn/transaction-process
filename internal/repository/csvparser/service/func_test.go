package csvparser

import (
	"bytes"
	"context"
	"testing"

	"flip/internal/repository/bus"
	"flip/internal/repository/storage"

	"github.com/golang/mock/gomock"
)

func Test_Process(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		storage *storage.MockRepository
		bus     *bus.MockRepository
	}
	type args struct {
		uploadID string
		content  string
	}

	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		args    args
		mock    func(f fields)
		wantErr bool
	}{
		{
			name: "SUCCESS stored only",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					storage: storage.NewMockRepository(ctrl),
					bus:     bus.NewMockRepository(ctrl),
				}
			},
			args: args{
				uploadID: "u1",
				content:  "1736530000,Alice,CREDIT,50000,SUCCESS,Salary\n",
			},
			mock: func(f fields) {
				f.storage.EXPECT().Create(gomock.Any(), "u1")
				f.storage.EXPECT().AddSuccess(gomock.Any(), "u1", gomock.Any())
				f.storage.EXPECT().Finish(gomock.Any(), "u1")
			},
		},
		{
			name: "FAILED published",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					storage: storage.NewMockRepository(ctrl),
					bus:     bus.NewMockRepository(ctrl),
				}
			},
			args: args{
				uploadID: "u2",
				content:  "1736530200,Bob,DEBIT,20000,FAILED,Card declined\n",
			},
			mock: func(f fields) {
				f.storage.EXPECT().Create(gomock.Any(), "u2")
				f.bus.EXPECT().Publish("u2", gomock.Any())
				f.storage.EXPECT().Finish(gomock.Any(), "u2")
			},
		},
		{
			name: "PENDING published",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					storage: storage.NewMockRepository(ctrl),
					bus:     bus.NewMockRepository(ctrl),
				}
			},
			args: args{
				uploadID: "u3",
				content:  "1736530300,Charlie,CREDIT,35000,PENDING,Review\n",
			},
			mock: func(f fields) {
				f.storage.EXPECT().Create(gomock.Any(), "u3")
				f.bus.EXPECT().Publish("u3", gomock.Any())
				f.storage.EXPECT().Finish(gomock.Any(), "u3")
			},
		},
		{
			name: "SUCCESS + FAILED combined",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					storage: storage.NewMockRepository(ctrl),
					bus:     bus.NewMockRepository(ctrl),
				}
			},
			args: args{
				uploadID: "u4",
				content: "" +
					"1736530000,Alice,CREDIT,50000,SUCCESS,Salary\n" +
					"1736530200,Bob,DEBIT,20000,FAILED,Card declined\n",
			},
			mock: func(f fields) {
				f.storage.EXPECT().Create(gomock.Any(), "u4")
				f.storage.EXPECT().AddSuccess(gomock.Any(), "u4", gomock.Any())
				f.bus.EXPECT().Publish("u4", gomock.Any())
				f.storage.EXPECT().Finish(gomock.Any(), "u4")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := tt.fields(ctrl)

			if tt.mock != nil {
				tt.mock(f)
			}

			p := NewParser(f.storage, f.bus, nil)

			err := p.Process(ctx, tt.args.uploadID, bytes.NewBufferString(tt.args.content))
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
