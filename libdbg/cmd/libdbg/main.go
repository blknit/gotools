package main

import "C"

//export hello
func hello(message string) string {
	return "hello, " + message
}

func main() {}
