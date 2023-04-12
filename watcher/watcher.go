package watcher

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

type restarter interface {
	RestartServer()
	ShutDownServer()
}

type Watcher struct {
	lastEdited   map[string]time.Time
	folder       string
	r            restarter
	startup      bool
	includeFiles *regexp.Regexp
}

func NewWatcher(folder string, r restarter, includeFiles *regexp.Regexp) *Watcher {
	w := &Watcher{folder: folder, r: r, startup: true, includeFiles: includeFiles}
	w.lastEdited = make(map[string]time.Time)
	w.updateFiles(w.folder)
	w.updateLastEdited()
	w.startup = false
	return w
}

func (w *Watcher) log(msg string) {
	if !w.startup {
		fmt.Println(msg)
	}
}

func (w *Watcher) updateFiles(currentFolder string) {
	entries, err := os.ReadDir(currentFolder)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			w.updateFiles(currentFolder + "/" + entry.Name())
		} else {
			fileName := currentFolder + "/" + entry.Name()
			if !w.includeFiles.MatchString(fileName) {
				continue
			}
			if _, ok := w.lastEdited[fileName]; !ok {
				w.log("Added file: " + fileName)
				w.lastEdited[fileName] = time.Time{}
			}
		}
	}
}

func (w *Watcher) updateLastEdited() bool {
	updated := false
	for file, modTime := range w.lastEdited {
		fileInfo, err := os.Lstat(file)
		if err != nil {
			w.log("Deleted file: " + file)
			delete(w.lastEdited, file)
			updated = true
			continue
		}
		lastMod := fileInfo.ModTime()
		if lastMod != modTime {
			w.log("Updated file: " + file)
			w.lastEdited[file] = lastMod
			updated = true
		}
	}
	return updated
}

func (w *Watcher) Watch(tickChannel chan bool, inputChannel chan string) {
	for {
		select {
		case <-tickChannel:
			w.updateFiles(w.folder)
			if w.updateLastEdited() {
				w.r.RestartServer()
			}
		case input := <-inputChannel:
			if input == "exit" || input == "quit" || input == "q" {
				w.r.ShutDownServer()
				return
			}
		}
	}
}
