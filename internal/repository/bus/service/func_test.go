package bus

import (
	"testing"
	"time"

	"flip/internal/model"
	"flip/internal/repository/bus"
)

func Test_Publish(t *testing.T) {
	type fields struct {
		b bus.Repository
	}
	type args struct {
		uploadID string
		tx       model.Transaction
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantUploadID string
		wantStatus   string
	}{
		{
			name: "publish event ok",
			fields: fields{
				b: NewDefault(),
			},
			args: args{
				uploadID: "u1",
				tx:       model.Transaction{Status: "FAILED"},
			},
			wantUploadID: "u1",
			wantStatus:   "FAILED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.b.Publish(tt.args.uploadID, tt.args.tx)

			select {
			case ev := <-tt.fields.b.Channel():
				if ev.UploadID != tt.wantUploadID {
					t.Fatalf("got uploadID=%v want=%v", ev.UploadID, tt.wantUploadID)
				}
				if ev.Tx.Status != tt.wantStatus {
					t.Fatalf("got status=%v want=%v", ev.Tx.Status, tt.wantStatus)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("no event received")
			}
		})
	}
}

func Test_ChannelNonBlocking(t *testing.T) {
	type fields struct {
		b bus.Repository
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "channel open and receives",
			fields: fields{b: NewDefault()},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.fields.b.Channel()

			select {
			case <-ch:
			case <-time.After(50 * time.Millisecond):
			}
		})
	}
}
