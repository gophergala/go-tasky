package tasky

import (
	"fmt"
	// "github.com/gorilla/mux"
	// "encoding/json"
	// "net/http"
	// "strings"
	"sync"
)

/*
the resource package handles registering the custom workers and setting up the routes.
*/

var (
	workerStore Resources
)

func init() {
	workerStore = Resources{
		Store: map[string]*Worker{},
	}
}

type Resources struct {
	sync.RWMutex
	Store map[string]*Worker
}

func (r *Resources) RegisterWorker(workerName string, worker Worker) error {
	r.Lock()
	id := fmt.Sprintf("%d", len(r.Store)) //replace this with database backing id, temp for now
	r.Store[id] = &worker
	r.Unlock()

	return nil
}

// func RespondJSON(w http.ResponseWriter, req *http.Request, v interface{}, code ...int) error {
// 	if code != nil {
// 		w.WriteHeader(code[0])
// 	}
// 	err := json.NewEncoder(w).Encode(v)
// 	if err != nil {
// 		//encoding failed
// 		panic(err)
// 	}
// 	return err
// }

// type BaseResource interface {
// 	Index(http.ResponseWriter, *http.Request)
// }

// type Resource struct {
// 	router           *mux.Router
// 	name             string      // name of resource, discovered from workerCollection
// 	path             string      // path is the URI to this resource
// 	workerCollection interface{} //worker object that implements Index interface
// }

// func (r *Resource) getPath(sub string) string {
// 	if strings.Contains(sub, r.path) {
// 		return sub
// 	}
// 	path := r.path
// 	if t := strings.Trim(sub, "/"); t != "" {
// 		path += "/" + t
// 	}

// 	return path
// }

// func (svc *Service) Resource(collection BaseResource) *Resource {
// 	// reflect name from objects type
// 	cs := fmt.Sprintf("%T", collection)
// 	name := strings.ToLower(cs[strings.LastIndex(cs, ".")+1:])
// 	if name == "" {
// 		panic("Tasky: Worker Resource naming failed: " + cs)
// 	}

// 	res := &Resource{
// 		router:           svc.router.PathPrefix(name).Subrouter(),
// 		name:             name,
// 		path:             svc.getPath(name, false),
// 		workerCollection: collection,
// 	}

// 	//set up base routes
// 	res.router.HandleFunc("", collection.Index).Methods("GET")

// 	//update service resources list
// 	svc.resources = append(svc.resources, res)

// 	return res

// }
