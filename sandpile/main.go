package main

import (
	"fmt"
	// "image"
	// "image/png"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

//CheckerBoard is the required game board type
type CheckerBoard [][]int

//parallel part
func SandpileMultiprocs(board CheckerBoard, numProcs int, input string) {
	c := make([]chan CheckerBoard, numProcs)
	for i := range c {
		c[i] = make(chan CheckerBoard, 1)
	}
	size := len(board)
	if input != "random" {
		board = Topple(board, size/2, size/2)
	}
	for !Indicator(board) {
		for i := 0; i < numProcs; i++ {
			if i == numProcs-1 {
				go SandpileProc(Fillin(board[i*(size/numProcs):], size), c[i], i, numProcs)
			} else {
				go SandpileProc(Fillin(board[i*(size/numProcs):(i+1)*(size/numProcs)], size), c[i], i, numProcs)
			}
		}
		for i := 0; i < numProcs; i++ {
			CombineBoards(board, c[i], i, numProcs)
		}
	}
}

func CopyBoard(board CheckerBoard) CheckerBoard {
	size := len(board)
	newBoard := make(CheckerBoard, size)

	for i := range newBoard {
		newBoard[i] = make([]int, size)
	}
	copy(newBoard, board)
	return newBoard
}

func Fillin(board CheckerBoard, size int) CheckerBoard {
	newBoard := make(CheckerBoard, 2)
	newBoard[0] = make([]int, size)
	newBoard[1] = make([]int, size)
	return append(append(newBoard[0:1], board...), newBoard[1])
}

func SandpileProc(board CheckerBoard, c chan CheckerBoard, k, numProcs int) {
	for !Indicator(board[1 : len(board)-1]) {
		for i := 1; i < len(board)-1; i++ {
			for j := range board[i] {
				board = Topple(board, i, j)
			}
		}
	}
	if k == 0 {
		c <- board[len(board)-1 : len(board)]
	} else if k == numProcs-1 {
		c <- board[0:1]
		return
	} else {
		c <- append(board[0:1], board[len(board)-1])
	}
}

func CombineBoards(board CheckerBoard, c chan CheckerBoard, i, numProcs int) {
	tmpBoard := <-c
	size := len(board)
	if i == 0 {
		for j := 0; j < size; j++ {
			board[(i+1)*(size/numProcs)][j] += tmpBoard[0][j]
		}
		return
	}
	if i == numProcs-1 {
		for j := 0; j < size; j++ {
			board[i*(size/numProcs)-1][j] += tmpBoard[0][j]
		}
		return
	}
	for j := 0; j < size; j++ {
		board[i*(size/numProcs)-1][j] += tmpBoard[0][j]
		board[(i+1)*(size/numProcs)][j] += tmpBoard[1][j]
	}
}

//this part is serial program

//InitialBoard takes size, pile and placement as input and returns a initialized board.
func InitialBoard(size, pile int, placement string) CheckerBoard {
	board := make(CheckerBoard, size)
	for i := range board {
		board[i] = make([]int, size)
	}

	if placement == "central" {
		board[size/2][size/2] = pile
	}

	if placement == "random" {
		randNums := make([]int, 3)
		remain := pile
		for i := 0; i < 100; i++ {
			randNums[0] = rand.Intn(size)
			randNums[1] = rand.Intn(size)
			randNums[2] = rand.Intn(remain + 1)
			board[randNums[0]][randNums[1]] = randNums[2]
			remain -= randNums[2]
			if i == 99 && remain != 0 {
				board[randNums[0]][randNums[1]] += remain
			}
		}
	}
	return board
}

//
func ToppleBoardOnce(board CheckerBoard) CheckerBoard {
	for i := range board {
		for j := range board[i] {
			board = Topple(board, i, j)
		}
	}
	return board
}

//ToppleBoard does topple on a CheckerBoard iteratively until every pile is smaller than 4
func ToppleBoard(board CheckerBoard) CheckerBoard {
	for !Indicator(board) {
		for i := range board {
			for j := range board[i] {
				board = Topple(board, i, j)
			}
		}
	}
	return board
}

func Indicator(board CheckerBoard) bool {
	for i := range board {
		for j := range board[i] {
			if board[i][j] >= 4 {
				return false
			}
		}
	}
	return true
}

//Topple is doing topple on a board cell
func Topple(board CheckerBoard, i, j int) CheckerBoard {
	row := len(board)
	col := len(board[0])
	for board[i][j] >= 4 {
		board[i][j] -= 4
		if Infield(row, col, i-1, j) {
			board[i-1][j]++
		}
		if Infield(row, col, i+1, j) {
			board[i+1][j]++
		}
		if Infield(row, col, i, j-1) {
			board[i][j-1]++
		}
		if Infield(row, col, i, j+1) {
			board[i][j+1]++
		}
	}
	return board
}

//Infield checks whether a index is inside the board
func Infield(row, col, i, j int) bool {
	if i >= 0 && i < row && j >= 0 && j < col {
		return true
	}
	return false
}

func CheckBoard(sBoard, pBoard CheckerBoard) bool {
	for i := range sBoard {
		for j := range pBoard {
			if sBoard[i][j] != pBoard[i][j] {
				return false
			}
		}
	}
	return true
}

func RunSerialParallel(size, pile int, placement string) (CheckerBoard, CheckerBoard) {
	numProcs := runtime.NumCPU()
	sBoard := InitialBoard(size, pile, placement)
	pBoard := CopyBoard(sBoard)
	ToppleBoard(sBoard)
	SandpileMultiprocs(pBoard, numProcs, placement)
	return sBoard, pBoard
}

func main() {
	rand.Seed(time.Now().Unix())
	sizeStr := os.Args[1]
	size, err1 := strconv.Atoi(sizeStr)
	if err1 != nil {
		fmt.Println("Something wrong in converting string to integer on size")
		os.Exit(1)
	}

	pileStr := os.Args[2]
	pile, err2 := strconv.Atoi(pileStr)
	if err2 != nil {
		fmt.Println("Something wrong in converting string to integer on pile")
		os.Exit(2)
	}

	placement := os.Args[3]
	fmt.Println("Arguments have been read")
	if !(placement == "central" || placement == "random") {
		fmt.Println("Wrong input in initialization mode")
		os.Exit(4)
	}

	numProcs := runtime.NumCPU()
	sBoard := InitialBoard(size, pile, placement)
	pBoard := CopyBoard(sBoard)
	fmt.Println("Board has been initialized")

	start := time.Now()

	ToppleBoard(sBoard)
	fmt.Println("End of toppling the board")

	elapsed := time.Since(start)
	fmt.Println("Serial program takes", elapsed)

	sBoardImg := DrawCheckerBoard(sBoard)
	fmt.Println("Image created")

	ImageToPNG(sBoardImg, "serial")
	fmt.Println("PNG created")

	//------------------------------------------------------------------------------------
	start2 := time.Now()

	SandpileMultiprocs(pBoard, numProcs, placement)
	fmt.Println("End of toppling the board")

	elapsed2 := time.Since(start2)
	fmt.Println("Parallel program takes", elapsed2)

	pBoardImg := DrawCheckerBoard(pBoard)
	fmt.Println("Image created")

	ImageToPNG(pBoardImg, "parallel")
	fmt.Println("PNG created")

	if CheckBoard(sBoard, pBoard) == true {
		fmt.Println("correct")
	} else {
		fmt.Println("something wrong with parallel")
	}
}
