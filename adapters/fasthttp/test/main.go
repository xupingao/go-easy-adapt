package main

import (
	fast "github.com/xupingao/go-easy-adapt/adapters/fasthttp"
)

func main() {
	server := fast.NewDefaultServer()
	err := server.Run()
	if err != nil {
		print(err.Error())
	}
}
