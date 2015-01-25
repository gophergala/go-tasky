package tasky

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

const (
	Enabled  = "Enabled"
	Disabled = "Disabled"
)

type Action uint64

const (
	Cancel Action = iota
	Pause
	Resume
	Restart
)

//WorkerDetails is an object that describes the worker
type WorkerDetails struct {
	// Name of the worker
	Name string `json:"name"`

	// human readable description of the worker
	Description string `json:"description"`

	//an config struct for all config values
	Config interface{} `json:"config_objects,omitempty"`
}

type Worker interface {
	// Details provides all metadata info about the worker and it's config
	Details() *WorkerDetails

	// Worker name
	Name() string

	// Description of the worker and it's usage
	Usage() string

	// Execute the task
	Perform([]byte, chan []byte, chan error, chan bool)

	// Worker status
	Status() string

	// Action to be taken on ongoing task
	Signal(Action) bool

	// Maximum number of simultaneous tasks allowed
	MaxNumTasks() uint64
}

type TaskyError struct {
	Error string
}

var (
	wMut    sync.RWMutex
	workers map[string]Worker

	tMut  sync.RWMutex
	tasks map[string]*taskyTask

	apiBase string
)

func init() {
	workers = make(map[string]Worker)
	tasks = make(map[string]*taskyTask)
	apiBase = "/tasky/v1"
}

func NewWorker(w Worker) (Worker, error) {
	tw := &taskyWorker{}
	tw.w = w

	name := strings.ToLower(w.Name())

	wMut.Lock()
	workers[name] = tw
	wMut.Unlock()

	return tw, nil
}

func uuid() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x%x%x%x%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type taskid struct {
	TaskId string
}

type tstat struct {
	TaskId   string
	Status   string
	Duration string `json:"Duration,omitempty"`
}

func handlerGetTaskOutput(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("id: ", id)

	tMut.RLock()
	t, ok := tasks[id]
	tMut.RUnlock()

	if !ok {
		e := TaskyError{"Could not found a task with given id"}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	output := string(t.result())

	if len(output) > 0 {
		js := make(map[string]interface{})

		js["TaskId"] = id
		js["Output"] = output

		jsonStr, err := json.Marshal(js)
		if err != nil {
			e := TaskyError{err.Error()}
			estr, _ := json.Marshal(e)
			log.Println("estr: ", estr)
			fmt.Fprintf(rw, "%s\n", estr)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(rw, "%s\n", string(jsonStr))
	}
}

func handlerCancelTask(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("id: ", id)

	tMut.RLock()
	t, ok := tasks[id]
	tMut.RUnlock()

	if !ok {
		e := TaskyError{"Could not found a task with given id"}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	t.cancel()
}

func handlerGetTaskStatus(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Println("id: ", id)

	tMut.RLock()
	t, ok := tasks[id]
	tMut.RUnlock()

	if !ok {
		e := TaskyError{"Could not found a task with given id"}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	s := t.status()

	var durStr string

	if s != Running {
		d := t.duration()
		if d > 0 {
			durStr = fmt.Sprintf("%v", d)
		}
	}

	ts := tstat{id, s, durStr}
	log.Println("ts: ", ts)
	jsonStr, err := json.Marshal(ts)
	if err != nil {
		e := TaskyError{err.Error()}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s\n", string(jsonStr))
}

func newTask(w Worker, job []byte) taskid {
	id := uuid()

	t := &taskyTask{}
	t.new(w)

	tMut.Lock()
	tasks[id] = t
	tMut.Unlock()

	go t.run(job)

	return taskid{id}
}

type ts struct {
	Tasks []tstat
}

func listTasks() ts {
	t := ts{}

	tMut.RLock()
	for k, v := range tasks {
		s := v.status()
		log.Println("s: ", s)

		durStr := ""

		if s != Running {
			d := v.duration()
			log.Println("d: ", d)
			if d > 0 {
				durStr = fmt.Sprintf("%v", d)
			}
		}

		if len(t.Tasks) <= 0 {
			t.Tasks = make([]tstat, 1)
			t.Tasks[0] = tstat{k, s, durStr}
		} else {
			t.Tasks = append(t.Tasks, tstat{k, s, durStr})
		}
	}
	tMut.RUnlock()

	return t
}

func handlerListTasks(rw http.ResponseWriter, r *http.Request) {
	t := listTasks()
	log.Println("tasks: ", t)

	jsonStr, err := json.Marshal(t)
	if err != nil {
		e := TaskyError{err.Error()}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s\n", jsonStr)
}

func handlerNewTask(rw http.ResponseWriter, r *http.Request) {
	job, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e := TaskyError{err.Error()}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("job: ", job)

	name := mux.Vars(r)["name"]
	log.Println("name: ", name)

	wMut.RLock()
	w, ok := workers[name]
	wMut.RUnlock()

	if !ok {
		e := TaskyError{"Could not found worker with given name"}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	t := newTask(w, job)
	jsonStr, err := json.Marshal(t)
	if err != nil {
		e := TaskyError{err.Error()}
		estr, _ := json.Marshal(e)
		log.Println("estr: ", estr)
		fmt.Fprintf(rw, "%s\n", estr)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s\n", string(jsonStr))
}

func handlerListWorkers(rw http.ResponseWriter, r *http.Request) {
	// jsonStr, _ := listWorkers()

	listWorkers, err := listWorkerDetails()
	if err != nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	RespondJSON(rw, r, listWorkers)
	return

	// fmt.Fprintf(rw, "%s\n", jsonStr)
}

func RegisterTaskyHandlers(r *mux.Router) {
	r.StrictSlash(true) //enables matching a route with or without a trailing slash
	//Handles /tasky/v1 routes. Create new subrouters off this
	tr := r.PathPrefix(apiBase).Subrouter()

	workersRtr := tr.PathPrefix("/workers").Subrouter()
	workersRtr.HandleFunc("/", handlerListWorkers).Methods("GET")
	workersRtr.HandleFunc("/{name}", handlerNewTask).Methods("POST")

	tasksRtr := tr.PathPrefix("/tasks").Subrouter()
	tasksRtr.HandleFunc("/", handlerListTasks).Methods("GET")
	// tasksRtr.HandleFunc("/{id:[0-9a-f]+}", handlerGetTaskInfo).Methods("GET")
	tasksRtr.HandleFunc("/{id:[0-9a-f]+}/status", handlerGetTaskStatus).Methods("GET")
	tasksRtr.HandleFunc("/{id:[0-9a-f]+}/cancel", handlerCancelTask).Methods("POST")
	tasksRtr.HandleFunc("/{id:[0-9a-f]+}/result", handlerGetTaskOutput).Methods("GET")
}

func RespondJSON(w http.ResponseWriter, req *http.Request, v interface{}, code ...int) error {
	if code != nil {
		w.WriteHeader(code[0])
	}
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		//encoding failed
		panic(err)
	}
	return err
}
