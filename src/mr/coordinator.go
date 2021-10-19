package mr

import (
	"container/list"
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

var workers map[string]WorkerInfo
var jobs map[string]map[string]Job
var mutex map[string]sync.RWMutex
var waitJobs list.List

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Register(args *RegisterRequest, reply *RegisterResponse) error {
	// 重复注册
	if _, ok := workers[args.workerId]; ok {
		reply = &RegisterResponse{
			success: false,
			code:    int16(400),
			message: string("Coordinators already have workers"),
		}
		return errors.New("Coordinators already have workers")
	}

	if len(args.workerId) != 0 {
		workers[args.workerId] = WorkerInfo{
			host:     args.host,
			port:     args.port,
			workerId: args.workerId,
		}
		// 初始化针对worker的锁
		mutex[args.workerId] = sync.RWMutex{}

		reply = &RegisterResponse{
			success: false,
			code:    int16(200),
			message: string("Coordinators already have workers"),
		}
	}

	return nil
}
func (c *Coordinator) Heartbeat(args *HeartbeatRequest, reply *HeartbeatRequest) error {
	if _, ok := workers[args.workerId]; ok {
		rwMutex := mutex[args.workerId]
		defer rwMutex.Unlock()
		rwMutex.Lock()
		jobs[args.workerId][args.jobId] = Job{}

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
				waitJobs.PushBack(Runnable{
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
