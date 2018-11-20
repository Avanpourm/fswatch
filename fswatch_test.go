package fswatch

import (
	"context"
	"fmt"
	"os"
	"testing"
)

type FsNotifyTest struct {
	FsNotifyBase
}

func (f *FsNotifyTest) Create(path string) error {
	fi, err := os.Stat(path)
	fmt.Printf("Create path:%v name:%v,%+v, %v\n", path, fi.Name(), fi, err)
	return nil
}

func TestWatcherRecursive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	WatcherRecursive(ctx, "./", &FsNotifyTest{})
	cancel()
}
