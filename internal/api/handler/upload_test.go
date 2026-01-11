package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"flip/internal/repository/bus"
	"flip/internal/repository/csvparser"
	"flip/internal/repository/storage"

	"github.com/golang/mock/gomock"
)

func TestHandler_upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockRepository(ctrl)
	mockBus := bus.NewMockRepository(ctrl)
	mockParser := csvparser.NewMockRepository(ctrl)
	h := New(mockStorage, mockParser, mockBus)
	body, ctype := buildMultipartFile("file", "ts,c,CR,100,SUCCESS,desc\n")
	type fields struct {
		store  storage.Repository
		parser csvparser.Repository
		bus    bus.Repository
	}
	type args struct {
		method string
		body   io.Reader
		setup  func()
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "method not allowed",
			fields: fields{
				store:  mockStorage,
				parser: mockParser,
				bus:    mockBus,
			},
			args: args{
				method: http.MethodGet,
				body:   nil,
				setup:  func() {},
			},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name: "missing file returns bad request",
			fields: fields{
				store:  mockStorage,
				parser: mockParser,
				bus:    mockBus,
			},
			args: args{
				method: http.MethodPost,
				body:   nil,
				setup:  func() {},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "valid upload calls parser and returns upload_id",
			fields: fields{
				store:  mockStorage,
				parser: mockParser,
				bus:    mockBus,
			},
			args: args{
				method: http.MethodPost,
				body:   body,
				setup: func() {
					mockParser.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"upload_id":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.setup()

			req := httptest.NewRequest(tt.args.method, "/upload", tt.args.body)
			if ctype != "" {
				req.Header.Set("Content-Type", ctype)
			}
			w := httptest.NewRecorder()
			h.async = false
			h.upload(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantBody != "" {
				body := w.Body.String()
				if !reflect.DeepEqual(true, contains(body, tt.wantBody)) {
					t.Errorf("response = %s, want contains %s", body, tt.wantBody)
				}
			}
		})
	}
}

func buildMultipartFile(fieldname, content string) (io.Reader, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile(fieldname, "test.csv")
	fw.Write([]byte(content))
	w.Close()
	return &buf, w.FormDataContentType()
}

func contains(haystack, needle string) bool {
	return bytes.Contains([]byte(haystack), []byte(needle))
}
