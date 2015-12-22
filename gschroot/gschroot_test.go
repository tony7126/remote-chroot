package gschroot

import (
	"testing"
	"utils"
	"time"
)

func TestNewGschroot(t *testing.T) {
	url := "http://localhost:8000/rootfs.tar"
	cmd := "sleep 10"
	name := "test"
	task, err := NewTask(url, cmd, name)
	dirExists, err := utils.Exists(task.path)
	if !dirExists {
		t.Error("Directory should exist at", task.path)
	}

	// Malformed URL protocol
	url = "htt://localhost:8000/rootfs.tar"
	task, err = NewTask(url, cmd, name)
	if err == nil {
		t.Error("Invalid URL should generate an error in NewTask function")
	}

	// 404
	url = "htt://localhost:8000/rootfs.tr"
	task, err = NewTask(url, cmd, name)
	if err == nil {
		t.Error("404 should create an error in NewTask function")
	}

	// Invalid Tarball
	url = "htt://localhost:8000/Vagrantfile"
	task, err = NewTask(url, cmd, name)
	if err == nil {
		t.Error("Invalid tarball should create an error in NewTask function")
	}


}

func TestBadCmd(t *testing.T) {
	url := "http://localhost:8000/rootfs.tar"
	cmd := "sleep1 10"
	name := "test"
	task, err := NewTask(url, cmd, name)

	go func(){
		for !task.Ts.Done {
			time.Sleep(100 * time.Millisecond) // give time to type in sudo password (don't have workaround for now)
		}
		
		err := task.kill()
		if err != nil {
			t.Error("Kill function should work on running task")
		}
	}()
	err = task.Run()
	if err != nil {
		t.Log(err)
	} else {
		t.Error("Should throw an error for invalid command")
	}

	task.Close()

}

func TestGoodCmd(t *testing.T) {
	url := "http://localhost:8000/rootfs.tar"
	cmd := "ls"
	name := "test"
	task, err := NewTask(url, cmd, name)
	server, err := NewServer(name)
	server.Register(task)
	syncChan := make(chan bool, 1)
	go func(){

		for !task.Ts.Done {
			time.Sleep(100 * time.Millisecond) // give time to type in sudo password (don't have workaround for now)
		}

		ts, err := QueryTask(name) //test RPC call
		if err != nil {
			t.Error("QueryTask has an error", err)
		}

		if ts.Pid == 0 {
			t.Error("TaskStatus should hold a valid pid")
		}
		
		if !ts.Done {
			t.Error("TaskStatus should indicate task is done")
		}

		err = task.kill() //test signal interrupt
		if err != nil {
			t.Error("Kill function should work on running task")
		}

		server.Close()

		exists, err := utils.Exists(SOCK_DIR + "/" + name)
		if exists {
			t.Error("Socket should be closed after RPC server is closed")
		}
		syncChan<-true
	}()

	err = task.Run()
	if err != nil {
		t.Error("Should be a valid task command")
	}

	exists, err := utils.Exists(task.path)
	if exists {
		t.Error("remote image filesystem should be deleted after task is closed")
	}

	<-syncChan

}

