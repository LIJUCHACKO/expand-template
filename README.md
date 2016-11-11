# expandtemplate
Create text file containing repeated/incremental/decremental pattern from a given template. It is very useful especially when you want create table for a relational database

 #USAGE
---------------


sample_template.txt
----------------------------------------------------------
    #SEQ1 srno=1:step:999
    #FOR entries=1:1:noofentries
    <<srno:3>>)  Entry no=<<entries>>
    #END_FOR
    #END_SEQ1
--------------------------------------------------------------

command line usage

./expandtemplate sample_template.txt output_file.txt step:2 noofentries:50
