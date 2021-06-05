package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

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
	fmt.Println("You have", numProcs, "CPUS available")

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

}

//CheckerBoard is the required game board type
type CheckerBoard [][]int

//parallel part

func SandpileMultiprocs(board CheckerBoard, numProcs int, input string) {
	size := len(board)
	var wg sync.WaitGroup
	// if input != "random" {
	// 	Topple(board, size/2, size/2)
	// }
	for !Indicator(board) {
		for i := 0; i < numProcs; i++ {
			if i == numProcs-1 {
				wg.Add(1)
				go SandpileProc(board[i*(size/numProcs):], i, numProcs)
				wg.Done()
			} else {
				wg.Add(1)
				go SandpileProc(board[i*(size/numProcs):(i+1)*(size/numProcs)], i, numProcs)
				wg.Done()
			}
		}
		wg.Wait()
		// for i := 0; i < numProcs; i++ {
		// 	CombineBoards(board, c)
		// }
		for j := 0; j < numProcs-1; j++ {
			// wg.Add(1)
			SandpileProc(board[(j+1)*(size/numProcs)-2:(j+1)*(size/numProcs)+2], -1, numProcs)
			// wg.Done()
		}
		// wg.Wait()
	}
}

func CopyBoard(board CheckerBoard) CheckerBoard {
	size := len(board)
	newBoard := make(CheckerBoard, size)

	for i := range newBoard {
		newBoard[i] = make([]int, size)
		for j := range newBoard[i] {
			newBoard[i][j] = board[i][j]
		}
	}
	return newBoard
}

func Fillin(board CheckerBoard, start, end, size int) CheckerBoard {
	newBoard := make(CheckerBoard, size)
	for i := range newBoard {
		newBoard[i] = make([]int, size)
	}
	for i := start; i < end; i++ {
		newBoard[i] = board[i-start]
	}
	return newBoard
}

func SandpileProc(board CheckerBoard, i, numProcs int) {
	// startt := time.Now()
	if i == 0 {
		for !Indicator(board[:len(board)-1]) {
			for k := 0; k < len(board)-1; k++ {
				for j := range board[k] {
					Topple(board, k, j)
				}
			}
		}
	} else if i == numProcs-1 {
		for !Indicator(board[1:]) {
			for k := 1; k < len(board); k++ {
				for j := range board[k] {
					Topple(board, k, j)
				}
			}
		}
	} else {
		for !Indicator(board[1 : len(board)-1]) {
			for k := 1; k < len(board)-1; k++ {
				for j := range board[k] {
					Topple(board, k, j)
				}
			}
		}
	}

	// elpased := time.Since(startt)
	// fmt.Println("Single proc process for toppling takes", elpased)
}

func CombineBoards(board CheckerBoard, c chan CheckerBoard) {
	tmpBoard := <-c
	//size := len(board)
	for i := range board {
		for j := range board {
			board[i][j] += tmpBoard[i][j]
		}
	}
	// if i == 0 {
	// 	for j := 0; j < size; j++ {
	// 		board[(i+1)*size/numProcs][j] += tmpBoard[0][j]
	// 	}
	// 	return
	// }
	// if i == numProcs-1 {
	// 	for j := 0; j < size; j++ {
	// 		board[i*size/numProcs-1][j] += tmpBoard[1][j]
	// 	}
	// 	return
	// }
	// for j := 0; j < size; j++ {
	// 	board[i*size/numProcs-1][j] += tmpBoard[0][j]
	// 	board[(i+1)*size/numProcs][j] += tmpBoard[1][j]
	// }
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
func ToppleBoardOnce(board CheckerBoard) {
	for i := range board {
		for j := range board[i] {
			Topple(board, i, j)
		}
	}
}

//ToppleBoard does topple on a CheckerBoard iteratively until every pile is smaller than 4
func ToppleBoard(board CheckerBoard) {
	for !Indicator(board) {
		for i := range board {
			for j := range board[i] {
				Topple(board, i, j)
			}
		}
	}
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
func Topple(board CheckerBoard, i, j int) {
	for board[i][j] >= 4 {
		board[i][j] -= 4
		row := len(board)
		col := len(board[0])
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
}

//Infield checks whether a index is inside the board
func Infield(row, col, i, j int) bool {
	if i >= 0 && i < row && j >= 0 && j < col {
		return true
	}
	return false
}
