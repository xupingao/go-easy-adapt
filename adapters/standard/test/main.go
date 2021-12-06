package main

import "github.com/xupingao/go-easy-adapt/adapters/standard"

func main() {
	server := standard.NewDefaultServer()
	err := server.Run()
	if err != nil {
		print(err.Error())
	}
}
