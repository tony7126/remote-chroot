package main

import (
    "flag"
	//"io/ioutil"
	"fmt"
    "gschroot"
    "log"
    "os"
	)

func main() {

	url := flag.String("url", "", "URL for image")
    cmd := flag.String("cmd", "", "Command to run")
    script := flag.String("script", "", "Path of script to run in chroot environment")
    name := flag.String("name", "", "Name of the running process so it can be queried later") //also used to create unix pipe
    taskName := flag.String("query", "", "Task name to look up status on")
    //stdout := flag.String("log.out", "", "file to log stdout to") // streams to command line by default
    //stderr := flag.String("log.err", "", "file to log stderr to") // streams to command line by default 

    flag.Parse()

    if *taskName != "" {
        ts, _ := gschroot.QueryTask(*taskName)
        fmt.Printf("%+v", ts)
        return
    } 
	if *url == "" {
		log.Fatal("Missing Image URL")
	}

    if *cmd == "" && *script == "" {
        log.Fatal("Missing command to run")
    }

    if *name == "" {
        log.Fatal("Missing name field")
    }

    server, err := gschroot.NewServer(*name)

    if err != nil {
        log.Fatal(err)
    }
    t, err := gschroot.NewTask(*url, *cmd, *name)
    if err != nil {
        fmt.Fprintln(os.Stderr, "err", err.Error())
        os.Exit(1)
    }
    
    server.Register(t) //so it can be queried
    err = t.Run()
    log.Println(err)
    //t.Close()
    server.Close()
}
