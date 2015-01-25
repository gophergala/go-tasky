package tasky

import (
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
)

type taskyWorker struct {
	w     Worker
	tasks []string
}

func (tw *taskyWorker) Details() *WorkerDetails {
	return tw.w.Details()
}

func (tw *taskyWorker) Name() string {
	return tw.w.Name()
}

func (tw *taskyWorker) Usage() string {
	return tw.w.Usage()
}

func (tw *taskyWorker) Perform(job []byte, dataCh chan []byte, errCh chan error, quitCh chan bool) {
	tw.w.Perform(job, dataCh, errCh, quitCh)
}

func (tw *taskyWorker) Status() string {
	return tw.w.Status()
}

func (tw *taskyWorker) Signal(act Action) bool {
	return tw.w.Signal(act)
}

func (tw *taskyWorker) MaxNumTasks() uint64 {
	return tw.w.MaxNumTasks()
}

type worker struct {
	Id    string `json:"name"`
	Usage string `json:"usage"`
}

type ws struct {
	Workers []worker `json:"workers,omitempty"`
}

func listWorkerDetails() ([]WorkerDetails, error) {
	wMut.RLock()
	if len(workers) == 0 {
		wMut.RUnlock()
		return nil, errors.New("No Registered workers were found.")
	}
	workersDetailList := make([]WorkerDetails, len(workers))
	i := 0
	for _, wdetails := range workers {
		spew.Dump(wdetails)
		workersDetailList[i] = *wdetails.Details()
		i++
	}
	wMut.RUnlock()

	return workersDetailList, nil
}

// methods for routes to invoke
func listWorkers() ([]byte, error) {
	w := ws{}

	wMut.RLock()
	for k, v := range workers {
		t := worker{}
		t.Id = k
		t.Usage = v.Usage()

		if len(w.Workers) <= 0 {
			w.Workers = make([]worker, 1)
			w.Workers[0] = t
		} else {
			w.Workers = append(w.Workers, t)
		}
	}
	wMut.RUnlock()

	jsonStr, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
