package commands
import (
	"sync"
	"strconv"
	"math"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)


func can_it_inc(num_to_conv *int64,add *int64)bool{
	if(*num_to_conv>0){
		if(*add>0){
			if((math.MaxInt64 - *num_to_conv)>=*add){
				return true
			}
			return false
		}
		return true;
	}else{
		if(*add<0){
			if((math.MinInt64 - *num_to_conv)<=*add){
				return true
			}
			return false
		}
		return true;
	}
}



func Handle_inc(rmu *sync.RWMutex , parts []string,mp *map[string]string)(string,inputoutput.WriteType){
	rmu.Lock()
	defer rmu.Unlock()

	//checking key correctness
	val , found := (*mp)[parts[1]];

	//if key not found
	if !found{
		mess:="key not found"
		return mess,inputoutput.KeyMiss ;
	}

	//if value is NAN
	conv_num , err := strconv.ParseInt(val,10,64)
	if err!=nil{
		mess:="value is NAN"
		return mess,inputoutput.Error;
	}

	//if more then 2 arguments then take 2nd indexed argument as amount by which you have to increase
	var inc_by string = "1"
	if len(parts)>2{
		inc_by=parts[2]
	}

	//converting number to be added to int
	add_num,err := strconv.ParseInt(inc_by,10,64)
	if err!=nil{                            //checking value correctness
		mess:="NAN not allowed/out of bound number"
		return mess,inputoutput.Error;
	}

	//checking for int64 bound and then returning
	if(can_it_inc(&conv_num,&add_num)){
		conv_num=conv_num+add_num
		(*mp)[parts[1]]= strconv.FormatInt(conv_num,10); 
		mess:="key: "+parts[1]+" increased by "+inc_by
		return mess, inputoutput.Normal
	}else{
		mess:="value going out of bound "
		return mess, inputoutput.Error
	}

}