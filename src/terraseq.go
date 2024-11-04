package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
	"strconv"
)

var Green = "\033[32m"
var Reset = "\033[0m"
var Red = "\033[31m"

// DNARecord represents a single DNA record regardless of input format
type DNARecord struct {
	RSID        string
	Chromosome  string
	Position    string
	Allele1     string
	Allele2     string
	RawGenotype string
}

// DNAData represents a collection of DNA records with metadata
type DNAData struct {
	Records []DNARecord
	Format  string
}

// ParseResult represents the outcome of parsing operations
type ParseResult struct {
	Data DNAData
	Err  error
}

type TemplateRecord struct {
	Chromosome  string
	RSID        string
	Value       float64
	Position    string
	ReferenceA1 string
	ReferenceA2 string
}

func ParseTemplate(filename string) ([]TemplateRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening template file: %v", err)
	}
	defer file.Close()

	var records []TemplateRecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			value, err := parseScientificNotation(fields[2])
			if err != nil {
				continue // Skip lines with invalid scientific notation
			}

			record := TemplateRecord{
				Chromosome:  fields[0],
				RSID:       fields[1],
				Value:      value,
				Position:   fields[3],
				ReferenceA1: fields[4],
				ReferenceA2: fields[5],
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading template file: %v", err)
	}

	return records, nil
}

// parseScientificNotation converts scientific notation strings to float64
func parseScientificNotation(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// AlignDNA aligns DNA records with a template and writes the result
func AlignDNA(data DNAData, templateRecords []TemplateRecord, outFile string, outFormat string) error {
	// Create a map for quick lookup of DNA records by RSID
	dnaMap := make(map[string]DNARecord)
	for _, record := range data.Records {
		dnaMap[record.RSID] = record
	}

	// Create output file
	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	// Write header based on format
	switch outFormat {
		case "23andme":
			output.WriteString("# rsid\tchromosome\tposition\tgenotype\n")
		case "ancestry":
			output.WriteString("# rsid\tchromosome\tposition\tallele1\tallele2\n")
		case "ftdnav2", "ftdnav1", "myheritage":
			output.WriteString("RSID,CHROMOSOME,POSITION,RESULT\n")
		default:
			return fmt.Errorf("unsupported output format: %s", outFormat)
	}

	// Track statistics
	var totalSnps, matchedSnps int

	// Process each template record
	for _, template := range templateRecords {
		totalSnps++
		var outputLine string

		if dnaRecord, exists := dnaMap[template.RSID]; exists {
			matchedSnps++
			// Use the actual DNA record data
			switch outFormat {
				case "23andme":
					outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\n",
								 template.RSID, template.Chromosome, template.Position, dnaRecord.RawGenotype)
				case "ancestry":
					outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
								 template.RSID, template.Chromosome, template.Position, dnaRecord.Allele1, dnaRecord.Allele2)
				case "ftdnav2":
					outputLine = fmt.Sprintf("%s,%s,%s,%s\n",
								 template.RSID, template.Chromosome, template.Position, dnaRecord.RawGenotype)
				case "ftdnav1", "myheritage":
					outputLine = fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"\n",
								 template.RSID, template.Chromosome, template.Position, dnaRecord.RawGenotype)
			}
		} else {
			// Use default values for missing SNPs
			defaultAllele1, defaultAllele2, defaultGenotype := getDefaultGenotype(template, outFormat)

				switch outFormat {
					case "23andme":
						outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\n",
									 template.RSID, template.Chromosome, template.Position, defaultGenotype)
					case "ancestry":
						outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
									 template.RSID, template.Chromosome, template.Position, defaultAllele1, defaultAllele2)
					case "ftdnav2":
						outputLine = fmt.Sprintf("%s,%s,%s,%s\n",
									 template.RSID, template.Chromosome, template.Position, defaultGenotype)
					case "ftdnav1", "myheritage":
						outputLine = fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"\n",
									 template.RSID, template.Chromosome, template.Position, defaultGenotype)
				}
		}

		output.WriteString(outputLine)
	}

	// Print statistics
	fmt.Printf("%s[INFO]%s Total SNPs in template: %d\n", Green, Reset, totalSnps)
	fmt.Printf("%s[INFO]%s Matched SNPs: %d (%.1f%%)\n", Green, Reset, matchedSnps, float64(matchedSnps)/float64(totalSnps)*100)
	fmt.Printf("%s[INFO]%s Missing SNPs: %d (%.1f%%)\n", Green, Reset, totalSnps-matchedSnps, float64(totalSnps-matchedSnps)/float64(totalSnps)*100)

	return nil
}


