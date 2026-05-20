package main

import "log"

func main() {
	log.Println("Event Consumer starting...")
	select {} // block forever
}
