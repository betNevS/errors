package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Frame wrap runtime.Frame
type Frame runtime.Frame

func (f Frame) file() string {
	if f.File == "" {
		return "unknown file"
	}
	return f.File
}

func (f Frame) line() int {
	return f.Line
}

func (f Frame) name() string {
	if f.Function == "" {
		return "unknown"
	}
	return f.Function
}

// Format formats frame, The following rules:
// %s -> source file
// %d -> source line
// %n -> function name
// %v -> equivalent to %s:%d
// %+s -> <function name>\n\t<path>
// %+v -> equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line()))
	case 'n':
		io.WriteString(s, funcname(f.name()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

type stack struct {
	frames *runtime.Frames
}

func (s stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for {
				f, more := s.frames.Next()
				frame := Frame(f)
				fmt.Fprintf(st, "\n%+v", frame)
				if !more {
					break
				}
			}
		}
	}
}

func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

func callers() stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	return stack{
		frames: frames,
	}
}
