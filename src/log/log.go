package log

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

var dFlag bool

func init() { dFlag = false }

func DebugToggle() bool {
	dFlag = !dFlag
	return dFlag
}

var buffer bytes.Buffer
var lock sync.Mutex

func getTime() string {
	buffer.Reset()
	buffer.WriteString("[")
	buffer.WriteString(time.Now().Format("15:04:05.000")) //time.StampMilli
	buffer.WriteString("] ")
	return buffer.String()
}

func toString(args ...interface{}) string {
	buffer.Reset()
	size := len(args) - 1
	for i, arg := range args {
		buffer.WriteString(fmt.Sprint(arg))
		if i == size {
			buffer.WriteString("\n")
		} else {
			// buffer.WriteString(" ")
		}
	}
	return buffer.String()
}

func I(args ...interface{}) {
	lock.Lock()
	fmt.Fprint(os.Stdout, getTime())
	fmt.Fprint(os.Stdout, "\033[39m", toString(args...), "\033[0m")
	lock.Unlock()
}

func E(args ...interface{}) {
	lock.Lock()
	fmt.Fprint(os.Stderr, getTime())
	fmt.Fprint(os.Stderr, "\033[31;1m", toString(args...), "\033[0m")
	lock.Unlock()
}

func D(args ...interface{}) {
	if !dFlag {
		return
	}
	lock.Lock()
	fmt.Fprint(os.Stdout, getTime())
	fmt.Fprint(os.Stdout, "\033[37m", toString(args...), "\033[0m")
	lock.Unlock()
}

func H(args ...interface{}) {
	lock.Lock()
	fmt.Fprint(os.Stdout, getTime())
	fmt.Fprint(os.Stdout, "\033[36;1m", toString(args...), "\033[0m")
	lock.Unlock()
}
