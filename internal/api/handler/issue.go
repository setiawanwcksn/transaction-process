package api

import (
	"encoding/json"
	"flip/internal/model"
	"flip/util"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) issues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("upload_id")
	statusStr := r.URL.Query().Get("status")

	if id == "" {
		http.Error(w, "missing upload_id", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	statuses := []string{"FAILED", "PENDING"}
	if statusStr != "" {
		parts := strings.Split(statusStr, ",")
		tmp := make([]string, 0)
		for _, p := range parts {
			val := strings.ToUpper(strings.TrimSpace(p))
			if val == "FAILED" || val == "PENDING" {
				tmp = append(tmp, val)
			}
		}
		if len(tmp) > 0 {
			statuses = tmp
		}
	}

	util.Log.InfoContext(ctx, "issues", "upload_id", id, "statuses", statuses)
	items, total, err := h.storage.IssuesFiltered(ctx, id, statuses, offset, limit)
	if err != nil {
		http.Error(w, "not found", http.StatusInternalServerError)
		return
	}

	util.Perf(ctx, "issues filtered", "upload_id", id, "statuses", statuses)

	resp := model.PaginatedResponse[model.Transaction]{
		Data: items,
		Pagination: model.Pagination{
			Offset:     offset,
			Limit:      limit,
			Total:      total,
			NextOffset: offset + limit,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
