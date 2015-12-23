package gschroot
import (
	"time"
	"fmt"
	"utils"
    "os"
    "os/exec"
    "os/signal"
    "net/http"
    "syscall"
    "bufio"
    "strings"
    "log"
)

const (
	INVALID_COMMAND_CODE = 127
    SOCK_DIR = "/tmp/gschroot_socks"
    LOG_DIR = "/tmp/gschroot_logs"
)
// Task to be run by chroot library
type Task struct {
	Name string
	Url string
	Cmd string
	ScriptPath string
	done chan bool
	path string
	Timeout int
	Ts *TaskStatus
	StdoutPath string
	StderrPath string
}

type TaskStatus struct {
	Pid int
	Running bool
	Done bool
	SignalsRecvd []int
    Err string
	TimeStarted time.Time
}

func NewTask(url, cmd, name string) (t *Task, err error) {
    t = new(Task)

    if e, err := utils.Exists(LOG_DIR); !e && err == nil {
        os.Mkdir(LOG_DIR, 0755)
    } else if err != nil {
        log.Fatal("Couldn't create log directory", err)
    }

    t.Url = url
    t.Cmd = cmd
    t.StdoutPath = LOG_DIR + "/" + name + ".out"
    t.StderrPath = LOG_DIR + "/" + name + ".error"
    t.path = utils.CreateRandDest()
    t.Ts = new(TaskStatus)
    os.MkdirAll(t.path, 0755)
    resp, err := http.Get(t.Url)
    if err != nil {
        return
    }

    switch resp.StatusCode {
    case 404:
        err = &TarNotFound{Url: t.Url}
    }

    if err != nil {
        return
    }

    defer resp.Body.Close()
    tr := resp.Body
    //s, _ := ioutil.ReadAll(resp.Body)

    err = utils.Untar(tr, t.path)
    err = utils.CopyFile("/home/tony/Dropbox/giant_swarm/go/src/cmd.sh", t.path + "/cmd.sh")
    if err != nil {
        return 
    }
    err = os.Chmod(t.path + "/cmd.sh", 0777)
    return
}

func (t *Task) Run() (err error) {
    //cmd := "./cmd.sh"
    defer t.Close()
    cmdArr := strings.Fields(t.Command())

    cmdArr = append([]string{"chroot", t.path}, cmdArr...)

    sigc := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    signal.Notify(sigc,
    	syscall.SIGKILL,
        syscall.SIGHUP,
        syscall.SIGINT,
        syscall.SIGUSR1,
        syscall.SIGTERM,
        syscall.SIGQUIT)
    go func() {
        s := <-sigc
        fmt.Println("signal", s)
        done<-true

    }()

    chrootCmd := exec.Command("sudo", cmdArr...)
    cmdOutReader, err := chrootCmd.StdoutPipe()
    err = chrootCmd.Start()
    t.Ts.Pid = chrootCmd.Process.Pid
    t.Ts.TimeStarted = time.Now()
    t.Ts.Running = true

    if err != nil {
        fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
        return
    }

    scanner := bufio.NewScanner(cmdOutReader)
    stdOutFile, err := os.OpenFile(t.StdoutPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    defer stdOutFile.Close()
   	if err != nil {
    	return err
	}

    go func() {
        for scanner.Scan() {
            fmt.Printf(scanner.Text() + "\n")
            if _, err := stdOutFile.WriteString(scanner.Text() + "\n"); err != nil {
            	fmt.Println(err)
            }
        }
    }()

    if err = chrootCmd.Wait(); err != nil {
        t.Ts.Err = err.Error()
        if exiterr, ok := err.(*exec.ExitError); ok {
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                es := status.ExitStatus() //geTs exit status code
                //t.Ts.SignalsRecvd = append(t.Ts.SignalsRecvd, es) TODO: make exit status field
                if es == INVALID_COMMAND_CODE {
                    err = &InvalidCommandError{Msg: "Command is invalid.", Code: INVALID_COMMAND_CODE}
                }
            }
        }
        
    }

    t.Ts.Done = true

    t.Ts.Running = false
    <-done
    return
}

// helper function used mainly for testing purposes
func (t *Task) kill() (err error){
	if proc, err := os.FindProcess(os.Getpid()); err == nil{
		proc.Signal(syscall.SIGTERM)
	}

	return

}



func (t *Task) Close() {
    os.RemoveAll(t.path)
    os.RemoveAll(t.StdoutPath)
    //t.closed = true
}

func (t *Task) Command() string {
	return t.Cmd
}

func QueryTask(processName string) (ts *TaskStatus, err error) {
	args := &QueryArgs{}
	var reply QueryReply

	call(SOCK_DIR + "/" + processName, "GsServer.GetTaskStatus", args, &reply)
	ts = reply.Ts
	//p := os.FindProcess(pid)
	return
}

func GetStdLogPath(processName string) (reply LogPathReply, err error) {
    args := &QueryArgs{}
    call(SOCK_DIR + "/" + processName, "GsServer.GetTaskStdLogPath", args, &reply)
    return

}