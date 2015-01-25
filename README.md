#Go-Tasky
<p align="center">
<img src="https://cloud.githubusercontent.com/assets/3473592/5893581/5a90b250-a4af-11e4-84a4-1a1b14ebc54d.png" alt="go-tasky-logo"/>
</p>

Go-Tasky is a simple golang package that makes it easy to create and execute tasks on a server or in a Docker container that are exposed with a RESTful api.  Think of Go-Tasky as your drop-in server companion tool, for when you need to expose information or run remote tasks via friendly APIs.  

Developed by [Mark Moudy](https://github.com/MarkMoudy) and [Mehul Choube](https://github.com/mchoube) during the 2015 Gophergala global hackathon.  

##Basic Principles

Register a custom worker with Go-Tasky that implements the worker interface and a set of routes will be created for that task. You can have any number of uniquely configured tasks for a worker. Go-Tasky runs on your server and gives you a customizable interface to build your toolkit for servers and containers. 

## Example Workers
- CopyFile - copy a file to a new location  
- Ifconfig - fetch network information read the contents of a file  
- Sleeper -  Sleep and wait

##Usage

###Getting Started
Install Go-Tasky in your $GOPATH with `go get`  
```Go
go get -u github.com/gophergala/go-tasky
```
Navigate to the examples directory and run  
```Go
// use the -addr flag to set server address(default is :8888) ex. -addr=:4444
go run main.go 
```

###Routes
Workers:  
GET /tasky/v1/workers/ - returns a list of available worker endpoints   
POST /tasky/v1/workers/{worker_name} - Creates a task to run.  

Tasks:  
GET /tasky/v1/tasks/ - Returns a list of all tasks.  
GET /tasky/v1/tasks/{id}/status - Returns the status of a task. 
GET /tasky/v1/tasks/{id}/result - Get the output of the task.  
POST /tasky/v1/tasks/{id}/cancel - Cancel the task. 

###Creating a Custom Worker

### Adding to your code


##Roadmap
- Serve a gui web interface alongside the api for point and click task management
- Expand selection of example workers
- Task Chaining so you can specify a sequence of tasks to run. 
- Authentication to secure access. 


## Example
List available workers:  
```go
curl http://localhost:8888/tasky/v1/workers/ | python -mjson.tool 
[
    {
        "config_objects": {
            "destination": "",
            "source": ""
        },
        "description": "CopyFile will copy a file on the server. You must specify the Source and Destination",
        "name": "copyfile"
    },
    {
        "description": "Ifconfig will return networking details from the server. No config is needed for this worker",
        "name": "ifconfig"
    },
    {
        "description": "Sleeps for a minute. This is to showcase long running tasks.",
        "name": "sleeper"
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
curl http://localhost:8888/tasky/v1/tasks/1acb0ad6c8fb4d4faa7aa0e1d0f5f0b6/result 
{
    "Output": {
        "Interfaces": [
            {
                "interface_name": "lo0",
                "ip_network": "::1/128",
                "mac_address": ""
            },
            {
                "interface_name": "lo0",
                "ip_network": "127.0.0.1/8",
                "mac_address": ""
            },
            {
                "interface_name": "lo0",
                "ip_network": "fe80::1/64",
                "mac_address": ""
            },
            {
                "interface_name": "en0",
                "ip_network": "fe80::6203:8ff:fe9b:814a/64",
                "mac_address": "60:03:08:9b:81:4a"
            },
            {
                "interface_name": "en0",
                "ip_network": "192.168.1.6/24",
                "mac_address": "60:03:08:9b:81:4a"
            },
            {
                "interface_name": "awdl0",
                "ip_network": "fe80::6089:27ff:fef1:b73c/64",
                "mac_address": "62:89:27:f1:b7:3c"
            },
            {
                "interface_name": "vboxnet9",
                "ip_network": "192.168.59.3/24",
                "mac_address": "0a:00:27:00:00:09"
            }
        ]
    },
    "TaskId": "6d8b4b6e3e0942faa2fcfe0f9a8757d7"
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
