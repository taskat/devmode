package restarter

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/taskat/devmode/config"
)

type Restarter struct {
	started bool
	pid     int
}

func NewRestarter() *Restarter {
	return &Restarter{}
}

func (r *Restarter) getPid() {
	fileName := config.PidFile()
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic("Could not read pid file: " + err.Error())
	}
	data = data[:len(data)-1]
	err = os.Remove(fileName)
	if err != nil {
		panic("Could not remove pid file: " + err.Error())
	}
	r.pid, err = strconv.Atoi(string(data))
	if err != nil {
		panic("Could not parse pid: " + err.Error())
	}
}

func (r *Restarter) killServer() {
	script := config.KillServerScript()
	cmd := exec.Command("sh", "-c", script+" "+strconv.Itoa(r.pid))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic("Could not kill server: " + err.Error())
	}
}

func (r *Restarter) RestartServer() {
	fmt.Println("Restarting server...")
	r.ShutDownServer()
	r.StartServer()
}

func (r *Restarter) ShutDownServer() {
	if r.started {
		r.killServer()
		r.waitForKill()
	}
}

func (r *Restarter) StartServer() {
	script := config.StartServerScript()
	cmd := exec.Command("sh", "-c", script)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting server: " + err.Error())
	}
	r.started = err == nil
	if r.started {
		r.getPid()
	}
}

func (r *Restarter) waitForKill() {
	timeout := config.WaitForServerKill()
	for {
		cmd := exec.Command("sh", "-c", "ps -p "+strconv.Itoa(r.pid))
		err := cmd.Run()
		if err != nil {
			_, isExitError := err.(*exec.ExitError)
			if !isExitError {
				panic("Could not check if server is running: " + err.Error())
			}
			break
		}
		time.Sleep(timeout)
		r.killServer()
		var waitStatus syscall.WaitStatus
		// It shows an error on windows, but it will run on linux, so there is no real problem
		_, err = syscall.Wait4(r.pid, &waitStatus, 0, nil)
		if err != nil {
			fmt.Println("Error:", err)
		}

	}
}
