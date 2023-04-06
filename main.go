package writeprogress

import (
	"sync"
	"sync/atomic"
)

type ProgressWriter struct {
	size    uint64
	written uint64

	watcherlock sync.RWMutex
	watchers []chan struct{}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	l := len(p)
	atomic.AddUint64(&pw.written, uint64(l))
	go func() {
			pw.watcherlock.RLock()
			defer pw.watcherlock.RUnlock()
			for _, w := range pw.watchers {
				if w == nil { continue }
				go func(w chan<- struct{}) { w<-struct{}{} }(w)
			}
	}()
	return l, nil
}


func (pw *ProgressWriter) GetProgress() float64 {
	var size, written = atomic.LoadUint64(&pw.size), atomic.LoadUint64(&pw.written)
	if size == 0 {
		return float64(written)
	} else if size == written {
		return 1.0
	}
	return float64(written) / float64(size)
}

func (pw *ProgressWriter) Resize(size uint64) {
	atomic.StoreUint64(&pw.size, size)
}

func NewProgressWriter(size uint64) *ProgressWriter {
	return &ProgressWriter{size: size}
}


func (pw *ProgressWriter)registerWatcher(w chan struct{}) {
	pw.watcherlock.Lock()
	defer pw.watcherlock.Unlock()
	pw.watchers = append(pw.watchers, w)
}

func (pw *ProgressWriter)deregisterWatcher(w chan struct{}) {
	pw.watcherlock.Lock()
	defer pw.watcherlock.Unlock()
	for i, x := range pw.watchers {
		if w == x { 
			pw.watchers[i] = nil
			return
		}
	}
}

func (pw *ProgressWriter)Watch(msg func(float64)) (done, cancel chan struct{}) {
	var w  = make(chan struct{})
	cancel = make(chan struct{})
	done = make(chan struct{})
	go func() {
		var progress float64
		pw.registerWatcher(w)
		loop:
		for {
			select {
			case <-w:
				progress = pw.GetProgress()
				msg(progress)
				if progress >= 1.0 {
					break loop
				}
			case <-cancel:
				break loop
			}
		}
		pw.deregisterWatcher(w)
		done <- struct{}{}
	}()
	return
}
