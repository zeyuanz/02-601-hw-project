package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var mutationRate = 1.0 / 100000.0 //default value from paper
var genIntability = 100.0         //default value rate from paper
var ranApop = 1.0 / 1000.0        //default value rate from paper
var evaApop = 0.0                 //default value rate from paper
var teloLen = 50                  //default value rate from paper
var ignor = 1.0 / 10.0            //default value rate from paper
var divid = 1.0                   //default value of a cell that could undergo division
var dif = 1.0                     //default prob of a stem cell differentiates into a differentiated cell
var ranDe = 1.0 / 1000.0          //default prob for necrosis because of physical, chemical etc.

//CellBoard is a struct to describe the distribution of cells
type CellBoard [][]Cell

//Cell is a struct to describe the relative parameteres and qualities of a cell
type Cell struct {
	types     Types
	parameter par
	hallmark  mark
	alive     bool
}

//Types is a struct to represent the type of the cells. It could be stem cell or differentiated cell
type Types struct {
	stem bool
}

//par is a struct to describe the parameteres of the cells
type par struct {
	baseMutationRate   float64 //initialize as mutationRate
	geneticInstability float64 //initialize as 100.0
	ranDeath           float64 //initialize as 1.0/1000.0
	ranApoptosis       float64 //initialize as 1.0/1000.0
	evaApoptosis       float64 //initialize as n/10 (extra likihood to death) unless nonApoptosis is true
	teloLength         int     //initialize as 100
	ignorance          float64 //initialize as 1.0/10.0
	division           float64 //initialized as 1.0
}

//mark is a struct to describe to relevant cancer hallmarks of cell
type mark struct {
	nonGrowthInhibit bool //initialize as false. If it is true, it could grow without limitation of growth factor.
	nonApoptosis     bool //initialize as false. If it is true, it will not die because of apoptosis.
	proliferation    bool //initialize as false. If it is true, it could kill other cell and will not stop division..
	instability      bool //initialize as false. If it is true, it will cause a increase in muatation rate.
	unDivid          bool //initialized as false. If it is true, it enables cell to undergo unlimited cell division process.
	otherMuatation   int  //initialize as zero. When obtain other mutation, it will cause extra likelihood to die unless nonApoptosis is obtained. Here we say must is 10, because it will die.
}

