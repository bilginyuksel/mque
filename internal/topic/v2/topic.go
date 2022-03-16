package topic

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Topic interface {
		CreateReader() Reader
		CreateWriter() Writer
	}

	Reader interface {
		Read(buffSize int64, offset int) ([]byte, error)
	}

	Writer interface {
		Write(msg []byte)
	}
)

type topic struct {
	filepath string
	file     *os.File
	locks    []chan<- struct{}
}

func New(name string) (Topic, error) {
	filepath := fmt.Sprintf("./data/%s.log", name)
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &topic{
		filepath: filepath,
		file:     f,
	}, nil
}

func (t *topic) read(readBuff int64, offset int) ([]byte, error) {
	if _, err := t.file.Seek(readBuff, 0); err != nil {
		log.Printf("could not seek, err: %v\n", err)
		return nil, err
	}

	length := make([]byte, 10)
	if _, err := t.file.Read(length); err != nil {
		log.Printf("could not read, err: %v\n", err)
		return nil, err
	}
	log.Printf("length: %v\n", string(length))

	lengthOfBytes, err := strconv.ParseInt(string(length), 10, 64)
	if err != nil {
		log.Printf("could not parse, err: %v\n", err)
		return nil, err
	}
	log.Printf("lengthOfBytes: %v\n", lengthOfBytes)

	msg := make([]byte, lengthOfBytes)
	_, err = t.file.Read(msg)
	return msg, err
}

func (t *topic) write(msg []byte) {
	contentLengthBytes := t.createByteContentArr(len(msg))

	var completeMsg []byte
	completeMsg = append(completeMsg, contentLengthBytes...)
	completeMsg = append(completeMsg, msg...)
	completeMsg = append(completeMsg, '\n')

	go func() {
		if _, err := t.file.Write(completeMsg); err != nil {
			log.Println(err)
		}

		t.file.Sync()
		t.releaseAll()
	}()
}

func (t *topic) createByteContentArr(msgLength int) []byte {
	ls := strconv.FormatInt(int64(msgLength), 10)
	ls = strings.Repeat("0", 10-len(ls)) + ls
	return []byte(ls)
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

func (r *reader) Read(buffLen int64, offset int) ([]byte, error) {
	stat, _ := r.t.file.Stat()
	if buffLen >= stat.Size() {
		<-r.lock
	}

	return r.t.read(buffLen, offset)
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
