package topic

const (
	NoAck = iota + 1
	HalfAck
)

type Reader interface {
	Get(offset int) []byte
	Size() int
	Listen() <-chan struct{}
}

type Writer interface {
	Push(msg []byte) error
}

type Topic interface {
	Push(msg []byte) error
	Get(offset int) []byte
	Size() int
	Listen() <-chan struct{}
}

func New(ack uint8) Topic {
	if ack == HalfAck {
		return &Consistent{
			FilePath: "data/queue.data",
			release:  make(chan struct{}),
		}
	}
	return &Basic{
		release: make(chan struct{}),
	}
}
