package tasky

// import (
// 	"encoding/json"
// 	"github.com/gorilla/mux"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"time"
// )

// type Service struct {
// 	// This is the full URI reference to the service and the base of all routes in tasky
// 	URI *url.URL

// 	//gorilla mux router for now
// 	router *mux.Router

// 	// resources stores a map of all managed resources
// 	resources []*Resource

// 	//uptime is when the service started
// 	uptime time.Time
// }

// //getPath returns the base path of the service, sub = subpath segment to append to path, absolute = bool return absolute path
// func (svc *Service) getPath(sub string, absolute bool) string {
// 	path := svc.URI.Path
// 	if absolute {
// 		path = svc.URI.String()
// 	}

// 	if sub != "" {
// 		path += sub
// 	}

// 	return path
// }

// func (svc *Service) baseHandler(w http.ResponseWriter, req *http.Request) {
// 	resources := make(map[string]string)
// 	for _, v := range svc.resources {
// 		resources[v.name] = svc.getPath(v.name, true)
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	RespondJSON(w, req, resources)
// }

// func NewService(uri string) *Service {
// 	u, err := url.Parse(uri)
// 	if err != nil {
// 		log.Panicln("Tasky: unable to parse Service URI: ", err.Error())
// 	}

// 	//add a trailing / so we dont get weird matches
// 	if u.Path == "" || u.Path[len(u.Path)-1] != '/' {
// 		u.Path += "/"
// 	}
// 	rtr := mux.NewRouter()
// 	svc := &Service{
// 		URI:       u,
// 		router:    rtr.PathPrefix(u.String()).Subrouter(),
// 		resources: make([]*Resource, 0),
// 		uptime:    time.Now(),
// 	}

// 	//set up default service routes, i.e for /tasky/v1/
// 	svc.router.HandleFunc("", svc.baseHandler).Methods("GET")

// 	// svc.router.
// 	log.Printf("Tasky: New Service %q", u.String())

// 	return svc
// }

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
