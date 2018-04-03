package sources

import (
	"github.com/mmcdole/gofeed"
)

func scrapeDilbert(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	c := make([]comicUpdate, 1)
	c[0] = comicUpdate{
				source: src.name,
				title: item.Title,
				url: item.Link,
			}
	return c, nil
}
