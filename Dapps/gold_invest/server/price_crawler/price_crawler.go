package price_crawler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	sourceURL = "https://giavang.org/"
)

func Crawl() (int, int, error) {
	res, err := http.Get(sourceURL)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, 0, fmt.Errorf(
			"request failed with status code: %d %s",
			res.StatusCode,
			res.Status,
		)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse HTML response: %w", err)
	}
	var buy, sell int
	doc.Find("span.gold-price").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			sell, err = strconv.Atoi(strings.ReplaceAll(strings.Split(s.Text(), " ")[0], ".", ""))
			if err != nil {
				err = fmt.Errorf(
					"no price data was found on the page, the website structure may have changed",
				)
			}
		}

		if i == 1 {
			buy, err = strconv.Atoi(strings.ReplaceAll(strings.Split(s.Text(), " ")[0], ".", ""))
			if err != nil {
				err = fmt.Errorf(
					"no price data was found on the page, the website structure may have changed",
				)
			}
		}
	})

	return buy, sell, err
}

// parsePrice removes commas and converts a price string to an int64.
func parsePrice(priceStr string) (int64, error) {
	// Remove the comma separator (e.g., "90,000" -> "90000")
	cleanedStr := strings.ReplaceAll(priceStr, ",", "")
	price, err := strconv.ParseInt(cleanedStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not convert '%s' to integer: %w", cleanedStr, err)
	}
	return price, nil
}
