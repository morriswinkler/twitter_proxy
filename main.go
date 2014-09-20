package main

import (
	"log"
	"os"
	"time"
)

var logfile = "/var/log/twitter_proxy/twitter_proxy.log" // log file
var host = "127.0.0.1"                                   // host
var port = "8585"                                        // port to listen
var templateDir = "/etc/twitter_proxy/templates"         // template folder

var ScreenName = "laydrop"                                   // twitter user
var TweetCacheFile = "/tmp/twitter_proxy_user_timeline.json" // local cache copy
var credentials = "/etc/twitter_proxy/CREDENTIALS"           // crediantials file

var returnMax = 100                 // max tweets to return
var refreshEvery = time.Minute * 10 // refresh in minutes

func main() {

	// log file
	logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file ", logFile.Name, ":", err)
	}
	logInit(logFile, logFile, logFile, logFile)

	// start httpd
	var httpd Server
	go httpd.Init(host, port, templateDir)

	//
	ticker := time.NewTicker(refreshEvery)

	getTweets()
	for _ = range ticker.C {
		getTweets()
	}
}
