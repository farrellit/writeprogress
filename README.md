# writeprogress
an `io.Writer` implementation that simply counts written bytes, meant to be used with `io.MultiWriter` to track write progress.

Example: 

[In this example](https://go.dev/play/p/KdKUfjV43SW) we copy 1e6 bytes from `/dev/urandom` to `/dev/null`.  We we register a watcher function
that writes progress percentage to stdout.  We then execute the `io.Copy`.  Afterwards we can make sure our watcher function was called for each write 
(`<-d`), only to ensure that our last write % doesn't come after our summary message.  


```
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/farrellit/writeprogress"
)

func main() {
	in, err := os.Open("/dev/urandom")
	if err != nil {
		panic(err)
	}
	out, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	defer out.Close()

	length := int64(1e6)
	wp := writeprogress.NewProgressWriter(uint64(length))
	d, _ := wp.Watch(func(p float64) { fmt.Printf("\r%2.0f%%", p*100) })

	if b, err := io.Copy(
		io.MultiWriter(out, wp),
		&io.LimitedReader{R: in, N: length},
	); err != nil {
		panic(err)
	} else {
		<-d
		fmt.Printf("\n%d/%d %2.0f%%\n", b, length, wp.GetProgress()*100)
	}

}

```
