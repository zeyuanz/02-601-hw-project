package main

import (
	"image/color"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type errPoints struct {
	plotter.XYs
	plotter.YErrors
	plotter.XErrors
}

//CalBoards takes a slice of boards and produce two figures.
//One is Proportion of different cell types among all cells.
//The other is Proportion of different cancer hallmarks among all tumor cells.
//Mean value is counted to eliminate randomness, variance is also considered, but the visulization effects are too bad.
func CalBoards(boards [][]CellBoard, filename string, interval int) {
	numCells := make([][][]float64, 3)
	for i := range numCells {
		numCells[i] = make([][]float64, len(boards[0]))
		for j := range numCells[i] {
			numCells[i][j] = make([]float64, len(boards))
		}
	}

	for i := range numCells[0] {
		for j := range numCells[0][i] {
			numCell := CalculateCells(boards[j][i], true)
			numCells[0][i][j] = numCell[0]
			numCells[1][i][j] = numCell[1]
			numCells[2][i][j] = numCell[2]
		}
	}

	meanNumCell := make([][]float64, 3)
	varianceNumCell := make([][]float64, 3)
	for i := range meanNumCell {
		meanNumCell[i] = make([]float64, len(boards[0]))
		varianceNumCell[i] = make([]float64, len(boards[0]))
	}

	for i := range numCells {
		for j := range numCells[0] {
			meanNumCell[i][j], varianceNumCell[i][j] = MeanVariance(numCells[i][j])
		}
	}
	PlotLineGraph(len(boards[0]), meanNumCell, varianceNumCell, filename+"NumCells", interval,"C")

	numCHs := make([][][]float64, 5)
	for i := range numCHs {
		numCHs[i] = make([][]float64, len(boards[0]))
		for j := range numCHs[i] {
			numCHs[i][j] = make([]float64, len(boards))
		}
	}

	for i := range numCHs[0] {
		for j := range numCHs[0][i] {
			numCH := CalCancerHallmarks(boards[j][i], true)
			numCHs[0][i][j] = numCH[0]
			numCHs[1][i][j] = numCH[1]
			numCHs[2][i][j] = numCH[2]
			numCHs[3][i][j] = numCH[3]
			numCHs[4][i][j] = numCH[4]
		}
	}

	meanCH := make([][]float64, 5)
	varianceCH := make([][]float64, 5)
	for i := range meanCH {
		meanCH[i] = make([]float64, len(boards[0]))
		varianceCH[i] = make([]float64, len(boards[0]))
	}

	for i := range numCHs {
		for j := range numCHs[0] {
			meanCH[i][j], varianceCH[i][j] = MeanVariance(numCHs[i][j])
		}
	}
	PlotLineGraph(len(boards[0]), meanCH, varianceCH, filename+"NumCHs", interval,"H")
}

//CalculateCells takes a cellboard as input and returns the number of stem cells, differentiated cells and cancer cells
//the order is stem, normal and cancer cells
func CalculateCells(board CellBoard, percentage bool) []float64 {
	countS, countN, countC := 0.0, 0.0, 0.0
	nums := make([]float64, 3)
	for i := range board {
		for j := range board[i] {
			if board[i][j].alive == true {
				if NumCancerHall(board[i][j]) != 0 {
					countC += 1.0
				} else if board[i][j].types.stem == true {
					countS += 1.0
				} else {
					countN += 1.0
				}
			}
		}
	}
	if percentage == false {
		nums[0] = countS // blue
		nums[1] = countN // red
		nums[2] = countC // green
	} else {
		row := float64(len(board))
		col := float64(len(board[0]))
		nums[0] = countS / (row * col)
		nums[1] = countN / (row * col)
		nums[2] = countC / (row * col)
	}
	return nums
}

//CalCancerHallmarks takes a cellboard as input and return a slice of int of number of cancer hallmarks of all cells
//the order is nonGrowthInhibit, onApoptosis, proliferation, instability, unDivid
func CalCancerHallmarks(board CellBoard, percentage bool) []float64 {
	nums := make([]float64, 5)
	totalCancer := 0.0
	for i := range board {
		for j := range board[i] {
			if board[i][j].alive == true {
				tmp := 0.0
				if board[i][j].hallmark.nonGrowthInhibit == true { // blue
					nums[0] += 1.0
					tmp += 1.0
				}
				if board[i][j].hallmark.nonApoptosis == true { // red
					nums[1] += 1.0
					tmp += 1.0
				}
				if board[i][j].hallmark.proliferation == true { // green
					nums[2] += 1.0
					tmp += 1.0
				}
				if board[i][j].hallmark.instability == true { // yellow
					nums[3] += 1.0
					tmp += 1.0
				}
				if board[i][j].hallmark.unDivid == true { // cyan
					nums[4] += 1.0
					tmp += 1.0
				}
				totalCancer += math.Min(1.0, tmp)
			}
		}
	}
	if percentage == false {
		return nums
	}

	if totalCancer != 0 {
		for i := range nums {
			nums[i] /= totalCancer
		}
	}

	return nums
}

//MeanVariance takes a slice of integers as input and returns the mean and variance of the slice
func MeanVariance(nums []float64) (float64, float64) {
	sum := 0.0
	squareSum := 0.0
	for i := range nums {
		sum += nums[i]
		squareSum += nums[i] * nums[i]
	}
	mean := sum / float64(len(nums))
	variance := squareSum/float64(len(nums)) - mean*mean
	return mean, variance
}

//PlotLineGraph takes x and a slice of y as input and plot a line graph with error bar. x is a integer, the x-axis is 1...x
func PlotLineGraph(x int, y [][]float64, variance [][]float64, filename string, interval int, types string) error {

	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
	green := MakeColor(0, 255, 0)
	yellow := MakeColor(255, 255, 0)
	//magenta := MakeColor(255, 0, 255)
	cyan := MakeColor(0, 255, 255)

	colors := make([]color.Color, 5)
	colors[0], colors[1], colors[2], colors[3], colors[4] = blue, red, green, yellow, cyan

	if x != len(y[0]) {
		log.Fatal("Input x is not equal to y in size.")
	}
	//this function takes y value and variance as input.
	//it returns the low and high error value of y value, which is y-variance and y+variance
	yError := func(y []float64, vari []float64, interval int) plotter.Errors {
		yErr := make(plotter.Errors, len(y)/interval+1)
		for i := range vari {
			if i%interval == 0 {
				yErr[i/interval].Low = 0.0 - vari[i]
				yErr[i/interval].High = vari[i]
			}
		}
		return yErr
	}

	p, err := plot.New()
	if err != nil {
		return err
	}

	pts := make([]plotter.XYs, len(y))

	for i := range y {
		if interval == 1 {
			pts[i] = make(plotter.XYs, x)
		} else {
			pts[i] = make(plotter.XYs, x/interval+1)

		}
		for j := range pts[i] {
			pts[i][j].X = float64(j * interval)
			pts[i][j].Y = y[i][j*interval]
		}
	}

	data := make([]errPoints, len(y))

	for i := range data {
		data[i] = errPoints{
			XYs:     pts[i],
			YErrors: plotter.YErrors(yError(y[i], variance[i], interval)),
		}
	}

	if types == "C" {
		for i := range pts {
			line, err := plotter.NewLine(pts[i])
			if err != nil {
				return err
			}

			//yerrs, err := plotter.NewYErrorBars(data[i])
			if err != nil {
				return err
			}
			switch i {
				case 0: line.Color = colors[0]
								p.Add(line)
								p.Legend.Add("Stem Cell",line)
				case 1: line.Color = colors[1]
								p.Add(line)
								p.Legend.Add("Differentiated Cell",line)
				case 2: line.Color = colors[2]
								p.Add(line)
								p.Legend.Add("Tumor Cell",line)
			}
		}

		p.Y.Max = 1.0
		p.Y.Min = 0.0

		p.Y.Label.Text = "Proportion"
		p.X.Label.Text = "Number of Generations"
		p.Title.Text = "Proportion of Different Types of Cells"
		if err := p.Save(13*vg.Inch, 5*vg.Inch, filename+".png"); err != nil {
			return err
		}
	}

	if types == "H" {
		for i := range pts {
			line, err := plotter.NewLine(pts[i])
			if err != nil {
				return err
			}

			//yerrs, err := plotter.NewYErrorBars(data[i])
			if err != nil {
				return err
			}
			switch i {
				case 0: line.Color = colors[0]
								p.Add(line)
								p.Legend.Add("nonGrowthInhibit",line)
				case 1: line.Color = colors[1]
								p.Add(line)
								p.Legend.Add("nonApoptosis",line)
				case 2: line.Color = colors[2]
								p.Add(line)
								p.Legend.Add("proliferation",line)
			  case 3: line.Color = colors[3]
								p.Add(line)
								p.Legend.Add("instability",line)
				case 4: line.Color = colors[4]
								p.Add(line)
								p.Legend.Add("unLimitedDivid",line)
			}
		}

		p.Y.Max = 1.0
		p.Y.Min = 0.0

		p.Y.Label.Text = "Proportion"
		p.X.Label.Text = "Number of Generations"
		p.Title.Text = "Proportion of Different Tumor Features in All Tumor Cells"
		if err := p.Save(13*vg.Inch, 5*vg.Inch, filename+".png"); err != nil {
			return err
		}
	}

	return err
}
