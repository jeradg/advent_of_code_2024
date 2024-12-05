package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

func main() {
	text, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(text)

	grid1 := buildGrid(reader)
	part1Total := totalForGrid(grid1, 1)

	_, err = reader.Seek(0, io.SeekStart)

	if err != nil {
		log.Fatal(err)
	}

	grid2 := buildGrid(reader)
	part2Total := totalForGrid(grid2, 2)

	fmt.Println("")
	fmt.Println("Total (Part 1):")
	fmt.Println(part1Total)

	fmt.Println("")
	fmt.Println("Total (Part 2):")
	fmt.Println(part2Total)
}

var directions []string = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

type Grid struct {
	FirstNode *GridNode
}

func (g *Grid) First() (*GridNode, error) {
	if g.FirstNode != nil {
		return g.FirstNode, nil
	}

	return &GridNode{}, &noNodeError{"grid does not have a first node"}
}

type GridNode struct {
	Value     string
	NorthNode *GridNode
	EastNode  *GridNode
	SouthNode *GridNode
	WestNode  *GridNode
	IsInMatch bool
}

type noNodeError struct {
	message string
}

func (e *noNodeError) Error() string {
	return e.message
}

func (gn *GridNode) NeighbourInDirection(direction string) (*GridNode, error) {
	// NOTE: reflection *seems* like a good idea, but I really hate this
	directionMeth := reflect.ValueOf(gn).MethodByName(direction)

	result := directionMeth.Call([]reflect.Value{})
	neighbour, err := result[0].Interface().(*GridNode), result[1].Interface()

	if err != nil {
		return &GridNode{}, err.(error)
	}

	return neighbour, nil
}

func (gn *GridNode) NeighbourInOppositeDirection(direction string) (*GridNode, error) {
	var oppositeDirection string
	dirIndex := -1

	for i, dir := range directions {
		if dir == direction {
			dirIndex = i
			break
		}
	}

	if dirIndex == -1 {
		return &GridNode{}, errors.New("direction not found")
	}

	numDirections := len(directions)
	oppositeIndex := ((dirIndex-(numDirections/2))%numDirections + numDirections) % numDirections
	oppositeDirection = directions[oppositeIndex]

	return gn.NeighbourInDirection(oppositeDirection)
}

func (gn *GridNode) matchesForNodePart1() int {
	if gn.Value != "X" {
		return 0
	}

	total := 0
	letters := []string{"M", "A", "S"}

	// If Value is "X",
	// this node might be the start of one or more matches
	for _, direction := range directions {
		node, err := gn.NeighbourInDirection(direction)

		if err != nil {
			continue
		}

		matchComponents := []*GridNode{gn}

		for i := 0; i < len(letters); i++ {
			testLetter := letters[i]

			if node.Value != testLetter {
				break
			} else if i == len(letters)-1 {
				matchComponents = append(matchComponents, node)

				for _, matchNode := range matchComponents {
					matchNode.IsInMatch = true
				}

				total++
			} else {
				nextNode, err := node.NeighbourInDirection(direction)

				if err != nil {
					break
				}

				matchComponents = append(matchComponents, node)
				node = nextNode
			}
		}
	}

	return total
}

func (gn *GridNode) matchesForNodePart2() int {
	if gn.Value != "A" {
		return 0
	}

	total := 0
	directions := []string{"NE", "SE", "SW", "NW"}
	matchCorners := []*GridNode{}

	// If Value is "A",
	// this node might be at the center of a match
	for _, direction := range directions {
		for _, match := range matchCorners {
			oppositeNode, err := gn.NeighbourInOppositeDirection(direction)

			if err != nil {
				continue
			}

			// Don't double-count
			if match == oppositeNode {
				continue
			}
		}

		matchComponents := []*GridNode{}

		node, err := gn.NeighbourInDirection(direction)

		if err != nil {
			continue
		}

		testLetter := "M"

		if node.Value != testLetter {
			continue
		}

		matchComponents = append(matchComponents, node)

		node, err = gn.NeighbourInOppositeDirection(direction)

		if err != nil {
			continue
		}

		testLetter = "S"

		if node.Value != testLetter {
			continue
		}

		matchComponents = append(matchComponents, node)
		matchCorners = append(matchCorners, matchComponents...)
	}

	if len(matchCorners) == 4 {
		gn.IsInMatch = true

		for _, matchNode := range matchCorners {
			matchNode.IsInMatch = true
		}

		total++
	}

	return total
}

