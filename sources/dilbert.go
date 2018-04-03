package sources

import (
	"telegram_webcomic_bot/configs"
	"time"
	"github.com/mmcdole/gofeed"
)

func scrapeDilbert(items []*gofeed.Item, conf *configs.Configs, src feedSrc) (comics []comicUpdate, lastItemTime time.Time) {
	comics = make([]comicUpdate, 0)

	for _, item := range(items) {
		updated, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			parseError(src.name, err)
			continue
		}

		if lastItemTime.Before(updated) {
			lastItemTime = updated
		}

		if conf.IsItemNew(src.name, updated) {
			comics = append(comics, comicUpdate{
				src.name,
				item.Title,
				item.Link,
				"",
				true,
			})
		}
	}

	return
}
