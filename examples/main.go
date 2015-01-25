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

	quitCh := make(chan bool)
	dataCh := make(chan []byte)
	errCh := make(chan error)

	go tw.Perform(nil, dataCh, errCh, quitCh)

	select {
	case o, ok := <-dataCh:
		if ok {
			fmt.Println("data: ", string(o))
		}

	case e := <-errCh:
		fmt.Println("error: ", e)
	}

	i := &Ifconfig{}

	tw2, err := tasky.NewWorker(i)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("i: ", tw2)

	fmt.Println("Info: ", string(tw2.Info()))

	go tw2.Perform(nil, dataCh, errCh, quitCh)

	select {
	case o, ok := <-dataCh:
		if ok {
			fmt.Println("data: ", string(o))
		}

	case e := <-errCh:
		fmt.Println("error: ", e)
	}
}
