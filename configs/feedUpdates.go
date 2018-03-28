package configs

import (
	"time"
)

func (c *Configs) UpdateFeed(url string, updateTime time.Time) {
	c.FeedUpdates[url] = updateTime
}

func (c *Configs) GetFeed(url string) (time.Time, bool) {
	t, ok := c.FeedUpdates[url]
	return t, ok
}

func (c *Configs) IsUpToDate(url string, lastEntry time.Time) (bool) {
	t, ok := c.GetFeed(url)

	uptodate := !ok
	if ok {
		uptodate = !(t.Before(lastEntry))
	}

	return uptodate
}
