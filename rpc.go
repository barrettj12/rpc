package rpc

// Promise represents a value to be returned at some point in the future.
// Based on the concept of promises from JavaScript.
type Promise[O any] interface {
	// Await waits for the promise to resolve, and returns the resolved value.
	Await() O
}

type promise[O any] struct {
	call *call[O]
}

func (p *promise[O]) Await() O {
	<-p.call.done
	return *p.call.result
}

// Call represents an RPC request from a remote client. A server handling this
// Call should call Do in the goroutine where you would like the work to be
// done.
type Call interface {
	// Do processes this Call.
	Do()
}

type call[O any] struct {
	f      func()
	result *O
	done   chan struct{}
}

func newCall[I, O any](f func(I) O, input I) *call[O] {
	c := &call[O]{
		done: make(chan struct{}),
	}
	c.f = func() {
		output := f(input)
		c.result = &output
		close(c.done)
	}
	return c
}

func (c *call[O]) Do() {
	c.f()
}

// Register returns a Client which can be used to make RPC calls for the given
// function f. RPC calls will be sent on the 'calls' channel, so it is expected
// that a 'server' is running, ready to receive calls from this channel and
// process them.
//
// The registered function f must have a single input and a single output. For
// functions which don't have this signature, you can still register them by
// creating input/output structs for them. For example, the function
//
//	func Set(key, val string) {
//		myMap[key] = val
//	}
//
// can be changed to
//
//	func Set(args SetArgs) SetResult {
//		myMap[args.key] = args.val
//		return SetResult{}
//	}
//
//	type SetArgs struct{ key, val string }
//	type SetResult struct{}
func Register[I, O any](calls chan Call, f func(I) O) Client[I, O] {
	return &client[I, O]{
		f:     f,
		calls: calls,
	}
}

// Client is a client used to make RPC calls for a single method.
type Client[I, O any] interface {
	// Call makes an RPC call with the given input.
	Call(I) Promise[O]
}

type client[I, O any] struct {
	f     func(I) O
	calls chan Call
}

func (c *client[I, O]) Call(input I) Promise[O] {
	call := newCall(c.f, input)
	c.calls <- call
	return &promise[O]{call}
}
