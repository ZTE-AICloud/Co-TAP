package main

import (
	"log"
	"os"
	"syscall"
)

var (
	DirToCheck = []string{"../uapregistry-works/logs"}
	ErrLogFile = "../uapregistry-works/logs/stderr.log"
)

func checkDirs(d string) {
	if _, err := os.Stat(d); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(d, os.FileMode(int(0750))); err != nil {
			log.Fatalf("Failed to mkdir:%s", d)
		}
	}
}

func main() {
	for _, d := range DirToCheck {
		checkDirs(d)
	}
	f, err := os.OpenFile(ErrLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(int(0640)))
	if err != nil {
		log.Fatalf("Failed to open file to log stderr:%v", err)
	}
	err = syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0)
	if err != nil {
		log.Fatalf("Failed to redirect stderr to regular file:%v", err)
	}
	cli := NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}
