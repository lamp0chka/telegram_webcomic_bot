package configs

import (
	"time"
)

func (c *Configs) UpdateFeed(name string, updateTime time.Time) {
	c.flock.Lock()
	c.contents.FeedUpdates[name] = updateTime
	c.flock.Unlock()
	c.Store()
}

func (c *Configs) GetFeed(name string) (time.Time, bool) {
	c.flock.RLock()
	t, ok := c.contents.FeedUpdates[name]
	c.flock.RUnlock()
	return t, ok
}

func (c *Configs) GetFeedSources() ([]string) {
	c.flock.RLock()
	srcs := make([]string, len(c.contents.FeedUpdates))
	i := 0
	for f, _ := range c.contents.FeedUpdates {
		srcs[i] = f
		i++
	}
	c.flock.RUnlock()
	return srcs
}

func (c *Configs) IsUpToDate(name string, lastEntry time.Time) (bool) {
	t, ok := c.GetFeed(name)

	uptodate := !ok
	if ok {
		uptodate = !(t.Before(lastEntry))
	}

	return uptodate
}
