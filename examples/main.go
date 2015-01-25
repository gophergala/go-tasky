package main

import (
	"fmt"
	"log"

	"github.com/gophergala/go-tasky/tasky"
)

func register() {
	cp := &CopyFile{}

	tw, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("cp: ", tw)

	i := &Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("i: ", tw2)

	s := &Sleeper{}

	tw3, err := tasky.NewWorker(s)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("s: ", tw3)
}

func main() {
	register()
}
