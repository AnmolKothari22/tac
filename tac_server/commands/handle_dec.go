package commands
import (
	"sync"
	"strconv"
	"math"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_dec(rmu *sync.RWMutex , parts []string,mp *map[string]string) (string,inputoutput.WriteType) {
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

	//if more then 2 arguments then take 2nd indexed argument as amount by which you have to decrease
	var dec_by string = "1"
	if len(parts)>2{
		dec_by=parts[2]
		if len(dec_by)>20{ // if lenght of stringy number is greater then 20 , then its already out of bound!!
			mess:="NAN or number out of bound(values should be >=INT64MIN )"
			return mess,inputoutput.Error;
		}
	}
	sub_num,err := strconv.ParseInt(dec_by,10,64)
	if err!=nil{                            //checking value correctness
		mess:="decreasing by a NAN not allowed"
		return mess,inputoutput.Error;
	}
	if((conv_num-math.MinInt64)>=sub_num){//int 64 out of bound checking 
		conv_num=conv_num-sub_num
	}else{
		mess:="value is going out of bound"
		return mess,inputoutput.Error;
	}

	(*mp)[parts[1]]= strconv.FormatInt(conv_num,10); 
	mess:="key: "+parts[1]+" decreased by "+dec_by
	return mess, inputoutput.Normal
}