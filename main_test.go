package writeprogress

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestProgressWriter(t *testing.T) {
	t.Run("NewProgressWriter", func(t *testing.T) {
		pw := NewProgressWriter(1)
		assert.Equal(t, uint64(1), pw.size, "size is initialized properly")
		assert.Equal(t, uint64(0), pw.written, "written is set properly")
	})
	t.Run("Write,Size,GetProgress", func(t *testing.T) {
		var data = []byte("hello, world!")
		var buf bytes.Buffer
		var pw = NewProgressWriter(0)
		assert.Equal(t, uint64(0), pw.size, "zero size init")
		w := io.MultiWriter(&buf, pw)
		if n, err := w.Write(data); err != nil {
			t.Error(err)
			return
		} else {
			assert.Equal(t, uint64(buf.Len()), pw.written, "correctly counted bytes written to buffer")
			assert.Equal(t, uint64(n), pw.written, "correctly counted bytes written")
			assert.Equal(t, pw.GetProgress(), float64(pw.written), "zero size progress is written bytes")
			pw.Resize(uint64(n))
			assert.Equal(t, pw.GetProgress(), float64(1.0), "correctly accounts for progress")
			pw.Resize(uint64(n*2))
			assert.Equal(t, pw.GetProgress(), float64(0.5), "correctly accounts for progress halfway through")
		}
	})
}