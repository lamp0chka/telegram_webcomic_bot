package sources

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"html"
)

func scrapeXkcd(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	imgUrl := item.Description[strings.Index(item.Description, "src=\"") + len("src=\""):]
	imgUrl = imgUrl[:strings.Index(imgUrl, "\"")]

	alt := item.Description[strings.Index(item.Description, "alt=\"") + len("alt=\""):]
	alt = alt[:strings.Index(alt, "\"")]
	alt = html.UnescapeString(alt)

	c := make([]comicUpdate, 1)
	c[0] = comicUpdate{
		source: src.name,
		title: item.Title,
		url: imgUrl,
		alt: alt,
	}
	return c, nil
}

func scrapeXkcdWhatIf(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	c := make([]comicUpdate, 1)
	c[0] = comicUpdate{
		source: src.name,
		title: item.Title,
		url: item.Link,
	}
	return c, nil
}
