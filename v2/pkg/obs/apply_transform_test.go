package obs_test

import (
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/obs"
)

type stubObserver struct {
	wrapped bool
}

func (s *stubObserver) Observe() any {
	s.wrapped = true
	return s
}

func TestApplyTransform_observer(t *testing.T) {
	stub := &stubObserver{}
	out := obs.ApplyTransform(stub)
	if !stub.wrapped {
		t.Fatal("expected Observe to be called")
	}
	if out != stub {
		t.Fatal("expected Observe return value")
	}
}

func TestApplyTransform_passthrough(t *testing.T) {
	if got := obs.ApplyTransform(42); got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
	v := struct{ n int }{n: 1}
	if got := obs.ApplyTransform(v); got != v {
		t.Fatalf("got %v, want %v", got, v)
	}
}
