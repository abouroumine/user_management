package main

import (
	s "./server"
)

var server = s.Server{}

func main() {
	server.Initialize()
}
