package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	tokenFlag        = flag.String("token", "", "Twitter API token")
	secretFlag       = flag.String("secret", "", "Twitter API secret")
	accessTokenFlag  = flag.String("access-token", "", "Twitter API token")
	accessSecretFlag = flag.String("access-secret", "", "Twitter API secret")
	sinceFlag        = flag.Int("since", 7, "Days since you want tweets from")
	fromFlag         = flag.String("from", "wallixcom", "Twitter account from wich to retrieve tweets")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	anaconda.SetConsumerKey(*tokenFlag)
	anaconda.SetConsumerSecret(*secretFlag)

	api := anaconda.NewTwitterApi(*accessTokenFlag, *accessSecretFlag)

	v := url.Values{}
	v.Set("screen_name", *fromFlag)
	v.Set("exclude_replies", "true")
	v.Set("include_rts", "true")
	v.Set("count", "100")

	allTweets, err := api.GetUserTimeline(v)
	if err != nil {
		log.Fatal(err)
	}

	var tweets [][2]string

	for _, tweet := range allTweets {
		date, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}

		daySince := time.Since(date).Hours() / 24
		if daySince > float64(*sinceFlag) {
			break
		}

		weekDay := fmt.Sprintf("%s %d", date.Weekday(), date.Day())
		tweets = append(tweets, [2]string{weekDay, tweet.Text})
	}

	w := csv.NewWriter(os.Stdout)

	w.Write([]string{fmt.Sprintf("Last %d days", *sinceFlag), "Tweets"})

	for _, tweet := range tweets {
		w.Write(tweet[:])
	}

	w.Flush()
}
