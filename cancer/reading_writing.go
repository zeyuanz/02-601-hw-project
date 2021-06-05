package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

//ReadBoardFromFile takes a filename as a string and reads in the data provided
//in this file, returning a game board.
func ReadBoardFromFile(filename string) (CellBoard, []string) {
	board := make(CellBoard, 0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var parLine []string
	scanner := bufio.NewScanner(file)
	parInput := false
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentArray := make([]Cell, 0)
		if parInput == false {
			parLine = strings.Split(currentLine, " ")
			parInput = true
			continue
		}
		for i := range currentLine {
			val := currentLine[i : i+1]
			var newCell Cell
			if val == "S" {
				newCell = InitialStemCell()
			} else if val == "D" {
				newCell = InitialDifCell()
			} else if val == "B" {
				newCell = InitialEmptyCell()
			}
			currentArray = append(currentArray, newCell)
		}
		board = append(board, currentArray)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return board,parLine
}
