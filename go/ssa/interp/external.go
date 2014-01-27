// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interp

// Emulated functions that we cannot interpret because they are
// external or because they use "unsafe" or "reflect" operations.

import (
	"math"
	"os"
	"runtime"
	"syscall"
	"time"
)

type externalFn func(fr *frame, args []value) value

// TODO(adonovan): fix: reflect.Value abstracts an lvalue or an
// rvalue; Set() causes mutations that can be observed via aliases.
// We have not captured that correctly here.

// Key strings are from Function.String().
var externals map[string]externalFn

func init() {
	// That little dot ۰ is an Arabic zero numeral (U+06F0), categories [Nd].
	externals = map[string]externalFn{
		"(*runtime.Func).Entry":           ext۰runtime۰Func۰Entry,
		"(*runtime.Func).FileLine":        ext۰runtime۰Func۰FileLine,
		"(*runtime.Func).Name":            ext۰runtime۰Func۰Name,
		"(*sync.Pool).Get":                ext۰sync۰Pool۰Get,
		"(*sync.Pool).Put":                ext۰sync۰Pool۰Put,
		"(reflect.Value).Bool":            ext۰reflect۰Value۰Bool,
		"(reflect.Value).CanAddr":         ext۰reflect۰Value۰CanAddr,
		"(reflect.Value).CanInterface":    ext۰reflect۰Value۰CanInterface,
		"(reflect.Value).Elem":            ext۰reflect۰Value۰Elem,
		"(reflect.Value).Field":           ext۰reflect۰Value۰Field,
		"(reflect.Value).Float":           ext۰reflect۰Value۰Float,
		"(reflect.Value).Index":           ext۰reflect۰Value۰Index,
		"(reflect.Value).Int":             ext۰reflect۰Value۰Int,
		"(reflect.Value).Interface":       ext۰reflect۰Value۰Interface,
		"(reflect.Value).IsNil":           ext۰reflect۰Value۰IsNil,
		"(reflect.Value).IsValid":         ext۰reflect۰Value۰IsValid,
		"(reflect.Value).Kind":            ext۰reflect۰Value۰Kind,
		"(reflect.Value).Len":             ext۰reflect۰Value۰Len,
		"(reflect.Value).MapIndex":        ext۰reflect۰Value۰MapIndex,
		"(reflect.Value).MapKeys":         ext۰reflect۰Value۰MapKeys,
		"(reflect.Value).NumField":        ext۰reflect۰Value۰NumField,
		"(reflect.Value).NumMethod":       ext۰reflect۰Value۰NumMethod,
		"(reflect.Value).Pointer":         ext۰reflect۰Value۰Pointer,
		"(reflect.Value).Set":             ext۰reflect۰Value۰Set,
		"(reflect.Value).String":          ext۰reflect۰Value۰String,
		"(reflect.Value).Type":            ext۰reflect۰Value۰Type,
		"(reflect.Value).Uint":            ext۰reflect۰Value۰Uint,
		"(reflect.error).Error":           ext۰reflect۰error۰Error,
		"(reflect.rtype).Bits":            ext۰reflect۰rtype۰Bits,
		"(reflect.rtype).Elem":            ext۰reflect۰rtype۰Elem,
		"(reflect.rtype).Field":           ext۰reflect۰rtype۰Field,
		"(reflect.rtype).Kind":            ext۰reflect۰rtype۰Kind,
		"(reflect.rtype).NumField":        ext۰reflect۰rtype۰NumField,
		"(reflect.rtype).NumMethod":       ext۰reflect۰rtype۰NumMethod,
		"(reflect.rtype).NumOut":          ext۰reflect۰rtype۰NumOut,
		"(reflect.rtype).Out":             ext۰reflect۰rtype۰Out,
		"(reflect.rtype).Size":            ext۰reflect۰rtype۰Size,
		"(reflect.rtype).String":          ext۰reflect۰rtype۰String,
		"bytes.Equal":                     ext۰bytes۰Equal,
		"bytes.IndexByte":                 ext۰bytes۰IndexByte,
		"hash/crc32.haveSSE42":            ext۰crc32۰haveSSE42,
		"math.Abs":                        ext۰math۰Abs,
		"math.Exp":                        ext۰math۰Exp,
		"math.Float32bits":                ext۰math۰Float32bits,
		"math.Float32frombits":            ext۰math۰Float32frombits,
		"math.Float64bits":                ext۰math۰Float64bits,
		"math.Float64frombits":            ext۰math۰Float64frombits,
		"math.Min":                        ext۰math۰Min,
		"reflect.New":                     ext۰reflect۰New,
		"reflect.TypeOf":                  ext۰reflect۰TypeOf,
		"reflect.ValueOf":                 ext۰reflect۰ValueOf,
		"reflect.init":                    ext۰reflect۰Init,
		"reflect.valueInterface":          ext۰reflect۰valueInterface,
		"runtime.Breakpoint":              ext۰runtime۰Breakpoint,
		"runtime.Caller":                  ext۰runtime۰Caller,
		"runtime.FuncForPC":               ext۰runtime۰FuncForPC,
		"runtime.GC":                      ext۰runtime۰GC,
		"runtime.GOMAXPROCS":              ext۰runtime۰GOMAXPROCS,
		"runtime.Gosched":                 ext۰runtime۰Gosched,
		"runtime.NumCPU":                  ext۰runtime۰NumCPU,
		"runtime.ReadMemStats":            ext۰runtime۰ReadMemStats,
		"runtime.SetFinalizer":            ext۰runtime۰SetFinalizer,
		"runtime.getgoroot":               ext۰runtime۰getgoroot,
		"strings.IndexByte":               ext۰strings۰IndexByte,
		"sync.runtime_Syncsemcheck":       ext۰sync۰runtime_Syncsemcheck,
		"sync.runtime_registerPool":       ext۰sync۰runtime_registerPool,
		"sync/atomic.AddInt32":            ext۰atomic۰AddInt32,
		"sync/atomic.CompareAndSwapInt32": ext۰atomic۰CompareAndSwapInt32,
		"sync/atomic.LoadInt32":           ext۰atomic۰LoadInt32,
		"sync/atomic.LoadUint32":          ext۰atomic۰LoadUint32,
		"sync/atomic.StoreInt32":          ext۰atomic۰StoreInt32,
		"sync/atomic.StoreUint32":         ext۰atomic۰StoreUint32,
		"syscall.Close":                   ext۰syscall۰Close,
		"syscall.Exit":                    ext۰syscall۰Exit,
		"syscall.Fstat":                   ext۰syscall۰Fstat,
		"syscall.Getpid":                  ext۰syscall۰Getpid,
		"syscall.Getwd":                   ext۰syscall۰Getwd,
		"syscall.Kill":                    ext۰syscall۰Kill,
		"syscall.Lstat":                   ext۰syscall۰Lstat,
		"syscall.Open":                    ext۰syscall۰Open,
		"syscall.ParseDirent":             ext۰syscall۰ParseDirent,
		"syscall.RawSyscall":              ext۰syscall۰RawSyscall,
		"syscall.Read":                    ext۰syscall۰Read,
		"syscall.ReadDirent":              ext۰syscall۰ReadDirent,
		"syscall.Stat":                    ext۰syscall۰Stat,
		"syscall.Write":                   ext۰syscall۰Write,
		"time.Sleep":                      ext۰time۰Sleep,
		"time.now":                        ext۰time۰now,
	}
}

