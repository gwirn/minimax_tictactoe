package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const boardSize = 3

/*fill board with 0 at the beginning*/
func clean_board(curBoard [][]int) {
	boardSize := len(curBoard)
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			curBoard[i][j] = 0
		}
	}
}

/*sum over a column of a 2D slice
* inSlice: the slice
* idx: the index of the column/row
* axis: the axis over which to calcualte (like numpy)
 */
func sumSlice2D(inSlice [][]int, idx int, axis int) int {
	sumVal := 0
	if axis == 0 {
		for _, i := range inSlice {
			sumVal += i[idx]
		}
	} else if axis == 1 {
		for _, i := range inSlice[idx] {
			sumVal += i
		}
	} else {
		panic(fmt.Sprintf("Invalid axis argument %d for slice with 2 dimensions", axis))
	}
	return sumVal
}

/* check the state of the game
* -1 X wins
* 1 O wins
* 0 tie
* 2 game is not over
 */
func checkGameState(curBoard [][]int) int {
	// check rows and columns for boardSize of the same symbol
	for ax := 0; ax < 2; ax++ {
		for ind := 0; ind < boardSize; ind++ {
			if lineVal := sumSlice2D(curBoard, ind, ax); lineVal == boardSize {
				return 1
			} else if lineVal == -boardSize {
				return -1
			}
		}
	}
	// check the diagonals for boardSize of the same symbol
	boardSize := len(curBoard)
	diag := 0
	revDiag := 0
	for i := 0; i < boardSize; i++ {
		diag += curBoard[i][i]
		revDiag += curBoard[i][boardSize-(i+1)]
	}
	if diag == boardSize || revDiag == boardSize {
		return 1
	} else if diag == -boardSize || revDiag == -boardSize {
		return -1
	}
	// see how many fields are left to play
	freeFields := 9
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			if curBoard[i][j] != 0 {
				freeFields--
			}
		}
	}
	if freeFields == 0 {
		return 0
	}
	return 2
}

/*print current state of the board*/
func printBoard(curBoard [][]int, tokenOne string, tokenNegOne string) {
	boardSize := len(curBoard)
	fmt.Printf("%s%s\n", strings.Repeat("----", boardSize), "-")
	for i := 0; i < boardSize; i++ {
		out := "| "
		for j := 0; j < boardSize; j++ {
			token := "0"
			if curBoard[i][j] == 1 {
				token = tokenOne
			} else if curBoard[i][j] == -1 {
				token = tokenNegOne
			} else {
				token = "0"
			}
			out = out + token + " | "
		}
		fmt.Println(out)
		fmt.Printf("%s%s\n", strings.Repeat("----", boardSize), "-")
	}
}

/*find the best value of the best move with the minimax algorithm*/
func minimax(curBoard [][]int, depth int, isMax bool) int {
	boardState := checkGameState(curBoard)
	// board states when win or tie
	switch {
	case boardState == 1:
		return 10
	case boardState == -1:
		return -10
	case boardState == 0:
		return 0
	}
	// find value of the best move
	if isMax {
		best := int(math.Inf(-1))
		for i := 0; i < boardSize; i++ {
			for j := 0; j < boardSize; j++ {
				if curBoard[i][j] == 0 {
					curBoard[i][j] = 1
					if mscore := minimax(curBoard, depth+1, !isMax) - depth; mscore > best {
						best = mscore
					}
					curBoard[i][j] = 0
				}
			}
		}
		return best
	} else {
		best := int(math.Inf(1))
		for i := 0; i < boardSize; i++ {
			for j := 0; j < boardSize; j++ {
				if curBoard[i][j] == 0 {
					curBoard[i][j] = -1
					if mscore := minimax(curBoard, depth+1, !isMax) + depth; mscore < best {
						best = mscore
					}
					curBoard[i][j] = 0
				}
			}
		}
		return best
	}
}

func bigger(x int, y int) bool {
	return x > y
}

func smaller(x int, y int) bool {
	return y > x
}

