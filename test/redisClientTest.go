package main

import (
    "goTools/clients"
    "fmt"
)
var (
    REDIS_CONNECT_HOST_AND_PORT string = "127.0.0.1:6379"
)

func main() {
    cli := clients.NewRedisClient()
    cli.Connect(REDIS_CONNECT_HOST_AND_PORT, "")
    fmt.Println(cli.Sadd("abc", "1"))
}
