package internal

import (
	"os"
	"fmt"
)

func WriteDNAData(data DNAData, outFile string, outFormat string) error {
	output, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer output.Close()

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
			// Ancestry format uses 0 for missing data
			return "0", "0", rawGenotype
		default:
			return allele1, allele2, rawGenotype
	}
}

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
	fmt.Printf("[INFO] Total SNPs in template: %d\n", totalSnps)
	fmt.Printf("[INFO] Matched SNPs: %d (%.1f%%)\n", matchedSnps, float64(matchedSnps)/float64(totalSnps)*100)
	fmt.Printf("[INFO] Missing SNPs: %d (%.1f%%)\n", totalSnps-matchedSnps, float64(totalSnps-matchedSnps)/float64(totalSnps)*100)

	return nil
}

