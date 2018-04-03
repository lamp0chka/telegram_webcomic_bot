package sources

import (
	"telegram_webcomic_bot/configs"
	"time"
	"github.com/mmcdole/gofeed"
	"strings"
	"html"
)

func scrapeXkcd(items []*gofeed.Item, conf *configs.Configs, src feedSrc) (comics []comicUpdate, lastItemTime time.Time) {
	comics = make([]comicUpdate, 0, len(items))

	for _, item := range(items){
		updated, err := time.Parse(time.RFC1123Z, item.Published)
		if err != nil {
			parseError(src.name, err)
			continue
		}

		if lastItemTime.Before(updated) {
			lastItemTime = updated
		}

		if conf.IsItemNew(src.name, updated) {
			imgUrl := item.Description[strings.Index(item.Description, "src=\"") + len("src=\""):]
			imgUrl = imgUrl[:strings.Index(imgUrl, "\"")]

			alt := item.Description[strings.Index(item.Description, "alt=\"") + len("alt=\""):]
			alt = alt[:strings.Index(alt, "\"")]
			alt = html.UnescapeString(alt)

			comics = append(comics, comicUpdate{
				src.name,
				item.Title,
				imgUrl,
				alt,
				false,
			})
		}
	}

	return
}
