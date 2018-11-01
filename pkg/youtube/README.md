# YouTube Uploader

## Upload a video

```
go run cmd/upload/main.go -file ~/Desktop/sample.mp4 -channelId UCAgn_0SnWiW8Inu5fbpjLEg -title "Example Video" -clientId xxx -clientSecret yyy
```

## List my videos

https://developers.google.com/youtube/v3/guides/auth/installed-apps

```
go run pkg/prototype/my_videos.go -secrets ~/Downloads/cs.json
```
