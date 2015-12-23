package gschroot

import (
	"net/rpc/jsonrpc"
	"net/rpc"
    "log"
    "fmt"
    "net"
    "os"
    "utils"

)


type GsServer struct {
	Addr string
	l net.Listener
	task *Task
	closed bool
}
func NewServer(addr string) (gs *GsServer, err error) {
    if e, err := utils.Exists(SOCK_DIR); !e && err == nil {
        os.Mkdir(SOCK_DIR, 0755)
    } else if err != nil {
        log.Fatal("Couldn't create socket directory", err)
    }
	gs = new(GsServer)
	gs.Addr = SOCK_DIR + "/" + addr
	gs.startServer()
	return
}

func (gs *GsServer) Register(t *Task) {
	gs.task = t
}

func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, errx := jsonrpc.Dial("unix", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	return false
}

func (gs *GsServer) Close() {
	gs.closed = true
	os.Remove(gs.Addr)
}

func (gs *GsServer) startServer()  {
	rpcServer := rpc.NewServer()
	rpcServer.Register(gs)

	l, e := net.Listen("unix", gs.Addr)
	gs.l = l
	if e != nil {
		log.Fatal("RegstrationServer", gs.Addr, " error: ", e)
	}

	
	go func() {
		for !gs.closed {
			conn, err := gs.l.Accept()
			if err == nil {
				go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
			} else {
				fmt.Println("RegistrationServer: accept error", err)
				break
			}
		}
		fmt.Println("RegistrationServer: done\n")
	}()
	
	return
}

type QueryArgs struct {
	Pid int
}

type QueryReply struct {
	Msg string
	Ts *TaskStatus
}

type LogPathReply struct {
	StdoutPath string
	StderrPath string
}

func (gs *GsServer) GetTaskStatus(args *QueryArgs, rep *QueryReply) error {
	rep.Msg = "here"
	rep.Ts = gs.task.Ts
	return nil
}

func (gs *GsServer) GetTaskStdLogPath(args *QueryArgs, rep *LogPathReply) error {
	rep.StdoutPath = gs.task.StdoutPath
	rep.StderrPath = gs.task.StderrPath
	return nil
}
