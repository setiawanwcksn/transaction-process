package bus

import (
	"flip/internal/model"

	"github.com/google/uuid"
)

func (b *Repository) Publish(uploadID string, tx model.Transaction) {
	b.ch <- model.Event{UploadID: uploadID, Tx: tx, ID: uuid.NewString()}
}

func (b *Repository) Channel() <-chan model.Event {
	return b.ch
}
