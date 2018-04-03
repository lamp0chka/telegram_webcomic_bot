package sources

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"regexp"
)

var comicsUrl = regexp.MustCompile("comicsthumbs")

func scrapeMa3(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	img := item.Description[strings.Index(item.Description, "<img src=\"")+len("<img src=\""):]
	img = img[:strings.Index(img, "\"")]
	img = comicsUrl.ReplaceAllString(img, "comics")

	c := make([]comicUpdate, 1)
	c[0] = comicUpdate{
		source: src.name,
		title: item.Title[strings.Index(item.Title, " - ") + len(" - "):],
		url: img,
	}
	return c, nil
}
