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
		client := clients.NewRedisClientForWeb(REDIS_CONNECT_HOST_AND_PORT, "")

		fmt.Println(client.Zscore("abc", "aac"))
		time.Sleep(time.Second)
	}

}
