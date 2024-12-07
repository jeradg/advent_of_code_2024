package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(text)

	grid := buildGrid(reader)
	part1Total := totalForGridPart1(grid)

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("Part 1 grid:")
	grid.Print(false)
	fmt.Println("")

	fmt.Println("")
	fmt.Println("Total (Part 1):")
	fmt.Println(part1Total)

	part2Total := totalForGridPart2(grid)

	fmt.Println("")
	fmt.Println("Part 2 grid:")
	grid.Reset()
	grid.Print(true)
	fmt.Println("")

	fmt.Println("")
	fmt.Println("Total (Part 2):")
	fmt.Println(part2Total)
}

var directions []string = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

type Grid struct {
	FirstNode           *GridNode
	GuardNode           *GridNode
	GuardStartingNode   *GridNode
	ObstacleAttemptNode *GridNode
	VisitedNodes        []*GridNode
	Done                bool
}

func (g *Grid) Reset() {
	for _, gn := range g.VisitedNodes {
		// Reset everything except for "IsObstacleCandidate"
		gn.Value = "."
		gn.VisitCount = 0
		gn.HasGuardFacedNorth = false
		gn.HasGuardFacedEast = false
		gn.HasGuardFacedSouth = false
		gn.HasGuardFacedWest = false
		gn.IsGuardInLoop = false
	}

	g.GuardNode = g.GuardStartingNode
	g.GuardNode.Value = "^"
	g.GuardNode.VisitCount = 1
	g.GuardNode.HasGuardFacedNorth = true
	g.GuardNode.HasGuardFacedEast = false
	g.GuardNode.HasGuardFacedSouth = false
	g.GuardNode.HasGuardFacedWest = false
	g.GuardNode.IsGuardInLoop = false
}

func (g *Grid) First() (*GridNode, bool) {
	if g.FirstNode != nil {
		return g.FirstNode, true
	}

	return &GridNode{}, false
}

func (g *Grid) Guard() (*GridNode, bool) {
	if g.GuardNode != nil {
		return g.GuardNode, true
	}

	return &GridNode{}, false
}

type GridNode struct {
	NorthNode           *GridNode
	EastNode            *GridNode
	SouthNode           *GridNode
	WestNode            *GridNode
	Value               string
	VisitCount          int
	HasGuardFacedNorth  bool
	HasGuardFacedEast   bool
	HasGuardFacedSouth  bool
	HasGuardFacedWest   bool
	IsGuardInLoop       bool
	IsObstacleCandidate bool
}

func (gn *GridNode) WalkGuard() (*GridNode, bool) {
	if !gn.IsGuardNode() {
		fmt.Fprintf(os.Stderr, "WalkGuard() called on non-guard node")
		os.Exit(1)
	}

	var newNode *GridNode
	var guardNode *GridNode
	var isInGrid bool

	switch gn.Value {
	case "^":
		newNode, isInGrid = gn.N()
	case ">":
		newNode, isInGrid = gn.E()
	case "v":
		newNode, isInGrid = gn.S()
	case "<":
		newNode, isInGrid = gn.W()
	}

	if !isInGrid {
		// Guard exited the grid!
		gn.Value = "X"
		return newNode, false
	}

	if newNode.IsObstacle() {
		// Rotate the guard
		switch gn.Value {
		case "^":
			gn.Value = ">"
		case ">":
			gn.Value = "v"
		case "v":
			gn.Value = "<"
		case "<":
			gn.Value = "^"
		}

		guardNode = gn
	} else {
		newNode.Value = gn.Value
		gn.Value = "X"
		newNode.VisitCount++
		guardNode = newNode
	}

	switch guardNode.Value {
	case "^":
		if guardNode.HasGuardFacedNorth {
			guardNode.IsGuardInLoop = true
		} else {
			guardNode.HasGuardFacedNorth = true
		}
	case ">":
		if guardNode.HasGuardFacedEast {
			guardNode.IsGuardInLoop = true
		} else {
			guardNode.HasGuardFacedEast = true
		}
	case "v":
		if guardNode.HasGuardFacedSouth {
			guardNode.IsGuardInLoop = true
		} else {
			guardNode.HasGuardFacedSouth = true
		}
	case "<":
		if guardNode.HasGuardFacedWest {
			guardNode.IsGuardInLoop = true
		} else {
			guardNode.HasGuardFacedWest = true
		}
	}

	return guardNode, true
}

