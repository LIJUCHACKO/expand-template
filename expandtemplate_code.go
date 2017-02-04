package main
import (
  "github.com/zdebeer99/goexpression"
  "bufio"
  "os"
  "fmt"
  "strings"
  "strconv"
  "regexp"
)

func readLines(path string) ([]string, error) {
	  file, err := os.Open(path)
	  if err != nil {
	    return nil, err
	  }
	  defer file.Close()

	  var lines []string
	  scanner := bufio.NewScanner(file)
	  for scanner.Scan() {
	    lines = append(lines, scanner.Text())
	  }
	  return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	  file, err := os.Create(path)
	  if err != nil {
	    return err
	  }
	  defer file.Close()

	  w := bufio.NewWriter(file)
	  for _, line := range lines {
	    fmt.Fprintln(w, line+"\r")
	  }
	  return w.Flush()
}

func ABSOLUTE(n int)(int){
      a:=0
      if n < 0 {
	  a = -n
      } else {
	  a = n
      }
      return a
 }
func extractreplaceportion(line string) []string{
        //get substring between << >>
	copy:=false
	substring:=""
	substringlist:=[]string{}
	for i,char:=range line{
		if i<(len(line)-1) {
			if char=='<' {
				if line[i+1]=='<'{
					copy=true
				}	
			}
			if copy {
				substring=substring+line[i:i+1]
			}
			if char=='>' {
				if line[i+1]=='>'{
					substring=substring+">"
					copy=false
					substringlist=append(substringlist,substring)
					substring=""
				}	
			}
		}
	}
	return substringlist
} 

func evaluateexpress(formula string) int{
        context := map[string]interface{}{
		
	}
	result:=0
	re := regexp.MustCompile("([a-zA-Z!<>@=;:~_{}\\[\\]]+)")
	matching:=re.FindStringSubmatch(formula)
	if len(matching)>0 {
		for _,word:=range matching {
			fmt.Printf("\nUnknown variable : %s\n", word)
		}
		var yes string
		fmt.Scan(&yes)
		os.Exit(0)
	}
	if len(formula)>0{
	  result=int(goexpression.Eval(formula, context))
	}
	return result
}

func replacewholeword(line string,orginal string,rplace string ) (string,bool){
  re := regexp.MustCompile("([a-zA-Z_]+)") 
  finalline:=""
  copiedupto:=-1
  present:=false
  if len(orginal)==0{
      return line,present
  }
  for i,_:=range line {
	if line[i]==orginal[0] {
	    aftermatching:=false
	    beforematching:=false
	    if (len(line)>len(orginal)+i){
	          aftermatching=re.MatchString(line[len(orginal)+i:len(orginal)+i+1])
	    }
	    if i>0 {
	          beforematching=re.MatchString(line[i-1:i])
	    }
	    if !aftermatching && !beforematching &&  (len(line)>len(orginal)+i-1) {
	        if line[i:i+len(orginal)]==orginal {
		  finalline=finalline+rplace
		  present=true
		  copiedupto=len(orginal)+i-1
		}
	    }
	}
	if i>copiedupto {
	  finalline=finalline+line[i:i+1]
	}
   }
   return finalline,present
}

func replacevariablename(line string,orginal string,rplace string ) (string,bool){
	orginal=strings.TrimSpace(orginal)
	rplace=strings.TrimSpace(rplace)
	present:=false
	final:=""
	if strings.Contains(line, "#IF:") ||strings.Contains(line, "#END_") || strings.Contains(line, "#FOR") || strings.Contains(line, "#SEQ") {
		final,present=replacewholeword(line,orginal,rplace)
	}else {
		final=line
		substringlist:=extractreplaceportion(line)
		
		for _,substring:=range substringlist {
			substringnew,yes:=replacewholeword(substring,orginal,rplace)
			if yes{
			    present=yes
			}
			final=strings.Replace(final,substring,substringnew,-1)
		}
	}
	
	return final,present
}




