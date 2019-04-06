# Interpolator

Reads a text or markdown file and replaces defined values with contents from previously defined content such as a dictionary or object definition.   

## Command Line API

### Sample Invocation

```
interpolate  -in=data -out=out glob=*sample*.md -search=./data/data-dict  -VarNames=desc,tech_desc  -keepNames=true -maxRec=99
```

* **-in** = path to input directory containing files to process .  Defaults to ./data
* **-out** = path to output directory where expanded files will be written.   defaults to ./out
* **-glob**= glob pattern to use when selecting files to process in input directory.  Defaults to *.md
* **-search** = Directory Base to search for files named in interpolated parameters.
* **-varNames** = Default variable name  matched in dictionary files.  Can be overridden if variable name is specified using #varname semantic.    May be common separated list to allow default lookup of backup fields such as look first in tech_desc then in desc.
* **-keepNames** = when set to true it will keep the supplied path as part of output text.   When not set or false will replace content of path with content.
* **-saveHtml**= when set to yes will convert the md file to Html and save it in the output directory. 



#### Saving md files converted to HTML

interpolate  -in=data -out=out glob=*sample*.md -search=./data/data-dict  -VarNames=desc,tech_desc  -keepNames=true -maxRec=99 -saveHtml=yes

Adding the -saveHtml=yes will cause the system to re-read the expanded markup and write a HTML version. 



# Data Format API

## Sample Input

[sample-api.md](data/sample-api.md)

```markdown
# Name Search

**path**: addrBook/search

**sample uri**: http://namesearch.com/addrBook/search?fname=joe&lname=jackson&maxRec=389

* **maxRec**={*maxRec}
* **fname**= {*person/fname}
* **lname**= {*person/lname#tech_desc}
  * **Type**={*person/lname#type}  **len**={*person/lname#len}

{*inc: inc/legal/copyright.txt}
```

* Any string contained inside of {} will be treated as a variable that needs to be resolved.  When first character after { is *.    The * was used to make it easier to avoid parsing and attempted interpolation when sample JSON or other curly brace languages are part of input. 

  * The system will first attempt to match the defined variable name in parameters passed in on the command line.     
  * It will search for a file at the path starting at base dir and will attempt to load that file. 
  * Any pathname that includes a # segment will treat everything before the # as the file path and anything after it as a matching path to find a segment within the file.
  * When -defaultVarName is set the system will use it to search inside the file content for a specified field.     
  * Any lookup key value in {} that starts with https:// or http:// will be treated as a URI and the system will return any text returned from that URI in the output document. 

* The system will always match the closest } whenever it encounters the opening {.  The closing } must not be separated by vertical white space such as /n.

* Any variable name not matched with a valid defined string or file will be replaced be included untouched.

* Any interpolated variable starting with "inc:" will be treated as a simple file include where the entire file located at the specified path will be read and inserted into the output replacing the {path}. 

* Dictionary YML must follow a simplified format.  The match pattern is simply the variable name that must start on first of line followed by :.   It will include text until it detects the next word followed by a : followed by a space.

  



## Sample Data Dictionary LOOKUP

**[person/fname.yml](data/data-dict/person/fname.yml)**

```
name: fname
domain: addressbook
table: person
desc: First name of person
type: string
len: 50
```

**[person/lname.yml](data/data-dict/person/lname.yml)**

```
name: lname
domain: addressbook
table: person
desc: Last Name of person
type: string
len: 50
```

**[legal/copyright.txt](data/data-dict/share/legal/copyright.txt)**

```
(C) Copyright Joseph Ellsworth Mar-2019
MIT License: https://opensource.org/licenses/MIT
Contact me on linkedin: https://www.linkedin.com/in/joe-ellsworth-68222/
```

## Sample Output

```
SAMPLE OUTPUT
```



# Build

Once you have the  [golang compiler](https://golang.org/dl/) installed.    

```
go get -u -t "github.com/joeatbayes/interpolator/interpolate"
```

It will create a executable in your GOPATH directory in bin/interpolate.  For windows it will be httpTest.exe.  If GOPATH is set the /tmp then the executable will be written to /tmp/bin/interpolate under Linux or /tmp/bin/interpolate.exe for windows. 

> Once you build the executable you can copy it to other computers without the compiler.   It only requires the single executable file.

HINT: set GOTPATH= your current working directory.  Or set it to your desired target directory and the go get command will create the executable in bin inside of that directory which is good because you may not have write privileges to the default GOPATH.

##### To Download all pre-built test cases, scripts and source code

```
git clone https://github.com/joeatbayes/interpolator interpolate
```

### To Build for Multiple OS Linux / Windows / Mac 

[make-go-all-os.bat](make-go-all-os.bat): Batch file to build executable for multiple OS.  runs on windows 

#### Windows example to build for other OS 

```bash
set GOPATH=%cd%

go get -u "github.com/joeatbayes/goutil/jutil"

set GOOS=darwin
set GOARCH=386
go build -o interpolate-darwin-386 interpolate/interpolate.go 

set GOOS=linux
set GOARCH=386
go build -o interpolate-linux-386 interpolate/interpolate.go 

set GOOS=windows
set GOARCH=386
go build -o interpolate-windows-386 interpolate/interpolate.go 

set GOOS=solaris
set GOARCH=amd64
go build -o interpolate-solaris-amd64 interpolate/interpolate.go 

```



