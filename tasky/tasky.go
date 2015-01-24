package tasky

type Action uint64

const (
	Cancel Action = iota
	Pause
	Resume
	Restart
)

type Worker interface {
	// Description of the worker and it's usage
	Info() string

	// List of available tasks it can service
	Services() string

	// Execute the task
	Perform() (string, bool)

	// Worker status
	Status() string

	// Action to be taken on ongoing task
	Signal(Action) bool

	// Worker statistics like number of tasks performed, failure rate,
	// average time per task etc
	Statistics() string
}

type WorkerInst struct {
	// SHA1 Id
	Id string
}

func Register(Worker) (*WorkerInst, error) {
	return nil, nil
}

func (w *WorkerInst) ListenAndServe() error {
	return nil
}
