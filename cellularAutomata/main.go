package main

import (
  "fmt"
  "os"
  "strconv"
)

type GameBoard [][]int

//PlayAutomaton takes an initial game board and a number of generations and returns the boards after playing automaton specified by a neighborhood and a collection of rule strings numGens+1 times.
func PlayAutomaton(initialBoard GameBoard, numGens int, neighborhood string, ruleStrings []string) []GameBoard {
	boards := make([]GameBoard, numGens+1)

	boards[0] = initialBoard

	for gen := 1; gen < numGens+1; gen++ {
		boards[gen] = UpdateBoard(boards[gen-1], neighborhood, ruleStrings)
	}

	return boards
}

func AssertBoardRectangular(board GameBoard) {
  if len(board) == 0 {
    fmt.Println("Game board has no rows.")
    os.Exit(2)
  }
  // range over rows, and ensure that they all have same length
  numCols := len(board[0])
  for i := 1; i < len(board); i++ {
    if len(board[i]) != numCols {
      fmt.Println("Board isn't rectangular. Row", i, "has length", len(board[i]), "and row 0 has length", len(board[0]))
      os.Exit(1)
    }
  }
}

//UpdateBoard takes a current game board and produces the board in the next generation according to whether the neighborhood is vonNeumann or Moore and a given set of rule strings.
//Assume input board is rectangular.
func UpdateBoard(currBoard GameBoard, neighborhood string, ruleStrings []string) GameBoard {
  AssertBoardRectangular(currBoard) // ensure board is rectangular

	numRows := len(currBoard)
	numCols := len(currBoard[0])
	newBoard := InitializeBoard(numRows, numCols)

	//everything in newBoard is 0

	// now we just need to range over cells of current board and fill values of new board

	for r := range currBoard {
		for c := range currBoard[r] {
      // I'm only going to update non-boundary cells.  We can make more sophisticated later.
      if r > 0 && r < numRows - 1 && c > 0 && c < numCols - 1 {
        newBoard[r][c] = UpdateCell(currBoard, r, c, neighborhood, ruleStrings)
      }
		}
	}

	return newBoard
}

func InitializeBoard(numRows, numCols int) GameBoard {
	board := make(GameBoard, numRows)
	for r := range board {
		// range over rows and make each row of length numCols
		board[r] = make([]int, numCols)
	}
	return board
}

/*
UpdateCell(board, r, c, nbrhood, rulestrings)
	for every integer i from 0 to #rulestrings – 1
		rule  rulestrings[i]
		if CellMatch(board, r, c, nbrhood, rule)
			return rule[length(rule) – 1]

*/

//UpdateCell takes a GameBoard, row/col indices, neighborhood, and rule strings. It finds the updated state of the cell at these indices in the next generation according to appropriate rule, and returns that state.
func UpdateCell(currBoard GameBoard, r, c int, neighborhood string, ruleStrings []string) int {
	for _, rule := range ruleStrings {
    if RuleMatch(currBoard, r, c, neighborhood, rule) {
      // updated state of current cell is final element of rule string
      newState, err := strconv.Atoi(rule[len(rule)-1:])
      if err != nil {
        fmt.Println("Error in accessing element at end of rule string")
        os.Exit(3)
      }
      return newState
    }
  }
  // we made it through entire for loop without a match :(
  fmt.Println("Error: No matching rule found.")
  os.Exit(4)
  return -1
}

//RuleMatch takes a current board along with row/col indices, a neighborhood type, and a single rule string as input. If the given neighborhood type of current cell matches first n-1 elements of rule string, returns true.  Returns false otherwise.
func RuleMatch(currBoard GameBoard, r, c int, neighborhood string, rule string) bool {
  if neighborhood == "vonNeumann" {
    return RuleMatchVN(currBoard, r, c, rule)
  } else if neighborhood == "Moore" {
    return RuleMatchMoore(currBoard, r, c, rule)
  }
  fmt.Println("Invalid neighborhood type given.")
  os.Exit(5)
  //can't ever hit this
  return false
}

//RuleMatchVN takes a game board, r and c indices, and a rule string. It returns true if the first n-1 elements of rule string match the VonNeumann neighborhood of board[r][c].
func RuleMatchVN(currBoard GameBoard, r, c int, rule string) bool {
  // later: we are doing this on a donut (torus)

  // first: let's convert rule string to an appropriate type (slice of ints)
  ruleSlice := ConvertRuleString(rule)
  //fmt.Println(rule)
  if currBoard[r][c] != ruleSlice[0] {
    return false
  }
  if currBoard[r-1][c] != ruleSlice[1] {
    return false
  }
  if currBoard[r][c+1] != ruleSlice[2] {
    return false
  }
  if currBoard[r+1][c] != ruleSlice[3] {
    return false
  }
  if currBoard[r][c-1] != ruleSlice[4] {
    return false
  }
  return true
}

