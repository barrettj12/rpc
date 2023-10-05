package rpc_test

import (
	"fmt"
	"time"

	"github.com/barrettj12/rpc/v1"
)

func Example() {
	l := lookup{
		data: map[string]string{
			"foo": "bar",
		},
	}

	get := rpc.NewRPC(l.Get)
	l.rpc = get

	go l.Serve()

	w := get.Send("foo")
	value := w.Wait()
	fmt.Println(value)
	// Output: bar
}

type lookup struct {
	data map[string]string
	rpc  *rpc.RPC[string, string]
}

func (l *lookup) Serve() {
	for {
		select {
		case call := <-l.rpc.Receive():
			call.Do()
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (l *lookup) Get(key string) string {
	return l.data[key]
}
