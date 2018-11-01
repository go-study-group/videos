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
	channelId := flag.String("channelId", "", "YouTube Channel Id")
	flag.Parse()

	fmt.Printf("Uploading %s to channel %s\n", *fileName, *channelId)

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("unable to load file: %v", err)
	}

	videoId, err := youtube.UploadToYoutube(file, channelId)
	if err != nil {
		log.Fatalf("unable to upload: %v", err)
	}

	fmt.Printf("Your video is here: https://www.youtube.com/watch?v=%s\n", videoId)
}
