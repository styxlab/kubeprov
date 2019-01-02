package main

import (
	"log"

	"github.com/styxlab/kubeprov/cmd"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("kubeprov: ")
}

func main() {
	cmd.Execute()
}
