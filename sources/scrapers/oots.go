package scrapers

import (
	"github.com/mmcdole/gofeed"
	"time"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/url"
	"github.com/pkg/errors"
	"telegram_webcomic_bot/configs"
	"strconv"
)

// OotS feed has only the global lastBuildDate tag, so the items must be kept in order
// using the comic number, to check if we are up to date, we check the comic number and store
// that additional info.

func TimeParserOots(item *gofeed.Item, feed *gofeed.Feed, src string) (time.Time, string, error) {
	conf := configs.GetConfigs()

	item_num := item.Title[:strings.Index(item.Title, ":")]
	item_i, err := strconv.Atoi(item_num)
	if err != nil {
		return time.Time{}, "", err
	}

	last_item, ok := conf.GetLastItem(src)
	if ok {
		last_i, err := strconv.Atoi(last_item)
		if err != nil {
			return time.Time{}, "", err
		}

		ok = last_i >= item_i
	}

	if !ok {
		t, err := time.Parse(time.RFC1123, feed.Updated)
		// add the comic number as microseconds to sort them correctly if more than one comic
		// has to be sent (the motifyNewSources function sorts the items based on update time)
		t = t.Add(time.Duration(item_i) * time.Microsecond)
		return t, item_num, err
	} else {
		t, ok := conf.GetFeed(src)
		if !ok {
			t = time.Time{}
		}
		return t, item_num, nil
	}
}

func ScrapeOots(item *gofeed.Item, src string) ([]ComicUpdate, error) {
	gq, err := goquery.NewDocument(item.Link)
	if err != nil {
		return nil, err
	}

	image := ""
	gq.Find("td > img").Each(func(i int, selection *goquery.Selection) {
		if len(image) <= 0 {
			img, ok := selection.Attr("src")
			if ok && strings.HasPrefix(img, "/comics/images/") {
				image = img
			}
		}
	})

	if len(image) <= 0 {
		return nil, errors.New("image link not found!")
	}

	u, err := url.Parse(item.Link)
	if err != nil {
		return nil, err
	}

	u.Path = image
	image = u.String()

	c := make([]ComicUpdate, 1)
	c[0] = ComicUpdate{
		Source: src,
		Title: item.Title,
		Url: image,
	}

	return c, nil
}
