# Up Next Approved

1. Finishing implementation of simple variable substitution from a data dictionary

# Under Consideration:

1. Implement a HTML Output Directory using the basic markdown parser.
2. When reading a directory need to write the output in a predictable output directory.  Should change semantic to always write file of same name as input in the specified output directory which makes output directory mandatory.
3. Modify default  input field a array to try first one path and then the next when looking up variables.
4. Add a notion to import specification such as + before actual text to allow specification of keeping path.
5. Implement file directory recursive walk for input.  Generate same output directory in the output directory.
6. Add Cache of input files read as key values.
7. Add Cache of full files plus var name segment 
8. Implement a full YAML parser for inputs.



# DONE:

1. 2019-04-07:JOE: Modify from input file / output file to input dir, output dir with a glob path to allow easier use when generating html output derived from input name.
2. 2019-04-07:JOE: I have checked in a new version:   Modify to use {* instead of { working.   -inc: to include entire file working.   Lookup content of single field from data dictionary file working.   Allow specification of specific field such as len: using  path#fieldname working.    Using default match field to extract field from data dictionary working.   Allowing a list of default match fields so it can try desc: then tech-desc working.   basedir specified on command line is working. I still need to test the nested recursion to ensure that {} settings in values loaded from included files are expanded correctly.   Also need to create a better use case for fallback when desc: is missing but tech_desc is present in the default variable names.   
3. 2019-04-05:JOE: Add File inclusion using {*INC
4. 2019-04-05:JOE: Basic implementation with interpolation from command line.



