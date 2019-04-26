package fswatch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

//
var watcher *fsnotify.Watcher

type FsNotifyBase struct {
}

func (f *FsNotifyBase) Create(path string) error {
	return nil
}
func (f *FsNotifyBase) Write(path string) error {
	return nil
}
func (f *FsNotifyBase) Remove(path string) error {
	return nil
}
func (f *FsNotifyBase) Rename(path string) error {
	return nil

}
func (f *FsNotifyBase) Chmod(path string) error {
	return nil

}
func (f *FsNotifyBase) Error(path string, err error) {
	return

}

type Ifsnotify interface {
	Create(path string) error
	Write(path string) error
	Remove(path string) error
	Rename(path string) error
	Chmod(path string) error
	Error(path string, err error)
}

func WatcherRecursive(ctx context.Context, path string, notify Ifsnotify) (err error) {
	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	// starting at the root of the project, walk each file/directory searching for
	// directories
	if err := filepath.Walk(path, watchDir); err != nil {
		// glog.Errorf("Walk Path:%v failed, %v", path, err.Error())
		return err
	}
	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			//glog.V(4).Infof("EVENT! %#v\n", event.String())
			//fmt.Printf("EVENT! %#v\n", event.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, err := os.Stat(event.Name)
				err = watchDir(event.Name, fi, err)
				if err != nil {
					notify.Error(event.Name, err)
				}
				notify.Create(event.Name)
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				notify.Write(event.Name)
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				watcher.Remove(event.Name)
				notify.Remove(event.Name)
			}

			if event.Op&fsnotify.Rename == fsnotify.Rename {
				notify.Rename(event.Name)
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				notify.Chmod(event.Name)
			}

			// watch for errors
		case err := <-watcher.Errors:
			notify.Error(path, err)
		case <-ctx.Done():
			return nil
		}
	}

}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if err != nil {
		return err
	}
	if fi == nil {
		return fmt.Errorf("path:%v, FileInfo is nil", path)
	}
	if fi.Mode().IsDir() {
		//fmt.Printf("watch add path:%v\n", path)
		return watcher.Add(path)

	}

	return nil
}
