package scrapers

import (
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// Oglaf has age check cookie and some comics span across multiple pages, but this information can only be desumed
// from the comic webpage, the feed only contains the link to the page, so we first need to get the comic page
// passing the age confirmation cookie, scrape the page to get the comic url and search for the "next page" button,
// if the button is present inform the main scraper function that it has to load another comic page
func getOglafStripAndNext(url string) (*ComicUpdate, string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	req.AddCookie(&http.Cookie{
		Name: "AGE_CONFIRMED",
		Value: "yes",
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}

	gq, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, "", err
	}

	img := gq.Find("#strip").First()
	link, ok := img.Attr("src")
	if !ok {
		return nil, "", errors.New("unable to find image url")
	}

	alt, ok := img.Attr("title")

	nextPage := gq.Find("#nx").First().Parent()
	np, _ := nextPage.Attr("href")

	return &ComicUpdate{
		Url: link,
		Alt: alt,
	}, np, nil
}

func ScrapeOglaf(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	c := make([]ComicUpdate, 0)

	var comic *ComicUpdate
	var err error
	link := item.Link
	nl := link

	for len(nl) > 0 {
		comic, nl, err = getOglafStripAndNext(link)
		if err != nil {
			return nil, err
		}

		comic.Title = item.Title
		comic.Source = src

		c = append(c, *comic)

		u, err := url.Parse(link)
		if err != nil {
			return nil, err
		}
		u.Path = nl
		link = u.String()
	}

	return c, nil
}
