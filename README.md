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
curl -O -L https://github.com/XinAlon/terraseq/releases/latest/download/terraseq.exe
# running the file
./terraseq.exe
```

#### **Linux**
```bash
# downloading the lastest version
wget https://github.com/XinAlon/terraseq/releases/latest/download/terraseq
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
```bash
usage: terraseq convert [--inFile INFILE] [--inFormat INFORMAT]
                        [--outFile OUTFILE] [--outFormat OUTFORMAT]

Parse optional command line arguments.

options:
  -h, --help                Display this help message and exit.
  --inFile INFILE           Specify the path to the input file (e.g., input.txt).
  --inFormat INFORMAT       Define the format of the input file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).
  --outFile OUTFILE         Specify the path for the output file (e.g., output.txt).
  --outFormat OUTFORMAT     Define the format of the output file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).
```


### Example: Extracting only snps from a .bim (plinkbed) file.
```bash
terraseq align --alignFile 1240K.bim --inFile myfile.txt --inFormat ftdnav1 --outFormat 23andme --outFile myfile_1240K.txt
```
#### Command Options: align
```bash
terraseq align -h
```
```bash
usage: terraseq align [--alignFile ALIGNFILE] [--inFile INFILE] [--inFormat INFORMAT]
                      [--outFile OUTFILE] [--outFormat OUTFORMAT]

Parse optional command line arguments.

options:
  -h, --help                Display this help message and exit.
  --alignFile ALIGNFILE     Specify the path to the alignment file (e.g., alignment.bim).
  --inFile INFILE           Specify the path to the input file (e.g., input.txt).
  --inFormat INFORMAT       Define the format of the input file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).
  --outFile OUTFILE         Specify the path for the output file (e.g., output.txt).
  --outFormat OUTFORMAT     Define the format of the output file (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage).
```
