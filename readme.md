# Interpolator

Reads a text or markdown file and replaces defined values with contents from previously defined content such as a dictionary or object definition.   

## Command Line API

### Sample Invocation

```
interpolate  -in=data/sample-api.md -out=out/sample-api.md -baseDir=./data/dict  -defaulVarPath=desc  -keepPaths=true
```

* **-in** = path of a input file to process.  May be a specific file or a glob pattern.

* **-out** = Location to write the output file once expanded
* **-baseDir** = Directory Base to search for files named in interpolated parameters.

* **-defaultVarPath** = string matched in predefined file to pull next string. 

* **-keepPaths** = when set to true it will keep the supplied path as part of output text.   When not set or false will replace content of path with content.



# Data Format API

## Sample Input

```markdown
= Name Search =
uri: http://namesearch.com?fname=joe&lname=jackson&maxRec=389

* maxRec={maxrecmaxRec}
* fname= {person/fname}
* lname= {person/lname}

{inc:legal/copyright.txt}
```

* Any string contained inside of {} will be treated as a variable that needs to be resolved. 

  * The system will first attempt to match the defined variable name in parameters passed in on the command line.     
  * It will search for a file at the path starting at base dir and will attempt to load that file. 
  * Any pathname that includes a # segment will treat everything before the # as the file path and anything after it as a matching path to find a segment within the file.
  * When -defaultVarName is set the system will use it to search inside the file content for a specified field.     
  * And value in {} that starts with https:// or http:// will be treated as a URI and the system will return any text returned from that URI in the output document. 

* The system will always match the closest } whenever it encounters the opening {.  The closing } must not be separated by vertical white space such as /n.

* Any variable name not matched with a valid defined string or file will be replaced be included untouched.

* Any interpolated variable starting with "inc:" will be treated as a simple file include where the entire file located at the specified path will be read and inserted into the output replacing the {path}. 

* Dictionary YML must follow a simplified format.  The match pattern is simply the variable name that must start on first of line followed by :.   It will include text until it detects the next word followed by a : followed by a space.

  



## Sample Data Dictionary LOOKUP

**person/fname.yml**

```
name: fname
domain: addressbook
table: person
desc: First name of person
type: string
len: 50
```

**person/lname.yml**

```
name: lname
domain: addressbook
table: person
desc: Last Name of person
type: string
len: 50
```

**legal/copyright.txt**

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

