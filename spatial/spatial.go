package main

import (
	"fmt"
	"strconv"
	"strings"
	"os"
	"log"
)
// The data stored in a single cell of a field
type Cell struct {
	strategy string  //represents "C" or "D" corresponding to the type of prisoner in the cell
	score    float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects
type GameBoard [][]Cell

func main() {
	fileName := os.Args[1] //enter fileName


	b,err2 := strconv.ParseFloat(os.Args[2],64) //enter points for D/C scenario
	if err2 != nil {
		log.Fatal(err2)
	}

	numGens,err1 := strconv.Atoi(os.Args[3]) //enter number of gens
	if err1 != nil {
		log.Fatal(err1)
	}

	outputFile := os.Args[4]
	fileContent := ReadBoardFromFiles(fileName)

	numRow, numCol, Board := ReadRowColStrategy(fileContent)
	currBoard := InitializeGameBoard(numRow, numCol, Board)

	fmt.Println("Files read and initialization successfully.")

	boards := UpdateGameBoards(currBoard, numGens, b)
	fmt.Println("Spatial game played.")

	cellWidth := 20

	imglist := DrawGameBoards(boards, cellWidth)
	fmt.Println("Imagelist Created. Now generate GIF.")

	//ImageToPNG(imglist[len(imglist)-1], outputFile)
	ImagesToGIF(imglist, outputFile)
	fmt.Println("GIF generated")
}

//UpdateGameBoards takes a current GameBoard, numGens as number of generations and b as the points for immunity.It returns a slice of GameBoard of numGens+1 number of generations.
func UpdateGameBoards(currBoard GameBoard, numGens int, b float64) []GameBoard {
	numGenBoards := make([]GameBoard, numGens+1)
	numGenBoards[0] = currBoard

	for i := 1; i < numGens+1; i++ {
		numGenBoards[i] = UpdateGameBoard(numGenBoards[i-1],b)
	}
	return numGenBoards
}

//UpdateGameBoard takes a current GameBoard and b as input. It returns the next generation GameBoard according to the neighboring strategy and scores (b).
func UpdateGameBoard(currBoard GameBoard, b float64) GameBoard {
	newBoard := InitializeGameBoard(len(currBoard), len(currBoard[0]), []string{})
	newBoard2 := InitializeGameBoard(len(currBoard), len(currBoard[0]), []string{})
	//Update score for every cell first
	for i := range currBoard {
		for j := range currBoard[i] {
			newBoard[i][j].score = UpdateCellScore(currBoard, i, j, b)
			newBoard[i][j].strategy = currBoard[i][j].strategy
			newBoard2[i][j].score = newBoard[i][j].score
		}
	}

	//Update strategy for every cell according to max score in the neighborhood
	for i := range newBoard {
		for j := range newBoard[i] {
			newBoard2[i][j].strategy = UpdateCellStrategy(newBoard, i, j, b)
		}
	}
	return newBoard2
}

//UpdateCellScore takes a current GameBoard and row, colunmn index and b as input. It returns the total score of a cell.
func UpdateCellScore(currBoard GameBoard, r, c int, b float64) float64 {
	newScore := 0.0
	//score update rules for "C" strategy
	if currBoard[r][c].strategy == "C" {
		for i := r-1; i <= r+1; i++ {
			for j := c-1; j <= c+1; j++ {
				if i == r && j == c {
					continue
				}
				if InField(currBoard,i,j) {
					if currBoard[i][j].strategy == "C" {
						newScore++
					}
				}
			}
		}
	}
	//score update strategy for "D" strategy
	if currBoard[r][c].strategy == "D" {
		for i := r-1; i <= r+1; i++ {
			for j := c-1; j <= c+1; j++ {
				if i == r && j == c {
					continue
				}
				if InField(currBoard,i,j) {
					if currBoard[i][j].strategy == "C" {
						newScore += b
					}
				}
			}
		}
	}
	return newScore
}

//InField takes currBoard and coordinates as input. It returns whether the coordinates is in the currBoard.
func InField(currBoard GameBoard, r,c int) bool {
	numRow := len(currBoard)
	numCol := len(currBoard[0])

	if r > numRow-1 || r < 0 || c > numCol-1 || c < 0 {
		return false
	}
	return true
}

//UpdateCellStrategy takes a current GameBoard and row, colunmn index and b as input. It returns the next strategy of a single cell.
func UpdateCellStrategy(currBoard GameBoard, r, c int, b float64) string {
	maxStrategy := currBoard[r][c].strategy
	maxScore := currBoard[r][c].score
  for i := r-1; i <= r+1; i++ {
		for j := c-1; j <= c+1; j++ {
			if InField(currBoard, i, j) {
				if currBoard[i][j].score > maxScore {
					maxScore = currBoard[i][j].score
					maxStrategy = currBoard[i][j].strategy
				}
			}
		}
	}
	return maxStrategy
}

//InitializeGameBoard takes initialBoard, numRow and numCol as input. It returns a GameBoard with row = numRow, col = numCol and every strategy in the initialBoard.
func InitializeGameBoard(numRow, numCol int, boardStr []string) GameBoard {
	initialBoard := make(GameBoard, numRow)
	for i := range initialBoard {
		initialBoard[i] = make([]Cell, numCol)
	}

	if len(boardStr) != 0 {
		for i := range initialBoard {
			for j := range initialBoard[i] {
				initialBoard[i][j].strategy = string(boardStr[i][j])
			}
		}
	} else {
		for i := range initialBoard {
			for j := range initialBoard[i] {
				initialBoard[i][j].strategy = "C" // set "C" as default value
			}
		}
	}

	return initialBoard
}

//ReadRowColStrategy takes fileContent as input. It returns the numRow, numCol and all strategies of the cells for the initial game board.
func ReadRowColStrategy(fileContent []string) (int,int,[]string) {
	numRow,errR := strconv.Atoi(strings.Split(fileContent[0]," ")[0])
	if errR != nil {
		fmt.Println("Wrong in transform string row to integer")
		os.Exit(2)
	}

	numCol,errC := strconv.Atoi(strings.Split(fileContent[0]," ")[1])
	if errC != nil {
		fmt.Println("Wrong in transform string column to integer")
		os.Exit(2)
	}
	return numRow,numCol,fileContent[1:]
}
