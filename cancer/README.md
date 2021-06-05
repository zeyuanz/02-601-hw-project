# 02-601 Programming Final Project
-------

# **Tumor Growth Simulation Based On Cellular Automata**
This project simulates the appearance and growth of tumor based on cellular automata, which is a locally, discrete mathmaticall model.

## Mode
Program has two types of mode: ```input mode``` and ```output mode`` and there are another two kinds of mode of each.

- ```Input mode```: ```default``` and ```user```
- ```Output mode```: ```graph``` and ```stats```

## Descirptions

- ```Default mode``` does not require user to provide a setting file. User only needs to provide the ```size of the board```, and ```number of generations```. Four stems cells will be generated to begin with at the center of the board. 
- ```User mode``` requires user to provide a setting file including the parameters for the first line and other parameters after that.

- ```Graph mode``` creates a GIF recording the board of every generation. Black represents tumor cells, purple represents stem cells and blue represents differentiated cells. 
- ```Stats mode``` creates two line graphs instead of GIF. One is the proportion of different types of cells among all the cells during all generations. The other is the proportion of different types of tumor features (cancer hallmark) among all the tumor cells during all generations.

## Run (4 different modes in combination)
1. default + graph

```
cancer.exe default graph numRows numCols numGens cellWidth outputFile
```
  * ```numRows``` is the number of rows of the board and numCols is the number of columns of the board. They should be positive integers;
  * ```numGens``` is the number of generations to simulate. It should be positive integers;
  * ```cellWidth``` is the width of the single grid shown in the final graph. Larger cellWidth means a single grid in the final graph is larger (more pixels). It should be positive integer;
  * ```ouputFile``` is the file name for output GIF graph. Outputfile will be named as outputFile.gif where "outputFile" is the name given by user.

2. default + stats

```
cancer.exe default stats numRows numCols numGens numReps outputFile interval
```
  * ```numRows, numCols, numGens``` and ```outputFile``` are same as those described above;
  * ```numReps``` is the number of repetitions you want to take. Since the model includes some randomness, the more you repeat the more precise it will be. It should be positive integers.
  * ```interval``` is the space of generation drawn on the graph. For example, if interval is 1, every generation will be drawn on the graph. If interval is 2, the first, third, fifth......generation will be drawn on the graph.
  * Format of outputFile. One graph is named as "outputFile"+numCells.png which is the proportion of different types of cells among all cells during the generations. The other is named as "outputFile"+numCHs.png which is the proportion of different types of tumor features (cancer hallmarks) among all tumor cells during the generation.

3. user + graph

```
cancer.exe user graph inputFile numGens cellWidth outputFile
```
  * ```inputFile``` is the name of the setting file given by the user. The first row of the setting should be the parameters, 8 in total, which are mutationRate, genInstability, ranApop, evaApop, teloLen, ignor, divid, dif and ranDe.
    * ```mutationRate:``` the basic mutation rate. It should be positive;
    * ```genInstability```: the multiplier for the mutation rate if genome instability is obtained by the cell. For example if it is set to be 100, when a cell obtains genome instability cancer hallmark, the mutation rate for the cell will becomes mutationRate * 100. It should be positive;
    * ```ranApop```: the basic probability for a cell to undergo apoptosis. It should be positive;
    * ```evaApop```: the extra likelihood for a cell to under apoptosis if mutations are acquired (if nonApoptosis is acquired then apoptosis will not happen). It should be set to zero, because it will be recalculated according to the specific condition of a cell;
    * ```teloLen```: the initial lenght of the telemore. It should be positive integers;
    * ```ignor```: the probability for a cell to kill the neighboring cells and perform movement and division if it has obtained proliferation feature;
    * ```divid```: the probability for a cell to divide;
    * ```dif```: the probability for a stem cell to differentiate into a differentiated cell;
    * ```ranDe```: the probability for a cell to undergo necrosis.
   * The details of the board should comes following the parameters. "B" represents blank (no cell), "S" represents stem cells and "D" represents differentiated cells. A row in the text is a row on the board. The board should be rectangle which means that the number of columns should be same in every row. Sample files are provided "default.txt".
   * ```numGens, cellWidth``` and ```outputFile``` are same as above.

4. user + stats

```
cancer.exe user stats inputFile numGens numReps outputFile interval
```
  * all parameters are same as described above.
