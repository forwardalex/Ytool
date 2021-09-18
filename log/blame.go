package log

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// BlameLine is a structure for a blame result for a specific user
type BlameLine struct {
	AuthorName  string
	AuthorEmail string
	AuthorDate  time.Time
	CommitName  string
	CommitEmail string
	CommitDate  time.Time
}

// Callback is called for each line
type Callback func(line BlameLine) error

var (
	authorPrefix        = "author"
	authorMailPrefix    = "author-mail"
	authorTimePrefix    = "author-time"
	committerPrefix     = "committer"
	committerMailPrefix = "committer-mail"
	committerTimePrefix = "committer-time"
)

// buffer pool to reduce GC
var bufferPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, MaxLineSize))
	},
}

// getBuffer fetches a buffer from the pool
func getBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// putBuffer returns a buffer to the pool
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// MaxLineSize is the maximum of one line of output. testing with 1K which seems OK
var MaxLineSize = 1024

//FindCommit  查询代码提交人
func FindCommit(ctx context.Context, fn string, line int, w io.Writer) (current BlameLine, err error) {
	cmd := exec.Command("git", "blame", "-e", "--root", "--line-porcelain", fn, fmt.Sprintf("-L %d,%d", line, line))
	r, err := cmd.StdoutPipe()
	if err != nil {
		Error("ERR ", err.Error())
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	defer r.Close()
	buf := getBuffer()
	defer putBuffer(buf)
	lr := bufio.NewReaderSize(r, MaxLineSize)
	s := bufio.NewScanner(lr)
	s.Buffer(buf.Bytes(), MaxLineSize)
	var writer *bufio.Writer
	if w != nil {
		writer = bufio.NewWriter(w)
		defer writer.Flush()
	}
	for s.Scan() {
		// make sure our context isn't done
		select {
		case <-ctx.Done():
			return current, nil
		default:
		}
		buf := s.Text()
		if writer != nil {
			_, err := writer.WriteString(buf)
			if err != nil {
				return current, fmt.Errorf("error writing buffer to output. %v", err)
			}
			_, err = writer.WriteString("\n")
			if err != nil {
				return current, fmt.Errorf("error writing buffer to output. %v", err)
			}
		}
		infos := strings.Split(buf, " ")
		switch infos[0] {
		case authorPrefix:
			current.AuthorName = infos[1]
		case authorMailPrefix:
			current.AuthorEmail = infos[1]
		case authorTimePrefix:
			i, err := strconv.ParseInt(infos[1], 10, 64)
			if err != nil {
				return current, err
			}
			current.AuthorDate = time.Unix(i, 0)
		case committerPrefix:
			current.CommitName = infos[1]
		case committerMailPrefix:
			current.CommitEmail = infos[1]
		case committerTimePrefix:
			i, err := strconv.ParseInt(infos[1], 10, 64)
			if err != nil {
				return current, err
			}
			current.CommitDate = time.Unix(i, 0)
		}
	}
	err = s.Err()
	if err != nil {
		if strings.Contains(s.Err().Error(), "file already closed") {
			return current, nil
		}
		return current, err
	}
	return current, nil
}