// align is the main function for DNA alignment
func align(inFile, inFormat, outFile, outFormat, alignFile string) error {
	// Parse input file
	var result ParseResult
	switch inFormat {
		case "ancestry":
			result = ParseAncestryDNA(inFile)
		case "23andme":
			result = Parse23andMe(inFile)
		case "ftdnav2":
			result = ParseFTDNA(inFile)
		case "ftdnav1", "myheritage":
			result = ParseMyHeritage(inFile)
		default:
			return fmt.Errorf("unsupported input format: %s", inFormat)
	}

	if result.Err != nil {
		return result.Err
	}

	// Parse template file
	templateRecords, err := ParseTemplate(alignFile)
	if err != nil {
		return fmt.Errorf("error parsing template file: %v", err)
	}

	// Align DNA with template and write output
	return AlignDNA(result.Data, templateRecords, outFile, outFormat)
}

func getDefaultGenotype(template TemplateRecord, outFormat string) (string, string, string) {
	// Default to the reference alleles
	allele1 := template.ReferenceA1
	allele2 := template.ReferenceA2
	rawGenotype := allele1 + allele2

	// For some formats, we might want to use different defaults
	switch outFormat {
		case "23andme", "ftdnav2", "ftdnav1", "myheritage":
			// These formats typically use "--" for missing data
			return "--", "--", "--"
		case "ancestry":
			// Ancestry format might prefer explicit reference alleles
			return "0", "0", rawGenotype
		default:
			return allele1, allele2, rawGenotype
	}
}

// ParseAncestryDNA reads and parses an Ancestry DNA file format
func ParseAncestryDNA(filename string) ParseResult {
	file, err := os.Open(filename)
	if err != nil {
		return ParseResult{Err: fmt.Errorf("error opening file: %v", err)}
	}
	defer file.Close()

	var records []DNARecord
	scanner := bufio.NewScanner(file)

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
			record := DNARecord{
				RSID:        fields[0],
				Chromosome:  fields[1],
				Position:    fields[2],
				Allele1:     fields[3],
				Allele2:     fields[4],
				RawGenotype: fields[3] + fields[4],
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return ParseResult{Err: fmt.Errorf("error reading file: %v", err)}
	}

	return ParseResult{
		Data: DNAData{
			Records: records,
			Format:  "ancestry",
		},
	}
}

// Parse23andMe reads and parses a 23andMe file format
func Parse23andMe(filename string) ParseResult {
	file, err := os.Open(filename)
	if err != nil {
		return ParseResult{Err: fmt.Errorf("error opening file: %v", err)}
	}
	defer file.Close()

	var records []DNARecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if line == "rsid\tchromosome\tposition\tgenotype" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) >= 4 {
			genotype := fields[3]
			allele1 := string(genotype[0])
			allele2 := allele1
			if len(genotype) > 1 {
				allele2 = string(genotype[1])
			}

			record := DNARecord{
				RSID:        fields[0],
				Chromosome:  fields[1],
				Position:    fields[2],
				Allele1:     allele1,
				Allele2:     allele2,
				RawGenotype: genotype,
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return ParseResult{Err: fmt.Errorf("error reading file: %v", err)}
	}

	return ParseResult{
		Data: DNAData{
			Records: records,
			Format:  "23andme",
		},
	}
}

