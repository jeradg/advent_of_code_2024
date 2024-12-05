package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var leftList []int
	var rightList []int
	totalDiff := 0
	timesLeftInRight := make(map[int]int)
	prevResultIdx := 0
	totalSimilarityScore := 0

	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		leftNum, err := strconv.Atoi(words[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading left num:", err)
		}

		rightNum, err := strconv.Atoi(words[1])

		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading right num:", err)
		}

		leftList = append(leftList, leftNum)
		rightList = append(rightList, rightNum)
	}

	slices.Sort(leftList)
	slices.Sort(rightList)

	for i := 0; i < len(leftList); i++ {
		diff := leftList[i] - rightList[i]
		absDiff := max(diff, -diff)

		if _, ok := timesLeftInRight[leftList[i]]; !ok {
			timesInRightList := 0

			for ; (prevResultIdx < len(leftList)) && (rightList[prevResultIdx] <= leftList[i]); prevResultIdx++ {
				if leftList[i] == rightList[prevResultIdx] {
					timesInRightList++
				}
			}

			timesLeftInRight[leftList[i]] = timesInRightList
		}

		totalSimilarityScore += leftList[i] * timesLeftInRight[leftList[i]]
		totalDiff += absDiff
	}

	fmt.Println("Diff is:")
	fmt.Println(totalDiff)

	fmt.Println("Total similarity score is:")
	fmt.Println(totalSimilarityScore)
}
