package topic

type Basic struct {
	Messages [][]byte

	release chan struct{}
}

func (q *Basic) Push(msg []byte) error {
	q.Messages = append(q.Messages, msg)

	select {
	case q.release <- struct{}{}:
	default:
	}
	return nil
}

func (q *Basic) Get(offset int) []byte {
	if offset >= len(q.Messages) {
		return nil
	}
	return q.Messages[offset]
}

func (q *Basic) Listen() <-chan struct{} {
	return q.release
}

func (q *Basic) Size() int {
	return len(q.Messages)
}
