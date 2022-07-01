package main

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"
)

func GetManagementConsolePassword() string {
	fmt.Print("Enter Management Console password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalf("read password from prompt: %s", err)
	}
	return string(bytepw)
}
