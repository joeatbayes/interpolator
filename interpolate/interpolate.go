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
	"path/filepath"
	"regexp"
	s "strings"
	"time"

	"github.com/joeatbayes/goutil/jutil"
	m2h "github.com/joeatbayes/goutil/mdtohtml"
)

type Interpolate struct {
	perf         *jutil.PerfMeasure
	pargs        *jutil.ParsedCommandArgs
	inPath       string
	outPath      string
	glob         string
	processExt   string // file name extension to use when processing directories.
	start        float64
	keepVarNames bool
	baseDir      string
	varPaths     []string
	saveHtml     bool
	loopDelay    float32
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
		fmt.Println("L205: pattern  error: specPath=", specPath, " rePatt=", rePatt, " parmErr=", parmErr, " data=\n", string(data), "\n\n")
	}
	m := lookPat.FindIndex([]byte(data))
	//fmt.Println("L208: specPath=", specPath, " rePatt=", rePatt, " parmErr=", parmErr, " m=", m, " data=\n", string(data), "\n\n")
	if m == nil {
		return ""
	} else {
		_, end := m[0], m[1]
		//fmt.Println("L213: rePatt=", rePatt, " parmErr=", parmErr, " start=", start, " end=", end)
		remaining := data[end:]
		var sb []string
		// accumulate line by line until
		//
		restArr := bytes.Split(remaining, nlByteArr)
		for _, tline := range restArr {
			mrest := MatchAnyTag.FindIndex(tline)
			//fmt.Println("L222: mrest=", mrest, " tline=", string(tline))
			if mrest == nil {
				sb = append(sb, string(tline))
				if len(tline) > 0 {
					sb = append(sb, "\n")
				}
			} else {
				break
			}
		}
		return s.Join(sb, "")
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
		fmt.Printf("L64: matchStr=%s original=%s\n", aMatchStr, origStr)
		if s.HasPrefix(aMatchStr, "inc:") {
			// Handle Include File
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
					if keepVarNames {
						sb = append(sb, varNameIncPath)
						sb = append(sb, " \n")
					}
					sb = append(sb, r.InterpolateStr(string(data)))
				}
			} else {
				sb = append(sb, origStr)
				// file to include not located.
			}

		} else if s.HasPrefix(aMatchStr, "http:") || s.HasPrefix(aMatchStr, "https:") {
			// Process as a URI to read
			// Handle URI Fetch

		} else if r.pargs.Exists(aMatchStr) {
			// substitute match string with parms value
			// or add it back in with the {} protecting it
			// TODO: Add lookup from enviornment variable
			//  if do not find it in the command line parms
			lookVal := r.pargs.Sval(aMatchStr, origStr) // "{*"+aMatchStr+"}")
			//fmt.Printf("L99: matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
			if keepVarNames {
				sb = append(sb, varNameIncPath)
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
			//fmt.Println("L158: fname=", tpath, "data=", string(data))
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
					if keepVarNames {
						sb = append(sb, varNameIncPath)
					}

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
	-`)
}

// Process a single input file
func (u *Interpolate) processFile(inFiName string, outFiName string) {
	inFile, err := os.Open(inFiName)
	if err != nil {
		fmt.Println("L423: error opening input file ", inFiName, " err=", err)
		os.Exit(3)
	}
	defer inFile.Close()

	outFile, sferr := os.Create(outFiName)
	if sferr != nil {
		fmt.Println("L430: Can not open out file ", outFiName, " sferr=", sferr)
		os.Exit(3)
	}

	scanner := bufio.NewScanner(inFile)
	//var b bytes.Buffer
	for scanner.Scan() {
		aline := scanner.Text()
		aline = s.TrimSpace(aline)
		if len(aline) < 1 {
			fmt.Fprintln(outFile, "")
			continue
		} else {
			outStr := u.InterpolateStr(aline)
			//fmt.Println("L444: outStr=", outStr)
			outFile.WriteString(outStr)
			fmt.Fprintln(outFile, "")
		}
	}
	outFile.Sync()
	outFile.Close()
	if u.saveHtml {
		m2h.SaveAsHTML(outFiName, u.loopDelay)
	}
}

// Process a single input file
func (u *Interpolate) processDir(inDirName string, outDirName string) {
	globPath := inDirName + "/*" + u.glob
	fmt.Println("L259: globPath=", globPath)
	files, err := filepath.Glob(globPath)
	if err != nil {
		fmt.Println("L289: ERROR processsing dir=", inDirName, " outDir=", outDirName, "globPath=", globPath, " err=", err)
	} else {
		fmt.Println("L264: files=", files)
		for _, fiPath := range files {

			fmt.Println("L301:  fiPath=", fiPath)
			dir, fname := filepath.Split(fiPath)
			fmt.Println("L271: fiPath=", fiPath, "dir=", dir, "fname=", fname)
			outName := filepath.Join(outDirName, fname)
			u.processFile(fiPath, outName)
		}
	}
}

func main() {
	startms := jutil.Nowms()
	const DefIn = "data"
	const DefOut = "out"
	parms := jutil.ParseCommandLine(os.Args)
	if parms.Exists("help") {
		PrintHelp()
		return
	}
	fmt.Println(parms.String())
	inName := parms.Sval("in", DefIn)
	outName := parms.Sval("out", DefOut)
	u := makeInterpolator()
	//fmt.Println("OutName=", outName)

	u.pargs = parms
	u.glob = parms.Sval("glob", "*.md")
	u.keepVarNames = parms.Bval("keepnames", false)
	u.varPaths = s.Split(parms.Sval("varnames", "desc"), ",")
	u.baseDir = s.TrimSpace(parms.Sval("search", "data/data-dict/"))
	u.saveHtml = parms.Bval("savehtml", false)
	u.loopDelay = parms.Fval("loopdelay", -1)
	jutil.EnsurDir(outName)
	//fmt.Println("u=", u, "saveHtml=", u.saveHtml)

	if jutil.IsDirectory(u.baseDir) == false {
		fmt.Println("L191: FATAL ERROR: baseDir ", u.baseDir, " must be a directory")
		os.Exit(3)
	}

	if u.loopDelay > 0 {
		for {
			u.processDir(inName, outName)
			jutil.Elap("L382: Finished Run", startms, jutil.Nowms())
			time.Sleep(time.Duration(u.loopDelay) * time.Second)
			startms = jutil.Nowms()
		}
	} else {
		u.processDir(inName, outName)
		jutil.Elap("L382: Finished Run", startms, jutil.Nowms())
	}
}
