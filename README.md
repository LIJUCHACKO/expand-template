# EXPANDTEMPLATE
This utility creates text file containing repeated/incremental/decremental patterns on the basis of a given template. It is very useful especially when you want to create table for a relational database

#USAGE

-sample_template.txt

----------------------------------------------------------
    #SEQ1 srno=1:step:99
    //* Heading- (this is a comment)
    #FOR entries=1:1:noofentries
    <<srno:3>>)  Entry no=<<entries>>
    #END_FOR
    #END_SEQ1
--------------------------------------------------------------

-command line usage

    ./expandtemplate sample_template.txt output_file.txt step:2 noofentries:50


    ./expandtemplate [input file] [output file] [variable name:value].......
    --include_comments will include comments also (for debugging purpose)
