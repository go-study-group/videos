package main

import (
	"flag"
	"fmt"
	"github.com/go-study-group/videos"
	"log"
	"os"
)

func main() {
	fmt.Println("YouTube Uploader")

	fileName := flag.String("file", "", "video file")
	title := flag.String("title", "", "video title")
	channelId := flag.String("channelId", "", "YouTube Channel Id")
	clientId := flag.String("clientId", "", "client id")
	clientSecret := flag.String("clientSecret", "", "client secret")
	flag.Parse()

	fmt.Printf("Uploading %s to channel %s\n", *fileName, *channelId)

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("unable to load file: %v", err)
	}

	tokenStore := youtube.NewTokenFileStore("cache.token")
	service := youtube.GetService(*clientId, *clientSecret, tokenStore)
	videoId, err := youtube.UploadToYoutube(service, file, *channelId, *title)
	if err != nil {
		log.Fatalf("unable to upload: %v", err)
	}

	fmt.Printf("Your video is here: https://www.youtube.com/watch?v=%s\n", videoId)
}