// ParseFTDNA reads and parses an FTDNA file format
func ParseFTDNA(filename string) ParseResult {
	file, err := os.Open(filename)
	if err != nil {
		return ParseResult{Err: fmt.Errorf("error opening file: %v", err)}
	}
	defer file.Close()

	var records []DNARecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if line == "RSID,CHROMOSOME,POSITION,RESULT" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 4 {
			genotype := fields[3]
			allele1 := string(genotype[0])
			allele2 := allele1
			if len(genotype) > 1 {
				allele2 = string(genotype[1])
			}

			record := DNARecord{
				RSID:        fields[0],
				Chromosome:  fields[1],
				Position:    fields[2],
				Allele1:     allele1,
				Allele2:     allele2,
				RawGenotype: genotype,
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return ParseResult{Err: fmt.Errorf("error reading file: %v", err)}
	}

	return ParseResult{
		Data: DNAData{
			Records: records,
			Format:  "ftdna",
		},
	}
}

// ParseMyHeritage reads and parses a MyHeritage file format
func ParseMyHeritage(filename string) ParseResult {
	file, err := os.Open(filename)
	if err != nil {
		return ParseResult{Err: fmt.Errorf("error opening file: %v", err)}
	}
	defer file.Close()

	var records []DNARecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if line == "RSID,CHROMOSOME,POSITION,RESULT" {
			continue
		}

		fields := strings.Split(line, ",")
		for i, field := range fields {
			fields[i] = strings.Trim(field, "\"")
		}

		if len(fields) >= 4 {
			genotype := fields[3]
			allele1 := string(genotype[0])
			allele2 := allele1
			if len(genotype) > 1 {
				allele2 = string(genotype[1])
			}

			record := DNARecord{
				RSID:        fields[0],
				Chromosome:  fields[1],
				Position:    fields[2],
				Allele1:     allele1,
				Allele2:     allele2,
				RawGenotype: genotype,
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return ParseResult{Err: fmt.Errorf("error reading file: %v", err)}
	}

	return ParseResult{
		Data: DNAData{
			Records: records,
			Format:  "myheritage",
		},
	}
}

// WriteDNAData writes DNA records to a file in the specified format
func WriteDNAData(data DNAData, outFile string, outFormat string) error {
	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

	// Write header
	switch outFormat {
		case "23andme":
			output.WriteString("# rsid\tchromosome\tposition\tgenotype\n")
		case "ancestry":
			output.WriteString("# rsid\tchromosome\tposition\tallele1\tallele2\n")
		case "ftdnav2", "ftdnav1", "myheritage":
			output.WriteString("RSID,CHROMOSOME,POSITION,RESULT\n")
		default:
			return fmt.Errorf("unsupported output format: %s", outFormat)
	}

	// Write records
	for _, record := range data.Records {
		var outputLine string
		switch outFormat {
			case "23andme":
				outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\n",
							 record.RSID, record.Chromosome, record.Position, record.RawGenotype)
			case "ancestry":
				outputLine = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
							 record.RSID, record.Chromosome, record.Position, record.Allele1, record.Allele2)
			case "ftdnav2":
				outputLine = fmt.Sprintf("%s,%s,%s,%s\n",
							 record.RSID, record.Chromosome, record.Position, record.RawGenotype)
			case "ftdnav1", "myheritage":
				outputLine = fmt.Sprintf("\"%s\",\"%s\",\"%s\",\"%s\"\n",
							 record.RSID, record.Chromosome, record.Position, record.RawGenotype)
		}
		output.WriteString(outputLine)
	}

	return nil
}

func convert(inFile, inFormat, outFile, outFormat string) error {
	var result ParseResult

	// Parse input file based on format
	switch inFormat {
		case "ancestry":
			result = ParseAncestryDNA(inFile)
		case "23andme":
			result = Parse23andMe(inFile)
		case "ftdnav2":
			result = ParseFTDNA(inFile)
		case "ftdnav1", "myheritage":
			result = ParseMyHeritage(inFile)
		default:
			return fmt.Errorf("unsupported input format: %s", inFormat)
	}

	if result.Err != nil {
		return result.Err
	}

	// Write output file
	return WriteDNAData(result.Data, outFile, outFormat)
}

