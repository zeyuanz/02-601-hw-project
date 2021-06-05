// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)

Our version of this program reads text from standard input, parsing it into a
Markov chain, and writes generated text to standard output.
The prefix and output lengths can be specified using the -prefix and -words
flags on the command-line.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
	"sort"
	"strconv"
	//"reflect"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	for i := range p {
		if len(p[i])==0 {
			p[i] = "\"\""
		}
	}
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string]Frequency
	prefixLen int
}
//Frequency is a type that contains a map with the string as key and the count as int vlaue
type Frequency map[string]int

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	newChain := new(Chain)
	(*newChain).chain = make(map[string]Frequency)
	(*newChain).prefixLen = prefixLen
	return newChain
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		_,exist := c.chain[key]
		if !exist {
			var tmpMap Frequency
			tmpMap = make(map[string]int)
			c.chain[key] = tmpMap
		}
		tmp := c.chain[key]
		tmp[s]++
		c.chain[key] = tmp
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	var next string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		//probSlice := make([]int, 0)
		tmpStr := make([]string, 0)
		for k := range choices {
			tmpStr = append(tmpStr, k)
			//probSlice = append(probSlice, v)
		}
		//total := Sum(probSlice)
		randNum := rand.Intn(len(tmpStr))
		/*for i := range probSlice {
			if randNum < Sum(probSlice[0:i+1]) {
				 next = tmpStr[i]
				 break
			}
		}*/
		next = tmpStr[randNum]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

//Sum takes a slice of integer as input and returns the sum of the slice.
func Sum(nums []int) int {
	sum := 0
	for i := range nums {
		sum += nums[i]
	}
	return sum
}

//PrintFrequencyTable takes prefixlen, chain and name of output file as input. It print the frequency table to the outputfile.
func PrintFrequencyTable(prefixLen int, c *Chain, fileName string) {
	outFile, err := os.Create(fileName)
	if err != nil {
		panic("something wrong in creating the file")
	}
	fmt.Fprintln(outFile,prefixLen)
	keyStr := make([]string, 0)
	for k := range c.chain {
		keyStr = append(keyStr, k)
	}
	sort.Strings(keyStr)
	for i := range keyStr {
		k := keyStr[i]
		v := c.chain[keyStr[i]]
		fmt.Fprint(outFile,k," ")
		for k1,v1 := range v {
			fmt.Fprint(outFile, k1," ", v1," ")
		}
		fmt.Fprintln(outFile)
	}
	outFile.Close()
}

//BuildFromFile takes a slice of files' name and Chain as input and produce the frequency table from those files
func (c *Chain) BuildFromFile (filesName []string) {
	for i := range filesName {
		file,err := os.Open(filesName[i])
		if err != nil {
			panic("something wrong in opening the file")
		}
		c.Build(file)
	}
}

func ReadFrequencyTable(fileName string) *Chain {
	var prefixLen int
	c := NewChain(prefixLen)
	file, err := os.Open(fileName)
	if err != nil {
		panic("something wrong in opening frequency table file")
	}
	defer file.Close() // close the file after the end of the function
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //until the end of the file
		currentLine := scanner.Text() //read by line
		if len(currentLine) == 1 {
			prefixLen,err = strconv.Atoi(currentLine)
			if err != nil {
				panic("something wrong in converting string to integer")
			}
			c.prefixLen = prefixLen
		} else {
			tmpLine := strings.Split(currentLine, " ")
			keys:= strings.Join(tmpLine[:prefixLen], " ")
			tmp := make(Frequency)
			for i := prefixLen; i < len(tmpLine)-1; i +=2 {
			 tmp[tmpLine[i]],_ = strconv.Atoi(tmpLine[i+1])
			}
			c.chain[keys]=tmp
		}
	}
	return c
}
//WriteStringsToFile takes a collection of strings and a filename and
//writes these strings to the given file, with each string on one line.
func WriteStringsToFile(patterns []string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("cannot create a file")
	}
	defer file.Close()

	for _, pattern := range patterns {
		fmt.Fprintln(file, pattern)
	}
}

func main() {
	// Register command-line flags.
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
	flag.Parse()                     // Parse command-line flags.
	augs := flag.Args()				 // Rest of the arguments after parsing
	if augs[0] == "read" {
		fileName := augs[2]
		prefixLen,err1 := strconv.Atoi(augs[1]) //defult value is 2
		if err1 != nil {
			panic("something wrong in converting string to integer")
		}
		filesName := augs[3:]
		c := NewChain(prefixLen)     // Initialize a new Chain.
		c.BuildFromFile(filesName)
		PrintFrequencyTable(prefixLen, c, fileName)
	}

	if augs[0] == "generate" {
		inFile:=augs[1]
		numWords,err1 := strconv.Atoi(augs[2])
		if err1 != nil {
			panic("something wrong in converting string to integer")
		}
		c := ReadFrequencyTable(inFile)
		text := c.Generate(numWords) // Generate text.
		fmt.Println(text)             // Write text to standard output.
	}
}
