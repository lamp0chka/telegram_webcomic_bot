package sources

import (
	"github.com/mmcdole/gofeed"
	"github.com/PuerkitoBio/goquery"
)

func scrapeDilbert(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	gq, err := goquery.NewDocument(item.Link)
	if err != nil {
		return nil, err
	}

	comic, ok := gq.Find(".img-comic").First().Attr("src")
	if !ok {
		comic = item.Link
	}

	c := make([]comicUpdate, 1)
	c[0] = comicUpdate{
				source: src.name,
				title: item.Title,
				url: comic,
			}
	return c, nil
}
