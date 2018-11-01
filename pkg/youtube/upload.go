package youtube

import (
	"google.golang.org/api/youtube/v3"
	"io"
)

func UploadToYoutube(service *youtube.Service, file io.Reader, channelId string, title string) (string, error) {
	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			CategoryId: "22",
			ChannelId:  channelId,
			Title:      title,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "private",
		},
	}
	call := service.Videos.Insert("snippet,status", video)

	response, err := call.Media(file).Do()
	if err != nil {
		return "", err
	}

	return response.Id, nil
}
