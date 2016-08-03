package main

// 19Jul14 JBH uses read IDs file of arg1, and reads through FastQ PE files arg2 arg3,
// which are expected to have same IDs at the same record (no error checking for this)

import ("os"; "bufio"; "fmt"; "strings"; "path"; "path/filepath"; "time"; "log"; "local/util")
const (CHECK_EVERY = 100000; WAIT_AT_LEAST = time.Millisecond*250)
var testParser = false
var makeFiles = false
var pe1_outname = ""
var pe2_outname = ""
var includeMachName = false // 06Jan15 JBH force machname to be included

func main() {
	
	idFile, pe1, pe2, output_suffix, excerptNonMatches, includeMach := getopts() // get filename to process and type of hdr fixup
	includeMachName = includeMach
	if pe1 == "" {
		usage()
		return
	}
	if pe2 != "" { // we need to write to files not just stdout
		makeFiles = true
		pe1_outname = makeOutfilename(pe1, output_suffix)
		pe2_outname = makeOutfilename(pe2, output_suffix)
		
		if testParser { // just show the output names
			fmt.Println(pe1_outname)
			fmt.Println(pe2_outname)
			return
		}
	}
	
	// fill idMap from the id file or stdin
	idMap := make(map[string]int) // we'll just use int 1 for our IDs that exist
	var id_scanner *bufio.Scanner
	if (idFile == "-" || idFile == "stdin") { // - as filename means read ID file from stdin, also string "stdin" does same
		id_scanner = bufio.NewScanner(os.Stdin)
	} else {
		ids, _ := os.Open(idFile)
		id_scanner = bufio.NewScanner(ids)
		defer ids.Close()
	}
	for id_scanner.Scan() {
		line := id_scanner.Text()
		if strings.Contains(line, ":") {
			id := hdrID(line, true)
			idMap[id] = 1
		}
	}
	fmt.Fprintln(os.Stderr, len(idMap), "IDs.")
	
	if testParser { // if we are just testing the parser, then do nothing else
		return 
	}
	
	excerptFromFile(pe1, idMap, output_suffix, !makeFiles, excerptNonMatches)
	
	if pe2 != "" {
		excerptFromFile(pe2, idMap, output_suffix, false, excerptNonMatches)
	}
}

func excerptFromFile(filename string, idMap map[string]int, output_suffix string, is_stdout bool, excerptNonMatching bool) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)
	
	fileSize := int64(-1)
	fi, err := file.Stat()
	if err == nil { fileSize = fi.Size() }
	
	// set up output
	var fileOut *os.File
	var outmsg = ""
	if is_stdout {
		fileOut = os.Stdout
	} else {
		outname := makeOutfilename(filename, output_suffix)
		fileOut, _ = os.Create(outname)
		defer fileOut.Close()
		outmsg = "to " + outname
	}
	
	fmt.Fprintln(os.Stderr, "Excerpting from", path.Base(filename), outmsg)
	
	numRecs := 0; numExcerpts := 0
	lastMsg := ""; newMsg := "" // progress msg written to stderr
	lastDisplayTime := time.Now()
	blanksToPad := "";
	
	for scanner.Scan() {
		hdr := scanner.Text()
		matched := idMap[ hdrID(hdr, false) ] == 1
		doExcerpt := matched // 01Jan15 JBH added option to excerpt records that DON'T match any of our IDs
		if (excerptNonMatching) {
			doExcerpt = !matched
		}
			
		if doExcerpt {
			numExcerpts++;
			fmt.Fprintln(fileOut, hdr)		
		}
		scanner.Scan()
		seq := scanner.Bytes()
		if doExcerpt {
			fmt.Fprintln(fileOut, string(seq))
		}
		scanner.Scan()
		plus := scanner.Bytes()
		if doExcerpt {
			fmt.Fprintln(fileOut, string(plus))
		}
		scanner.Scan()
		qual := scanner.Bytes()
		if doExcerpt {
			fmt.Fprintln(fileOut, string(qual))
		}
		
		numRecs++
		if ( (numRecs % CHECK_EVERY == 0) && (time.Now().Sub(lastDisplayTime) > WAIT_AT_LEAST)) {
			fpos, _ := file.Seek(0, os.SEEK_CUR) // baroque equivalent of tell() for the file offset used here
			pct := float64(fpos) / float64(fileSize)
			newMsg = util.Comma(int64(numRecs)) + " records " + util.FloatToPctStr(pct, 2) + " " + util.Comma(int64(numExcerpts)) + " excerpts."
			if len(lastMsg) > len(newMsg) {
				btp := len(lastMsg) - len(newMsg)
				blanksToPad = strings.Repeat(" ", btp)
			}
			newMsg += blanksToPad
			fmt.Fprint(os.Stderr, strings.Repeat("\b", len(lastMsg)) + newMsg)
			lastMsg = newMsg //strings.TrimRight(newMsg, " ")
			lastDisplayTime = time.Now()
		}
		
	}
	if err := scanner.Err(); err != nil { log.Fatal(err) }
	
	fmt.Fprintln(os.Stderr, strings.Repeat("\b", len(lastMsg)) + util.Comma(int64(numRecs)) + " records total, " +util.Comma(int64(numExcerpts))+ " excerpted.")
}

