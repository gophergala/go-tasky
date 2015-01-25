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

## Example Workers
- CopyFile - copy a file to a new location  
- Ifconfig - fetch network information read the contents of a file  
- sleeper -  Sleep and wait?
- check the values of system settings  

##Routes

Workers:  
GET /tasky/v1/workers/ - returns a list of available worker endpoints   
POST /tasky/v1/workers/{worker_name} - Creates a task to run.  

TODO:
GET /tasky/v1/workers/{worker_name}/info - returns a description of the worker and it's usage   
GET /tasky/v1/{worker_name}/statistics - returns statistics for the worker like number of tasks performed, failure rate, average time take per task etc  

Tasks:  
GET /tasky/v1/tasks/ - Returns a list of all tasks.  
GET /tasky/v1/tasks/{id}/status - Returns the status of a task. 
POST /tasky/v1/tasks/{id}/cancel - Cancel the task.   

TODO:
PUT /tasky/v1/task/{task_id} - Update the configuration of the task.  
POST /tasky/v1/task/{task_id}/actions - Modify the state of the task (cancel, pause, resume, run)  
GET /tasky/v1/task/{task_id}/statistics - returns the statistics about the task, such as time to complete task  

RuleChains:  
For later, but used to chain multiple tasks together in an ordered fashion.  

## Example
List of workers available:  
```go
curl http://localhost:8888/tasky/v1/workers/ | python -mjson.tool 
[
    {
        "config_objects": {
            "destination": "",
            "source": ""
        },
        "description": "CopyFile will copy a file on the server. You must specify the Source and Destination",
        "name": "CopyFile"
    },
    {
        "description": "Ifconfig will return networking details from the server. No config is needed for this worker",
        "name": "Ifconfig"
    }
]
```

## Worker Interface
Create a custom worker by implementing the worker interface:  
```go
type Worker interface {
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
```