func main() {
	/*  Arguments :
			1. 	"graph" or "stats"
	    2.  number of rows
	    3.  number of cols
	    4.  number of generations
	    5.  width of the grid (if graph) repetition times (if stats)
	    6.  output file name
			7.	intervals in the graph (if stats)
	*/
	rand.Seed(time.Now().Unix())

	// dealing with input arguments
	inputMode := os.Args[1]
	if inputMode != "default" && inputMode != "user" {
		panic("Wrong mode!")
	}
	if inputMode == "default" {
		mode := os.Args[2]
		if mode != "graph" && mode != "stats" {
			panic("Wrong mode!")
		}
		row := os.Args[3]
		numRow, err1 := strconv.Atoi(row)
		if err1 != nil {
			panic("something wrong in converting row to integer")
		}

		col := os.Args[4]
		numCol, err2 := strconv.Atoi(col)
		if err2 != nil {
			panic("something wrong in converting colunmn to integer")
		}

		Gens := os.Args[5]
		numGens, err3 := strconv.Atoi(Gens)
		if err3 != nil {
			panic("something wrong in converting number of generations to integer")
		}

		var cellWidth, numRep int
		var err4, err5 error

		if mode == "graph" {
			width := os.Args[6]
			cellWidth, err4 = strconv.Atoi(width)
			if err4 != nil {
				panic("something wrong in converting width of cell to integer")
			}
		}

		if mode == "stats" {
			rep := os.Args[6]
			numRep, err5 = strconv.Atoi(rep)
			if err5 != nil {
				panic("something wrong in converting number of repetitions to integer")
			}
		}

		fileName := os.Args[7]

		fmt.Println("Parameters has been set to simulate the tumorigenesis and tumor growth.")

		initialBoard := InitialBoard(numRow, numCol)
		for i := numRow/2 - 2; i <= numRow/2+2; i++ {
			for j := numCol/2 - 2; j <= numCol/2+2; j++ {
				initialBoard[i][j] = InitialStemCell()
			}
		}

		// two modes, here is graph mode
		if mode == "graph" {

			fmt.Println("InitialBoard has been created.")

			boards := UpdateBoards(initialBoard, numGens)

			fmt.Println("Simulation has done. Now create the images.")

			images := DrawCellBoards(boards, cellWidth)

			fmt.Println("Images has been created. Now transfer to GIF.")

			ImagesToGIF(images, fileName)

			fmt.Println("GIF has been created.")

		}

		// here is stats mode
		if mode == "stats" {
			interval, err6 := strconv.Atoi(os.Args[8])
			if err6 != nil {
				log.Fatal(err6)
			}
			boards := make([][]CellBoard, numRep)
			for i := 0; i < numRep; i++ {
				boards[i] = UpdateBoards(initialBoard, numGens)
			}
			fmt.Println(numRep, "times repetition Finished.")

			CalBoards(boards, fileName, interval)

			fmt.Println("Line graph has been produced!")
		}
	}

	if inputMode == "user" {
		mode := os.Args[2]
		if mode != "graph" && mode != "stats" {
			panic("Wrong mode!")
		}

		settings := os.Args[3]
		initialBoard, para := ReadBoardFromFile(settings)
		fmt.Println(len(initialBoard), len(initialBoard[0]))
		mutationRate, _ = strconv.ParseFloat(para[0], 64)
		genIntability, _ = strconv.ParseFloat(para[1], 64)
		ranApop, _ = strconv.ParseFloat(para[2], 64)
		evaApop, _ = strconv.ParseFloat(para[3], 64)
		teloLen, _ = strconv.Atoi(para[4])
		ignor, _ = strconv.ParseFloat(para[5], 64)
		divid, _ = strconv.ParseFloat(para[6], 64)
		dif, _ = strconv.ParseFloat(para[7], 64)
		ranDe, _ = strconv.ParseFloat(para[8], 64)

		Gens := os.Args[4]
		numGens, err3 := strconv.Atoi(Gens)
		if err3 != nil {
			panic("something wrong in converting number of generations to integer")
		}

		var cellWidth, numRep int
		var err4, err5 error

		if mode == "graph" {
			width := os.Args[5]
			cellWidth, err4 = strconv.Atoi(width)
			if err4 != nil {
				panic("something wrong in converting width of cell to integer")
			}
		}

		if mode == "stats" {
			rep := os.Args[5]
			numRep, err5 = strconv.Atoi(rep)
			if err5 != nil {
				panic("something wrong in converting number of repetitions to integer")
			}
		}

		fileName := os.Args[6]

		// two modes, here is graph mode
		if mode == "graph" {

			fmt.Println("InitialBoard has been created.")

			boards := UpdateBoards(initialBoard, numGens)

			fmt.Println("Simulation has done. Now create the images.")

			images := DrawCellBoards(boards, cellWidth)

			fmt.Println("Images has been created. Now transfer to GIF.")

			ImagesToGIF(images, fileName)

			fmt.Println("GIF has been created.")

		}

		// here is stats mode
		if mode == "stats" {
			interval, err6 := strconv.Atoi(os.Args[7])
			if err6 != nil {
				log.Fatal(err6)
			}
			boards := make([][]CellBoard, numRep)
			for i := 0; i < numRep; i++ {
				boards[i] = UpdateBoards(initialBoard, numGens)
			}
			fmt.Println(numRep, "times repetition Finished.")

			CalBoards(boards, fileName, interval)

			fmt.Println("Line graph has been produced!")
		}
	}
}