func hdrID(hdr string, bCanBe2ndFld bool) string {
	// 31Jul14 JBH make this so it handles m8 style records (either tabbed or comma delimited fields)
	// we're presuming an Illumina Casava 1.8 style ID in first field (or entirety) of line
	// 25Aug14 add option for checking 2nd field for ":" formatted ID if no colons in 1st field
	// 06Jan15 JBH add option to handle older Illumina IDs if there are 5 fields and last one has '#' or '/' in it
	// 06Jan15 JBH if includeMachName golbal is set prefix it to recID
	
	ixFldDelim := strings.IndexAny(hdr, "\t, ") // use first found of tab comma or space as a delimited
	if ixFldDelim > 0 { // returns -1 if not found but we don't want a delimiter as first field either
		if !bCanBe2ndFld {
			hdr = hdr[:ixFldDelim]
		} else { // if first field doesn't have colons, try second
			fld1 := hdr[:ixFldDelim]
			if strings.Contains(fld1, ":") {
				hdr = fld1
			} else {
				hdr = hdr[ixFldDelim+1:]
				ix2ndDelim := strings.IndexAny(hdr, "\t, ")
				if ix2ndDelim > 0 {
					hdr = hdr[:ix2ndDelim]
				}
			}
		}
	}
	components := strings.Split(hdr, ":")
	
	var recID string
	if len(components) == 5 { // presume old-style Illumina header (06Jan15 JBH)
		ixIndex := strings.IndexAny(components[4], "#/")
		if ixIndex > 0 {
			recID = components[1] + ":" + components[2] + ":" + components[3] + ":" + components[4][:ixIndex]
		}
	} else {
		recID = components[3] + ":" + components[4] + ":" + components[5] + ":" + components[6]
	}
	
	if includeMachName && recID != "" {
		recID = components[0] + ":" + recID
	}
	
	if testParser {
		fmt.Println(recID)
	}
	
	return recID
}

func makeOutfilename(filename string, out_suff string) string {
	// 03Aug14 JBH puts out_suff before the file suffix (e.g, .fq or .fastq or .fq.gz)	
	extension := filepath.Ext(filename)
	basename := filename[0:len(filename)-len(extension)]
	if extension == ".gz" { // look back for another extension
		zip_extension := extension
		extension = filepath.Ext(basename)
		if extension != "" {
			basename = basename[0:len(basename)-len(extension)]
		}
		extension += zip_extension
	}
	if out_suff == "" {
		out_suff = "_extract"
	} else if out_suff[0] != '_' {
		out_suff = "_" + out_suff
	}
	return basename + out_suff + extension
}

func getopts() (string, string, string, string, bool, bool) { // return 3 filenames and the output_suffix // 01Jan15 added excerptNonMatches boolean
	idFile := ""; pe1 := ""; pe2 := ""; output_suffix := ""; excerptNonMatches := false; includeMach := false
	
	for ixarg := 1; ixarg < len(os.Args); ixarg++ {
		arg := os.Args[ixarg]
		if arg == "-test" { // test ID parser
			testParser = true
		} else if arg == "-ext" { // set an output suffix other than 'excerpt' (which we will prefix with an underscore)
			ixarg++
			output_suffix = os.Args[ixarg]
		} else if arg == "-v" {	// 01Jan15 JBH added -v argument to allow non-matching records to be output (-v similiar usage in grep)
			excerptNonMatches = true
		} else if arg == "-mach" {
			includeMach = true
		} else if idFile == "" {
			idFile = arg
		} else if pe1 == "" {
			pe1 = arg
		} else if pe2 == "" {
			pe2 = arg
		}
	}
	
	return idFile, pe1, pe2, output_suffix, excerptNonMatches, includeMach
}

func usage() {
        fmt.Println("Usage: excerptByIDs <IDfile>|- <PE_file_1> [<PE_file_2>] [-ext <output_suffix>] [-v] [-mach] [-test]\n\n" +
					"       Outputs records from the FastQ files that match one of the IDs in IDfile.\n" +
					"       (The same ID can be present more than once but is used only once.)\n" +
					"       To read the IDs from stdin you can use the hyphen - as the first parameter.\n\n" +
					"       PE_file_1 records are written to stdout if it is the only file.\n" +
					"       If PE_file_1 and PE_file_2 are present, then new files are written where\n" +
					"       the file name has _extract inserted before the file suffix. Or, you can use\n" +
					"       -ext <output_suffix> to specify a string other than 'extract'.\n" +
					"       PE_file_1 is processed then PE_file_2.\n" +
					"\n       -test option, with just 1 PE_file arg writes IDs to stdout, with 2 PE_file args shows output names.\n" +
					"       -mach option includes the complete machine name in the ID (only needed if PE outputs combined).\n" +
					"       -v option inverts meaning of a matched record. Records NOT in the IDs are output.\n")
}
