package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
)

func DrawCheckerBoard(board CheckerBoard) image.Image {
	height := len(board)
	width := len(board[0])
	c := CreateNewCanvas(width, height)
	// if height != width {
	// 	fmt.Println("Board's height is not equal to width")
	// 	os.Exit(3)
	// }
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

func DrawGridLines(pic Canvas, cellWidth int) {
	w, h := pic.width, pic.height
	// first, draw vertical lines
	for i := 1; i < pic.width/cellWidth; i++ {
		y := i * cellWidth
		pic.MoveTo(0.0, float64(y))
		pic.LineTo(float64(w), float64(y))
	}
	// next, draw horizontal lines
	for j := 1; j < pic.height/cellWidth; j++ {
		x := j * cellWidth
		pic.MoveTo(float64(x), 0.0)
		pic.LineTo(float64(x), float64(h))
	}
	pic.Stroke()
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
