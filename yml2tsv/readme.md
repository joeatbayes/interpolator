# yml2csv

> **Extract Named Fields from set of YML files and convert to a TSV (Tab Delimited) file**



## Example Command

```bash
yml2csv  -in=../data/data-dict/db -out=tmp.tsv glob=*.yml  -vars=type,len,desc
```

```
-in  =  Directory to read. Files will be read recursively.

-out =  output file to write TSV file to.

-glob=  glob path filter used to filter set up input files
        processed.
        
-vars=  list of variable names to be included in output
        the relative path to input file is always first
        column.
```



## Example Output

```
path    type    size    desc
person\fname.yml        string          First name of person
person\lname.yml        string          Last Name of person
```

## Example Output in Excel

![sample-out-yml2tsv-excel-01.jpg](../docs/sample-out-yml2tsv-excel-01.jpg)

## Build

```
go get -u -t "github.com/joeatbayes/interpolator/yml2tsv"
```

