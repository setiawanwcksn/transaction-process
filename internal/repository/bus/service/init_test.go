package bus

import "testing"

func TestNewDefault(t *testing.T) {
	svc := NewDefault()
	if svc == nil {
		t.Fatalf("NewDefault() error")
	}
}
