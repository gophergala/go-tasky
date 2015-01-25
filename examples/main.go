package main

import (
	"fmt"
	"log"

	"github.com/gophergala/go-tasky/tasky"
)

func main() {
	cp := &CopyFile{}

	tw, err := tasky.NewWorker(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("cp: ", tw)

	fmt.Println("Info: ", string(tw.Info()))

	e, b := tw.Perform(nil)
	fmt.Println("e: ", string(e))
	fmt.Println("b: ", b)

	i := &Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("i: ", tw2)

	fmt.Println("Info: ", string(tw2.Info()))

	e, b = tw2.Perform(nil)
	fmt.Println("e: ", string(e))
	fmt.Println("b: ", b)
}
