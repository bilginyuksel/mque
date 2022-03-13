package topic

type (
	Topic interface {
		CreateReader() Reader
		CreateWriter() Writer
	}

	Reader interface {
		Read(offset int) []byte
	}

	Writer interface {
		Write(msg []byte)
	}
)

type topic struct {
	msgs  [][]byte
	locks []chan<- struct{}
}

func New() Topic {
	return &topic{}
}

func (t *topic) read(offset int) []byte {
	return t.msgs[offset]
}

func (t *topic) write(msg []byte) {
	t.msgs = append(t.msgs, msg)

	t.releaseAll()
}

func (t *topic) releaseAll() {
	for _, lock := range t.locks {
		select {
		case lock <- struct{}{}:
		default:
		}
	}
}

func (t *topic) CreateReader() Reader {
	lock := make(chan struct{})
	t.locks = append(t.locks, lock)

	return &reader{
		t:    t,
		lock: lock,
	}
}

type reader struct {
	t    *topic
	lock <-chan struct{}
}

func (r *reader) Read(offset int) []byte {
	if offset >= len(r.t.msgs) {
		<-r.lock
	}
	return r.t.read(offset)
}

func (t *topic) CreateWriter() Writer {
	return &writer{t}
}

type writer struct {
	t *topic
}

func (w *writer) Write(msg []byte) {
	w.t.write(msg)
	w.t.releaseAll()
}
