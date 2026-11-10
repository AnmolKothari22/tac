package commands
import (
	"sync"
	"strconv"
	"math"
//	"time"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)


func can_it_dec(num_to_conv *int64,sub *int64)bool{
	if(*num_to_conv>0){
		if(*sub<0){
			if((*num_to_conv - math.MaxInt64)<=*sub){
				return true
			}
			return false
		}
		return true;
	}else{
		if(*sub>0){
			if((*num_to_conv - math.MinInt64)>=*sub){
				return true
			}
			return false
		}
		return true;
	}
}


func Handle_dec(rmu *sync.RWMutex , parts []string,mp *map[string]shards.Stored_values) (string,inputoutput.WriteType) {
	rmu.Lock()
	defer rmu.Unlock()

	//checking key correctness
	comb , found := (*mp)[parts[1]];

	//if key not found
	if !found{
		mess:="key not found"
		return mess,inputoutput.KeyMiss ;
	}

	//checking if key is expired 
	if comb.IsExpired(){//checking if exp is before current time
		delete((*mp),parts[1])
		mess:="Key expired"
		return mess,inputoutput.KeyMiss;
	}

	//extracting value from key
	val := comb.Value

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
	}

	//converting number to be subtracted to int
	sub_num,err := strconv.ParseInt(dec_by,10,64)
	if err!=nil{                            //checking value correctness
		mess:="decreasing by a NAN not allowed"
		return mess,inputoutput.Error;
	}

	//checking for int64 bound and then returning
	if(can_it_dec(&conv_num,&sub_num)){
		conv_num=conv_num-sub_num
		(*mp)[parts[1]]= shards.Stored_values{Value:strconv.FormatInt(conv_num,10) , Created_at: comb.Created_at , Expire_after: comb.Expire_after } 
		mess:="key: "+parts[1]+" decreased by "+dec_by
		return mess, inputoutput.Normal
	}else{
		mess:="value going out of bound "
		return mess, inputoutput.Error
	}

}