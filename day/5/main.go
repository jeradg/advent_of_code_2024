package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Key: a page number, Value: a map with keys of the pages that this page must appear *before*, each of which has the val `true`
// E.g., { "5": { "99": true }, "7", { "6": true, "54": true } }
type Rules map[string](map[string]bool)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	rulePairs := make([][]string, 0)
	updatePages := make([][]string, 0)
	haveSeenBlankLine := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			haveSeenBlankLine = true
		} else if !haveSeenBlankLine {
			rulePairs = append(rulePairs, strings.Split(line, "|"))
		} else {
			updatePages = append(updatePages, strings.Split(line, ","))
		}
	}

	fmt.Println("rulePairs:")
	fmt.Printf("%v\n\n", rulePairs)

	fmt.Println("updatePages:")
	fmt.Printf("%v\n\n", updatePages)

	rules := rulesForRulePairs(rulePairs)

	fmt.Printf("%+v\n", rules)

	validPages, invalidPages := validAndInvalidUpdatePages(updatePages, rules)

	fmt.Println("Valid pages:")
	fmt.Println(validPages)

	fmt.Println("Part 1 answer:")
	fmt.Println(sumMiddlePages(validPages))

	fmt.Println("")

	fmt.Println("Invalid pages:")
	fmt.Println(invalidPages)

	fmt.Println("")

	fixedInvalidPages := fixInvalidPages(invalidPages, rules)

	fmt.Println("Part 2 answer:")
	fmt.Println(sumMiddlePages(fixedInvalidPages))
}

func rulesForRulePairs(rulePairs [][]string) Rules {
	rules := Rules{}

	for _, rulePair := range rulePairs {
		left, right := rulePair[0], rulePair[1]

		if afterPages, pageExists := rules[left]; pageExists {
			afterPages[right] = true
		} else {
			afterPages := make(map[string]bool)
			afterPages[right] = true
			rules[left] = afterPages
		}
	}

	return rules
}

func validAndInvalidUpdatePages(updatePages [][]string, rules Rules) ([][]string, [][]string) {
	validUpdatePages := make([][]string, 0)
	invalidUpdatePages := make([][]string, 0)

	for _, pages := range updatePages {
		valid := arePagesValid(pages, rules)

		if valid {
			validUpdatePages = append(validUpdatePages, pages)
		} else {
			invalidUpdatePages = append(invalidUpdatePages, pages)
		}
	}

	return validUpdatePages, invalidUpdatePages
}

func arePagesValid(pages []string, rules Rules) bool {
	valid := true

	pagesSeen := make([]string, 0)

	for _, page := range pages {
		for _, pageSeen := range pagesSeen {
			if rules[page][pageSeen] {
				valid = false
				break
			}
		}

		if !valid {
			break
		}

		pagesSeen = append(pagesSeen, page)
	}

	return valid
}

func fixInvalidPages(invalidPages [][]string, rules Rules) [][]string {
	fixedPages := make([][]string, 0)

	for _, pages := range invalidPages {
		attempt := tryFixPages(pages, rules)

		if arePagesValid(attempt, rules) {
			fmt.Printf("Success! Fixed pages with attempt: %v\n", attempt)

			fixedPages = append(fixedPages, attempt)
		} else {
			fmt.Fprintf(os.Stderr, "Attempt didn't work, WHY NOT???\n")
			fmt.Fprintf(os.Stderr, "attempt: %v\n", attempt)
			os.Exit(1)
		}
	}

	if len(fixedPages) != len(invalidPages) {
		fmt.Fprintf(os.Stderr, "Length of fixedPages doesn't match length of invalidPages, WHY NOT???\n")
		fmt.Fprintf(os.Stderr, "len(fixedPages): %v, len(invalidPages): %v\n", len(fixedPages), len(invalidPages))
		os.Exit(1)
	}

	return fixedPages
}

func tryFixPages(pages []string, rules Rules) []string {
	invalidPage := ""
	invalidPageIndex := -1
	mustGoBeforeIndex := -1

	pagesSeen := make([]string, 0)

	for i, page := range pages {
		for j, pageSeen := range pagesSeen {
			if rules[page][pageSeen] {
				invalidPage = page
				invalidPageIndex = i
				mustGoBeforeIndex = j
				break
			}
		}

		if invalidPage != "" {
			break
		}

		pagesSeen = append(pagesSeen, page)
	}

	if invalidPage == "" {
		fmt.Fprintf(os.Stderr, "Didn't find invalid page, WHY NOT???\n")
		os.Exit(1)
	}

	if invalidPageIndex == -1 {
		fmt.Fprintf(os.Stderr, "No invalid page index, WHY NOT???\n")
		os.Exit(1)
	}

	if mustGoBeforeIndex == -1 {
		fmt.Fprintf(os.Stderr, "No must-go-before index, WHY NOT???\n")
		os.Exit(1)
	}

	fmt.Printf("invalidPage: %v, invalidPageIndex: %v, mustGoBeforeIndex: %v\n", invalidPage, invalidPageIndex, mustGoBeforeIndex)

	attempt := append(make([]string, 0), pages[:mustGoBeforeIndex]...)
	attempt = append(attempt, invalidPage)
	attempt = append(attempt, pages[mustGoBeforeIndex:invalidPageIndex]...)
	attempt = append(attempt, pages[invalidPageIndex+1:]...)

	if arePagesValid(attempt, rules) {
		fmt.Printf("Success! Fixed pages with attempt: %v\n", attempt)

		return attempt
	} else {
		fmt.Println("Trying again!")

		return tryFixPages(attempt, rules)
	}
}

func sumMiddlePages(updatePages [][]string) int {
	total := 0

	for _, pages := range updatePages {
		middleIndex := len(pages) / 2

		num, err := strconv.Atoi(pages[middleIndex])

		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse middle num, WHAT???\n")
			os.Exit(1)
		}

		total += num
	}

	return total
}
