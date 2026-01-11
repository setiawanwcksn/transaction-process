package api

import (
	"flip/util"
	"fmt"
	"net/http"
)

func (h *Handler) balance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("upload_id")

	defer func() {
		util.Log.InfoContext(ctx, "balance requested", "upload_id", id)
		util.Perf(ctx, "balance", "upload_id", id)
	}()

	if id == "" {
		http.Error(w, "missing upload_id", http.StatusBadRequest)
		return
	}

	bal, err := h.storage.Balance(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"balance":%d}`, bal)))
}
