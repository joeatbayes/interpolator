package main

// interpolate.go

import (
	"bufio"
	"bytes"

	//"bytes"
	//"encoding/json"
	"fmt"
	"io/ioutil"

	//"net/http"
	"os"
	//"path/filepath"
	//m2h "mdtohtml"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	s "strings"
	"time"

	"github.com/joeatbayes/goutil/jutil"
	m2h "github.com/joeatbayes/goutil/mdtohtml"
)

type Interpolate struct {
	perf           *jutil.PerfMeasure
	pargs          *jutil.ParsedCommandArgs
	inPath         string
	inName         string
	outName        string
	outPath        string
	glob           string
	processExt     string // file name extension to use when processing directories.
	start          float64
	keepVarNames   bool
	searchDirs     []string // Renamed from baseDir to searchDirs. Also changed type from string to array of strings
	varPaths       []string
	saveHtml       bool
	loopDelay      float32
	recurseDir     bool
	crossRefFiName string
	crossRefFi     *os.File
	currInFiName   string
	currFiFldUsg   map[string]string
}

func makeInterpolator(parms *jutil.ParsedCommandArgs) *Interpolate {
	const DefIn = "data"
	const DefOut = "out"
	r := Interpolate{}
	r.start = jutil.Nowms()
	r.perf = jutil.MakePerfMeasure(25000)
	r.inName = parms.Sval("in", DefIn)
	r.outName = parms.Sval("out", DefOut)
	r.pargs = parms
	r.currFiFldUsg = make(map[string]string)
	r.glob = parms.Sval("glob", "*.md")
	r.keepVarNames = parms.Bval("keepnames", false)
	r.varPaths = s.Split(parms.Sval("varnames", "desc"), ",")
	r.searchDirs = s.Split(s.TrimSpace(parms.Sval("search", "data/data-dict/")), ",") // Use strings.Split to make the search directories into an array
	r.saveHtml = parms.Bval("savehtml", false)
	r.loopDelay = parms.Fval("loopdelay", -1)
	r.recurseDir = parms.Bval("r", false)
	jutil.EnsurDir(r.outName)
	r.makeCrossRefFi()

	for _, direct := range r.searchDirs { // Changed this to make it a loop to check if each search directory exists
		if jutil.IsDirectory(direct) == false {
			fmt.Println("L73: FATAL ERROR: searchDir ", direct, " must be a directory")
			os.Exit(3)
		}
	}
	return &r
}

func (r *Interpolate) makeCrossRefFi() {
	r.crossRefFiName = filepath.Join(r.outName, "usage_cross_ref.tsv")
	fi, err := os.Create(r.crossRefFiName)
	if err != nil {
		fmt.Println("L84: Error opening cross ref file=", r.crossRefFiName, " err=", err)
		os.Exit(2)
	}
	r.crossRefFi = fi
	r.crossRefFi.WriteString("referenced\treferenced_by\n")
}

func (r *Interpolate) elapSec() float64 {
	return jutil.CalcElapSec(r.start)
}

// Pattern used to find values we are interpolating
var ParmMatch, ParmErr = regexp.Compile(`\{\*.*?\}`)

// Pattern used to find any named tag in yml
var MatchAnyTag, perr2 = regexp.Compile(`^\s*\w+?\:`)

var nlByte = byte('\n')
var nlByteArr = []byte("\n")

// find index of the matching field and then take text until
// we find the next tag indicating starting next element.
// or hit the end of string.
func (r *Interpolate) GetFieldSingle(data []byte, specPath string) string {
	rePatt := `\s*?` + specPath + `\:`
	lookPat, parmErr := regexp.Compile(rePatt)
	if parmErr != nil {
		fmt.Println("L111: pattern  error: specPath=", specPath, " rePatt=", rePatt, " parmErr=", parmErr, " data=\n", string(data), "\n\n")
	}
	m := lookPat.FindIndex([]byte(data))
	//fmt.Println("L114: specPath=", specPath, " rePatt=", rePatt, " parmErr=", parmErr, " m=", m, " data=\n", string(data), "\n\n")
	if m == nil {
		return ""
	} else {
		_, end := m[0], m[1]
		//fmt.Println("L119: rePatt=", rePatt, " parmErr=", parmErr, " end=", end)
		remaining := data[end:]
		var sb []string
		// accumulate line by line until
		//
		restArr := bytes.Split(remaining, nlByteArr)
		for _, tline := range restArr {
			mrest := MatchAnyTag.FindIndex(tline)
			//fmt.Println("127: mrest=", mrest, " tline=", string(tline))
			if mrest == nil {
				sb = append(sb, string(tline))
			} else {
				break
			}
		}
		if len(sb) > 1 {
			return s.Join(sb, "\n")
		} else {
			return s.Join(sb, "")
		}
	}
}

