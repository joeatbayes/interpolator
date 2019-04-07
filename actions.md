# Up Next Approved

1. Fix the sample from lname.yml where it is not detecting end of pre-formated markup

2. Add support example showing how to use primary and backup lookup varnames

   

# Under Consideration:

1. Save a list of the markdown and html generated and save as a index page.
2. Add a notion to import specification such as + before actual text to allow specification of keeping path.   EG:  {*person/name#desc} would just substitute the label with path with the content of the person/name yml with desc field.  Where {*+*person/name#desc} would keep the label path and add it to the output.
3. Implement file directory recursive walk for input.  Generate same output directory in the output directory.
4. Add Cache of input files read as key values.
5. Add Cache of full files plus var name segment 
6. Implement a full YAML parser for inputs.



# DONE:

1. 2019-04-07: Modify generated HTML to reload the page to coincide with the loopDelay specified on command line.  This will cause the browser to reload the page as new versions are saved.
2. 2019-04-07: Add feature to reprocess once every x seconds to regenerate output files. 
3. 2019-04-07: Add CSS Styling to the generated Markdown.
4. 2019-04-07: Modify default  input field a array to try first one path and then the next when looking up variables.
5. 2019-04-07: When reading a directory need to write the output in a predictable output directory.  Should change semantic to always write file of same name as input in the specified output directory which makes output directory mandatory.
6. 2019-04-07:JOE: Implement a HTML Output Directory using the basic markdown parser.
7. 2019-04-07:JOE: Modify from input file / output file to input dir, output dir with a glob path to allow easier use when generating html output derived from input name.
8. 2019-04-07:JOE: I have checked in a new version:   Modify to use {* instead of { working.   -inc: to include entire file working.   Lookup content of single field from data dictionary file working.   Allow specification of specific field such as len: using  path#fieldname working.    Using default match field to extract field from data dictionary working.   Allowing a list of default match fields so it can try desc: then tech-desc working.   basedir specified on command line is working. I still need to test the nested recursion to ensure that {} settings in values loaded from included files are expanded correctly.   Also need to create a better use case for fallback when desc: is missing but tech_desc is present in the default variable names.   
9. 2019-04-05:JOE: Add File inclusion using {*INC
10. 2019-04-05:JOE: Basic implementation with interpolation from command line.



