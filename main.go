package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
)

func main() {

	inotifyFd, err := unix.InotifyInit()
	fmt.Println(inotifyFd)
	if err != nil {
		fmt.Print(err)
	}

	filePath := "./api.log"
	watcher, err := unix.InotifyAddWatch(inotifyFd, filePath, unix.IN_ALL_EVENTS)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(watcher, " | err=", err)

	/*
		type InotifyEvent struct {
		    Wd     int32
		    Mask   uint32
		    Cookie uint32
		    Len    uint32
		    Name   [0]uint8
		}
	*/
	go func() {
		var eventBuf [4096 * unix.SizeofInotifyEvent]byte
		for {
			readLength, err := unix.Read(inotifyFd, eventBuf[:])
			if err != nil {
				log.Fatal(err)
				continue
			}
			var offset uint32
			for offset <= uint32(readLength-unix.SizeofInotifyEvent) {

			}
		}
	}()

}