//InitialBoard takes number of rows and number of colunms as input and return a initialized CellBoard with defualt paratmeters.
func InitialBoard(numRow, numCol int) CellBoard {
	newBoard := make(CellBoard, numRow)
	for i := range newBoard {
		newBoard[i] = make([]Cell, numCol)
	}

	for i := range newBoard {
		for j := range newBoard[i] {
			newBoard[i][j] = InitialEmptyCell()
		}
	}
	return newBoard
}

//CopyBoard takes a CellBoard as input and returns a copy of the new Board
func CopyBoard(currBoard CellBoard) CellBoard {
	newBoard := InitialBoard(len(currBoard), len(currBoard[0]))
	for i := range newBoard {
		for j := range newBoard[i] {
			newBoard[i][j] = currBoard[i][j]
		}
	}
	return newBoard
}

//InitialEmptyCell takes nothing as input and returns an empyt Cell with all parameters initialized as defualt.
func InitialEmptyCell() Cell {
	var newCell Cell

	newCell.alive = false
	newCell.hallmark.nonApoptosis = false
	newCell.hallmark.proliferation = false
	newCell.hallmark.nonGrowthInhibit = false
	newCell.hallmark.instability = false
	newCell.hallmark.unDivid = false
	newCell.hallmark.otherMuatation = 0

	newCell.types.stem = false

	newCell.parameter.baseMutationRate = mutationRate
	newCell.parameter.geneticInstability = genIntability
	newCell.parameter.ranApoptosis = ranApop
	newCell.parameter.evaApoptosis = evaApop
	newCell.parameter.teloLength = teloLen
	newCell.parameter.ignorance = ignor
	newCell.parameter.division = divid
	newCell.parameter.ranDeath = ranDe

	return newCell
}

//InitialDifCell takes nothing as input and returns a differentiated cell.
func InitialDifCell() Cell {
	newCell := InitialEmptyCell()
	newCell.alive = true
	newCell.types.stem = false
	newCell.parameter.division = divid / 4.0
	return newCell
}

//InitialStemCell takes nothing as input and returns a stem cell as default.
func InitialStemCell() Cell {
	newCell := InitialEmptyCell()
	newCell.alive = true
	newCell.types.stem = true
	return newCell
}

//UpdateBoards takes a currBoard and numGens as input. It returns a slice of CellBoard of numGens + 1.
func UpdateBoards(currBoard CellBoard, numGens int) []CellBoard {
	cellBoards := make([]CellBoard, numGens+1)
	cellBoards[0] = CopyBoard(currBoard)
	for i := 1; i < numGens+1; i++ {
		cellBoards[i] = UpdateBoard(cellBoards[i-1])
	}
	return cellBoards
}

//UpdateBoard takes a CellBoard as input. It returns a CellBoard updated based on the event and parameteres of each cell.
//Here we use concurrency to implement some randomness, however it could still be optimized to more randomness which fit into the reality better.
//WaitGroup are used to make sure everything finishes before next event happen and to avoid weird things happen on the board.
func UpdateBoard(currBoard CellBoard) CellBoard {
	numRows := len(currBoard)
	numCols := len(currBoard[0])
	var wg sync.WaitGroup

	//update death event
	for i := 0; i < numRows; i++ {
		wg.Add(1)
		for j := 0; j < numCols; j++ {
			go func(i, j int) {
				currBoard[i][j] = UpdateDeath(currBoard[i][j])
			}(i, j)
		}
		wg.Done()
	}

	//update movement event
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			currBoard = UpdateMove(currBoard, i, j)
		}
	}

	//update division event
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			currBoard = UpdateDivision(currBoard, i, j)
		}
	}

	//update differentiation event
	for i := 0; i < numRows; i++ {
		wg.Add(1)
		for j := 0; j < numCols; j++ {
			if currBoard[i][j].types.stem == true && NumHall(currBoard[i][j]) == 0 {
				go func(i, j int) {
					currBoard[i][j] = UpdateDif(currBoard, i, j)
				}(i, j)
			}
		}
		wg.Done()
	}

	//update mutation event
	for i := 0; i < numRows; i++ {
		wg.Add(1)
		for j := 0; j < numCols; j++ {
			go func(i, j int) {
				currBoard[i][j] = UpdateMutation(currBoard[i][j])
			}(i, j)
		}
		wg.Done()
	}
	wg.Wait()
	return currBoard
}

