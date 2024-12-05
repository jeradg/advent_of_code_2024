package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

func main() {
	text, err := io.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(text)

	part1Total := parse(reader, false)

	_, err = reader.Seek(0, io.SeekStart)

	if err != nil {
		log.Fatal(err)
	}

	part2Total := parse(reader, true)

	fmt.Println("")
	fmt.Println("Total (Part 1):")
	fmt.Println(part1Total)

	fmt.Println("")
	fmt.Println("Total (Part 2):")
	fmt.Println(part2Total)
}

func parse(reader *bytes.Reader, useToggleInstructions bool) int {
	enabled := true
	total := 0
	statements := make([]string, 0)
	readingArgs := false
	statement := ""
	args := make([]int, 0)
	rawArg := ""

	fmt.Printf("useToggleInstructions: %v\n", useToggleInstructions)

	for {
		r, _, err := reader.ReadRune()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard input\n")
			os.Exit(1)
		}

		char := string(r)

		if useToggleInstructions &&
			(char == "d" ||
				(char == "o" && statement == "d") ||
				(char == "(" && statement == "do") ||
				(char == "n" && statement == "do") ||
				(char == "'" && statement == "don") ||
				(char == "t" && statement == "don'") ||
				(char == "(" && statement == "don't")) {
			statement += char
		} else if useToggleInstructions && (char == ")" && statement == "do(") {
			statement += char
			statements = append(statements, statement)
			enabled = true
			statement = ""
			args = make([]int, 0)
		} else if useToggleInstructions && (char == ")" && statement == "don't(") {
			statement += char
			statements = append(statements, statement)
			enabled = false
			statement = ""
			args = make([]int, 0)
		} else if enabled && readingArgs && (char == "," || char == ")") && len(rawArg) == 0 {
			readingArgs = false
			statement = ""
			args = make([]int, 0)
			rawArg = ""
		} else if enabled && readingArgs && (char == "," || char == ")") && len(rawArg) > 0 {
			statement += char

			num, err := strconv.Atoi(rawArg)

			if err != nil {
				fmt.Fprintln(os.Stderr, "error parsing numeric argument:", err)
				break
			}

			args = append(args, num)
			rawArg = ""

			if char == ")" {
				if enabled {
					total += multiply(args...)
					statements = append(statements, statement)
				}
				statement = ""
				args = make([]int, 0)
				readingArgs = false
			}
		} else if enabled && readingArgs && unicode.IsDigit(r) {
			statement += char
			rawArg += char
		} else if enabled && char == "m" ||
			(char == "u" && statement == "m") ||
			(char == "l" && statement == "mu") {
			statement += char
		} else if enabled && char == "(" && statement == "mul" {
			statement += char
			readingArgs = true
		} else {
			readingArgs = false
			statement = ""
			args = make([]int, 0)
			rawArg = ""
		}
	}

	fmt.Printf("Statements: %v\n", statements)

	return total
}

func multiply(nums ...int) int {
	total := 1

	for _, num := range nums {
		total *= num
	}

	return total
}
