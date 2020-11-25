package main

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func main() {

	inotifyFd, err := unix.InotifyInit()
	fmt.Println(inotifyFd)
	if err != nil {
		fmt.Print(err)
	}

}
