package sources

import (
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

func getOglafStripAndNext(url string) (*comicUpdate, string, error) {
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

	return &comicUpdate{
		url: link,
		alt: alt,
	}, np, nil
}

func scrapeOglaf(item *gofeed.Item, src feedSrc) ([]comicUpdate, error) {
	c := make([]comicUpdate, 0)

	var comic *comicUpdate
	var err error
	link := item.Link
	nl := link

	for len(nl) > 0 {
		comic, nl, err = getOglafStripAndNext(link)
		if err != nil {
			return nil, err
		}

		comic.title = item.Title
		comic.source = src.name

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