//UpdateMove takes currBoard and the index i, j as input. It returns the updated status of a cell movement.
//It considers the feature of cancer hallmark and all the possible neighborhoods and choose one randomly.
func UpdateMove(currBoard CellBoard, i, j int) CellBoard {
	nextMove := 0
	newBoard := CopyBoard(currBoard)
	if currBoard[i][j].alive == true {
		//consider the feature of proliferation
		if currBoard[i][j].hallmark.proliferation == true {
			if GenerateProb(currBoard[i][j].parameter.ignorance) {
				//find the possible movement
				possMove := NextMove(currBoard, i, j, true)
				ranMove := rand.Intn(len(possMove))
				nextMove = possMove[ranMove]
				newBoard[i][j] = InitialEmptyCell()
				switch nextMove {
				case 0:
					newBoard[i][j] = currBoard[i][j]
				case 1:
					newBoard[i-1][j] = currBoard[i][j]
				case 2:
					newBoard[i][j+1] = currBoard[i][j]
				case 3:
					newBoard[i+1][j] = currBoard[i][j]
				case 4:
					newBoard[i][j-1] = currBoard[i][j]
				}
			} else {
				//find possible movement
				possMove := NextMove(currBoard, i, j, false)
				ranMove := rand.Intn(len(possMove))
				nextMove = possMove[ranMove]
				newBoard[i][j] = InitialEmptyCell()
				switch nextMove {
				case 0:
					newBoard[i][j] = currBoard[i][j]
				case 1:
					newBoard[i-1][j] = currBoard[i][j]
				case 2:
					newBoard[i][j+1] = currBoard[i][j]
				case 3:
					newBoard[i+1][j] = currBoard[i][j]
				case 4:
					newBoard[i][j-1] = currBoard[i][j]
				}
			}
		}
		//another case if without proliferation feature
		if currBoard[i][j].hallmark.proliferation == false {
			possMove := NextMove(currBoard, i, j, false)
			ranMove := rand.Intn(len(possMove))
			nextMove = possMove[ranMove]
			newBoard[i][j] = InitialEmptyCell()
			switch nextMove {
			case 0:
				newBoard[i][j] = currBoard[i][j]
			case 1:
				newBoard[i-1][j] = currBoard[i][j]
			case 2:
				newBoard[i][j+1] = currBoard[i][j]
			case 3:
				newBoard[i+1][j] = currBoard[i][j]
			case 4:
				newBoard[i][j-1] = currBoard[i][j]
			}
		}
	}
	return newBoard
}

//GenerateProb takes a float64 as input and regard it as the probability of being true in a single trial. Then perform the simulation and produce true or false based on the simulation.
func GenerateProb(prob float64) bool {
	ranNum := rand.Float64()
	if ranNum <= prob {
		return true
	}
	return false
}

