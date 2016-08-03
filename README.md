# excerptByIDs

excerptByIDs is a tool to output records in a fastq or a pair of fastq files in which the
records match one of the IDs in IDfile, which is a list of IDs separated by newline
characters.


#### Contents

README.md : this README

excerptByIDs.go : Go source code


#### How to compile

$ git clone <git-repository-path>
$ cd <git-repository-name>
$ go build excerptByIDs.go


#### Usage

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


#### Authorship

author: James Henderson, jhenderson@calacademy.org
