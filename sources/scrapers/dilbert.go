package scrapers

import (
	"github.com/mmcdole/gofeed"
	"github.com/PuerkitoBio/goquery"
)

func ScrapeDilbert(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	gq, err := goquery.NewDocument(item.Link)
	if err != nil {
		return nil, err
	}

	comic, ok := gq.Find(".img-comic").First().Attr("src")
	if !ok {
		comic = item.Link
	}

	title := gq.Find(".comic-title-name").First().Text()
	if len(title) <= 0 {
		title = item.Title
	}

	c := make([]ComicUpdate, 1)
	c[0] = ComicUpdate{
				Source: src,
				Title: title,
				Url: comic,
			}
	return c, nil
}