func main() {
	// Create a new flag set
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)
	alignCmd := flag.NewFlagSet("align", flag.ExitOnError)

	// Define flags for the convert command
	inFile := convertCmd.String("inFile", "", "Path to input file")
	inFormat := convertCmd.String("inFormat", "", "Format of input file")
	outFile := convertCmd.String("outFile", "", "Path to output file")
	outFormat := convertCmd.String("outFormat", "", "Format of output file")

	alignInFile := alignCmd.String("inFile", "", "Path to input file")
	alignInFormat := alignCmd.String("inFormat", "", "Format of input file")
	alignOutFile := alignCmd.String("outFile", "", "Path to output file")
	alignOutFormat := alignCmd.String("outFormat", "", "Format of output file")
	alignFile := alignCmd.String("alignFile", "", "Path to your alignment file")

	// Help flags
	convertHelp := convertCmd.Bool("help", false, "Show convert command help")
	convertHelpShort := convertCmd.Bool("h", false, "Show convert command help")
	alignHelp := alignCmd.Bool("help", false, "Show align command help")
	alignHelpShort := alignCmd.Bool("h", false, "Show align command help")

	if len(os.Args) < 2 {
		fmt.Println("Usage: terraseq [command]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("   convert\t     Converts a DNA file to another format.")
		fmt.Println("   align\t     Aligns DNA sequences with a reference.")
		return
	}

	switch os.Args[1] {
		case "convert":
			convertCmd.Parse(os.Args[2:])

			if *convertHelp || *convertHelpShort {
				printConvertHelp()
				return
			}

			var missingFlags []string
			if *inFile == "" {
				missingFlags = append(missingFlags, "--inFile")
			}
			if *inFormat == "" {
				missingFlags = append(missingFlags, "--inFormat")
			}
			if *outFile == "" {
				missingFlags = append(missingFlags, "--outFile")
			}
			if *outFormat == "" {
				missingFlags = append(missingFlags, "--outFormat")
			}

			if len(missingFlags) > 0 && len(missingFlags) < 4 {
				printConvertHelp()
				fmt.Println()
				fmt.Println(Red + "[Warning]" + Reset + " Missing Flags: \n", strings.Join(missingFlags, ", "))
				fmt.Println()
				return
			} else if len(missingFlags) == 4 {
				printConvertHelp()
				return
			}

			fmt.Println(Green + "[INFO]" + Reset + " Converting...")
			if err := convert(*inFile, *inFormat, *outFile, *outFormat); err != nil {
				fmt.Printf(Red + "[Warning]" + Reset + " Error during conversion: %v\n", err)
				return
			}

			fmt.Println(Green + "[INFO]" + Reset + " Conversion completed successfully!")

		case "align":
			alignCmd.Parse(os.Args[2:])

			if *alignHelp || *alignHelpShort {
				printAlignHelp()
				return
			}

			var missingFlags []string
			if *alignFile == "" {
				missingFlags = append(missingFlags, "--alignFile")
			}
			if *alignInFile == "" {
				missingFlags = append(missingFlags, "--inFile")
			}
			if *alignInFormat == "" {
				missingFlags = append(missingFlags, "--inFormat")
			}
			if *alignOutFile == "" {
				missingFlags = append(missingFlags, "--outFile")
			}
			if *alignOutFormat == "" {
				missingFlags = append(missingFlags, "--outFormat")
			}

			if len(missingFlags) > 0 && len(missingFlags) < 5 {
				printAlignHelp()
				fmt.Println()
				fmt.Println(Red + "[Warning]" + Reset + " Missing Flags: \n", strings.Join(missingFlags, ", "))
				fmt.Println()
				return
			} else if len(missingFlags) == 5 {
				printAlignHelp()
				return
			}

			fmt.Println(Green + "[INFO]" + Reset + " Aligning DNA...")
			if err := align(*alignInFile, *alignInFormat, *alignOutFile, *alignOutFormat, *alignFile); err != nil {
				fmt.Printf(Red + "[Warning]" + Reset + " Error during alignment: %v\n", err)
				return
			}
			fmt.Println(Green + "[INFO]" + Reset + " Alignment completed successfully!")

		default:
			fmt.Println("Usage: terraseq [command]")
			fmt.Println()
			fmt.Println("Commands:")
			fmt.Println("   convert\t     Converts a DNA file to another format.")
			fmt.Println("   align\t     Aligns DNA sequences with a reference.")
			return
	}
}

func printConvertHelp() {
	fmt.Println("usage: terraseq convert [--inFile INFILE] [--inFormat INFORMAT]")
	fmt.Println("                        [--outFile OUTFILE] [--outFormat OUTFORMAT]")
	fmt.Println()
	fmt.Println("Parse optional command line arguments.")
	fmt.Println()
	fmt.Println("options:")
	fmt.Println("  -h, --help                Display this help message and exit.")
	fmt.Println("  --inFile INFILE           Specify the path to the input file (e.g., input.txt).")
	fmt.Println("  --inFormat INFORMAT       Define the format of the input file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).")
	fmt.Println("  --outFile OUTFILE         Specify the path for the output file (e.g., output.txt).")
	fmt.Println("  --outFormat OUTFORMAT     Define the format of the output file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).")
}

func printAlignHelp() {
	fmt.Println("usage: terraseq align [--alignFile ALIGNFILE] [--inFile INFILE] [--inFormat INFORMAT]")
	fmt.Println("                      [--outFile OUTFILE] [--outFormat OUTFORMAT]")
	fmt.Println()
	fmt.Println("Parse optional command line arguments.")
	fmt.Println()
	fmt.Println("options:")
	fmt.Println("  -h, --help                Display this help message and exit.")
	fmt.Println("  --alignFile ALIGNFILE     Specify the path to the alignment file (e.g., alignment.bim).")
	fmt.Println("  --inFile INFILE           Specify the path to the input file (e.g., input.txt).")
	fmt.Println("  --inFormat INFORMAT       Define the format of the input file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).")
	fmt.Println("  --outFile OUTFILE         Specify the path for the output file (e.g., output.txt).")
	fmt.Println("  --outFormat OUTFORMAT     Define the format of the output file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).")
}

