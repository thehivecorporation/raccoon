package main

import "log"

func LaunchLogger(c chan string) {
	for msg := range c {
		log.Println(msg)
	}
}