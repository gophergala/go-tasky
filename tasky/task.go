package tasky

import (
	"fmt"
	"time"
)

const (
	Started   = "Started"
	Running   = "Running"
	Canceled  = "Canceled"
	Completed = "Completed"
	Failed    = "Failed"
)

type taskyTask struct {
	id   string
	w    Worker
	stat string
	quit chan bool
	out  []byte
	err  error
	dur  time.Duration
}

func (t *taskyTask) new(w Worker) {
	t.id = uuid()
	t.w = w
	t.stat = Started
	t.quit = nil
	t.out = nil
	t.err = nil
	t.dur = time.Duration(0)
}

func (t *taskyTask) run(job []byte) {
	if t.stat != Started {
		return
	}

	t.stat = Running
	t.quit = make(chan bool)
	q := make(chan bool)
	data := make(chan []byte)
	err := make(chan error)

	start := time.Now()

	go t.w.Perform(job, data, err, q)

	select {
	case o, ok := <-data:
		if ok {
			t.out = o
			t.dur = time.Since(start)
			t.stat = Completed
		} else {
			t.err = fmt.Errorf("Received invalid output from the worker.")
			t.stat = Failed
		}
		return

	case e := <-err:
		t.err = e
		t.stat = Failed
		return

	case <-t.quit:
		q <- true
		t.err = fmt.Errorf("Tasked canceled")
		t.stat = Canceled
		return
	}
}

func (t *taskyTask) status() string {
	return t.stat
}

func (t *taskyTask) cancel() {
	if t.stat == Running {
		t.quit <- true
	}
}

func (t *taskyTask) result() []byte {
	return t.out
}

func (t *taskyTask) error() error {
	return t.err
}

func (t *taskyTask) duration() time.Duration {
	return t.dur
}
