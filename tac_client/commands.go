package client

func (c *Client) Set(key, value string) error { //SET key value
    _, err := c.execute("SET", key, value)
    return err
}

func (c *Client) Set_TTL(key, value , timetype,time string) error { // SET key value ms|s time
    _, err := c.execute("SET", key, value,timetype,time)
    return err
}

func (c *Client) Get(key string) (string, error) { // GET key
    return c.execute("GET", key)
}