//NextMove takes a currBoard and index i, j as input. It returns all the possible next movement position of the cell(i,j). Considering Von Neumann neighborhood, if any cell in the neighborhood is alive, then the position is invalid unless cell(i,j) is a cancer cell. This will be considered in the Update Move function.
//state is the bool value to decide whether a cell has proliferation feature
func NextMove(currBoard CellBoard, i, j int, state bool) []int {
	possMove := make([]int, 0)
	possMove = append(possMove, 0)
	//the case when a cell does not have proliferation feature
	if state == false {
		//find all possible place
		if InField(currBoard, i-1, j) && currBoard[i-1][j].alive == false {
			possMove = append(possMove, 1)
		}
		if InField(currBoard, i, j+1) && currBoard[i][j+1].alive == false {
			possMove = append(possMove, 2)
		}
		if InField(currBoard, i+1, j) && currBoard[i+1][j].alive == false {
			possMove = append(possMove, 3)
		}
		if InField(currBoard, i, j-1) && currBoard[i][j-1].alive == false {
			possMove = append(possMove, 4)
		}
	} else {
		//the case when a cell has proliferation feature
		if InField(currBoard, i-1, j) && currBoard[i-1][j].hallmark.proliferation == false {
			possMove = append(possMove, 1)
		}
		if InField(currBoard, i, j+1) && currBoard[i][j+1].hallmark.proliferation == false {
			possMove = append(possMove, 2)
		}
		if InField(currBoard, i+1, j) && currBoard[i+1][j].hallmark.proliferation == false {
			possMove = append(possMove, 3)
		}
		if InField(currBoard, i, j-1) && currBoard[i][j-1].hallmark.proliferation == false {
			possMove = append(possMove, 4)
		}
	}
	return possMove
}

//InField takes a currBoard and index i j as input. It returns whether index i j is within the currBoard.
func InField(currBoard CellBoard, i, j int) bool {
	numRow := len(currBoard)
	numCol := len(currBoard[0])
	if i >= 0 && i < numRow && j >= 0 && j < numCol {
		return true
	}
	return false
}

//UpdateDeath takes currBoard and the index i, j as input. It returns the updated status of a cell apoptosis.
//apoptosis is induced by the random number based on the probability of the cell.
//nonApoptosis cells can evade apoptosis process but they can still undergo necrosis.
func UpdateDeath(cell Cell) Cell {
	newCell := cell
	if cell.alive == true {
		//consider nonApoptosis feature, only prob for necrosis
		if cell.hallmark.nonApoptosis == true {
			if GenerateProb(cell.parameter.ranDeath) {
				newCell = InitialEmptyCell()
			}
			return newCell
		}
		//another case when nonApoptosis is false, prob for necrosis and apoptosis
		if cell.hallmark.nonApoptosis == false {
			cell.parameter.ranApoptosis = UpdateApoptosisProb(cell)
			prob := cell.parameter.ranApoptosis
			if GenerateProb(prob + cell.parameter.ranDeath) {
				newCell = InitialEmptyCell()
			}
		}
	}
	return newCell
}

//UpdateApoptosisProb takes a cell as input. It returns the udpate value of ranApoptosis accroding to the number of hallmarks.
//it updates the probability of apoptosis based on teloLength of each cell and number of mutations.
func UpdateApoptosisProb(cell Cell) float64 {
	numHall := NumHall(cell)
	cell.parameter.evaApoptosis = float64(numHall) / 10.0
	cell.parameter.ranApoptosis = ranApop + cell.parameter.evaApoptosis + 1.0 - float64(cell.parameter.teloLength)/float64(teloLen)
	return math.Min(cell.parameter.ranApoptosis, 1.0)
}

//NumHall takes a cell as input and return the number of hallmarks that equals true
func NumHall(cell Cell) int {
	count := 0
	if cell.hallmark.nonApoptosis == true {
		count++
	}
	if cell.hallmark.proliferation == true {
		count++
	}
	if cell.hallmark.nonGrowthInhibit == true {
		count++
	}
	if cell.hallmark.instability == true {
		count++
	}
	count += cell.hallmark.otherMuatation
	return count
}

