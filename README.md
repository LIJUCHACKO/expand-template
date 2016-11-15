# EXPANDTEMPLATE
This utility creates text file containing repeated/incremental/decremental patterns on the basis of a given template. It is very useful especially when you want to create table for a relational database

#USAGE

##INPUT FILE
- sample_template.txt

----------------------------------------------------------
    #SEQ1 srno=1:step:99
    //* Heading- (this is a comment)
    #FOR entries=1:1:noofentries
    <<srno:3>>)  Entry no=<<entries>>
    #END_FOR
    #END_SEQ1
----------------------------------------------------------

- Note: Read "writing expand template script.pdf" for help

##EXPANDING
- command line usage

    ./expandtemplate sample_template.txt output_file.txt step:2 noofentries:10


    ./expandtemplate [input file] [output file] [variable name:value].......
    --include_comments will include comments also (for debugging purpose)
- Note: All computations are integer based

##OUTPUT
- output_file.txt

----------------------------------------------------------
    001)    Entry no=1
    003)    Entry no=2
    005)    Entry no=3
    007)    Entry no=4
    009)    Entry no=5
    011)    Entry no=6
    013)    Entry no=7
    015)    Entry no=8
    017)    Entry no=9
    019)    Entry no=10
----------------------------------------------------------
