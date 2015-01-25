package main

import (
	"time"

	"github.com/gophergala/go-tasky/tasky"
)

type Sleeper struct {
}

func (d *Sleeper) Info() []byte {
	s := `{
		"Usage": {
		}
	}`

	return []byte(s)
}

func (d *Sleeper) Services() []byte {
	return nil
}

func (d *Sleeper) Perform(job []byte, dataCh chan []byte, errCh chan error, quitCh chan bool) {
	done := make(chan bool)
	go func() {
		time.Sleep(5 * time.Minute)
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

func (d *Sleeper) Status() []byte {
	return nil
}

func (d *Sleeper) Signal(act tasky.Action) bool {
	return true
}

func (d *Sleeper) Statistics() []byte {
	return nil
}