// wrapError returns an interpreted 'error' interface value for err.
func wrapError(err error) value {
	if err == nil {
		return iface{}
	}
	return iface{t: errorType, v: err.Error()}
}

func ext۰runtime۰Func۰Entry(fr *frame, args []value) value {
	return 0
}

func ext۰runtime۰Func۰FileLine(fr *frame, args []value) value {
	return tuple{"unknown.go", -1}
}

func ext۰runtime۰Func۰Name(fr *frame, args []value) value {
	return "unknown"
}

func ext۰sync۰Pool۰Get(fr *frame, args []value) value {
	if New := (*args[0].(*value)).(structure)[4]; New != nil {
		return call(fr.i, fr, 0, New, nil)
	}
	return nil
}

func ext۰sync۰Pool۰Put(fr *frame, args []value) value {
	return nil
}

func ext۰bytes۰Equal(fr *frame, args []value) value {
	// func Equal(a, b []byte) bool
	a := args[0].([]value)
	b := args[1].([]value)
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ext۰bytes۰IndexByte(fr *frame, args []value) value {
	// func IndexByte(s []byte, c byte) int
	s := args[0].([]value)
	c := args[1].(byte)
	for i, b := range s {
		if b.(byte) == c {
			return i
		}
	}
	return -1
}

func ext۰crc32۰haveSSE42(fr *frame, args []value) value {
	return false
}

func ext۰math۰Float64frombits(fr *frame, args []value) value {
	return math.Float64frombits(args[0].(uint64))
}

func ext۰math۰Float64bits(fr *frame, args []value) value {
	return math.Float64bits(args[0].(float64))
}

func ext۰math۰Float32frombits(fr *frame, args []value) value {
	return math.Float32frombits(args[0].(uint32))
}

func ext۰math۰Abs(fr *frame, args []value) value {
	return math.Abs(args[0].(float64))
}

func ext۰math۰Exp(fr *frame, args []value) value {
	return math.Exp(args[0].(float64))
}

func ext۰math۰Float32bits(fr *frame, args []value) value {
	return math.Float32bits(args[0].(float32))
}

func ext۰math۰Min(fr *frame, args []value) value {
	return math.Min(args[0].(float64), args[1].(float64))
}

func ext۰runtime۰Breakpoint(fr *frame, args []value) value {
	runtime.Breakpoint()
	return nil
}

func ext۰runtime۰Caller(fr *frame, args []value) value {
	// TODO(adonovan): actually inspect the stack.
	return tuple{0, "somefile.go", 42, true}
}

func ext۰runtime۰FuncForPC(fr *frame, args []value) value {
	// TODO(adonovan): actually inspect the stack.
	return (*value)(nil)
	//tuple{0, "somefile.go", 42, true}
	//
	//func FuncForPC(pc uintptr) *Func
}

func ext۰runtime۰getgoroot(fr *frame, args []value) value {
	return os.Getenv("GOROOT")
}

func ext۰strings۰IndexByte(fr *frame, args []value) value {
	// func IndexByte(s string, c byte) int
	s := args[0].(string)
	c := args[1].(byte)
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func ext۰sync۰runtime_registerPool(fr *frame, args []value) value {
	return nil
}

func ext۰sync۰runtime_Syncsemcheck(fr *frame, args []value) value {
	return nil
}

func ext۰runtime۰GOMAXPROCS(fr *frame, args []value) value {
	return runtime.GOMAXPROCS(args[0].(int))
}

func ext۰runtime۰GC(fr *frame, args []value) value {
	runtime.GC()
	return nil
}

func ext۰runtime۰Gosched(fr *frame, args []value) value {
	runtime.Gosched()
	return nil
}

func ext۰runtime۰NumCPU(fr *frame, args []value) value {
	return runtime.NumCPU()
}

func ext۰runtime۰ReadMemStats(fr *frame, args []value) value {
	// TODO(adonovan): populate args[0].(Struct)
	return nil
}

func ext۰atomic۰LoadUint32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	return (*args[0].(*value)).(uint32)
}

func ext۰atomic۰StoreUint32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	*args[0].(*value) = args[1].(uint32)
	return nil
}

