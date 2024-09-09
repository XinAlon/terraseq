package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
)

func AncestryDNA(inFile, outFile, outFormat string) error {
	inputF, err := os.Open(inFile)
	if err != nil {
		return fmt.Errorf("Error opening input file: %v", err)
	}
	defer inputF.Close()

	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer output.Close()

	if outFormat == "23andme" {
		output.WriteString("# rsid\tchromosome\tposition\tgenotype\n")
	} else if outFormat == "ancestry" {
		output.WriteString("# rsid\tchromosome\tposition\tallele1\tallele2\n")
	} else {
		return fmt.Errorf("Unsupported output format: %s", outFormat)
	}

	scanner := bufio.NewScanner(inputF)
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

			if outFormat == "23andme" {
				outputLine := fmt.Sprintf("%s\t%s\t%s\t%s\n", rsid, chromosome, position, genotype)
				output.WriteString(outputLine)
			}
			if outFormat == "ancestry" {
				outputLine := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", rsid, chromosome, position, allele1, allele2)
				output.WriteString(outputLine)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading input file: %v", err)
	}
	return nil
}

func main() {
	// Create a new flag set
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)

	// Define flags for the convert command
	inFile := convertCmd.String("inFile", "", "Path to input file")
	inFormat := convertCmd.String("inFormat", "", "Format of input file")
	outFile := convertCmd.String("outFile", "", "Path to output file")
	outFormat := convertCmd.String("outFormat", "", "Format of output file")

	if len(os.Args) < 2 {
		fmt.Println("Usage: terraseq [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("   convert\t     Converts a DNA-file to another format.")
		return
	}

	if os.Args[1] != "convert" {
		fmt.Println("Usage: terraseq [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("   convert\t        Converts a DNA-file to another format.")
		return
	}

	convertCmd.Parse(os.Args[2:])

	if *inFile == "" || *inFormat == "" || *outFile == "" || *outFormat == "" {
		fmt.Println("Missing Flags.")
		fmt.Println()
		fmt.Println("Usage: terraseq convert --inFile --inFormat --outFile --outFormat")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("   --inFile\t        Path to your input file.")
		fmt.Println("   --inFormat\t        Format of your input file.")
		fmt.Println("   --outFile\t        Path to your output file.")
		fmt.Println("   --outFormat\t        Format of your output file.")
		return
	}

	var err error
	if *inFormat == "ancestry" {
		err = AncestryDNA(*inFile, *outFile, *outFormat)
	} else {
		fmt.Println("Unsupported format conversion")
		return
	}

	// Handle conversion errors
	if err != nil {
		fmt.Printf("Error during conversion: %v\n", err)
		return
	}

	fmt.Println("Conversion completed successfully")
}
