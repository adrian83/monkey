package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/adrian83/monkey/pkg/repl"
)

func main() {
	sysUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s!\n", sysUser.Username)
	fmt.Println("This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
