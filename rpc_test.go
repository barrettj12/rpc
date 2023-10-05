package rpc_test

import (
	"fmt"
	"time"

	"github.com/barrettj12/rpc"
)

func Example() {
	l := lookup{
		data: map[string]string{
			"foo": "bar",
		},
		rpcCalls: make(chan rpc.Call),
	}
	go l.Serve()

	get := rpc.Register(l.rpcCalls, l.Get)
	value := get.Call("foo").Await()
	fmt.Println(value) // bar

	setClient := rpc.Register(l.rpcCalls, l.Set)
	set := func(key, val string) {
		setClient.Call(SetArgs{key, val}).Await()
	}

	set("foo", "baz")
	value = get.Call("foo").Await()
	fmt.Println(value) // baz

	// Output:
	// bar
	// baz
}

type lookup struct {
	data     map[string]string
	rpcCalls chan rpc.Call
}

func (l *lookup) Serve() {
	for {
		select {
		case call := <-l.rpcCalls:
			call.Do()
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (l *lookup) Get(key string) string {
	return l.data[key]
}

func (l *lookup) Set(args SetArgs) SetResult {
	l.data[args.key] = args.val
	return SetResult{}
}

type SetArgs struct{ key, val string }
type SetResult struct{}

func (l *lookup) Delete(key string) {
	delete(l.data, key)
}
