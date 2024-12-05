package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	part1SafeReportsCount := 0
	part2SafeReportsCount := 0

	for scanner.Scan() {
		levels := list.New()
		line := scanner.Text()

		fmt.Println("")

		for _, level := range strings.Fields(line) {
			num, err := strconv.Atoi(level)

			if err != nil {
				fmt.Fprintln(os.Stderr, "error reading current num:", err)
			}

			levels.PushBack(num)
		}

		if isSafePart1(levels) {
			fmt.Println("Safe! for part 1")
			part1SafeReportsCount++
		} else {
			fmt.Println("Unsafe :( for part 1")
		}

		if isSafePart2(levels) {
			fmt.Println("Safe! for part 2")
			part2SafeReportsCount++
		} else {
			fmt.Println("Unsafe :( for part 2")
		}
	}

	fmt.Println("")
	fmt.Println("Number of safe reports (Part 1):")
	fmt.Println(part1SafeReportsCount)

	fmt.Println("")
	fmt.Println("Number of safe reports (Part 2):")
	fmt.Println(part2SafeReportsCount)
}

func isSafePart1(levels *list.List) bool {
	safe := true
	direction := "none"
	var prevDirection string

	for level := levels.Front(); level != nil; level = level.Next() {
		fmt.Print(level.Value)
		fmt.Print(" ")
	}
	fmt.Print("\n")

	for level := levels.Front(); safe && level.Next() != nil; level = level.Next() {
		prevDirection = direction

		prev, ok := level.Value.(int)

		if !ok {
			fmt.Fprintln(os.Stderr, "prev value not an int:", prev)
		}

		current, ok := level.Next().Value.(int)

		if !ok {
			fmt.Fprintln(os.Stderr, "current value not an int:", current)
		}

		fmt.Printf("prev: %d, current: %d\n", prev, current)

		if current == prev {
			safe = false
			fmt.Println("Failed because of no change in level")
			break
		} else if current > prev {
			direction = "inc"
		} else {
			direction = "dec"
		}

		if prevDirection != "none" && prevDirection != direction {
			safe = false
			fmt.Println("Failed because of change in direction")
			break
		}

		diff := current - prev
		absDiff := max(diff, -diff)

		if (absDiff == 0) || (absDiff > 3) {
			safe = false
			fmt.Printf("Failed because of diff out of bounds (%d)\n", absDiff)
			break
		}
	}

	if safe {
		return true
	} else {
		return false
	}
}

func isSafePart2(levels *list.List) bool {
	if isSafePart1(levels) {
		return true
	}

	isSafeWithoutOneLevel := false

	for i := 0; i < levels.Len(); i++ {
		levelsMinusALevel := list.New()
		levelsMinusALevel.PushBackList(levels)

		level := levelsMinusALevel.Front()

		for j := 0; j < i; j++ {
			level = level.Next()
		}

		levelsMinusALevel.Remove(level)

		if isSafePart1(levelsMinusALevel) {
			isSafeWithoutOneLevel = true
			break
		}
	}

	return isSafeWithoutOneLevel
}
