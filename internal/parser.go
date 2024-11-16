package internal

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"path/filepath"
)

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

func ParseTemplate(filename string) ([]TemplateRecord, error) {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".bim" && ext != ".snp" {
		return nil, fmt.Errorf("unsupported file extension: %s. Only .bim and .snp files are supported", ext)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening alignFile: %v", err)
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

		if ext == ".bim" {
			if len(fields) < 6 {
				continue // Skip invalid lines
			}
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
		} else { // .snp file
			if len(fields) < 6 {
				continue // Skip invalid lines
			}
			value, err := parseScientificNotation(fields[2])
			if err != nil {
				continue // Skip lines with invalid scientific notation
			}
			record := TemplateRecord{
				Chromosome:  fields[1],
				RSID:       fields[0],
				Value:      value,
				Position:   fields[3],
				ReferenceA1: fields[4],
				ReferenceA2: fields[5],
			}
			records = append(records, record)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading alignFile: %v", err)
	}

	return records, nil
}

func parseScientificNotation(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
