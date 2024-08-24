package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: convertAncestry input.txt output.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	input, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	outputHeader := "# rsid\tchromosome\tposition\tgenotype\n"
	output.WriteString(outputHeader)

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if line == "rsid\tchromosome\tposition\tallele1\tallele2" {
			continue
		}

		fields := strings.Split(line, "\t")

		if len(fields) >= 5 {
			rsid := fields[0]
			chromosome := fields[1]
			position := fields[2]
			allele1 := fields[3]
			allele2 := fields[4]

			genotype := allele1 + allele2

			outputLine := fmt.Sprintf("%s\t%s\t%s\t%s\n", rsid, chromosome, position, genotype)

			output.WriteString(outputLine)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
	}
}
