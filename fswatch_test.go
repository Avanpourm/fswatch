package main

import (
	"context"
	"testing"
)

func TestWatcherRecursive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	WatcherRecursive(ctx, "./", &FsNotifyBase{})
	cancel()
}
