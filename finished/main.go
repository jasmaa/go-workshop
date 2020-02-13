package main

import (
	"fmt"
	"time"
)

type entry interface {
	getKey() string
}

type human struct {
	name string
	age  int
}

type dog struct {
	name  string
	owner *human
	age   int
}

func (h human) getKey() string {
	return h.name
}

func (d dog) getKey() string {
	return d.name
}

func main() {

	title := "Awesome Lookup"
	fmt.Printf("Welcome to %s!\n", title)

	// Setup
	alice := human{"Alice", 56}
	bob := human{"Bob", 9}
	missile := dog{"Missile", &bob, 56}
	entryList := []entry{alice, bob, missile}

	lookup := make(map[string]entry)
	for _, entry := range entryList {
		lookup[entry.getKey()] = entry
	}

	// Go-routines
	c := make(chan string)

	go func() {
		time.Sleep(10 * time.Second)
		c <- "\nTimed out!"
	}()

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("You have %d queries left...\n", 10-i)
			fmt.Print("Enter: ")

			var input string
			fmt.Scanln(&input)

			if entry, ok := lookup[input]; ok {
				fmt.Println(entry)
			} else {
				fmt.Println("No entry found")
			}
		}
		c <- "Session ended"
	}()

	msg := <-c
	fmt.Println(msg)
}
