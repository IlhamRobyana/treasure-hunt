package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

//the coords used for the map coordinates
type coords struct {
	x int
	y int
}

// check if there's an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// load map from the text file
func loadMap() (treasureMap string) {
	content, err := ioutil.ReadFile("map.txt")
	check(err)
	treasureMap = string(content)
	return
}

// get the possible coords of the treasure
func getPossibleCoords() (possibleCoords []coords) {
	possibleCoords = []coords{
		{
			x: 5,
			y: 2,
		},
		{
			x: 6,
			y: 2,
		},
		{
			x: 5,
			y: 3,
		},
		{
			x: 5,
			y: 4,
		},
		{
			x: 3,
			y: 4,
		},
	}
	return
}

// place the possible coords on the map
func placePossibleCoords(treasureMap *string, possibleCoords []coords) {
	temp := []byte(*treasureMap)
	for i := 0; i < len(possibleCoords); i++ {
		x := possibleCoords[i].x
		y := possibleCoords[i].y
		temp[x+10*y] = '$'
	}
	*treasureMap = string(temp)
}

// make the map into 2d array
func generate2DMap(treasureMap string) [][]byte {
	var map2D [][]byte
	for i := 0; i < 6; i++ {
		temp := treasureMap[i*10 : i*10+8]
		map2D = append(map2D, []byte(temp))
	}
	return map2D
}

// prints the map on console
func printMap(map2d [][]byte) {
	fmt.Println(" 01234567")
	for i := 0; i < 6; i++ {
		fmt.Print(i)
		fmt.Println(string(map2d[i]))
	}
}

// prints the instructions and map at the beginning
func printInit(map2D [][]byte, possibleCoords []coords) {
	fmt.Println("Find a treasure within this map!")
	printMap(map2D)
	fmt.Println("\n# are walls")
	fmt.Println(". are walkable spaces")
	fmt.Println("X is you")
	fmt.Println("$ are possible treasure spots\n")
	fmt.Println("Possible treasure coords are following:")
	fmt.Println(possibleCoords)
	fmt.Println("You are currently at {1,4}")
	fmt.Println("\nYou need to move north, and then east, and then south in order")
}

// get the randomized treasure spot
func getTreasureSpot(possibleCoords []coords) coords {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	element := r1.Intn(len(possibleCoords))
	return possibleCoords[element]
}

// move the player according to cycle and steps
func move(cycle, steps int, currentSpot *coords, map2D [][]byte) ([][]byte, bool) {
	prevX := currentSpot.x
	prevY := currentSpot.y
	if steps < 1 {
		fmt.Println("steps must be greater than zero")
		return map2D, false
	}
	var path []byte
	switch cycle {
	case 0:
		for i := 1; i <= steps; i++ {
			path = append(path, map2D[prevY-i][prevX])
			currentSpot.y--
		}
	case 1:
		for i := 1; i <= steps; i++ {
			path = append(path, map2D[prevY][prevX+i])
			currentSpot.x++
		}
	case 2:
		for i := 1; i <= steps; i++ {
			path = append(path, map2D[prevY+i][prevX])
			currentSpot.y++
		}
	}
	if strings.Contains(string(path), "#") {
		if path[0] == '#' {
			fmt.Println("Can't go anywhere else in that direction, Game Over.")
			return nil, false
		} else {
			currentSpot.x = prevX
			currentSpot.y = prevY
			fmt.Println("There's a wall in the way with number of steps taken, try shorther steps")
			return map2D, false
		}
	}
	map2D[prevY][prevX] = '.'
	map2D[currentSpot.y][currentSpot.x] = 'X'
	return map2D, true
}

// the game loop
func hunt(map2D [][]byte, possibleCoords []coords, treasureSpot coords) bool {
	currentSpot := coords{1, 4}
	directions := [3]string{"north", "east", "south"}
	for i := 0; i < 3; i++ {
		for {
			printMap(map2D)
			fmt.Println("Type the number of steps to the " + directions[i] + " direction")
			var steps int
			fmt.Scanln(&steps)
			map2D, valid := move(i, steps, &currentSpot, map2D)
			if valid {
				break
			} else if map2D == nil {
				return false
			}
		}
	}
	printMap(map2D)
	if currentSpot == treasureSpot {
		return true
	} else {
		return false
	}
}

// prints the result of the game
func printResult(result bool, treasureSpot coords) {
	if result {
		fmt.Println("Congratulations! You found the treasure!")
	} else {
		fmt.Println("Too bad! The treasure is in another spot!")
		fmt.Print("It's at ")
		fmt.Println(treasureSpot)
	}
}

func main() {
	treasureMap := loadMap()
	possibleCoords := getPossibleCoords()
	placePossibleCoords(&treasureMap, possibleCoords)
	map2D := generate2DMap(treasureMap)
	printInit(map2D, possibleCoords)
	treasureSpot := getTreasureSpot(possibleCoords)
	result := hunt(map2D, possibleCoords, treasureSpot)
	printResult(result, treasureSpot)
}
