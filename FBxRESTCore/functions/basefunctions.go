package _functions

import(
	"strings"
)

func ReplaceAritmetics(str string) (string) {
	var s string = strings.Replace(str," .gt. "," > ",-1)
	s = strings.Replace(s," .lt. "," < ",-1)
	s = strings.Replace(s," .gte. "," >= ",-1)
	s = strings.Replace(s," .lte. "," <= ",-1)
	return s
}

func HasOperator(value string) bool {
	
	if(strings.HasPrefix(value,">")){
		 return true
	}
	if(strings.HasPrefix(value,"<")){
		return true
	}
	if(strings.HasPrefix(value,"=")){
		 return true
	}
	if(strings.HasPrefix(value,"!=")){
		return true
    }
	return false
}

//Splits
func SplitPars(str string, csplit string) (cmd []string) {
	var inComment = false
	var erg = ""
	var serg []string


	for _, v := range str {		
		if(string(v) == "'"){
			inComment = !inComment			
			erg += "'"
		} else if(string(v) == csplit && !inComment){						
			if(len(erg) > 0){
				serg = append(serg, erg)
			} 
			erg = ""
		} else{
			
			erg += string(v)
		}
	}
	
	if(len(erg) > 0){
		serg = append(serg, erg)
	}
	return serg
}