//UpdateDivision takes currBoard and the index i, j as input. It returns the updated status of a cell division.
//Idea is the same as NextMove function. However, for division we do not need to erase the cell in original place. Every divsion, for normal cell, teloLength get shortened by 1.
//Still we need to consider the special cancer hallmark.
func UpdateDivision(currBoard CellBoard, i, j int) CellBoard {
	nextMove := 0
	newBoard := CopyBoard(currBoard)
	if currBoard[i][j].alive == true {
		//consider nonGrowthInhibit feature and growth factor limitation
		if currBoard[i][j].hallmark.nonGrowthInhibit == true || GrowthLimit(currBoard, i, j) {
			//consider proliferation feature
			if currBoard[i][j].hallmark.proliferation == true {
				//consider unDivid feature and it is related to the change of teloLength
				if currBoard[i][j].hallmark.unDivid == true || (currBoard[i][j].hallmark.unDivid == false && currBoard[i][j].parameter.teloLength > 0) {
					if GenerateProb(currBoard[i][j].parameter.division) {
						//find all the possible movement
						possMove := NextMove(currBoard, i, j, true)
						if len(possMove) == 1 {
							return newBoard
						}
						ranMove := rand.Intn(len(possMove))
						for ranMove == 0 {
							ranMove = rand.Intn(len(possMove))
						}
						nextMove = possMove[ranMove]
						if currBoard[i][j].hallmark.unDivid == false && currBoard[i][j].parameter.teloLength > 0 {
							currBoard[i][j].parameter.teloLength--
							currBoard[i][j] = UpdateMutationRate(currBoard[i][j])
						}
						switch nextMove {
						case 1:
							newBoard[i-1][j] = currBoard[i][j]
						case 2:
							newBoard[i][j+1] = currBoard[i][j]
						case 3:
							newBoard[i+1][j] = currBoard[i][j]
						case 4:
							newBoard[i][j-1] = currBoard[i][j]
						}
					}
				}
			} else if currBoard[i][j].parameter.teloLength > 0 {
				if GenerateProb(currBoard[i][j].parameter.division) {
					possMove := NextMove(currBoard, i, j, false)
					if len(possMove) == 1 {
						return newBoard
					}
					ranMove := rand.Intn(len(possMove))
					for ranMove == 0 {
						ranMove = rand.Intn(len(possMove))
					}
					nextMove = possMove[ranMove]
					//consider whether stem cell or not and whether unDivid feature
					if currBoard[i][j].types.stem == false && currBoard[i][j].hallmark.unDivid == false {
						currBoard[i][j].parameter.teloLength--
						currBoard[i][j] = UpdateMutationRate(currBoard[i][j])
					}
					switch nextMove {
					case 1:
						newBoard[i-1][j] = currBoard[i][j]
					case 2:
						newBoard[i][j+1] = currBoard[i][j]
					case 3:
						newBoard[i+1][j] = currBoard[i][j]
					case 4:
						newBoard[i][j-1] = currBoard[i][j]
					}
				}
			}
		}
	}
	return newBoard
}

//GrowthLimit takes the growth factor and nutrient limitations into consideration.
//A cell can only divide if there is a growth factor signal, but the number of cells cannot go beyong the limitations by the overall nutrients.
//Unless it is a cancer cell.
func GrowthLimit(currBoard CellBoard, i, j int) bool {
	count := CountNCells(currBoard, i, j, "N")
	if count >= 3 {
		return true
	}
	return false
}

//CountNCells counts the number of alive normal cells or stem cells in the Moore neighborhood.
//Based on whether the input string is N or S.
func CountNCells(currBoard CellBoard, i, j int, s string) int {
	count := 0

	//count normal cells
	if s == "N" {
		for m := i - 1; m <= i+1; m++ {
			for k := j - 1; k <= j+1; k++ {
				if InField(currBoard, m, k) && !(m == i && k == j) && currBoard[m][k].alive == true {
					count++
				}
			}
		}
		return count
	}

	//count stem cells
	if s == "S" {
		for m := i - 1; m <= i+1; m++ {
			for k := j - 1; k <= j+1; k++ {
				if InField(currBoard, m, k) && !(m == i && k == j) && currBoard[m][k].types.stem == true {
					count++
				}
			}
		}
	}
	return count
}

