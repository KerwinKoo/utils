package filemonitor

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// These are the generalized file operations that can trigger a notification.
const (
	Create fsnotify.Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// EventChangeDo call fu func when event of target (defined by arg path,
// file or files in directory)
// return:
//		done----channel, if done-chan is readable, stop the monitor
func EventChangeDo(path string, monitorOp fsnotify.Op, fu func()) (chan bool, error) {
	done := make(chan bool, 1)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return done, err
	}

	err = watcher.Add(path)
	if err != nil {
		return done, err
	}

	go func(watcher *fsnotify.Watcher, done chan bool) {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)

				if event.Op&Write == monitorOp {
					fu()
				}
				if event.Op&Create == monitorOp {
					fu()
				}
				if event.Op&Remove == monitorOp {
					fu()
				}
				if event.Op&Rename == monitorOp {
					fu()
				}
				if event.Op&Chmod == monitorOp {
					fu()
				}

			case <-done:
				return
			case err = <-watcher.Errors:
				return
			}
		}
	}(watcher, done)

	return done, nil
}
