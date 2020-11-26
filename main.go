package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"time"
	"unsafe"
)

func main() {

	inotifyFd, err := unix.InotifyInit()
	fmt.Println(inotifyFd)
	if err != nil {
		fmt.Print(err)
		log.Fatal(err)
	}

	filePath := "./"
	watcher, err := unix.InotifyAddWatch(inotifyFd, filePath, unix.IN_ALL_EVENTS)
	if err != nil {
		fmt.Print(err)
		log.Fatal(err)
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
		var bufferIndex uint32
		var readEventBufferLength uint32
		// 每个SizeofInotifyEvent的大小是0x10，也就是16字节，但是还需要加上name长度
		singleInotifyEventSize := uint32((unix.SizeofInotifyEvent + unix.NAME_MAX + 1))
		// 20是按照最大量预估一个文件或者文件夹能发生的事件数量
		eventBufferLength := singleInotifyEventSize * 20
		eventBuffer := make([]byte, eventBufferLength)
		for {
			bufferIndex = 0
			rawReadEventBufferLength, err := unix.Read(inotifyFd, eventBuffer[:])
			readEventBufferLength = uint32(rawReadEventBufferLength)
			if err != nil {
				log.Fatal(err)
				continue
			}
			for bufferIndex < readEventBufferLength {
				// 从eventBuffer中开始按照顺序拿取第一个event
				singleEvent := (*unix.InotifyEvent)(unsafe.Pointer(&eventBuffer[bufferIndex]))
				// 不知道为什么在监控目录时候，InotifyEvent的Name总是为空，连Go语法都过不了
				// 所以只能通过这种方式委婉地获取文件名
				if singleEvent.Len > 0 {
					fileNameByte := (*[unix.PathMax]byte)(unsafe.Pointer(&eventBuffer[bufferIndex+unix.SizeofInotifyEvent]))[:int(singleEvent.Len):int(singleEvent.Len)]
					fmt.Println(string(fileNameByte))
				}
				bufferIndex = uint32(bufferIndex) + uint32(unix.SizeofInotifyEvent) + singleEvent.Len
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
