# terraseq

1. [Installation](#installation)
   - [Windows](#windows)
   - [Linux](#linux)
2. [Usage](#usage)

---
## **Installation**

#### **Windows**
```bash
# downloading the lastest version
curl -O -L https://github.com/enelsr/terraseq/releases/latest/download/terraseq.exe
# running the file
./terraseq.exe
```

#### **Linux**
```bash
# downloading the lastest version
wget https://github.com/enelsr/terraseq/releases/latest/download/terraseq
# make the file executable
chmod +x terraseq
# running the file
./terraseq
```

---
## **Usage**

### Example: Converting file to another format.
```bash
terraseq convert --inFile myfile.txt --inFormat ancestry --outFormat 23andme --outFile myfile_converted.txt
```
#### Command Options: convert
```bash
terraseq convert -h
```
```
usage: terraseq convert [-i|--inFile FILE] [-f|--inFormat FORMAT]
                      [-o|--outFile FILE] [-t|--outFormat FORMAT]

Parse optional command line arguments.

options:
  -h, --help                  Display this help message and exit
  -i, --inFile FILE           Specify the path to the input file
                              (e.g., input.txt)
  -f, --inFormat FORMAT       Define the input file format
                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)
  -o, --outFile FILE          Specify the path for the output file
                              (e.g., output.txt)
  -t, --outFormat FORMAT      Define the output file format
                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)
```


### Example: Extracting only snps from a .bim (plink) or .snp (eigenstrat|packedancestrymap) file.
```bash
terraseq align --alignFile 1240K.bim --inFile myfile.csv --inFormat ftdnav1 --outFormat 23andme --outFile myfile_1240K.txt
```
#### Command Options: align
```bash
terraseq align -h
```
```
usage: terraseq align [-a|--alignFile FILE] [-i|--inFile FILE] [-f|--inFormat FORMAT]
                      [-o|--outFile FILE] (-t|--outFormat FORMAT) (--flip)

Parse optional command line arguments.

options:
  -h, --help                  Display this help message and exit
  -a, --alignFile FILE        Specify the path to the alignment file
                              (e.g., 1240K.bim, 1240K.snp)
  -i, --inFile FILE           Specify the path to the input file
                              (e.g., input.txt)
  -f, --inFormat FORMAT       Define the format of the input file
                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)
  -o, --outFile FILE          Specify the path for the output file
                              (e.g., output.txt)
  -t, --outFormat FORMAT      Define the format of the output file
                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)
  --flip                      Flips the alleles in accordance with the reference
```
