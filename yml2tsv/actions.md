# yml2tsv - Features & Actions



## Planned:

1. Provide better support for nested yml attributes.  By processing the leading space we should be able to process contained attributes such as conversion containing name containing 3 fields desc,  edi-doc,  edi-field,  conversion_logic t and output them as conversion.desc, conversion.edi-doc, conversion.edi-field. person.name.first.

## Under Consideration:

1. Provide a feature to allow -sep to allow production of CSV.     Provide -sepRep to replace separator with different character. 
2. Provide a GREP feature to further filter file names to include that can not be filtered with GREP.
3. Provide Wildcard support for variable names such as  \*desc\* which would detect desc:  techDesc: and description:   This could be complex because we can not predict the set of matches and without that we could not build our output set.  We may need to scan the data set to find the matches first and rescan to build output.
4. Provide feature to make variable name matching case insensitive.   

## Completed:

1. * DONE:JOE:2019-04-10: Convert \n and TAB to SPACE in embedded comment to make compatible for loading with excel.
2. * DONE:JOE:2019-04-10:  Properly handled multi-lines of text after initial YML variable TAG name.
3. * DONE:JOE:2019-04-10: Basic implementation of tool to read a set of simple YML files.  Extract a set of fields and produce a TSV (Tab Delimited File) containing those files. 