// parse input data as a file and attempt to extract the content
// for a specific field.   Tries to find the requested field first
// then works through options until it finds a match or runs out
func (r *Interpolate) GetField(data []byte, specPath string, defPaths []string) string {
	// When a field is specified it must be matched with
	// no fallback.
	if specPath > " " {
		return r.GetFieldSingle(data, specPath)
	}

	// try the default paths in order specified
	if len(defPaths) > 0 {
		for _, tmpVar := range defPaths {
			tres := r.GetFieldSingle(data, tmpVar)
			if tres > " " {
				return tres
			}
		}
	}
	return " "
}

func (r *Interpolate) InterpolateStr(str string, direct string) (string, bool) {
	// Added an input paramater (string) for the search file path (to replace r.baseDir), and output paramater (bool) to say if we are writing to the output file
	//fmt.Println("166: Interpolate atr=", str)
	if len(str) < 3 || str < " " {
		return str, true
	}
	ms := ParmMatch.FindAllIndex([]byte(str), -1)
	if len(ms) < 1 {
		return str, true // no match found
	}

	writing := false // Created variable to say if we need to write sb to the file
	//sb := strings.Builder
	var sb []string
	last := 0
	slen := len(str)
	for _, m := range ms {
		keepVarNames := r.keepVarNames
		origStr := str[m[0]:m[1]]
		start, end := m[0]+2, m[1]-1
		//fmt.Printf("m[0]=%d m[1]=%d match = %q\n", m[0], m[1], str[start:end])
		if start > last-1 {
			// add the string before the match to the buffer
			sb = append(sb, str[last:start-2])
		}
		aMatchStr := s.ToLower(str[start:end])
		varNameIncPath := " *(" + aMatchStr + ")* "
		fmt.Printf("191: matchStr=%s original=%s\n", aMatchStr, origStr)
		tfa := s.Split(aMatchStr, "#")
		r.currFiFldUsg[tfa[0]] = r.currInFiName
		if s.HasPrefix(aMatchStr, "inc:") {
			newSb, flag := r.handleInclude(aMatchStr, sb, origStr, varNameIncPath, direct) // Took the chunk of code in here and made it a separate function: handleInclude
			writing = flag
			sb = newSb

		} else if s.HasPrefix(aMatchStr, "http:") || s.HasPrefix(aMatchStr, "https:") {
			// Process as a URI to read
			// Handle URI Fetch

		} else if r.pargs.Exists(aMatchStr) {
			// substitute match string with parms value
			// or add it back in with the {} protecting it
			// TODO: Add lookup from enviornment variable
			//  if do not find it in the command line parms
			lookVal := r.pargs.Sval(aMatchStr, origStr) // "{*"+aMatchStr+"}")
			fmt.Printf("L209: matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
			if keepVarNames {
				sb = append(sb, varNameIncPath)
			}
			interp, _ := r.InterpolateStr(lookVal, direct) // Changed this line to accomadate changes to InterpolateStr inputs and outputs
			sb = append(sb, interp)

		} else {
			newSb, flag := r.readParse(aMatchStr, sb, origStr, varNameIncPath, direct) // Took the chunk of code in here and made it a separate function: readParse
			writing = flag
			sb = newSb
		}
		last = end + 1
	}
	if last < slen-1 {
		// append any remaining characters after
		// end of the last match
		sb = append(sb, str[last:slen])
	}
	joinedSb := s.Join(sb, "")
	return s.ReplaceAll(joinedSb, " \\n", "<BR>"), writing
}

func (r *Interpolate) handleInclude(aMatchStr string, sb []string, origStr string, varNameIncPath string, direct string) ([]string, bool) {
	// Handle Include File
	// Process simple file include
	//fmt.Println("L234:found inc: prefix")
	matchPort := s.TrimSpace(aMatchStr[4:])
	matchPort = r.mergePaths(direct, matchPort) // Added function to merge two file paths for specificity if desired
	tpath := filepath.Join(direct, matchPort)
	flag := false // Created flag to say if we need to write sb to the file
	//fmt.Println("L238: matchPort=", matchPort, " tpath=", tpath)
	if jutil.Exists(tpath) {
		data, err := ioutil.ReadFile(tpath)
		if err != nil {
			//fmt.Println("Error reading ", tpath, " err=", err)
			// could not read file so copy original path
			// into output file
			sb = append(sb, origStr)
		} else {
			// save file read into our output buffer
			//fmt.Println("L248: Add file from buffer data=\n", string(data), "\n\n")
			flag = true // If successfully reads the file at the specified path, we want to write to the output file
			if r.keepVarNames {
				sb = append(sb, varNameIncPath)
				sb = append(sb, " \n")
			}
			interp, _ := r.InterpolateStr(string(data), direct) // Changed this line to accomadate changes to InterpolateStr inputs and outputs
			sb = append(sb, interp)
		}
	} else {
		sb = append(sb, origStr)
		// file to include not located.
	}
	return sb, flag
}

func (r *Interpolate) readParse(aMatchStr string, sb []string, origStr string, varNameIncPath string, direct string) ([]string, bool) {
	// Try read file and parse out the requested variable name
	varSpecPath := ""
	pathFrag := s.SplitN(aMatchStr, "#", 2)
	basicPath := s.TrimSpace(pathFrag[0])
	if len(pathFrag) == 2 {
		varSpecPath = pathFrag[1]
	}
	basicPath = r.mergePaths(direct, basicPath) // Added function to merge two file paths
	tpath := filepath.Join(direct, basicPath) + ".yml"
	tpath = s.Replace(tpath, ".yml.yml", ".yml", 1)
	data, err := ioutil.ReadFile(tpath)
	flag := false // Flag to say if we are going to write the string to the output file
	// fmt.Println("L313: fname=", tpath, "data=", string(data))
	if err != nil {
		fmt.Println("L315: Error reading ", tpath, " err=", err)
		// could not read file so copy original path
		// into output file
		sb = append(sb, origStr)
		sb = append(sb, " *FILE NOT FOUND "+tpath+" *")
	} else {
		flag = true // If the file is found, set the writing flag to true
		// save file read into our output buffer
		extractStr := r.GetField(data, varSpecPath, r.varPaths)
		if extractStr > " " {
			if r.keepVarNames {
				sb = append(sb, varNameIncPath)
			}
			interp, _ := r.InterpolateStr(extractStr, direct) // Changed this line to accomadate changes to InterpolateStr inputs and outputs
			sb = append(sb, interp)
		} else {
			// Could not find a match so output default
			sb = append(sb, origStr)
			sb = append(sb, " *VARIABLE NOT FOUND* ")
		}
	}
	//sb = append(sb, aMatchStr)
	return sb, flag
}

func (r *Interpolate) mergePaths(direct string, basicPath string) string {
	frontSplit := s.Split(direct, "\\")  // Splits the search directory into an array of folders
	backSplit := s.Split(basicPath, "/") // Splits the path to the .yml file into an array of folders
	specific := false                    // Flag to say if the search directory contains part of the basicPath
	i := len(frontSplit) - 1             // iterator int for frontSplit
	j := 0                               // iterator int for backSplit
	for i >= 0 {
		if frontSplit[i] == backSplit[j] {
			specific = true // If the least specific folder in backSplit is found in frontSplit, we want to see if the file path continues to be the same
			break
		}
		i--
	}
	if specific {
		i++
		j++
		for i < len(frontSplit) { // Continue seeing if the next specific folders are the same in both arrays
			if j >= len(backSplit) {
				specific = false // If the search directory is more specific than even the yml file location, then we won't change anything
				break
			}
			if frontSplit[i] != backSplit[j] {
				specific = false // If the yml file location is in a different folder than the search directory, don't change anything
				break
			}
			i++
			j++
		}
	}
	if specific { // This will be true if the end of frontSplit is the same as the beginning of backSplit
		lenstr := 0
		for k := 0; k < j; k++ {
			lenstr = lenstr + len(backSplit[k]) // Add up the length of all the strings at the beginning of backSplit that are the same as the end of frontSplit
			lenstr++                            // Add one for each slash
		}
		basicPath = basicPath[lenstr:] // Cut off the beginning of basicPath (yml location) so it can be found by joining it with search directory
	}
	return basicPath
}

// Process a single input file
func (u *Interpolate) processFile(inFiName string, outFiName string) {
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("L344: error opening input file ", inFiName, " err=", err)
		os.Exit(3)
	}
	defer inFile.Close()
	u.currInFiName = inFiName
	outFile, sferr := os.Create(outFiName)
	if sferr != nil {
		fmt.Println("L351: Can not open out file ", outFiName, " sferr=", sferr)
		os.Exit(3)
	}

	scanner := bufio.NewScanner(inFile)
	//var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Text()
		if len(aline) < 1 {
			fmt.Fprintln(outFile, "")
			continue
		} else {
			for searchDex, direct := range u.searchDirs { // Added this for loop to go through each search directory every time we reach a new line
				outStr, writing := u.InterpolateStr(aline, direct) // Changed this line to accomadate changes to InterpolateStr inputs and outputs
				//fmt.Println("L365: outStr=", outStr)
				if writing || searchDex == len(u.searchDirs)-1 {
					outFile.WriteString(outStr) // Changed to only write to the output file if we find the file (writing flag is set to true), or we are on the last search directory
					fmt.Fprintln(outFile, "")
					break
				}
			}
		}
	}
	outFile.Sync()
	outFile.Close()
	if u.saveHtml {
		m2h.SaveAsHTML(outFiName, u.loopDelay)
	}

	// save references used in this file and setup for next file
	for refUsed, useIn := range u.currFiFldUsg {
		refUsed = s.TrimSpace(s.Replace(refUsed, "inc:", "", 1))
		u.crossRefFi.WriteString(refUsed + "\t" + useIn + "\n")
	}
	u.currFiFldUsg = make(map[string]string)
}

