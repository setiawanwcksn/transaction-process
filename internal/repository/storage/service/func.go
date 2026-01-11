package storage

import (
	"context"
	"flip/internal/model"
	"flip/util"
	"fmt"
	"strings"
)

func (s *Repository) Create(ctx context.Context, id string) {
	defer func() {
		util.Log.InfoContext(ctx, "storage create", "upload_id", id)
	}()

	s.mu.Lock()
	s.data[id] = &model.Statement{UploadID: id}
	s.mu.Unlock()
}

func (s *Repository) AddSuccess(ctx context.Context, id string, tx model.Transaction) error {
	defer func() {
		util.Log.InfoContext(ctx, "storage add success", "upload_id", id)
		util.Perf(ctx, "storage add success", "upload_id", id)
	}()

	s.mu.Lock()
	stmt, ok := s.data[id]
	if !ok {
		s.mu.Unlock()
		return ErrNotFound
	}

	stmt.Transactions = append(stmt.Transactions, tx)
	stmt.Lines++

	switch tx.Type {
	case "CREDIT":
		stmt.Balance += tx.Amount
	case "DEBIT":
		stmt.Balance -= tx.Amount
	}
	fmt.Println("worker commit", "amount", tx.Amount, "new_balance", stmt.Balance)

	s.mu.Unlock()
	return nil
}

func (s *Repository) AddFailed(ctx context.Context, id string, tx model.Transaction) error {
	defer func() {
		util.Log.InfoContext(ctx, "storage add failed", "upload_id", id)
		util.Perf(ctx, "storage add failed", "upload_id", id)
	}()

	s.mu.Lock()
	stmt, ok := s.data[id]
	if !ok {
		s.mu.Unlock()
		return ErrNotFound
	}

	stmt.Transactions = append(stmt.Transactions, tx)
	stmt.Lines++
	stmt.Issues = append(stmt.Issues, tx)

	s.mu.Unlock()
	return nil
}

func (s *Repository) Finish(ctx context.Context, id string) {
	defer func() {
		util.Log.InfoContext(ctx, "storage finish", "upload_id", id)
		util.Perf(ctx, "storage finish", "upload_id", id)
	}()

	s.mu.Lock()
	if stmt, ok := s.data[id]; ok {
		stmt.Done = true
	}
	s.mu.Unlock()
}

func (s *Repository) Get(ctx context.Context, id string) (*model.Statement, error) {
	defer func() {
		util.Log.InfoContext(ctx, "storage get", "upload_id", id)
	}()

	s.mu.RLock()
	stmt, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return stmt, nil
}

func (s *Repository) Balance(ctx context.Context, id string) (int64, error) {
	defer func() {
		util.Log.InfoContext(ctx, "storage balance", "upload_id", id)
	}()

	s.mu.RLock()
	stmt, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		return 0, ErrNotFound
	}
	return stmt.Balance, nil
}

func (s *Repository) Progress(ctx context.Context, id string) (*model.Statement, error) {
	s.mu.RLock()
	stmt, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return stmt, nil
}

func (s *Repository) IssuesFiltered(ctx context.Context, id string, statuses []string, offset, limit int) ([]model.Transaction, int, error) {
	s.mu.RLock()
	stmt, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		return nil, 0, ErrNotFound
	}

	set := make(map[string]bool)
	for _, st := range statuses {
		set[strings.ToUpper(st)] = true
	}

	filtered := make([]model.Transaction, 0)
	for _, tx := range stmt.Issues {
		if set[tx.Status] {
			filtered = append(filtered, tx)
		}
	}

	total := len(filtered)

	if offset >= total {
		return []model.Transaction{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}
