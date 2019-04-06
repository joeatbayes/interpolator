package main

// interpolate.go

import (
	"bufio"
	//"bytes"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	//"path/filepath"
	"regexp"
	s "strings"
	//"time"
	"path/filepath"

	"github.com/joeatbayes/goutil/jutil"
)

type Interpolate struct {
	perf       *jutil.PerfMeasure
	pargs      *jutil.ParsedCommandArgs
	inPath     string
	processExt string // file name extension to use when processing directories.
	start      float64
	keepPaths  bool
	baseDir    string
	varPaths   []string
}

func makeInterpolator() *Interpolate {
	r := Interpolate{}
	r.perf = jutil.MakePerfMeasure(25000)
	r.start = jutil.Nowms()
	return &r
}

func (r *Interpolate) elapSec() float64 {
	return jutil.CalcElapSec(r.start)
}

// Pattern used to find values we are interpolating
var ParmMatch, ParmErr = regexp.Compile(`\{\*.*?\}`)

// Pattern used to find any named tag in yml
var MatchAnyTag, perr2 = regexp.Compile(`\v+\s*\w+?\:`)

// find index of the matching field and then take text until
// we find the next tag indicating starting next element.
// or hit the end of string.
func (r *Interpolate) GetFieldSingle(data []byte, specPath string) string {
	rePatt := `\s*?` + specPath + `\:`
	lookPat, parmErr := regexp.Compile(rePatt)
	m := lookPat.FindIndex([]byte(data))
	fmt.Println("L57: specPath=", specPath, " rePatt=", rePatt, " parmErr=", parmErr, " m=", m, " data=\n", string(data), "\n\n")
	if m == nil {
		return ""
	} else {
		start, end := m[0], m[1]
		fmt.Println("L62: rePatt=", rePatt, " parmErr=", parmErr, " start=", start, " end=", end)
		remaining := data[end:]
		mrest := MatchAnyTag.FindIndex(remaining)
		//fmt.Println("L65: mrest=", mrest, " remaining=\n", string(remaining), "\n\n\n")
		if mrest == nil {
			return string(remaining)
		} else {
			restStart := mrest[0]
			varMatch := string(remaining[0:restStart])
			//fmt.Println("L71: varMatch=", varMatch)
			return varMatch
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

func (r *Interpolate) InterpolateStr(str string) string {
	//fmt.Println("L246: Interpolate atr=", str)
	if len(str) < 3 || str < " " {
		return str
	}
	ms := ParmMatch.FindAllIndex([]byte(str), -1)
	if len(ms) < 1 {
		return str // no match found
	}
	//sb := strings.Builder
	var sb []string
	last := 0
	slen := len(str)
	for _, m := range ms {
		origStr := str[m[0]:m[1]]
		start, end := m[0]+2, m[1]-1
		//fmt.Printf("m[0]=%d m[1]=%d match = %q\n", m[0], m[1], str[start:end])
		if start > last-1 {
			// add the string before the match to the buffer
			sb = append(sb, str[last:start-2])
		}
		aMatchStr := s.ToLower(str[start:end])
		fmt.Printf("L64: matchStr=%s original=%s\n", aMatchStr, origStr)
		if s.HasPrefix(aMatchStr, "inc:") {
			// Process simple file include
			//fmt.Println("L69:found inc: prefix")
			matchPort := s.TrimSpace(aMatchStr[4:])
			tpath := filepath.Join(r.baseDir, matchPort)
			//fmt.Println("L71: matchPort=", matchPort, " tpath=", tpath)
			if jutil.Exists(tpath) {
				data, err := ioutil.ReadFile(tpath)
				if err != nil {
					//fmt.Println("Error reading ", tpath, " err=", err)
					// could not read file so copy original path
					// into output file
					sb = append(sb, origStr)
				} else {
					// save file read into our output buffer
					//fmt.Println("L137 Add file from buffer data=\n", string(data), "\n\n")
					sb = append(sb, r.InterpolateStr(string(data)))
				}
			} else {
				sb = append(sb, origStr)
				// file to include not located.
			}

		} else if s.HasPrefix(aMatchStr, "http:") || s.HasPrefix(aMatchStr, "https:") {
			// Process as a URI to read

		} else if r.pargs.Exists(aMatchStr) {
			// substitute match string with parms value
			// or add it back in with the {} protecting it
			// TODO: Add lookup from enviornment variable
			//  if do not find it in the command line parms
			lookVal := r.pargs.Sval(aMatchStr, origStr) // "{*"+aMatchStr+"}")
			//fmt.Printf("L99: matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
			if r.keepPaths {
				sb = append(sb, "*"+aMatchStr+"* ")
			}
			sb = append(sb, r.InterpolateStr(lookVal))

		} else {
			// Try read file and parse out the requested variable name
			varSpecPath := ""
			pathFrag := s.SplitN(aMatchStr, "#", 2)
			basicPath := s.TrimSpace(pathFrag[0])
			if len(pathFrag) == 2 {
				varSpecPath = pathFrag[1]
			}
			tpath := filepath.Join(r.baseDir, basicPath) + ".yml"
			tpath = s.Replace(tpath, ".yml.yml", ".yml", 1)
			data, err := ioutil.ReadFile(tpath)
			fmt.Println("L158: fname=", tpath, "data=", string(data))
			if err != nil {
				fmt.Println("L160: Error reading ", tpath, " err=", err)
				// could not read file so copy original path
				// into output file
				sb = append(sb, origStr)
				sb = append(sb, " *FILE NOT FOUND "+tpath+" *")
			} else {
				// save file read into our output buffer
				extractStr := r.GetField(data, varSpecPath, r.varPaths)
				if extractStr > " " {
					sb = append(sb, r.InterpolateStr(extractStr))
				} else {
					// Could not find a match so output default
					sb = append(sb, origStr)
					sb = append(sb, " *VARIABLE NOT FOUND* ")
				}
			}
			//sb = append(sb, aMatchStr)
		}
		last = end + 1
	}
	if last < slen-1 {
		// append any remaining characters after
		// end of the last match
		sb = append(sb, str[last:slen])
	}
	return s.Join(sb, "")
}

func PrintHelp() {
	fmt.Println(
		` interpolate  -in=data/sample-api.md -out=out/sample-api.md -baseDir=./data/dict  -defaulVarPath=desc  -keepPaths=true

  -in = path of a input file to process.  May be a specific file or a glob pattern.

  -out = Location to write the output file once expanded
  -baseDir = Directory Base to search for files named in interpolated parameters.

  -defaultVarPath = string matched in predefined file to pull next string. 

  -keepPaths = when set to true it will keep the supplied path as part of output text.   When not set or false will replace content of path with content.

	-`)
}

// Process a single input file
func (u *Interpolate) processFile(inFiName string, outFiName string) {
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("error opening input file ", inFiName, " err=", err)
		os.Exit(3)
	}
	defer inFile.Close()

	outFile, sferr := os.Create(outFiName)
	if sferr != nil {
		fmt.Println("Can not open out file ", outFiName, " sferr=", sferr)
		os.Exit(3)
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	//var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Text()
		aline = s.TrimSpace(aline)
		if len(aline) < 1 {
			continue
		} else {
			outStr := u.InterpolateStr(aline)
			fmt.Println("L160: outStr=", outStr)
			fmt.Fprintln(outFile, outStr)
		}
	}
	outFile.Sync()
}

func main() {
	startms := jutil.Nowms()
	const DefInFiName = "data/sample.tst"
	const DefOutFiName = "out/sample.txt"
	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}
	fmt.Println(parms.String())
	inFiName := parms.Sval("in", DefInFiName)
	outFiName := parms.Sval("out", DefOutFiName)
	u := makeInterpolator()
	fmt.Println("OutFileName=", outFiName)

	u.pargs = parms
	u.keepPaths = parms.Bval("keeppaths", false)
	u.varPaths = s.Split(parms.Sval("varpaths", "desc"), ",")
	u.baseDir = s.TrimSpace(parms.Sval("basedir", "data/data-dict/"))
	if jutil.IsDirectory(u.baseDir) == false {
		fmt.Println("L191: FATAL ERROR: baseDir ", u.baseDir, " must be a directory")
		os.Exit(3)
	}

	u.processFile(inFiName, outFiName)
	jutil.Elap("L382: Finished Run", startms, jutil.Nowms())
}
