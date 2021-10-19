package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type HeartbeatRequest struct {
	workerId string
	jobId    string
	host     string
	port     int16
}
type HeartbeatResponse struct {
	success bool
	message string
	code    int8
}

type RegisterRequest struct {
	workerId string
	host     string
	port     int16
}

type RegisterResponse struct {
	success bool
	code    int16
	message string
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	return "127.0.0.1:9707"
}
