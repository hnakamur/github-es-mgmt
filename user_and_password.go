package main

import (
	"bufio"
	"os"
)

func getOrReadUserAndPassword(argUser, argPassword string) (user, password string, err error) {
	if argUser != "" && argPassword != "" {
		return argUser, argPassword, nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	user = argUser
	if user == "" && scanner.Scan() {
		user = scanner.Text()
	}

	password = argPassword
	if password == "" && scanner.Scan() {
		password = scanner.Text()
	}

	if err = scanner.Err(); err != nil {
		return "", "", err
	}
	return user, password, nil
}