func ext۰atomic۰LoadInt32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	return (*args[0].(*value)).(int32)
}

func ext۰atomic۰StoreInt32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	*args[0].(*value) = args[1].(int32)
	return nil
}

func ext۰atomic۰CompareAndSwapInt32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	p := args[0].(*value)
	if (*p).(int32) == args[1].(int32) {
		*p = args[2].(int32)
		return true
	}
	return false
}

func ext۰atomic۰AddInt32(fr *frame, args []value) value {
	// TODO(adonovan): fix: not atomic!
	p := args[0].(*value)
	newv := (*p).(int32) + args[1].(int32)
	*p = newv
	return newv
}

func ext۰runtime۰SetFinalizer(fr *frame, args []value) value {
	return nil // ignore
}

func ext۰runtime۰funcname_go(fr *frame, args []value) value {
	// TODO(adonovan): actually inspect the stack.
	return (*value)(nil)
	//tuple{0, "somefile.go", 42, true}
	//
	//func FuncForPC(pc uintptr) *Func
}

func ext۰time۰now(fr *frame, args []value) value {
	nano := time.Now().UnixNano()
	return tuple{int64(nano / 1e9), int32(nano % 1e9)}
}

func ext۰time۰Sleep(fr *frame, args []value) value {
	time.Sleep(time.Duration(args[0].(int64)))
	return nil
}

func ext۰syscall۰Exit(fr *frame, args []value) value {
	panic(exitPanic(args[0].(int)))
}

func ext۰syscall۰Getwd(fr *frame, args []value) value {
	s, err := syscall.Getwd()
	return tuple{s, wrapError(err)}
}

func ext۰syscall۰Getpid(fr *frame, args []value) value {
	return syscall.Getpid()
}

func valueToBytes(v value) []byte {
	in := v.([]value)
	b := make([]byte, len(in))
	for i := range in {
		b[i] = in[i].(byte)
	}
	return b
}
