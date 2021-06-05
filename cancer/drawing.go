package main

import (
	"image"
)

//DrawCellBoards takes a slice of cellboards and cellWidth and returns a slice of image
func DrawCellBoards(boards []CellBoard, cellWidth int) []image.Image {
	numGenerations := len(boards)
	imageList := make([]image.Image, numGenerations)
	for i := range boards {
		imageList[i] = DrawCellBoard(boards[i], cellWidth)
	}
	return imageList
}

//DrawCellBoard takes a cellBoard and cellWidth as input. It returns a image of the board
func DrawCellBoard(board CellBoard, cellWidth int) image.Image {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewCanvas(width, height)

	// declare colors
	white := MakeColor(255, 255, 255)
	gray1 := MakeColor(200, 200, 200)
	gray2 := MakeColor(130, 130, 130)
	gray3 := MakeColor(80, 80, 80)
	gray4 := MakeColor(40, 40, 40)
	black := MakeColor(0, 0, 0)
	purple := MakeColor(150, 150, 255)
	lightBlue := MakeColor(100, 200, 255)

	// fill in colored squares
	for i := range board {
		for j := range board[i] {
			if board[i][j].alive == true {
				numHall := NumCancerHall(board[i][j])
				if numHall == 0 {
					if board[i][j].types.stem == true {
						c.SetFillColor(purple)
					} else {
						c.SetFillColor(lightBlue)
					}
				} else if numHall == 1 {
					c.SetFillColor(gray1)
				} else if numHall == 2 {
					c.SetFillColor(gray2)
				} else if numHall == 3 {
					c.SetFillColor(gray3)
				} else if numHall == 4 {
					c.SetFillColor(gray4)
				} else if numHall == 5 {
					c.SetFillColor(black)
				} else {
					panic("Error: Out of range value " + string(numHall) + " in board when drawing board.")
				}
			} else {
				c.SetFillColor(white)
			}
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}
	return c.img
}

//NumCancerHall takes a cell as input and calculate the number of hallmarks related to cancer
func NumCancerHall(cell Cell) int {
	count := 0
	if cell.hallmark.nonApoptosis == true {
		count++
	}
	if cell.hallmark.proliferation== true {
		count++
	}
	if cell.hallmark.nonGrowthInhibit == true {
		count++
	}
	if cell.hallmark.instability == true {
		count++
	}
	if cell.hallmark.unDivid == true {
		count++
	}
	return count
}

//DrawGridLines returns a bunch of lines the separate the grids
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
