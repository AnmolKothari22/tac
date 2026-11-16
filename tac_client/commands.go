package client


func (c *Client) Set(key, value string) error { //SET key value
    _, err := c.execute("SET", key, value)
    return err
}

//timetype is either ms or s
func (c *Client) Set_TTL(key, value , timetype,time string) error { 
    _, err := c.execute("SET", key, value,timetype,time)
    return err
}

// remaining_time_unit can only be ms or s 
// ms for milliseconds
// s for seconds
func (c *Client) TTL(key ,remaining_time_unit string)(string,error){
    return c.execute("TTL",key,remaining_time_unit)
}

func (c *Client) Get(key string) (string, error) { // GET key
    return c.execute("GET", key)
}

func (c *Client) Del(key string)(error){
    _, err:= c.execute("DEL",key)
    return err
}

func (c *Client) Dec(key ,decrease_by string)error{
    _,err := c.execute("DEC",key,decrease_by)
    return err
}

func (c *Client) Inc(key ,increase_by string)error{
    _,err :=  c. execute("INC",key,increase_by)
    return err
}