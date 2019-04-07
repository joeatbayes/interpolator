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
	"github.com/shurcooL/github_flavored_markdown"
	//"gopkg.in/russross/blackfriday.v2"
)

var htmlPrefix = `
  <html>
    <head>
	  <meta charset="utf-8">
	<style>
	
body {
    background-color: #FFF;
    color: #172B4D;
    font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen,Ubuntu,"Fira Sans","Droid Sans","Helvetica Neue",sans-serif;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 1.42857142857143;
    -ms-overflow-style: -ms-autohiding-scrollbar;
    text-decoration-skip: ink;
}
	 p {
    display: block;
	magin-left: 1em;
    margin-block-start: 1em;
    margin-block-end: 1em;
    margin-inline-start: 0px;
    margin-inline-end: 0px;
}

h1 {
    font-size: 1.74285714em;
    font-style: inherit;
    font-weight: 600;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h2 {
    font-size: 1.52857143em;
    font-style: inherit;
    font-weight: 500;
    letter-spacing: -.008em;
    line-height: 1.2;
    margin-top: 28px;
    color: #172B4D;
}

 h3 {
    font-size: 1.34285714em;
    font-style: inherit;
    font-weight: 600;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h4 {
    font-size: 1.1;
    font-style: inherit;
    font-weight: 500;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}

h5 {
    font-size: 1.0;
    font-style: inherit;
    font-weight: 400;
    line-height: 1.25;
    letter-spacing: -.006em;
    margin-top: 24px;
    color: #172B4D;
}


pre {
    background: #F4F5F7;
    border: 1px solid #DFE1E6;
    border-radius: 3px;
    overflow-x: auto;
    padding: 5px 10px;
    word-wrap: normal;
	font-family: monospace; 
	font-size: 12px;
	font-weight: 300;
	
}


blockquote  {
    margin: 12px 12px 12px 12px;
	color: Green;
	font-weight: 700;
}

a:link {
  text-decoration: none;
  color: blue;
  font-weight: 500;
}

a:visited {
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
  font-weight: 550;
}

a:active {
  text-decoration: underline;
}

	</style>
    </head>
	<body>
	  <article class="markdown-body entry-content" style="padding: 30px;">
  `
var htmlSuffix = `
    </article>
	</body>
	</html>
  `

type Interpolate struct {
	perf       *jutil.PerfMeasure
	pargs      *jutil.ParsedCommandArgs
	inPath     string
	outPath    string
	glob       string
	processExt string // file name extension to use when processing directories.
	start      float64
	keepPaths  bool
	baseDir    string
	varPaths   []string
	saveHtml   bool
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

func (u *Interpolate) saveAsHTML(srcFile string) {
	// Save File as HTML if the HTML save directory
	// has been specified.
	fext := filepath.Ext(srcFile)
	fmt.Println("L226: fext=", fext, "u.saveHtml=", u.saveHtml)
	if u.saveHtml && fext == ".md" {
		data, err := ioutil.ReadFile(srcFile)
		//fmt.Println("L229: fname=", srcFile, "data=", string(data))
		if err != nil {
			fmt.Println("L260: Error reading ", srcFile, " err=", err)
		} else {
			hname := s.Replace(srcFile, ".md", ".html", 1)
			fmt.Println("L234: html filename=", hname)
			f, err := os.Create(hname)
			if err != nil {
				fmt.Println("error writing to ", hname, " err=", err)
			} else {
				f.WriteString(htmlPrefix)
				f.Write(github_flavored_markdown.Markdown(data))
				f.WriteString(htmlSuffix)
				f.Close()
			}
		}
	}
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
			//fmt.Println("L160: outStr=", outStr)
			outFile.WriteString(outStr)
			fmt.Fprintln(outFile, "")
		}
	}
	outFile.Sync()
	outFile.Close()
	u.saveAsHTML(outFiName)
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
	fmt.Println("OutName=", outName)

	u.pargs = parms
	u.glob = parms.Sval("glob", "*.md")
	u.keepPaths = parms.Bval("keepnames", false)
	u.varPaths = s.Split(parms.Sval("varnames", "desc"), ",")
	u.baseDir = s.TrimSpace(parms.Sval("search", "data/data-dict/"))
	u.saveHtml = parms.Bval("savehtml", false)
	jutil.EnsurDir(outName)
	fmt.Println("u=", u, "saveHtml=", u.saveHtml)

	if jutil.IsDirectory(u.baseDir) == false {
		fmt.Println("L191: FATAL ERROR: baseDir ", u.baseDir, " must be a directory")
		os.Exit(3)
	}

	u.processDir(inName, outName)
	jutil.Elap("L382: Finished Run", startms, jutil.Nowms())
}
