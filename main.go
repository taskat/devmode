package main

import (
	"fmt"
	"time"

	"github.com/taskat/devmode/config"
	"github.com/taskat/devmode/restarter"
	"github.com/taskat/devmode/watcher"
)

func createWatcher() *watcher.Watcher {
	r := restarter.NewRestarter()
	r.StartServer()
	folder := config.WatchFolder()
	include := config.IncludeFiles()
	return watcher.NewWatcher(folder, r, include)
}

func timer(tickChannel chan bool) {
	timeout := config.TimeoutBetweenChecks()
	for {
		time.Sleep(timeout)
		tickChannel <- true
	}
}

func readInput(inputChannel chan string) {
	for {
		var input string
		fmt.Scanln(&input)
		inputChannel <- input
	}
}

func main() {
	tickChannel := make(chan bool, 1)
	go timer(tickChannel)
	inputChannel := make(chan string, 1)
	go readInput(inputChannel)
	w := createWatcher()
	w.Watch(tickChannel, inputChannel)
}