// Process a single input file
func (u *Interpolate) processDir(inDirName string, outDirName string) {
	globPath := inDirName + "/*" + u.glob
	fmt.Println("L391: globPath=", globPath)
	files, err := filepath.Glob(globPath)
	if err != nil {
		fmt.Println("L394: ERROR processsing dir=", inDirName, " outDir=", outDirName, "globPath=", globPath, " err=", err)
	} else {
		fmt.Println("L396: files=", files)
		for _, fiPath := range files {

			fmt.Println("L399:  fiPath=", fiPath)
			dir, fname := filepath.Split(fiPath)
			fmt.Println("L401: fiPath=", fiPath, "dir=", dir, "fname=", fname)
			outName := filepath.Join(outDirName, fname)
			u.processFile(fiPath, outName)
		}
	}
}

func (u *Interpolate) processDirRecursive(inDir string, outDir string) {
	inDirPref := path.Clean(inDir)
	fmt.Println("L410: Recursive Directory Walk inDir=", inDir, "inDirClean=", inDirPref)
	err := filepath.Walk(inDirPref,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() == false {
				// We are going to use our glob path semantics to process directories
				// so we will skip
				return nil
			} else {
				// It is a directory so compute it's position relative to
				// our starting directory so we can compute our correct
				// output directory.
				//fmt.Println("L425: path=", path, "name=", info.Name(), " info=", info, " size=", info.Size())
				pathRelToIn, err := filepath.Rel(inDirPref, path)
				if err != nil {
					fmt.Println("L428: Error finding relative path inDirPref=", inDirPref, " path=", path, " err=", err)
				} else {
					relOutDir := filepath.Join(outDir, pathRelToIn)
					jutil.EnsurDir(relOutDir)
					//fmt.Println("L432: pathRelToIn=", pathRelToIn, " relOutDir=", relOutDir)
					u.processDir(path, relOutDir)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println("L439: processDirRecursive dir=", inDir, " err:", err)
	}
}

func PrintHelp() {
	fmt.Println(
		`interpolate  -in=data -out=out glob=*sample*.md -search=./data/data-dict  -VarNames=desc,tech_desc  -keepNames=true -maxRec=99

  -in = path to input directory containing files to process. 
        Defaults to ./data
  -out = path to output directory where expanded files will be
         written.   defaults to ./out
  -glob= glob pattern to use when selecting files to process in 
         input directory.  Defaults to *.md
  -search = Directory Base to search for files named in
         interpolated parameters.
  -varNames = Default variable name  matched in dictionary files.
         Can be overridden if variable name is specified using 
		 #varname semantic.    May be common separated list to 
		 allow default lookup of backup fields such as look first 
		 in tech_desc then in desc. -varNames=desc,tech_desc - 
		 Causes the system to search first in the desc: field 
		 then in the tech_desc.   This would use the business 
		 description to be used first and then filled in tech 
		 desc if desc is not found.   Just reverse the order to 
		 cause it to use the technical description first. It 
		 will use the first one found.   When the varname is 
		 specified using the # semantic it will use the 
		 specified var name and ignore the default varNames.
  -keepNames = when set to true it will keep the supplied path as 
         part of output text.   When not set or false will 
		 replace content of path with content.
  -saveHtml=yes when set to yes will convert the md file to Html
         and save it in the output directory. 
  -maxRec this is a variable defined on command line that is 
         being interpolated.  Resolution of variables defined 
		 on command line take precedence over those  resolved 
		 in files. 
  -loopDelay - When set the system will process the input.  
         Sleep for a number of seconds and then re-process.
		 This is intended to keep a generated file available to 
		 easily reload.  eg:  -loopDelay=30 will cause the system to 
		 reprocess the input files once every 30 seconds.
  -r=yes - When set to yes the system will recursively walk
         all directories contained in the -in directory and process
		 every file that matches the glob patttern.  If not set will
		 only walk the named directory.  When recurseDir is set to 
		 yes then it will create directories in the output directory
		 that mirror the input directory path whenever a matching input
		 file is found. 
  -makeCrossRef=yes - when set to yes the system will generate a .md file
         and a sorted file containing contents of files referencing a 
		 specific data dictionary item.  Defaults to false.
	-`)
}

/* Load a text file into slice of strings
Return the first line as header and an slice
containing the strings.  Supresses empty lines
And lines beginning with # */
func LoadFileWithHeader(inFiName string) (string, []string) {
	fmt.Println("L500: LoadFileWithHeader inFiName=", inFiName)
	start := jutil.Nowms()
	tarr := make([]string, 0, 1000)
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("L505: ERROR: loadFi ", inFiName, " err=", err)
		return "", nil
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	lineCnt := 0
	scanner.Scan()
	headers := scanner.Text()
	for scanner.Scan() {
		aline := scanner.Text()
		lineCnt++
		if len(aline) <= 0 {
			continue
		}
		aline = s.TrimSpace(aline)
		if s.HasPrefix(aline, "#") {
			continue
		}
		tarr = append(tarr, aline)
		fmt.Println("L524 aline=", aline)
	}
	jutil.Elap("L526: loadFi "+inFiName, start, jutil.Nowms())
	return headers, tarr
}

/* Open a File containing a single line of header
   and and a set of strings.  Sort the strings and
   write a new file containing the same header line
   and all the strings */
//func sortFileWithHeader(inFiName string, outFiName string) (string, string[]) {

func sortFileWithHeader(inFiName string, outFiName string) (string, []string) {

	headers, tarr := LoadFileWithHeader(inFiName)
	fmt.Println("tarr=", tarr)
	sort.Strings(tarr)
	fmt.Println("tarr=", tarr)
	fi, err := os.Create(outFiName)
	if err != nil {
		fmt.Println("L544: Error opening output sort file=", outFiName, " err=", err)
		os.Exit(2)
	}
	defer fi.Close()
	fmt.Println("L548: sortFileWithHeader inFiName=", inFiName, " outFiName=", outFiName, " #Rec=", len(tarr))
	fi.WriteString(headers)
	fi.WriteString("\n")
	for _, aline := range tarr {
		fmt.Println("L552: aline=", aline)
		fi.WriteString(aline)
		fi.WriteString("\n")
	}
	fi.Sync()

	return headers, tarr
}

func PadRightFixed(tstr string, targLen int, padChar string) string {
	numPad := targLen - len(tstr)
	if numPad <= 0 {
		return tstr // tstr[0:targLen]
	} else {
		var b s.Builder
		b.WriteString(tstr)
		for i := 0; i < numPad; i++ {
			b.WriteString(padChar)
		}
		return b.String()
	}
} // func

const tableColLen = 50

func makeCrossRef(inFiName string, outFiName string) {
	srtFiName := s.Replace(inFiName, ".tsv", ".srt.tsv", 1)
	headers, tarr := sortFileWithHeader(inFiName, srtFiName)
	outFi, err := os.Create(outFiName)
	if err != nil {
		fmt.Println("L582: Error opening makeCrossRef file=", outFiName, " err=", err)
		os.Exit(2)
	}
	defer outFi.Close()
	fmt.Println("L586:  infiName=", inFiName, " outFiName=", outFiName, "srtFiName=", srtFiName, "#Rec=", len(tarr))

	currFld := ""
	headSegArr := s.SplitN(headers, "\t", 2)

	seg1Pad := PadRightFixed(headSegArr[0], tableColLen, " ")
	seg2Pad := PadRightFixed(headSegArr[1], tableColLen, " ")
	outFi.WriteString("|" + seg1Pad + "|" + seg2Pad + "|\n")

	dashes := PadRightFixed("-", tableColLen, "-")
	spaces := PadRightFixed(" ", tableColLen, " ")
	outFi.WriteString("|" + dashes + " | " + dashes + " |\n")

	for _, aline := range tarr {

		fldArr := s.SplitN(aline, "\t", 2)
		fmt.Println("aline=", aline, " fldArr=", fldArr)
		fldRef := s.TrimSpace(fldArr[0])
		refBy := s.TrimSpace(fldArr[1])
		if fldRef > " " && refBy > " " {
			fldRefPad := PadRightFixed(fldRef, tableColLen, " ")
			refByPad := PadRightFixed(fldRef, tableColLen, " ")
			if fldRef != currFld {
				outFi.WriteString("|" + fldRefPad + "|" + refByPad + "|\n")
				currFld = fldRef
			} else {
				outFi.WriteString("|" + spaces + " |" + refBy + "|\n")
			}
		}
	}
	outFi.WriteString("\n")
	outFi.Sync()
}

func main() {
	startms := jutil.Nowms()

	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}
	fmt.Println(parms.String())
	u := makeInterpolator(parms)
	//fmt.Println("OutName=", outName)
	for {
		u.makeCrossRefFi()
		if u.recurseDir == false {
			// process single file path
			u.processDir(u.inName, u.outName)
		} else {
			u.processDirRecursive(u.inName, u.outName)
		}
		jutil.Elap("L382: Finished Run", startms, jutil.Nowms())

		if parms.Bval("makecrossref", false) {
			u.crossRefFi.Close()
			crossMDOut := s.Replace(u.crossRefFiName, ".tsv", ".md", 1)
			fmt.Println("crossMDOut=", crossMDOut)
			makeCrossRef(u.crossRefFiName, crossMDOut)
		}
		if u.loopDelay <= 0 {
			break
		} else {
			time.Sleep(time.Duration(u.loopDelay) * time.Second)
			startms = jutil.Nowms()
			u.makeCrossRefFi()
		}

	}
}