func main() {
        fmt.Printf("---expandtemplate--VER. 15/12/2016 \n -------------------------------")
	argsWithProg := os.Args
         if len(argsWithProg)<3 {
	    fmt.Printf("\n   usage- [--include_comments][parameter1:value] [parameter2:value].. [template-file] [output-file]")
		var yes string
		fmt.Scan(&yes)
		return

	 }
	started:=0
	filename:=""
	finaloutputfile:=""
	include_comments:=false
	for i,argument:= range argsWithProg {
		if i>0 {
			if !strings.Contains(argument, ":") && !strings.Contains(argument, "--") {
			        if len(filename)==0 {
					filename=argument
				}else{
					finaloutputfile=argument
				}
			} else if strings.Contains(argument, "--include_comments"){
			      include_comments=true
			}
		}
	}
	
	
	processedcontent:=[]string{}
	file_content,eerr := readLines(filename)
	if eerr != nil {
		fmt.Printf("Read : %s\n", eerr)
		var yes string
		fmt.Scan(&yes)
		os.Exit(0)
	}
	for j,argument:= range argsWithProg {
		if j>0 {
			if strings.Contains(argument, ":") {
				temp:=strings.Split(argument, ":")
				for i,_:= range file_content {
					file_content[i],_=replacevariablename(file_content[i],temp[0],temp[1])
					if (strings.Contains(file_content[i], "#DECREMENT_SEQ") || strings.Contains(file_content[i], "#INCREMENT_SEQ") || strings.Contains(file_content[i], "#RESTART_SEQ") || strings.Contains(file_content[i], "#IF:") ||strings.Contains(file_content[i], "#SEQ") || strings.Contains(file_content[i], "#FOR")|| strings.Contains(file_content[i], "#END_")) {
					      file_content[i]=strings.Replace(file_content[i],",","",-1)
					      file_content[i]=strings.Replace(file_content[i],";","",-1)
					      file_content[i]=strings.Replace(file_content[i],"\"","",-1)					
					}

				}
			} 	
		}
	}

	
	///////FOR LOOPS
	fmt.Printf("\nExpanding FOR LOOPS")
	continueexcu:=true
	for continueexcu {
	       continueexcu=false
		started=0
		copyformulti:=false
		partcontent:=[]string{}
		variablename:=""
		startedval:=0
		step:=0
		afterverynoline:=1
		endval:=0
		for _,line:= range file_content {
			if strings.Contains(line, "#END_FOR") && started>0 {
				if started==1 {
					copyformulti=false
					continueexcu=true
					fmt.Printf("...OVER")
					if step>0 {
						for i:=startedval; i<=endval;i=i+step {
							for k:=0;k<afterverynoline;k++ {
							      for _,tline:= range partcontent {
								      replacement,_:=replacevariablename(tline,variablename,strconv.Itoa(int(i)))
								      processedcontent=append(processedcontent,replacement)
							      }
							}
							
						}
					}
					if step<0 {
						for i:=startedval; i>=endval;i=i+step {
							for k:=0;k<afterverynoline;k++ {
							      for _,tline:= range partcontent {
								      replacement,_:=replacevariablename(tline,variablename,strconv.Itoa(int(i)))
								      processedcontent=append(processedcontent,replacement)
							      }
							}
						}
					}
					partcontent=partcontent[:0]
				}
				if (started>0){
				  started=started-1	
				}
			}

			if copyformulti{
				partcontent=append(partcontent,line)
				
			} else {
				if !strings.Contains(line, "#FOR") && !strings.Contains(line, "#END_FOR") {
					processedcontent=append(processedcontent,line)
				}
			}

			if strings.Contains(line, "#FOR") {
				if started==0 {
					fmt.Printf("\n--Expanding %s",line)
					copyformulti=true
					temp2:=strings.Split(line, "=")
					seqvar:=strings.Split(temp2[0], " ")
					variablename=strings.TrimSpace(seqvar[1])
					limits:=strings.Split(temp2[1], ":")
					startedval=evaluateexpress(limits[0])
					if strings.Contains(limits[1], "/") {
					    temp3:=strings.Split(limits[1], "/")
					    step=evaluateexpress(temp3[0])
					    afterverynoline=evaluateexpress(temp3[1])
					} else {
					    afterverynoline=1
					    step=evaluateexpress(limits[1])
					}
					endval=evaluateexpress(limits[2])
					//if endval>=startedval {
					//     step=ABSOLUTE(step)
					//}else{
					//     step=-ABSOLUTE(step)
					//}
				} 
				started=started+1	
				

			}
	      
		 }
		 
		 
		 file_content=file_content[:0]
		 for _,line:= range processedcontent {
		    file_content=append(file_content,line)
		  }
		 //fmt.Println(file_content)
		 processedcontent=processedcontent[:0]
		 if started>0 {
			fmt.Printf("\nERROR:#END_FOR missing [%s]",filename)
			var yes string
			fmt.Scan(&yes)
			os.Exit(0)
		}
	 }
	 
	 processedcontent=processedcontent[:0]
	 fmt.Printf("\nApplying IF conditions")
	 donotcopy:=false
	 hierarchy:=0
	 reablecpyat:=0
	 for _,line:= range file_content {
		if strings.Contains(line, "#IF:")   {
		      hierarchy=hierarchy+1
		      reablecpyat=hierarchy
		      if !donotcopy {
			     
			      temp2:=strings.Split(line, ":")
			      condition:=temp2[1]
			      if strings.Contains(condition, ">") {
				  temp3:=strings.Split(condition, ">")
				  if !(evaluateexpress(temp3[0])>evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
				  
			      }
			      if strings.Contains(condition, ">=") {
				  temp3:=strings.Split(condition, ">=")
				  if !(evaluateexpress(temp3[0])>=evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
			      }
			      if strings.Contains(condition, "<") {
				  temp3:=strings.Split(condition, "<")
				  if !(evaluateexpress(temp3[0])<evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
			      }
			      if strings.Contains(condition, "<=") {
				  temp3:=strings.Split(condition, "<=")
				  if !(evaluateexpress(temp3[0])<=evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
			      }
			       if strings.Contains(condition, "==") {
				  temp3:=strings.Split(condition, "==")
				  if !(evaluateexpress(temp3[0])==evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
			      }
			       if strings.Contains(condition, "!=") {
				  temp3:=strings.Split(condition, "!=")
				  if !(evaluateexpress(temp3[0])!=evaluateexpress(temp3[1])){
				      donotcopy=true
				  }
			      }
		               if donotcopy {
				  fmt.Printf("\n--Applying %s (CONDITIONS FALSE)",line)
			      }else {
                                  fmt.Printf("\n--Applying %s (CONDITIONS TRUE)",line)
			      }
		      }
		}
		if !donotcopy {			      
					processedcontent=append(processedcontent,line)
		}
		if strings.Contains(line, "#END_IF")   {
		    if reablecpyat==hierarchy {
			donotcopy=false
		    }
		    hierarchy=hierarchy-1
		    //fmt.Printf("...OVER")
		}
	 }
	 if donotcopy {
			fmt.Printf("\nERROR:#END_IF missing [%s]",filename)
			var yes string
			fmt.Scan(&yes)
			os.Exit(0)
	 }
	 file_content=file_content[:0]
	for _,line:= range processedcontent {
	    file_content=append(file_content,line)
	}
	 processedcontent=processedcontent[:0]
	 
	/////SEQ 
	fmt.Printf("\nNumbering SEQ")
	continueexcu=true
	ENDSEQ:="#END_SEQ999"
	RESTARTSEQ:="#RESTART_SEQ**"
	INCREMSEQ:="#INCREMENT_SEQ**"
	DECREMSEQ:="#DECREMENT_SEQ**"
	for continueexcu {
	        continueexcu=false
		started=0
		variablename:=""
		startedval:=0
		step:=0
		endval:=0
		value:=0
		afterverynoline:=1
		countline:=0
		header:=""
		for _,line:= range file_content {
			if strings.Contains(line, ENDSEQ)   {
				if started==1 {
					continueexcu=true
					step=0
					started=0
					fmt.Printf("...OVER")
					header=""
				} 
			}
			if strings.Contains(line, INCREMSEQ) || strings.Contains(line, DECREMSEQ) {
			       if(countline>0)  {  
				      countline=0
				    value=value+step 
				      if step<0 {
					      if value<endval{
						      value=startedval;
					      }
				      }
				      if step>0 {
					      if value>endval{
						      value=startedval;
					      }
				      }
			       }
			} else if strings.Contains(line, RESTARTSEQ) { 
			       value=startedval;
			} else if (!strings.Contains(line, "#SEQ") && !strings.Contains(line, "#END_SEQ")) {
				if step!=0{
				        replacement,present:=replacevariablename(line,variablename,strconv.Itoa(int(value)))
					processedcontent=append(processedcontent,replacement)
					if present {
						countline=countline+1
					        if countline>=afterverynoline{
						    value=value+step 
						    countline=0
						}

						if step<0 {
							if value<endval{
								value=startedval;
							}
						}
						if step>0 {
							if value>endval{
								value=startedval;
							}
						}
					}
				} else {
				        	processedcontent=append(processedcontent,line)
				}
				
			}else if ((strings.Contains(line, "#SEQ")&& started>0 || strings.Contains(line, "#END_SEQ") || strings.Contains(line, "#RESTART_SEQ") || strings.Contains(line, "#INCREMENT_SEQ") || strings.Contains(line, "#DECREMENT_SEQ")  )  && !strings.Contains(line, ENDSEQ)) {
					processedcontent=append(processedcontent,line)
			}
			
			if strings.Contains(line, "#SEQ")   {
				if started==0 {
					fmt.Printf("\n--Numbering %s",line)
					temp2:=strings.Split(line, "=")
					seqvar:=strings.Split(temp2[0], " ")
					variablename=strings.TrimSpace(seqvar[1])
					limits:=strings.Split(temp2[1], ":")
					startedval=evaluateexpress(limits[0])
					value=startedval
					if strings.Contains(limits[1], "/") {
					    temp3:=strings.Split(limits[1], "/")
					    step=evaluateexpress(temp3[0])
					    afterverynoline=evaluateexpress(temp3[1])
					} else {
					    afterverynoline=1
					    step=evaluateexpress(limits[1])
					}
					endval=evaluateexpress(limits[2])
					header=strings.TrimSpace(seqvar[0])
					ENDSEQ="#END_"+strings.Replace(header,"#","",-1)
					RESTARTSEQ="#RESTART_"+strings.Replace(header,"#","",-1)
					INCREMSEQ="#INCREMENT_"+strings.Replace(header,"#","",-1)
					DECREMSEQ="#DECREMENT_"+strings.Replace(header,"#","",-1)
					started=1
				        if endval>=startedval {
					     step=ABSOLUTE(step)
					}else{
					     step=-ABSOLUTE(step)
					}
					countline=0
				} else if strings.Contains(line, header) {
				      fmt.Printf("\nERROR: %s  use numbering scheme  for overlapping SEQs (eg: #SEQ1, #SEQ2..or #SEQa, #SEQb.. ) [%s]",header,filename)
				      var yes string
				      fmt.Scan(&yes)
				      os.Exit(0)
				  
				}
			}
		 }
		 file_content=file_content[:0]
		 for _,line:= range processedcontent {
		    file_content=append(file_content,line)
		  }
		  
		 //fmt.Println(file_content)
		 processedcontent=processedcontent[:0]
		 if started>0 {
		      fmt.Printf("\nERROR: %s missing [%s]",ENDSEQ,filename)
		      var yes string
		      fmt.Scan(&yes)
		      os.Exit(0)
	      } 
	 }
	// fmt.Println(file_content)
	 
	 fmt.Printf("\nSolving expressions")
	 
	 //process all expressions in statements
	 for i,_:= range file_content {
		substringlist:=extractreplaceportion(file_content[i])
		for _,substring:=range substringlist {
			temp:=strings.Replace(substring, "<<", "", -1)
			temp=strings.Replace(temp, ">>", "", -1)
			temp2:=strings.Split(temp, ":")
			expresult:=strconv.Itoa(evaluateexpress(temp2[0]))
			if len(temp2)>1 {
				str, _ := strconv.Atoi(temp2[1])
				for len(expresult)<str{
					expresult="0"+expresult
				}
			}
			file_content[i]=strings.Replace(file_content[i],substring,expresult,-1)
		}
	 }
	 
	 processedcontent=processedcontent[:0]
	 for _,line:= range file_content {
	      if (!strings.Contains(line, "#RESTART_SEQ") && !(strings.Contains(line, "//*") && !include_comments) && !strings.Contains(line, "#IF:") && !strings.Contains(line, "#DECREMENT_SEQ") && !strings.Contains(line, "#INCREMENT_SEQ") && !strings.Contains(line, "#SEQ") && !strings.Contains(line, "#FOR")&& !strings.Contains(line, "#END_")) {
		    processedcontent=append(processedcontent,line)
	      }
	 }
	 
	 
	 
	if err := writeLines(processedcontent, finaloutputfile); err != nil {
		fmt.Printf("writeLines: %s", err)
		var yes string
		fmt.Scan(&yes)
		return
	    
	}
	fmt.Printf("\nOVER: saved to %s\n",finaloutputfile)
}