//ConvertRuleString takes a rule string and converts it to a slice of integers
//of the same length.
func ConvertRuleString(rule string) []int {
  ruleSlice := make([]int, len(rule))

  // range through rule, convert current symbol to slice
  for i := range rule {
    //let's convert rule[i] to an integer
    num, err := strconv.Atoi(rule[i:i+1])
    if err != nil {
      panic("Issue with converting rule string to integers.")
    }
    ruleSlice[i] = num
  }

  return ruleSlice
}

func SliceMatches(a1, a2 []int) bool {
  // assumes a1 and a2 have same length
  if len(a1) != len(a2) {
    panic("Error: Two slices have unequal length.")
  }
  for i := range a1 {
    if a1[i] != a2[i] {
      return false
    }
  }
  return true
}

//RuleMatchMoore takes a game board, r and c indices, and a rule string. It returns true if the first n-1 elements of rule string match the Moore neighborhood of board[r][c].
func RuleMatchMoore(currBoard GameBoard, r, c int, rule string) bool {
  ruleSlice := ConvertRuleString(rule)

  if currBoard[r][c] != ruleSlice[0] { // central element
    return false
  }
  topRow := currBoard[r-1][c-1:c+2]
  bottomRow := currBoard[r+1][c-1:c+2]
  // check if top bottom row match ... if either doesn't, return false
  if SliceMatches(topRow, ruleSlice[1:4]) == false || SliceMatches(bottomRow, ruleSlice[6:9]) == false {
    return false
  }

  // now check squares to left and right of central element
  if currBoard[r][c-1] != ruleSlice[4] {
    return false
  }
  if currBoard[r][c+1] != ruleSlice[5] {
    return false
  }

  //default: return true
  return true
}

func AddVNRotations(ruleStrings []string) []string {
  // range over rule strings, append each of three rotations
  for _, rule := range ruleStrings {
    for _, rule2 := range ThreeRotations(rule) {
      ruleStrings = append(ruleStrings, rule2)
    }
  }

  return ruleStrings
}

func ThreeRotations(rule string) []string {
  ruleStrings := make([]string, 3)
  for i := 0; i < 3; i++ {
    rule = SingleRotation(rule)
    ruleStrings[i] = rule
  }

  return ruleStrings
}

func SingleRotation(rule string) string {
  rule2 := rule[:1]
  rule2 += rule[2:5]
  rule2 += rule[1:2]
  rule2 += rule[5:]
  return rule2
}

func main() {
  fmt.Println("Self-replicating systems ?")

  //I want to be able to take neighborhood type, filename containing rule strings, filename containing initial board, and output file name for animated GIF as input from the user.
  // We also should take cell width, and number of generations to play the game.

  // command line arguments go into an array of strings called os.Args[].  The arrays has length n+1, where n = # input parameters.
  // os.Args[0] is the file name of the program.
  // os.Args[1] is the first parameter given, etc.

  neighborhood := os.Args[1] // either "vonNeumann" or "Moore"

  ruleFile := os.Args[2] // tells us where to look for rules

  initialBoardFile := os.Args[3] // tells us where to look for starting board

  outputFile := os.Args[4] // what we want to call the GIF

  cellWidth, err1 := strconv.Atoi(os.Args[5]) // gives width of cell for drawing board
  if err1 != nil {
    panic("Issue in converting cellwidth command line parameter")
  }

  numGens, err2 := strconv.Atoi(os.Args[6]) // gives number of generations for game
  if err2 != nil {
    panic("Issue in converting cellwidth command line parameter")
  }

  fmt.Println("All command line arguments read successfully.")

  // read in a board
  initialBoard := ReadBoardFromFile(initialBoardFile)
  fmt.Println("We have read board from file.")

  // read in a collection of rules
  ruleStrings := ReadRulesFromFile(ruleFile)
  fmt.Println("We have read rules from file.")

  //let's add all rotations if neighborhood is vonNeumann
  if neighborhood == "vonNeumann" {
    ruleStrings = AddVNRotations(ruleStrings)
  }

  // play the automaton numGenerations generations
  boards := PlayAutomaton(initialBoard, numGens, neighborhood, ruleStrings)

  fmt.Println("Automaton played successfully! Now we draw the boards.")

  // produce animated GIF corresponding to automaton

  imglist := DrawGameBoards(boards, cellWidth)
  fmt.Println("Image lists created. Now Transfer to GIF.")

  ImagesToGIF(imglist, outputFile)

}
