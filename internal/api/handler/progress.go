package api

import (
	"encoding/json"
	"flip/util"
	"net/http"
)

func (h *Handler) progress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("upload_id")

	if id == "" {
		http.Error(w, "missing upload_id", http.StatusBadRequest)
		return
	}

	stmt, err := h.storage.Progress(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := map[string]any{
		"upload_id": stmt.UploadID,
		"lines":     stmt.Lines,
		"issues":    len(stmt.Issues),
		"balance":   stmt.Balance,
		"done":      stmt.Done,
	}

	util.Perf(ctx, "progress", "upload_id", id)

	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
