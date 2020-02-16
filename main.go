package main

import (
	"fmt"
	"os"
)

const store = "store.txt"

func readNotes() ([]byte, error) {
}

func writeNote(note string) (int, error) {
	file, err := os.OpenFile(store, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return 0, err
	}
	n, err := file.WriteString(note + "\n")
	if err != nil {
		return 0, err
	}
	return n, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected note or '--read/-r' flag.")
		os.Exit(1)
	}
	read := os.Args[1] == "-r" || os.Args[1] == "--read"
	if read {
		notes, err := readNotes()
		if err != nil {
			fmt.Printf("Error reading notes: %v", err)
			os.Exit(1)
		}
		fmt.Println(notes)
		// Pass to fzf
		// Msg
		os.Exit(0)
	}
	for i, n := range os.Args[1:] {
		bytecount, err := writeNote(n)
		if err != nil {
			fmt.Printf("Error writing note '%d': %v", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("Wrote %d bytes\n", bytecount)
	}
}
