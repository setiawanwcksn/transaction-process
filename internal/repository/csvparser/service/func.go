package csvparser

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"flip/internal/model"
	"flip/util"
)

func (p *Repository) Process(ctx context.Context, uploadID string, r io.Reader) error {
	defer func() {
		util.Log.InfoContext(ctx, "csvparser process", "upload_id", uploadID)
	}()

	p.storage.Create(ctx, uploadID)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			continue
		}

		tsInt, _ := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		amount, _ := strconv.ParseInt(strings.TrimSpace(parts[3]), 10, 64)

		tx := model.Transaction{
			Timestamp:    time.Unix(tsInt, 0),
			Counterparty: strings.TrimSpace(parts[1]),
			Type:         strings.TrimSpace(parts[2]),
			Amount:       amount,
			Status:       strings.ToUpper(strings.TrimSpace(parts[4])),
			Description:  strings.TrimSpace(parts[5]),
		}

		switch tx.Status {
		case model.StatusSuccess:
			_ = p.storage.AddSuccess(ctx, uploadID, tx)
		case model.StatusFailed:
			_ = p.storage.AddFailed(ctx, uploadID, tx)
		case model.StatusPending:
			p.bus.Publish(uploadID, tx)
		}
	}

	p.storage.Finish(ctx, uploadID)
	return scanner.Err()
}
