package main

import (
	"encoding/json"
	"fmt"

	"tekito/lib/logging"
	"tekito/lib/utils"
)

type Message struct {
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

func (msg Message) String() string {
	b, _ := json.Marshal(&msg)
	return string(b)
}

func main() {
	fmt.Println("Hello world")

	messages := []Message{}
	utils.ReadLines("test.json", func(i int64, line string) bool {
		// var msg map[string]interface{}
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			logging.Log.Warn(err.Error())
			return false
		}
		fmt.Println(msg)
		messages = append(messages, msg)
		return true
	})
	fmt.Println(messages)

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		ch2 <- 1
	}()

LO:
	for {
		select {
		case v1 := <-ch1:
			fmt.Println(v1)
		case <-ch2:
			break LO
		}
	}
}
