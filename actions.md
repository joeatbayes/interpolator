# Actions & Feature Planning

## Up Next Approved

1. Add support example showing how to use primary and backup lookup varnames

   

## Under Consideration:

1. Add support for relative files by looking for ./ as first path segment.  When that is encountered use the  directory of the last read file rather than search as starting path.
2. Add a notion to import specification such as + before actual text to allow specification of keeping path.   EG:  {*person/name#desc} would just substitute the label with path with the content of the person/name yml with desc field.  Where {*+*person/name#desc} would keep the label path and add it to the output.
3. Save a list of the markdown and html generated and save as a index page.
4. Cache File names and timestamp changed and only re-read if the file date modified has changed.
5. Modify to use the BlackFriday markdown generator more directly.  The current md to html tool is too large and slows down download.    When I tried it at first it did not produce what expected but I think it may be able to work with correct options and it is much smaller and more commonly used. 
6. Add Cache of input files read as key values.
7. Add Cache of full files plus var name segment 
8. Implement a full YAML parser for inputs.
9. Add support for inDir for comma separated list of names process in order received.



## DONE:

- 2019-04-09:JOE: Fixed bug causing indention of leading space to be lost.  This caused improper indenting of bullets for resulting markup.
- 2019-04-08:JOE: Modify all sample commands in readme to properly reference -glob 
- 2019-04-08:JOE: Add support for recursive walk of Dir Tree and generate new file with same name in relative tree in output. Implement file directory recursive walk for input.  Generate same output directory in the output directory.
- 2019-04-08:JOE: Fix keepVarnames so it works reliably and upgraded to work with yes or true on command line. 
- 2019-04-07:JOE:  Add example better showing a shared component.
- 2019-04-07:JOE: Fix the sample from lname.yml where it is not detecting end of pre-formatted markup
- 2019-04-07:JOE: Modify generated HTML to reload the page to coincide with the loopDelay specified on command line.  This will cause the browser to reload the page as new versions are saved.
- 2019-04-07: Add feature to reprocess once every x seconds to regenerate output files. 
- 2019-04-07: Add CSS Styling to the generated Markdown.
- 2019-04-07: Modify default  input field a array to try first one path and then the next when looking up variables.
- 2019-04-07: When reading a directory need to write the output in a predictable output directory.  Should change semantic to always write file of same name as input in the specified output directory which makes output directory mandatory.
- 2019-04-07:JOE: Implement a HTML Output Directory using the basic markdown parser.
- 2019-04-07:JOE: Modify from input file / output file to input dir, output dir with a glob path to allow easier use when generating html output derived from input name.
- 2019-04-07:JOE: I have checked in a new version:   Modify to use {* instead of { working.   -inc: to include entire file working.   Lookup content of single field from data dictionary file working.   Allow specification of specific field such as len: using  path#fieldname working.    Using default match field to extract field from data dictionary working.   Allowing a list of default match fields so it can try desc: then tech-desc working.   basedir specified on command line is working. I still need to test the nested recursion to ensure that {} settings in values loaded from included files are expanded correctly.   Also need to create a better use case for fallback when desc: is missing but tech_desc is present in the default variable names.   
- 2019-04-05:JOE: Add File inclusion using {*INC
- 2019-04-05:JOE: Basic implementation with interpolation from command line.



