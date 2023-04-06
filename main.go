package writeprogress

import( "sync/atomic" )

type ProgressWriter struct {
	size    uint64
	written uint64
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	l := len(p)
	atomic.AddUint64(&pw.written, uint64(l))
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
