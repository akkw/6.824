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
	WorkerId string
	JobId    string
	Host     string
	Port     int16
}
type HeartbeatResponse struct {
	Success bool
	Message string
	Code    int16
}

type RegisterRequest struct {
	WorkerId string
	Host     string
	Port     int16
}

type RegisterResponse struct {
	Success bool
	Code    int16
	Message string
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	return "127.0.0.1:9707"
}
