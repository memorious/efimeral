package main

import (
	"efimeral/cmd/command"
)

func main() {
	commandObject := new(command.EfimeralObject)
	commandObject.Entry()
}
