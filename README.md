# excerptByIDs

excerptByIDs is a tool to output records in a fastq or a pair of fastq files in which the
records match one of the IDs in IDfile, which is a list of IDs separated by newline
characters.

### Contents

README.md : this README

excerptByIDs.go : Go source code


### How to compile  
  
\# We have provided an executable compiled for use on 64-bit Linux systems in each release.  
\# If you need to compile the code for use on other architectures, you will need Go tools.  
\# If you do not have them, download and install Go tools as described here <https://golang.org/doc/install>  
$ git clone https://github.com/calacademy-research/excerptByIDs.git  
$ cd excerptByIDs  
$ go build excerptByIDs.go  


### Usage

$ excerptByIDs

       Usage: excerptByIDs <IDfile>|- <PE_file_1> [<PE_file_2>] [-ext <output_suffix>] [-v] [-mach] [-test]

       Outputs records from the FastQ files that match one of the IDs in IDfile.
       (The same ID can be present more than once but is used only once.)
       To read the IDs from stdin you can use the hyphen - as the first parameter.

       PE_file_1 records are written to stdout if it is the only file.
       If PE_file_1 and PE_file_2 are present, then new files are written where
       the file name has _extract inserted before the file suffix. Or, you can use
       -ext <output_suffix> to specify a string other than 'extract'.
       PE_file_1 is processed then PE_file_2.

       -test option, with just 1 PE_file arg writes IDs to stdout, with 2 PE_file args shows output names.
       -mach option includes the complete machine name in the ID (only needed if PE outputs combined).
       -v option inverts meaning of a matched record. Records NOT in the IDs are output.

### Citing

#### Authorship

Code author: James B. Henderson  
README.md authors: <a href="https://orcid.org/0000-0002-0210-7261" target="orcid.widget" rel="noopener noreferrer" style="vertical-align:top;"><img src="https://orcid.org/sites/default/files/images/orcid_16x16.png" style="width:1em;margin-right:.5em;" alt="ORCID iD icon">Zachary R. Hanna</a>, James B. Henderson  

#### Version 1.0.2
[![DOI](https://zenodo.org/badge/24128/calacademy-research/excerptByIDs.svg)](https://zenodo.org/badge/latestdoi/24128/calacademy-research/excerptByIDs)
