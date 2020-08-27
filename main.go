package main

import (
	"github.com/OmerKahani/tagChanger/cmd"
	"os"
)

func main() {
	if err := tagChanger.GetCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
