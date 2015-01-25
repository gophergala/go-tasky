package tasky

import (
	"fmt"
	"testing"
)

type dumyWorker struct {
}

func (d *dumyWorker) Info() []byte {
	return nil
}

func (d *dumyWorker) Services() []byte {
	return nil
}

func (d *dumyWorker) Perform(job []byte) ([]byte, bool) {
	return nil, true
}

func (d *dumyWorker) Status() []byte {
	return nil
}

func (d *dumyWorker) Signal(a Action) bool {
	return true
}

func (d *dumyWorker) Statistics() []byte {
	return nil
}

func TestNewWorker(t *testing.T) {
	d := &dumyWorker{}

	tw, err := NewWorker(d)
	if err != nil {
		t.Error("expected NewWorker() to work got error %v", err)
	}

	fmt.Println(tw)
}
