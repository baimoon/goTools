package main

import (
	"fmt"
	"goTools/clients"
	"time"
)

var (
	REDIS_CONNECT_HOST_AND_PORT string = "127.0.0.1:6379"
)

func main() {

	for {
		client := clients.NewRedisClient("", "")
		fmt.Println(client.Get("a"))
		time.Sleep(50 * time.Nanosecond)
	}

}
