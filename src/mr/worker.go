package mr

import (
	"fmt"
	"net"
	"os"
)
import "log"
import "net/rpc"
import "hash/fnv"

type KeyValue struct {
	Key   string
	Value string
}
type WorkerInfo struct {
	workerId string
	host     string
	port     int16
}

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	register()
	ch := make(chan bool)
	go heartbeat(ch)
}

func heartbeat(ch chan bool) {
	request := &HeartbeatRequest{
		WorkerId: getLocalHost() + "1",
		JobId:    "job_id",
		Host:     getLocalHost(),
		Port:     0,
	}
	response := &HeartbeatResponse{}
	call("Coordinator.Heartbeat", request, response)
}

func register() {
	request := RegisterRequest{
		Host:     getLocalHost(),
		WorkerId: getLocalHost() + "1",
	}
	response := RegisterResponse{}
	call("Coordinator.Register", request, &response)
}

func getLocalHost() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("tcp", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
