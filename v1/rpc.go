package rpc

type RPC[I, O any] struct {
	f  func(I) O
	ch chan Op
}

func NewRPC[I, O any](f func(I) O) *RPC[I, O] {
	return &RPC[I, O]{
		f:  f,
		ch: make(chan Op),
	}
}

func (r *RPC[I, O]) Send(input I) Waiter[O] {
	myOp := newOp(r.f, input)
	r.ch <- myOp
	return &waiter[O]{myOp.(*op[O])}
}

func (r *RPC[I, O]) Receive() <-chan Op {
	return r.ch
}

type Op interface {
	Do()
}

type op[O any] struct {
	f      func()
	done   chan struct{}
	result *O
}

func newOp[I, O any](f func(I) O, input I) Op {
	myOp := &op[O]{
		done: make(chan struct{}),
	}
	myOp.f = func() {
		output := f(input)
		myOp.result = &output
		close(myOp.done)
	}
	return myOp
}

func (o *op[O]) Do() {
	o.f()
}

type Waiter[O any] interface {
	Wait() O
}

type waiter[O any] struct {
	op *op[O]
}

func (w *waiter[O]) Wait() O {
	<-w.op.done
	return *w.op.result
}
