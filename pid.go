package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func handlePidFile() {
	var err error
	var pidfile *os.File
	pidfile, err = os.OpenFile(pidFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error opening pidfile: %s: %s", pidFile, err.Error())
	}
	if err = syscall.Flock(int(pidfile.Fd()), syscall.LOCK_NB|syscall.LOCK_EX); err != nil {
		log.Printf("Error locking  pidfile: %s: %s", pidFile, err.Error())
	}
	syscall.Ftruncate(int(pidfile.Fd()), 0)
	syscall.Write(int(pidfile.Fd()), []byte(fmt.Sprintf("%d", os.Getpid())))
}
