package bus

import "flip/internal/model"

// mockgen -source=internal/repository/bus/bus.go -destination=internal/repository/bus/mock_bus.go -package=bus
type Repository interface {
	Publish(uploadID string, tx model.Transaction)
	Channel() <-chan model.Event
}
