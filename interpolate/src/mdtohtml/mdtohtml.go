package mdtohtml

// Library to support producing HTML with CSS styling from MD
// files that look similar to that used by Github.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	s "strings"
	//"github.com/joeatbayes/goutil/jutil"
	"github.com/shurcooL/github_flavored_markdown"
	//"gopkg.in/russross/blackfriday.v2"
)

func SaveAsHTML(srcFile string, loopDelay float32) string {
	// Save File as HTML if the HTML save directory
	// has been specified.
	fext := filepath.Ext(srcFile)
	fmt.Println("L226: fext=", fext, "srcFi=", srcFile)
	if fext == ".md" {
		data, err := ioutil.ReadFile(srcFile)
		//fmt.Println("L229: fname=", srcFile, "data=", string(data))
		if err != nil {
			fmt.Println("L260: Error reading ", srcFile, " err=", err)
		} else {
			hname := s.Replace(srcFile, ".md", ".html", 1)
			//fmt.Println("L234: html filename=", hname)
			f, err := os.Create(hname)
			if err != nil {
				fmt.Println("error writing to ", hname, " err=", err)
			} else {
				f.WriteString(HtmlPrefix)
				if loopDelay > 0 {
					delayStr := fmt.Sprintf("%f", (loopDelay * 1000))
					//fmt.Println("L407: delayStr=", delayStr, "reloadStr=", htmlJavscriptReload)
					tmp := s.Replace(string(HtmlJavscriptReload), "loopDelay", delayStr, 1)
					f.WriteString(tmp)
				}
				f.Write(github_flavored_markdown.Markdown(data))
				f.WriteString(HtmlSuffix)
				f.Close()
				return hname
			}
		}
	}
	return ""
}
