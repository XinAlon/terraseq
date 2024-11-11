package internal

type DNARecord struct {
	RSID        string
	Chromosome  string
	Position    string
	Allele1     string
	Allele2     string
	RawGenotype string
}

type DNAData struct {
	Records []DNARecord
	Format  string
}

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
