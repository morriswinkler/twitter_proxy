package main

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var results *twittergo.Timeline
var cache *twittergo.Timeline

func LoadCredentials() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile(credentials)
	if err != nil {
		return
	}
	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return
}

func getTweets() {
	var (
		err    error
		client *twittergo.Client
		req    *http.Request
		resp   *twittergo.APIResponse
		max_id uint64
		out    *os.File
		query  url.Values

		text []byte
	)
	if client, err = LoadCredentials(); err != nil {
		ERROR.Println("Could not parse CREDENTIALS file:", err)
		os.Exit(1)
	}
	if out, err = os.Create(TweetCacheFile); err != nil {
		ERROR.Println("Could not create output file:", TweetCacheFile)
		os.Exit(1)
	}
	defer out.Close()
	const (
		count   int = 100
		urltmpl     = "/1.1/statuses/user_timeline.json?%v"
		minwait     = time.Duration(10) * time.Second
	)
	query = url.Values{}
	query.Set("count", fmt.Sprintf("%v", count))
	query.Set("screen_name", ScreenName)
	total := 0

	if max_id != 0 {
		query.Set("max_id", fmt.Sprintf("%v", max_id))
	}
	endpoint := fmt.Sprintf(urltmpl, query.Encode())
	if req, err = http.NewRequest("GET", endpoint, nil); err != nil {
		ERROR.Println("Could not parse request:", err)
	}
	if resp, err = client.SendRequest(req); err != nil {
		ERROR.Println("Could not send request:", err)
	}

	results = &twittergo.Timeline{}
	if err = resp.Parse(results); err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			dur := rle.Reset.Sub(time.Now()) + time.Second
			if dur < minwait {
				// Don't wait less than minwait.
				dur = minwait
			}
			WARNING.Println("Rate limited. Reset at", rle.Reset, "Waiting for", dur)
			time.Sleep(dur)

		} else {
			ERROR.Println("Problem parsing response:", err)
		}
	}
	batch := len(*results)
	if batch == 0 {
		INFO.Println("No more results, end of timeline.")
	} else {
		cache = results
	}
	for _, tweet := range *results {
		if text, err = json.Marshal(tweet); err != nil {
			ERROR.Println("Could not encode Tweet:", err)
			os.Exit(1)
		}
		out.Write(text)
		out.Write([]byte("\n"))
		max_id = tweet.Id() - 1
		total += 1
	}
	INFO.Println("Got", batch, "Tweets")
	if resp.HasRateLimit() {
		INFO.Println(resp.RateLimitRemaining(), " calls available")
	}

	INFO.Println("--------------------------------------------------------")
	INFO.Println("Wrote", total, "Tweets to", TweetCacheFile)
}
