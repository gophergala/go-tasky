package workers

import (
	"time"

	"github.com/gophergala/go-tasky/tasky"
)

type Sleeper struct {
}

func (d *Sleeper) Name() string {
	return "Sleeper"
}

func (d *Sleeper) Usage() string {
	s := "{\"Usage\":{}}"

	return s
}

func (d *Sleeper) Perform(job []byte, dataCh chan []byte, errCh chan error, quitCh chan bool) {
	done := make(chan bool)
	go func() {
		time.Sleep(1 * time.Minute)
		dataCh <- []byte("Done sleeping.")
		done <- true
	}()

	select {
	case <-done:
		return

	case <-quitCh:
		return
	}
}

func (d *Sleeper) Status() string {
	return tasky.Enabled
}

func (d *Sleeper) Signal(act tasky.Action) bool {
	return true
}

func (d *Sleeper) MaxNumTasks() uint64 {
	return 10
}
