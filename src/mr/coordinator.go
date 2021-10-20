package mr

import (
	"errors"
	"log"
	"sync"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	// Your definitions here.
}
type Runnable struct {
	jobId       string
	fileName    string
	startOffset int64
	endOffset   int64
}
type Job struct {
	offset int64
}

var workers = make(map[string]WorkerInfo)
var jobs = make(map[string]map[string]Job)
var mutex = make(map[string]sync.RWMutex)
var waitJobs = make([]Runnable, 10)

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Register(args RegisterRequest, reply *RegisterResponse) error {
	// 重复注册
	if _, ok := workers[args.WorkerId]; ok {
		reply = &RegisterResponse{
			Success: false,
			Code:    int16(400),
			Message: string("Coordinators already have workers"),
		}
		return errors.New("Coordinators already have workers")
	}

	if len(args.WorkerId) != 0 {
		workers[args.WorkerId] = WorkerInfo{
			host:     args.Host,
			port:     args.Port,
			workerId: args.WorkerId,
		}
		// 初始化针对worker的锁
		mutex[args.WorkerId] = sync.RWMutex{}
		//response
		reply.Message = "Coordinators already have workers"
		reply.Success = true
		reply.Code = int16(200)
	}

	return nil
}
func (c *Coordinator) Heartbeat(args HeartbeatRequest, reply *HeartbeatResponse) error {
	if _, ok := workers[args.WorkerId]; ok {
		rwMutex := mutex[args.WorkerId]

		rwMutex.Lock()
		if _, ok := jobs[args.WorkerId]; !ok {
			jobs[args.WorkerId] = make(map[string]Job)
		}
		jobs[args.WorkerId][args.JobId] = Job{offset: 0}
		rwMutex.Unlock()
		reply.Message = "Success"
		reply.Success = true
		reply.Code = int16(200)
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()

	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("tcp", ":9707")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	// Your code here.

	for _, s := range files {
		file, err := os.Open(s)
		if err != nil {
			stat, err := file.Stat()
			if err != nil {
				size := stat.Size()
				_ = append(waitJobs, Runnable{
					fileName:    stat.Name(),
					startOffset: 0,
					endOffset:   size,
					jobId:       stat.Name(),
				})
			}
		}
	}

	c.server()
	return &c
}