//UpdateDif takes a stem cell as input. It returns whether the cell will perform differentiation into a normal cell.
func UpdateDif(currBoard CellBoard, i, j int) Cell {
	if CountNCells(currBoard, i, j, "S") >= 5 {
		if GenerateProb(dif) {
			currBoard[i][j] = InitialDifCell()
		}
	}
	return currBoard[i][j]
}

//UpdateMutation takes currBoard and the index i, j as input. It returns the updated status of a cell mutation.
func UpdateMutation(cell Cell) Cell {
	newCell := cell
	prob := newCell.parameter.baseMutationRate
	if GenerateProb(prob) {
		newCell = GenerateMutation(newCell)
		newCell = UpdateMutationRate(newCell)
	}
	return newCell
}

//UpdateMutationRate takes cell as input. It update the muation rate according to the hallmarks and return the new cell after updating.
func UpdateMutationRate(cell Cell) Cell {
	if cell.hallmark.instability == true {
		cell.parameter.baseMutationRate = mutationRate * cell.parameter.geneticInstability * math.Log2(float64(teloLen-cell.parameter.teloLength+2))
	} else {
		cell.parameter.baseMutationRate = mutationRate * math.Log2(float64(teloLen-cell.parameter.teloLength+2))
	}
	return cell
}

//GenerateMutation takes a cell as input. It returns the random mutation in the hallmarks.
//It considers all the possible mutations, and recovering is not considered.
func GenerateMutation(cell Cell) Cell {
	probArray := make([]int, 0)
	//find the possible mutaton
	if cell.alive == true {
		if cell.hallmark.instability == false {
			probArray = append(probArray, 1)
		}
		if cell.hallmark.nonApoptosis == false {
			probArray = append(probArray, 1)
		}
		if cell.hallmark.proliferation == false {
			probArray = append(probArray, 1)
		}
		if cell.hallmark.nonGrowthInhibit == false {
			probArray = append(probArray, 1)
		}
		if cell.hallmark.unDivid == false {
			probArray = append(probArray, 1)
		}

		//generate random numbers based number of possible mutations
		probArray = append(probArray, cell.hallmark.otherMuatation)
		totalCount := SumArray(probArray)
		ranNum := rand.Intn(totalCount)
		index := SumIndex(ranNum, probArray)

		//match the random numbers to the mutations
		if index == probArray[len(probArray)-1] {
			cell.hallmark.otherMuatation++
			return cell
		}
		count := 1
		if cell.hallmark.instability == false {
			if count != index {
				count++
			} else {
				cell.hallmark.instability = true
				cell.types.stem = false
				return cell
			}
		}
		if cell.hallmark.nonApoptosis == false {
			if count != index {
				count++
			} else {
				cell.hallmark.nonApoptosis = true
				cell.types.stem = false
				return cell
			}
		}
		if cell.hallmark.proliferation == false {
			if count != index {
				count++
			} else {
				cell.hallmark.proliferation = true
				cell.types.stem = false
				return cell
			}
		}
		if cell.hallmark.nonGrowthInhibit == false {
			if count != index {
				count++
			} else {
				cell.hallmark.nonGrowthInhibit = true
				cell.types.stem = false
				return cell
			}
		}
		if cell.hallmark.unDivid == false {
			if count == index {
				cell.hallmark.unDivid = true
				cell.types.stem = false
				return cell
			}
		}
	}
	return cell
}

//SumArray takes a slice of integer as input and returns the sum of the array.
func SumArray(array []int) int {
	if len(array) == 0 {
		panic("The length of the input array should not be zero")
	}
	sum := 0
	for i := range array {
		sum += array[i]
	}
	return sum
}

//SumIndex takes an integer and a slice of integer as input. It returns the minimum index such that sum of integers in the array that smaller than or equal to the index is bigger than or equal the input integer.
func SumIndex(sum int, array []int) int {
	if len(array) == 0 {
		panic("The length of the input array should not be zero")
	}
	for i := range array {
		if sum <= SumArray(array[:i+1]) {
			return i
		}
	}
	panic("The result should have been returned")
}
