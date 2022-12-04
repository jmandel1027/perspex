package main

import "github.com/jmandel1027/perspex/services/backend/pkg/server"

func main() {

	go server.Serve()

	select {}
}
