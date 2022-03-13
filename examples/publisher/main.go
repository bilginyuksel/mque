package main

import (
	"github.com/bilginyuksel/mque/pkg/mq"
)

func main() {
	writer, err := mq.NewWriter(mq.WriterConfig{
		URL:         ":8080",
		Topic:       "test",
		MaxByteSize: 20000,
		MinByteSize: 1,
	})
	if err != nil {
		panic(err)
	}

	writer.Write()
}
