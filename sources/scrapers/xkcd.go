package scrapers

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"html"
)

func ScrapeXkcd(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	imgUrl := item.Description[strings.Index(item.Description, "src=\"") + len("src=\""):]
	imgUrl = imgUrl[:strings.Index(imgUrl, "\"")]

	alt := item.Description[strings.Index(item.Description, "alt=\"") + len("alt=\""):]
	alt = alt[:strings.Index(alt, "\"")]
	alt = html.UnescapeString(alt)

	c := make([]ComicUpdate, 1)
	c[0] = ComicUpdate{
		Source: src,
		Title: item.Title,
		Url: imgUrl,
		Alt: alt,
	}
	return c, nil
}

func ScrapeXkcdWhatIf(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	c := make([]ComicUpdate, 1)
	c[0] = ComicUpdate{
		Source: src,
		Title: item.Title,
		Url: item.Link,
	}
	return c, nil
}
