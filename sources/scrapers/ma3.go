package scrapers

import (
	"github.com/mmcdole/gofeed"
	"strings"
	"regexp"
)

var comicsThumb = regexp.MustCompile("comicsthumbs")
var comicsTb = regexp.MustCompile("tn[.]png")

func ScrapeMa3(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	img := item.Description[strings.Index(item.Description, "<img src=\"")+len("<img src=\""):]
	img = img[:strings.Index(img, "\"")]
	img = comicsThumb.ReplaceAllString(img, "comics")
	img = comicsTb.ReplaceAllString(img, ".png")

	c := make([]ComicUpdate, 1)
	c[0] = ComicUpdate{
		Source: src,
		Title: item.Title[strings.Index(item.Title, " - ") + len(" - "):],
		Url: img,
	}
	return c, nil
}
