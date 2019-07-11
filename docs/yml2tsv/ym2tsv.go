package main

// Yml2tsv.go - Harvest fields from YML files and convert to
// into a TSV File. Any \n\t will be converted to space.
// List of fields may be supplied on command line.  Glob pattern on command
// line will be used to filter input data set.

import (
	"bufio"
	//"bytes"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	//"path/filepath"
	//m2h "mdtohtml"
	"path"
	"path/filepath"
	"regexp"
	s "strings"
	//"time"

	"github.com/joeatbayes/goutil/jutil"
)

type ctl struct {
	perf     *jutil.PerfMeasure
	pargs    *jutil.ParsedCommandArgs
	inDir    string
	outPath  string
	glob     string
	start    float64
	vars     []string
	varsLook map[string]int
	varsCnt  int
}

func makeCtl(parms *jutil.ParsedCommandArgs) *ctl {
	const DefIn = "data/data-dict/db"
	const DefOut = "out/data-dict.tsv"
	const DefVars = "type,len,desc"
	r := ctl{}
	r.pargs = parms
	r.perf = jutil.MakePerfMeasure(25000)
	r.start = jutil.Nowms()

	fmt.Println(parms.String())
	r.inDir = parms.Sval("in", DefIn)
	r.outPath = parms.Sval("out", DefOut)
	r.glob = parms.Sval("glob", "*.yml")
	r.vars = s.Split(parms.Sval("vars", DefVars), ",")
	r.varsLook = make(map[string]int)
	r.varsCnt = len(r.vars)
	for ndx, vname := range r.vars {
		r.varsLook[vname] = ndx
	}
	// TODO:Convert VARS into fast lookup Set

	if jutil.IsDirectory(r.inDir) == false {
		fmt.Println("L191: FATAL ERROR: InDir ", r.inDir, " must be a directory")
		os.Exit(3)
	}
	return &r
}

// Pattern used to find values we are interpolating
var ParmMatch, ParmErr = regexp.Compile(`\{\*.*?\}`)

// Pattern used to find any named tag in yml
var MatchAnyTag, perr2 = regexp.Compile(`^\s*\w+?\:`)

var nlByte = byte('\n')
var nlByteArr = []byte("\n")

func CountLeadSpaceB(astr []byte) int {
	for ndx, tchar := range astr {
		if tchar != ' ' {
			return ndx
		}
	}
	return len(astr)
}

var emptyStr = ""

func ClearStrArr(arr []string) {
	numEle := len(arr)
	for i := 0; i < numEle; i++ {
		arr[i] = emptyStr
	}
}

// Process a single input file
// scan for starting yml tokens.  Once it is found then it becomes the field name
// then scan for next token start respecting space indents.
// accumulating text as you go.   When next field or end of file is found then
// then finish accumulating and save.
//
// Take the list of fields and use them in order specified for the list of
// of fields.
//
// Must save as a dict of fields from pass 1 to build a total list of fields
//
// a list of records until we are done because we
func (u *ctl) processFile(inFiName string, outFi *os.File) {
	pathRelToIn, err := filepath.Rel(u.inDir, inFiName) // Need this to produce nicer looking fiPath with dir Prefix removed
	fmt.Println("L109: pathRelToIn=", pathRelToIn)
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("L423: error opening input file ", inFiName, " err=", err)
		os.Exit(3)
	}
	defer inFile.Close()
	var sb []string
	currRec := make([]string, u.varsCnt)
	var currFldName = ""

	flushField := func() {
		// New field defined so need to save to keeper Dict
		// if field is in the specified output set. output file.
		ndx, fnd := u.varsLook[currFldName]
		if fnd {
			outVal := s.Join(sb, " ")
			outVal = s.Replace(outVal, "\t", " ", -1)
			outVal = s.TrimSpace(outVal)
			currRec[ndx] = outVal
		}
		//fmt.Println("L129: flushField ndx=", ndx, " currFldName=", currFldName, " currRec=", currRec)
		// zero out for next field
		sb = sb[:0]
		currFldName = ""
	}

	flushRec := func() {
		// new Rec reached so need to save to output file
		//outFile.WriteString(outStr)
		outFi.WriteString(pathRelToIn)
		outFi.WriteString("\t")
		outStr := s.Join(currRec, "\t")
		outFi.WriteString(outStr)
		outFi.WriteString("\n")
		// setup for next rec
		ClearStrArr(currRec)
	}

	scanner := bufio.NewScanner(inFile)
	//var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Bytes()
		//aline = s.TrimSpace(aline)
		if len(aline) < 1 {
			continue
		} else {
			m := MatchAnyTag.FindIndex(aline)

			//fmt.Println("L222: m=", m, " aline=", string(aline))
			if m == nil {
				sb = append(sb, string(aline))
			} else {
				if len(sb) > 0 {
					flushField()
				}

				// Setup Contents for next Record
				start, end := m[0], m[1]
				fldName := aline[start : end-1]
				// leadSpace := CountLeadSpaceB(fldName) // Use this to build up nested yml semantics.
				currFldName = s.TrimSpace(string(fldName))
				rest := aline[end:]
				sb = append(sb, string(rest))
			}

		}
	}
	// Flush any leftover
	flushField()
	flushRec()

}

// Process a single input file
func (u *ctl) processDir(inDirName string, outFi *os.File) {
	globPath := inDirName + "/*" + u.glob
	fmt.Println("L259: globPath=", globPath)
	files, err := filepath.Glob(globPath)
	if err != nil {
		fmt.Println("L289: ERROR processsing dir=", inDirName, "globPath=", globPath, " err=", err)
	} else {
		//fmt.Println("L264: files=", files)
		for _, fiPath := range files {
			fmt.Println("L333:  fiPath=", fiPath)
			u.processFile(fiPath, outFi)
		}
	}
}

func (u *ctl) processDirRecursive(inDir string, outFi *os.File) {
	inDirPref := path.Clean(inDir)
	fmt.Println("L344: Recursive Directory Walk inDir=", inDir, "inDirClean=", inDirPref)
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

				if err != nil {
					fmt.Println("L362: Error finding relative path inDirPref=", inDirPref, " path=", path, " err=", err)
				} else {
					//fmt.Println("L359: pathRelToIn=", pathRelToIn, " relOutDir=", relOutDir)
					u.processDir(path, outFi)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println("L365: processDirRecursive dir=", inDir, " err:", err)
	}
}

func printHelp() {
	fmt.Println(
		`yml2csv  -in=data/data-dict/db -out=out glob=*.yml  -vars=type,len,desc

  -in = path to input directory containing files to process. 
        Defaults to ./data
  -out = filename containg generated TSV. 
  -glob= glob pattern to use when selecting files to process in 
         input directory.  Defaults to *.md
  -vars= List of variables to retain from the yml files. 
  
	-`)
}

func main() {
	startms := jutil.Nowms()

	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		printHelp()
		return
	}
	u := makeCtl(parms)
	outFi, sferr := os.Create(u.outPath)
	if sferr != nil {
		fmt.Println("L430: Can not open out file ", u.outPath, " sferr=", sferr)
		os.Exit(3)
	}
	defer outFi.Close()
	headers := "path\t" + s.Join(u.vars, "\t") + "\n"
	outFi.WriteString(headers)

	u.processDirRecursive(u.inDir, outFi)
	jutil.Elap("L382: Finished Run", startms, jutil.Nowms())

}
