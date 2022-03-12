package topic

import "os"

type Consistent struct {
	Messages [][]byte
	FilePath string

	release chan struct{}
}

func (cq *Consistent) Push(msg []byte) error {
	cq.Messages = append(cq.Messages, msg)

	if err := os.WriteFile(cq.FilePath, msg, os.ModeAppend); err != nil {
		return err
	}

	select {
	case cq.release <- struct{}{}:
	default:
	}
	return nil
}

func (cq *Consistent) Get(offset int) []byte {
	if offset >= len(cq.Messages) {
		return nil
	}
	return cq.Messages[offset]
}

func (cq *Consistent) Listen() <-chan struct{} {
	return cq.release
}

func (cq *Consistent) Size() int {
	return len(cq.Messages)
}