func (gn *GridNode) N() (*GridNode, bool) {
	if gn.NorthNode != nil {
		return gn.NorthNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) NE() (*GridNode, bool) {
	if gn.NorthNode != nil && gn.NorthNode.EastNode != nil {
		return gn.NorthNode.EastNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) E() (*GridNode, bool) {
	if gn.EastNode != nil {
		return gn.EastNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) SE() (*GridNode, bool) {
	if gn.SouthNode != nil && gn.SouthNode.EastNode != nil {
		return gn.SouthNode.EastNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) S() (*GridNode, bool) {
	if gn.SouthNode != nil {
		return gn.SouthNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) SW() (*GridNode, bool) {
	if gn.SouthNode != nil && gn.SouthNode.WestNode != nil {
		return gn.SouthNode.WestNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) W() (*GridNode, bool) {
	if gn.WestNode != nil {
		return gn.WestNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) NW() (*GridNode, bool) {
	if gn.NorthNode != nil && gn.NorthNode.WestNode != nil {
		return gn.NorthNode.WestNode, true
	}

	return &GridNode{}, false
}

func (gn *GridNode) IsGuardNode() bool {
	switch gn.Value {
	case "^", ">", "v", "<":
		return true
	default:
		return false
	}
}

func (gn *GridNode) IsObstacle() bool {
	return gn.Value == "#" || gn.Value == "0"
}

func buildGrid(reader *bytes.Reader) *Grid {
	grid := Grid{Done: false}
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
		var newNode GridNode

		if i == 0 {
			if j == 0 {
				newNode = GridNode{Value: char}
				current = &newNode
				grid.FirstNode = current
			} else {
				newNode = GridNode{Value: char, NorthNode: currentLineFirst}
				current = &newNode
				currentLineFirst.SouthNode = current
			}

			currentLineFirst = current
		} else {
			newNode = GridNode{Value: char, WestNode: current}

			north, ok := current.NE()

			// Errors for control flow... I'm guessing this isn't really the "Go way"?
			if ok {
				newNode.NorthNode = north
				north.SouthNode = &newNode
			}

			current.EastNode = &newNode

			current = &newNode
		}

		if newNode.IsGuardNode() {
			newNode.VisitCount = 1
			// NOTE: ...assuming the guard can only start facing north!!!
			newNode.HasGuardFacedNorth = true
			grid.GuardStartingNode = &newNode
			grid.GuardNode = &newNode
		}

		i++
	}

	return &grid
}

func totalForGridPart1(grid *Grid) int {
	// The current guard node has already been visited,
	// so we start at 1
	total := 1

	for {
		guardNode, ok := grid.Guard()

		if !ok {
			fmt.Fprintf(os.Stderr, "grid does not seem to have a guard node")
			os.Exit(1)
		}

		if !guardNode.IsGuardNode() {
			fmt.Fprintf(os.Stderr, "grid.GuardNode is not a valid guard node")
			os.Exit(1)
		}

		newGuardNode, ok := guardNode.WalkGuard()

		// Guard changed nodes (or left the grid)
		if newGuardNode != guardNode {
			grid.GuardNode = newGuardNode

			if newGuardNode.VisitCount == 1 {
				grid.VisitedNodes = append(grid.VisitedNodes, newGuardNode)
				total++
			}
		}

		if !ok {
			// Guard left the grid
			grid.Done = true
			break
		}
	}

	return total
}

func totalForGridPart2(grid *Grid) int {
	total := 0

	nodesToAttempt := make([]*GridNode, len(grid.VisitedNodes))
	copy(nodesToAttempt, grid.VisitedNodes)

	for _, gn := range nodesToAttempt {
		grid.Reset()
		gn.Value = "0"

		for {
			guardNode, ok := grid.Guard()

			if !ok {
				fmt.Fprintf(os.Stderr, "grid does not seem to have a guard node")
				os.Exit(1)
			}

			if !ok {
				fmt.Fprintf(os.Stderr, "grid does not seem to have a guard node")
				os.Exit(1)
			}

			if !guardNode.IsGuardNode() {
				fmt.Fprintf(os.Stderr, "grid.GuardNode is not a valid guard node")
				os.Exit(1)
			}

			newGuardNode, ok := guardNode.WalkGuard()

			// Guard changed nodes (or left the grid)
			if newGuardNode != guardNode {
				grid.GuardNode = newGuardNode

				if newGuardNode.VisitCount == 1 {
					grid.VisitedNodes = append(grid.VisitedNodes, newGuardNode)
				}
			}

			if newGuardNode.IsGuardInLoop {
				gn.IsObstacleCandidate = true

				total++
				break
			}

			if !ok {
				// Guard left the grid
				grid.Done = true
				break
			}
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
func (g *Grid) Print(showObstacleCandidates bool) {
	first, ok := g.First()

	if !ok {
		fmt.Fprintf(os.Stderr, "grid is empty")
		os.Exit(1)
	}

	current, currentLineFirst := first, first

	for {
		if showObstacleCandidates && current.IsObstacleCandidate {
			fmt.Print("0")
		} else {
			fmt.Print(current.Value)
		}

		east, ok := current.E()

		if ok {
			// There's an east, keep going east on this line
			current = east
		} else {
			// There's no east, go to the next line
			fmt.Print("\n")

			nextLineFirst, ok := currentLineFirst.S()

			if !ok {
				break
			}

			currentLineFirst = nextLineFirst
			current = nextLineFirst
		}
	}

	fmt.Print("\n")
}