func (gn *GridNode) N() (*GridNode, error) {
	if gn.NorthNode != nil {
		return gn.NorthNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction N"}
}
func (gn *GridNode) NE() (*GridNode, error) {
	if gn.NorthNode != nil && gn.NorthNode.EastNode != nil {
		return gn.NorthNode.EastNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction NE"}
}
func (gn *GridNode) E() (*GridNode, error) {
	if gn.EastNode != nil {
		return gn.EastNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction E"}
}
func (gn *GridNode) SE() (*GridNode, error) {
	if gn.SouthNode != nil && gn.SouthNode.EastNode != nil {
		return gn.SouthNode.EastNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction SE"}
}
func (gn *GridNode) S() (*GridNode, error) {
	if gn.SouthNode != nil {
		return gn.SouthNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction S"}
}
func (gn *GridNode) SW() (*GridNode, error) {
	if gn.SouthNode != nil && gn.SouthNode.WestNode != nil {
		return gn.SouthNode.WestNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction SW"}
}
func (gn *GridNode) W() (*GridNode, error) {
	if gn.WestNode != nil {
		return gn.WestNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction W"}
}
func (gn *GridNode) NW() (*GridNode, error) {
	if gn.NorthNode != nil && gn.NorthNode.WestNode != nil {
		return gn.NorthNode.WestNode, nil
	}

	return &GridNode{}, &noNodeError{"node does not exist in direction NW"}
}

func buildGrid(reader *bytes.Reader) *Grid {
	grid := Grid{}
	var current *GridNode
	var currentLineFirst *GridNode

	i := 0
	j := 0

	// Build the grid left to right, top to bottom
	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard input\n")
			os.Exit(1)
		}

		if isLineTerminator(r) {
			// Start a new line
			i = 0
			j++
			continue
		}

		char := string(r)

		if i == 0 {
			if j == 0 {
				newNode := GridNode{Value: char}
				current = &newNode
				grid.FirstNode = current
			} else {
				newNode := GridNode{Value: char, NorthNode: currentLineFirst}
				current = &newNode
				currentLineFirst.SouthNode = current
			}

			currentLineFirst = current
		} else {
			newNode := GridNode{Value: char, WestNode: current}

			north, err := current.NE()

			// Errors for control flow... I'm guessing this isn't really the "Go way"?
			if err == nil {
				newNode.NorthNode = north
				north.SouthNode = &newNode
			}

			current.EastNode = &newNode

			current = &newNode
		}

		i++
	}

	return &grid
}

func totalForGrid(grid *Grid, part int) int {
	total := 0

	first, err := grid.First()

	if err != nil {
		fmt.Fprintf(os.Stderr, "grid is empty")
		os.Exit(1)
	}

	current, currentLineFirst := first, first

	for {
		if part == 1 {
			total += current.matchesForNodePart1()
		} else if part == 2 {
			total += current.matchesForNodePart2()
		}

		east, err := current.E()

		if err == nil {
			// There's an east, keep going east on this line
			current = east
		} else {
			// There's no east, go to the next line
			nextLineFirst, err := currentLineFirst.S()

			if err != nil {
				break
			}

			currentLineFirst = nextLineFirst
			current = nextLineFirst
		}
	}

	return total
}

func isLineTerminator(r rune) bool {
	switch r {
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	}
	return false
}

// debugging util
func (g *Grid) Print(onlyMatches bool) {
	first, err := g.First()

	if err != nil {
		fmt.Fprintf(os.Stderr, "grid is empty")
		os.Exit(1)
	}

	current, currentLineFirst := first, first

	for {
		if onlyMatches && !current.IsInMatch {
			fmt.Print(".")
		} else {
			fmt.Print(current.Value)
		}

		east, err := current.E()

		if err == nil {
			// There's an east, keep going east on this line
			current = east
		} else {
			// There's no east, go to the next line
			fmt.Print("\n")

			nextLineFirst, err := currentLineFirst.S()

			if err != nil {
				break
			}

			currentLineFirst = nextLineFirst
			current = nextLineFirst
		}
	}

	fmt.Print("\n")
}

// debugging util
func (g *GridNode) PrintNeighbourhood() {
	dir, err := g.NW()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	dir, err = g.N()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	dir, err = g.NE()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	fmt.Print("\n")

	dir, err = g.W()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	fmt.Print(g.Value)

	dir, err = g.E()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	fmt.Print("\n")

	dir, err = g.SW()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	dir, err = g.S()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	dir, err = g.SW()
	if err == nil {
		fmt.Print(dir.Value)
	} else {
		fmt.Print("-")
	}

	fmt.Print("\n")
}
