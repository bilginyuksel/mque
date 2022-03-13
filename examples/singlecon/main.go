package main

import (
	"github.com/bilginyuksel/mque/pkg/mq"
)

func main() {
	reader, err := mq.NewReader(mq.ReaderConfig{
		URL:         ":8080",
		Topic:       "test",
		MaxByteSize: 232322,
		MinByteSize: 1,
	})
	if err != nil {
		panic(err)
	}

	reader.Read()
}
