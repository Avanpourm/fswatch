package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
)

//
var watcher *fsnotify.Watcher

type FsNotifyBase struct {
}

func (f *FsNotifyBase) Create(path string) error {
	glog.V(4).Infof("Create %v", path)
	return nil
}
func (f *FsNotifyBase) Write(path string) error {

	glog.V(4).Infof("Write %v", path)
	return nil
}
func (f *FsNotifyBase) Remove(path string) error {
	glog.V(4).Infof("Remove %v", path)
	return nil
}
func (f *FsNotifyBase) Rename(path string) error {
	glog.V(4).Infof("Rename %v", path)
	return nil

}
func (f *FsNotifyBase) Chmod(path string) error {
	glog.V(4).Infof("Chmod %v", path)
	return nil

}

type Ifsnotify interface {
	Create(path string) error
	Write(path string) error
	Remove(path string) error
	Rename(path string) error
	Chmod(path string) error
}

func WatcherRecursive(ctx context.Context, path string, notify Ifsnotify) {
	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(path, watchDir); err != nil {
		glog.Errorf("Walk Path:%v failed, %v", path, err.Error())
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			fmt.Printf("EVENT! %#v\n", event.String())
			glog.V(4).Infof("EVENT! %#v\n", event.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, err := os.Stat(event.Name)
				watchDir(event.Name, fi, err)
				notify.Create(event.Name)
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				notify.Write(event.Name)
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				notify.Remove(event.Name)
				fmt.Printf("watch remove path:%v\n", event.Name)
				watcher.Remove(event.Name)
			}

			if event.Op&fsnotify.Rename == fsnotify.Rename {
				notify.Rename(event.Name)
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				notify.Chmod(event.Name)
			}

			// watch for errors
		case err := <-watcher.Errors:
			glog.Errorf("get a watcher Error Path:%v failed, %v", path, err.Error())
		case <-ctx.Done():
			return
		}
	}

}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		fmt.Printf("watch add path:%v\n", path)
		return watcher.Add(path)

	}

	return nil
}
