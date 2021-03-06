package configs

import (
	"time"
)

func (c *Configs) UpdateFeed(name string, updateTime time.Time) {
	c.flock.RLock()
	_, ok := c.contents.FeedUpdates[name]
	c.flock.RUnlock()
	c.flock.Lock()
	if !ok {
		c.newFeeds = append(c.newFeeds, name)
	}
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

func (c *Configs) IsItemNew(name string, itemTime time.Time) (bool) {
	t, ok := c.GetFeed(name)

	newer := !ok
	if ok {
		newer = t.Before(itemTime)
	}

	return newer
}

func (c *Configs) GetNewFeedSources() ([]string) {
	c.flock.RLock()
	s := c.newFeeds[:]
	c.flock.RUnlock()
	return s
}

func (c *Configs) ClearNewFeedSources() {
	c.flock.Lock()
	c.newFeeds = c.newFeeds[:0]
	c.flock.Unlock()
}

func (c *Configs) PopNewFeedSources() ([]string) {
	c.flock.Lock()
	sources := c.newFeeds[:]
	c.newFeeds = c.newFeeds[:0]
	c.flock.Unlock()
	return sources
}

func (c *Configs) StoreLastItem(source, item string) {
	c.flock.Lock()
	c.contents.LastItem[source] = item
	c.flock.Unlock()
	c.Store()
}

func (c *Configs) GetLastItem(source string) (string, bool) {
	c.flock.RLock()
	item, ok := c.contents.LastItem[source]
	c.flock.RUnlock()
	return item, ok
}