/*find the best move coordinates for all possible positions*/
func findMove(curBoard [][]int, symbol int) (int, int) {
	moveScore := int(math.Inf(symbol * -1))
	moveCoordX, moveCoordY := -1, -1
	symBool := symbol == -1
	compFunc := smaller
	if symbol == -1 {
		compFunc = bigger
	}
	// of all possible moves search them and chose the best
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			if curBoard[i][j] == 0 {
				curBoard[i][j] = symbol
				curMoveScore := minimax(curBoard, 0, symBool)
				curBoard[i][j] = 0
				if compFunc(moveScore, curMoveScore) {
					moveScore = curMoveScore
					moveCoordX = i
					moveCoordY = j
				}
			}
		}
	}
	return moveCoordX, moveCoordY
}

/*scan user input and convert it to board coordinates*/
func userInputScan() (int, int) {
	var userInput string
	fmt.Print("Your input: ")
	fmt.Scanln(&userInput)
	inAsNum, err  := strconv.Atoi(userInput)
	if err != nil {
		fmt.Println("Game terminated")
		os.Exit(0)
	}

	chosenY := (inAsNum - 1) % boardSize
	chosenX := int(math.Floor(float64(inAsNum - 1) / boardSize))
	return chosenX, chosenY
}

func main() {
	fmt.Println("")
	fmt.Println(strings.Repeat("+", 65))
	fmt.Println("+ Board coordinates to set your symbol", strings.Repeat(" ", 24), "+")
	fmt.Println("+ Enter the coordinate and press enter/return to make your move +")
	fmt.Println("+ To end the game input any other character apart from 1-9", strings.Repeat(" ", 4), "+")
	fmt.Println(strings.Repeat("+", 65))

	fmt.Println("")
	fmt.Println("Board coordinates:")
	for i := 1; i <= boardSize*boardSize; i++ {
		fmt.Printf("%-2d ", i)
		if i % boardSize == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")

	// get user symbol for the board display
	var userChar string
	for {
		fmt.Print("Enter your symbol (one character): ")
		fmt.Scanln(&userChar)
		if len(userChar) == 1 {
			break
		}
	}

	board := [][]int{}
	for i := 0; i < boardSize; i++ {
		row := []int{}
		for j := 0; j < boardSize; j++ {
			row = append(row, 0)
		}
		board = append(board, row)
	}

	boardState := 2
	// user board coordinates
	chosenX, chosenY := -1, -1
	starter := 1
	rountCount := 1
	for {
		fmt.Printf("*** ROUND %d ***\n", rountCount)
		if starter == 1 {
			printBoard(board, "*", userChar)
			for {
				// ask until valid input
				for {
					userX, userY := userInputScan()
					if board[userX][userY] == 0 {
						chosenX, chosenY = userX, userY
						break
					}
				}
				// play and check state after each move
				board[chosenX][chosenY] = -1
				printBoard(board, "*", userChar)
				boardState = checkGameState(board)
				if boardState != 2 {
					break
				}

				xMove, yMove := findMove(board, 1)
				board[xMove][yMove] = 1
				printBoard(board, "*", userChar)
				boardState = checkGameState(board)
				if boardState != 2 {
					break
				}
			}
		} else {
			for {
				xMove, yMove := findMove(board, 1)
				board[xMove][yMove] = 1
				printBoard(board, "*", userChar)
				boardState = checkGameState(board)
				if boardState != 2 {
					break
				}

				// ask until valid input
				for {
					userX, userY := userInputScan()
					if board[userX][userY] == 0 {
						chosenX, chosenY = userX, userY
						break
					}
				}
				// play and check state after each move
				board[chosenX][chosenY] = -1
				printBoard(board, "*", userChar)
				boardState = checkGameState(board)
				if boardState != 2 {
					break
				}
			}
		}
		rountCount++
		if rountCount%2 == 0 {
			starter = -1
		} else {
			starter = 1
		}
		switch {
		case boardState == 1:
			fmt.Println("-_-_- You lose -_-_-")
		case boardState == -1:
		fmt.Println("-_-_- You win -_-_-")
		case boardState == 0:
		fmt.Println("-_-_- Tie -_-_-")

		}
		// reset board for the next round
		clean_board(board)
	}
}
