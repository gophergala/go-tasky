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
    },
    {
        "description": "Sleeps for a minute. This is to showcase long running tasks.",
        "name": "Sleeper"
    }
]
```

Submit a task:  
```go
curl -X POST -H "Content-Type: application/json" -d '{}' http://localhost:8888/tasky/v1/workers/sleeper | python -mjson.tool
{
    "TaskId": "fe965d7a2b14471c8339b555f0f51014"
}

curl -X POST -H "Content-Type: application/json" -d '{}' http://localhost:8888/tasky/v1/workers/sleeper | python -mjson.tool
{
    "TaskId": "594a8e2fad2247328521a9d85066c627"
}

curl -X POST -H "Content-Type: application/json" -d '{}' http://localhost:8888/tasky/v1/workers/sleeper | python -mjson.tool
{
    "TaskId": "96ac00d2f3c343e9a301cd9236497cf8"
}
```

List all tasks:  
```go
curl http://localhost:8888/tasky/v1/tasks/ | python -mjson.tool 
{
    "Tasks": [
        {
            "Duration": "1m0.000260498s",
            "Status": "Completed",
            "TaskId": "749ef2a6c47c4493836bebe412f2cd8a"
        },
        {
            "Duration": "1m0.00034707s",
            "Status": "Completed",
            "TaskId": "3d8946e02fa34517a093ed46f84323fc"
        },
        {
            "Duration": "1m0.00013859s",
            "Status": "Completed",
            "TaskId": "ab6b3d680e324ecb8c60bcec118d2a7c"
        },
        {
            "Status": "Running",
            "TaskId": "fe965d7a2b14471c8339b555f0f51014"
        },
        {
            "Status": "Running",
            "TaskId": "594a8e2fad2247328521a9d85066c627"
        },
        {
            "Status": "Running",
            "TaskId": "96ac00d2f3c343e9a301cd9236497cf8"
        }
    ]
}
```

Cancel a task:  
```go
curl -X POST http://localhost:8888/tasky/v1/tasks/96ac00d2f3c343e9a301cd9236497cf8/cancel | python -mjson.tool

curl http://localhost:8888/tasky/v1/tasks/ | python -mjson.tool 
{
    "Tasks": [
        {
            "Status": "Running",
            "TaskId": "fe965d7a2b14471c8339b555f0f51014"
        },
        {
            "Status": "Running",
            "TaskId": "594a8e2fad2247328521a9d85066c627"
        },
        {
            "Status": "Canceled",
            "TaskId": "96ac00d2f3c343e9a301cd9236497cf8"
        }
    ]
}
```

Get Result of a task:  
```go
curl -X POST -H "Content-Type: application/json" -d '{}' http://localhost:8888/tasky/v1/workers/ifconfig | python -mjson.tool
{
    "TaskId": "ac2d1192bd4d46e48207956554aa230f"
}

curl http://localhost:8888/tasky/v1/tasks/ac2d1192bd4d46e48207956554aa230f/result | python -mjson.tool
{
    "Output": {
        "Interfaces": [
            {
                "interface_name": "lo",
                "ip_network": "127.0.0.1/8",
                "mac_address": ""
            },
            {
                "interface_name": "lo",
                "ip_network": "::1/128",
                "mac_address": ""
            },
            {
                "interface_name": "eth0",
                "ip_network": "172.24.20.26/23",
                "mac_address": "f0:de:f1:4c:3a:ee"
            },
            {
                "interface_name": "eth0",
                "ip_network": "fe80::f2de:f1ff:fe4c:3aee/64",
                "mac_address": "f0:de:f1:4c:3a:ee"
            }
        ]
    },
    "TaskId": "ac2d1192bd4d46e48207956554aa230f"
}
```

## Worker Interface
Create a custom worker by implementing the worker interface:  
```go
type Worker interface {
    // Worker name
    Name() string

	// Details provides all metadata info about the worker and it's config
	Details() *WorkerDetails

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
