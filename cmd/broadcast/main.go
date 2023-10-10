package main

import (
	"github.com/donovandicks/godistsys/pkg/server"
)

func main() {
	server := server.NewBroadcastServer[float64]()
	server.Run()
}
