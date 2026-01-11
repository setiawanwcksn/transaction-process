package api

import (
	"flip/util"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := r.Context()
	uploadID := uuid.NewString()

	util.Log.InfoContext(ctx, "upload received", "upload_id", uploadID)
	util.Perf(ctx, "upload", "upload_id", uploadID)

	if h.async {
		go func() {
			util.Log.InfoContext(ctx, "parser goroutine start", "upload_id", uploadID)
			_ = h.parser.Process(ctx, uploadID, file)
			util.Log.InfoContext(ctx, "parser goroutine end", "upload_id", uploadID)
		}()
	} else {
		err = h.parser.Process(ctx, uploadID, file)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"upload_id":"` + uploadID + `"}`))
}
