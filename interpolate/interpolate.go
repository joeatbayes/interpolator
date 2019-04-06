package main

// interpolate.go

import (
	"bufio"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	//"path/filepath"
	"regexp"
	s "strings"
	//"time"

	"github.com/joeatbayes/goutil/jutil"
)

type Interpolate struct {
	perf       *jutil.PerfMeasure
	pargs      *jutil.ParsedCommandArgs
	inPath     string
	processExt string // file name extension to use when processing directories.
	start      float64
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

var ParmMatch, ParmErr = regexp.Compile("\\{.*?\\}")

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
		start, end := m[0]+1, m[1]-1
		//fmt.Printf("m[0]=%d m[1]=%d match = %q\n", m[0], m[1], str[start:end])
		if start > last-1 {
			// add the string before the match to the buffer
			sb = append(sb, str[last:start-1])
		}
		aMatchStr := s.ToLower(str[start:end])
		// substitute match string with parms value
		// or add it back in with the {} protecting it
		// TODO: Add lookup from enviornment variable
		//  if do not find it in the command line parms
		lookVal := r.pargs.Sval(s.ToLower(aMatchStr), "{"+aMatchStr+"}")
		//fmt.Printf("matchStr=%s  lookVal=%s\n", aMatchStr, lookVal)
		sb = append(sb, lookVal)
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

	outFile, sferr := os.Create(outFiName)
	if sferr != nil {
		fmt.Println("Can not open out file ", outFiName, " sferr=", sferr)
		os.Exit(3)
	}

	defer inFile.Close()
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
			fmt.Println(outFile, outStr)
		}
	}
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
	u.processFile(inFiName, outFiName)
	jutil.Elap("L382: Finished Run", startms, jutil.Nowms())
}
