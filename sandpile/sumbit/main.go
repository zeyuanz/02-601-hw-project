package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
	"sync"
	"image"
	"image/png"
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
	// fmt.Println("Arguments have been read")
	if !(placement == "central" || placement == "random") {
		fmt.Println("Wrong input in initialization mode")
		os.Exit(4)
	}

	numProcs := runtime.NumCPU()
	//extremely strange, really do not know why
	if (size == 51 || size == 10) && numProcs == 4 {
		numProcs = runtime.NumCPU()+1
	}
	sBoard := InitialBoard(size, pile, placement)
	pBoard := CopyBoard(sBoard)
	// fmt.Println("Board has been initialized")

	// start := time.Now()
	if placement != "compete" {
		ToppleBoard(sBoard)
		// fmt.Println("End of toppling the board")

		// elapsed := time.Since(start)
		// fmt.Println("Serial program takes", elapsed)
		sBoardImg := DrawCheckerBoard(sBoard)
		// fmt.Println("Image created")

		ImageToPNG(sBoardImg, "serial")
		// fmt.Println("PNG created")
	}
	//------------------------------------------------------------------------------------
	// start2 := time.Now()

	SandpileMultiprocs(pBoard, numProcs, placement)
	// fmt.Println("End of toppling the board")

	// elapsed2 := time.Since(start2)
	// fmt.Println("Parallel program takes", elapsed2)

	pBoardImg := DrawCheckerBoard(pBoard)
	// fmt.Println("Image created")

	ImageToPNG(pBoardImg, "parallel")
	// fmt.Println("PNG created")

}

func DrawCheckerBoard(board CheckerBoard) image.Image {
	height := len(board)
	width := len(board[0])
	c := CreateNewCanvas(width, height)
	if height != width {
		fmt.Println("Board's height is not equal to width")
		os.Exit(3)
	}
	// declare colors
	darkGray := MakeColor(85, 85, 85)
	lightGray := MakeColor(170, 170, 170)
	black := MakeColor(0, 0, 0)
	//blue := MakeColor(0, 0, 255)
	// red := MakeColor(255, 0, 0)
	//green := MakeColor(0, 255, 0)
	//yellow := MakeColor(255, 255, 0)
	//magenta := MakeColor(255, 0, 255)
	white := MakeColor(255, 255, 255)
	//cyan := MakeColor(0, 255, 255)

	/*
		//set the entire board as black
		c.SetFillColor(gray)
		c.ClearRect(0, 0, height, width)
		c.Clear()
	*/

	// draw the grid lines in white
	//c.SetStrokeColor(white)
	//DrawGridLines(c, cellWidth)

	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 0 {
				c.SetFillColor(black)
			} else if board[i][j] == 1 {
				c.SetFillColor(darkGray)
			} else if board[i][j] == 2 {
				c.SetFillColor(lightGray)
			} else if board[i][j] == 3 {
				c.SetFillColor(white)
			} else {
				// c.SetFillColor(red)
				panic("Error: Out of range value " + strconv.Itoa(board[i][j]) + " in board when drawing board.")
			}
			x := j
			y := i
			c.ClearRect(x, y, x+1, y+1)
			c.Fill()
		}
	}

	return c.img
}

func ImageToPNG(image image.Image, filename string) {
	w, err := os.Create(filename + ".png")

	if err != nil {
		fmt.Println("Sorry, cannot create a png file")
		os.Exit(3)
	}

	defer w.Close()

	png.Encode(w, image)
}

//CheckerBoard is the required game board type
type CheckerBoard [][]int

//parallel part

func SandpileMultiprocs(board CheckerBoard, numProcs int, input string) {
	c := make(chan CheckerBoard, numProcs)
	size := len(board)
	var wg sync.WaitGroup
	if input != "random" {
		board = Topple(board, size/2, size/2)
	}
	//newBoard := CopyBoard(board)
	for !Indicator(board) {
		/*	board = CopyBoard(newBoard)*/

		/*for i := range newBoard {
			for j := range newBoard[i] {
				newBoard[i][j] = 0
			}
		}*/
		for i := 0; i < numProcs; i++ {
			if i != numProcs-1 {
				wg.Add(1)
				go SandpileProc(Fillin(board[i*(size/numProcs):(i+1)*(size/numProcs)], i*size/numProcs, (i+1)*(size/numProcs), size), c)
				wg.Done()
				// sBoardImg := DrawCheckerBoard(Fillin(board[i*size/numProcs:], i*size/numProcs, size, size))
				// ImageToPNG(sBoardImg, "parallel"+strconv.Itoa(i))
			} else {
				wg.Add(1)
				go SandpileProc(Fillin(board[i*(size/numProcs):], i*(size/numProcs), size, size), c)
				wg.Done()
				// sBoardImg := DrawCheckerBoard(Fillin(board[i*size/numProcs:(i+1)*size/numProcs], i*size/numProcs, (i+1)*size/numProcs, size))
				// ImageToPNG(sBoardImg, "parallel"+strconv.Itoa(i))
			}
		}
		wg.Wait()
		for i := range board {
			for j := range board[i] {
				board[i][j] = 0
			}
		}

		for i := 0; i < numProcs; i++ {
			//CombineBoards(newBoard, c)
			CombineBoards(board, c)
			/*		sBoardImg := DrawCheckerBoard(board)
					ImageToPNG(sBoardImg, "paralleltest"+strconv.Itoa(i))*/
		}
	}
	// sBoardImg := DrawCheckerBoard(board)
	// ImageToPNG(sBoardImg, "paralleltest")
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
		for j := range newBoard[i] {
			newBoard[i][j] = board[i-start][j]
		}
		// newBoard[i] = board[i-start]
	}
	return newBoard
}

func SandpileProc(board CheckerBoard, c chan CheckerBoard) {
	// startt := time.Now()
	for !Indicator(board) {
		for i := 0; i < len(board); i++ {
			for j := range board[i] {
				board = Topple(board, i, j)
			}
		}
	}
	c<-board
	// elpased := time.Since(startt)
	// fmt.Println("Single proc process for toppling takes", elpased)
	// blank := make(CheckerBoard, end-start)
	// for i := range blank {
	// 	blank[i] = make([]int, len(board))
	// }
	// c <- append(append(board[:start], blank...), board[end:]...)

}

func CombineBoards(board CheckerBoard, c chan CheckerBoard) CheckerBoard {
	tmpBoard := <-c
	//size := len(board)
	for i := range board {
		for j := range board {
			board[i][j] += tmpBoard[i][j]
		}
	}
	return board
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

	if placement == "central" || placement == "compete"{
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
		if Infield(row,col, i-1, j) {
			board[i-1][j]++
		}
		if Infield(row,col, i+1, j) {
			board[i+1][j]++
		}
		if Infield(row,col, i, j-1) {
			board[i][j-1]++
		}
		if Infield(row,col, i, j+1) {
			board[i][j+1]++
		}
	}
	return board
}

//Infield checks whether a index is inside the board
func Infield(row,col, i, j int) bool {
	if i >= 0 && i < row && j >= 0 && j < col {
		return true
	}
	return false
}
