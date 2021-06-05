package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ReadBoardFromFiles(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileContent := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println("Something wrong with the scanner")
		os.Exit(1)
	}

	return fileContent
}
