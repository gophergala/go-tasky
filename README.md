#Go-Tasky
Go-Tasky is a simple go package that makes it easy to execute tasks on your server that are exposed with a RESTful api. Developed for the gophergala global hackathon 2015.  

##Basic Principles

Register a task with go-tasky that implements the worker interface and a set of routes will be created for that task. You can have any number of uniquely configured tasks for a worker. Go-Tasky runs on your server and expose api endpoints. 

##Getting Started
Install Go-Tasky in your $GOPATH with `go get`  
```Go
go get -u github.com/gophergala/go-tasky
```
Navigate to the examples directory and run  
```Go
// use the -addr flag to set server address(default is :8888) ex. -addr=:4444
go run main.go 
```

##Workers
- copy a file to a new location  
- read the contents of a file  
- check the values of system settings 
- !TODO add more!  

##Use Case Example
Use this tool when debugging or checking values on your system.

##Routes
Workers:  
POST /tasky/v1/worker/register - registers a worker with go-tasky   
GET /tasky/v1/workers - returns a list of available worker endpoints   
GET /tasky/v1/workers/{worker_name} - returns a list of available tasks to run  
GET /tasky/v1/workers/{worker_name}/info - returns a description of the worker and it's usage   
POST /tasky/v1/workers/{worker_name} - Creates a new task to run with the worker and returns a unique task id  
GET /tasky/v1/{worker_name}/statistics - returns statistics for the worker like number of tasks performed, failure rate, avearge time take per task etc

Tasks:  
GET /tasky/v1/task/{task_id} - Fetch details of a single task.  
PUT /tasky/v1/task/{task_id} - Update the configuration of the task.  

POST /tasky/v1/task/{task_id}/actions - Modify the state of the task (cancel, pause, resume, run)  
GET /tasky/v1/task/{task_id}/status - returns the status of the task  
GET /tasky/v1/task/{task_id}/statistics - returns the statistics about the task, such as time to complete task  


PUT /tasky/v1/task/{task_id} - Modify the state of the task (cancel, pause, resume, rerun)   
GET /tasky/v1/task/{task_id}/status - returns the status of the task   
GET /tasky/v1/task/{task_id}/statistics - returns the statistics about a task like time it took to complete the task etc   

RuleChains:  
For later, but used to chain multiple tasks together in an ordered fashion.  

## Worker Interface
The worker interface corresponds to an individual worker type
```go
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
	// avaerage time per task etc
	Statistics() string
}
```


## Task Interface
The task interface corresponds to details about a specific task running with a worker.  
```go
type Task interface {
    // logic to create a new task, called from the worker create endpoint
    Create()

    // returns the details of an individual task 
    Read(*Context)

    // modify the configuration of an individual task
    Update(*Context)

    // delete an individual task
    Delete(*Context)

}

```
