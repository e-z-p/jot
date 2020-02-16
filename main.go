package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const store = "store.txt"

func readNotes() ([]string, error) {
	var notes []string

	file, err := os.Open(store)
	defer file.Close()
	if err != nil {
		return []string{}, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		notes = append(notes, string(scanner.Text()))
	}
	if err = scanner.Err(); err != nil {
		return []string{}, err
	}

	return notes, nil
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

func fuzzyFind(notes []string) error {
	fzf := exec.Command("fzf")
	stdin, err := fzf.StdinPipe()
	defer stdin.Close()
	if err != nil {
		log.Fatalf("Error obtaining stdin: %s", err.Error())
	}
	for _, n := range notes {
		stdin.Write([]byte(n))
	}
	fzf.Start()
	err = fzf.Wait()
	return err
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
		err = fuzzyFind(notes)
		if err != nil {
			fmt.Printf("Error fuzzy finding notes: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	for i, n := range os.Args[1:] {
		bytecount, err := writeNote(n)
		if err != nil {
			fmt.Printf("Error writing note '%d': %v", i+1, err)
			os.Exit(1)
		}
		fmt.Printf("%d bytes written\n", bytecount)
	}
	fmt.Printf("jotted %d note(s)\n", len(os.Args[1:]))
}
