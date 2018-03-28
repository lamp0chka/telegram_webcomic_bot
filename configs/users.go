package configs

func (c *Configs) UserAddSource(uid, source string) {
	u, ok := c.Users[uid]
	if !ok {
		u = make([]string, 1)
	}
	add := true
	for _, s := range(u) {
		if s == source {
			add = false
			break
		}
	}
	if add {
		c.Users[uid] = append(u, source)
	}
}

func (c *Configs) UserDelSource(uid, source string) {
	u, ok := c.Users[uid]
	if ok {
		index := -1
		for i, s := range(u) {
			if s == source {
				index = i
				break
			}
		}
		if index >= 0 {
			c.Users[uid] = append(u[:index], u[index+1:]...)
		}
	}
}

func (c *Configs) GetUser(uid string) ([]string, bool) {
	u, ok := c.Users[uid]
	return u, ok
}

func (c *Configs) GetUsers() ([]string) {
	uids := make([]string, len(c.Users))
	for u := range c.Users{
		uids = append(uids, u)
	}
	return uids
}
