package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const directory = "/migrations"
const acceptedChars = "abcdefghijklmnopqrstuvwxyz_"

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	namePtr := flag.String("name", "", "the name of the migration (only lowercase and underscore)")

	pubPtr := flag.Bool("isPublic", false, "set true if this migration is ran at the public level")

	flag.Parse()

	if *namePtr == "" {
		fmt.Printf("A name must be provided\n")
		return
	}

	for _, c := range *namePtr {
		if !strings.Contains(acceptedChars, string(c)) {
			fmt.Printf("A name must only include lower case letters and underscores\n")
			return
		}
	}

	if err := Create(dir+directory, *namePtr, *pubPtr); err != nil {
		fmt.Printf("Error while trying to make migration: %+v\n", err)
	}
}
