package youtubeAPI

import (
	"log"
	"net/http"

	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var logger *log.Logger

var YoutubeToken = "###"
var service *youtube.Service

func InitLog(lg *log.Logger) {
	logger = lg
	logger.Println("Log set")
}

func InitYoutubeAPI() error {
	logger.Println("Initialising API...")
	client := &http.Client{
		Transport: &transport.APIKey{Key: YoutubeToken},
	}
	var err error
	service, err = youtube.New(client)
	return err
}

func SearchForSong(query string) (string, error) {
	call := service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		return "", err
	}
	if len(resp.Items) == 0 {
		return "Not found", err
	}
	res := resp.Items[0]
	songID := res.Id.VideoId
	return "https://music.youtube.com/watch?v=" + songID, err
}

func SongQueryByID(id string) (string, error) {
	call := service.Search.List([]string{"id", "snippet"}).Q(id).MaxResults(1)
	resp, err := call.Do()
	if err != nil {
		logger.Println(err)
		return "", err
	}
	if len(resp.Items) == 0 {
		return "", err
	}
	res := resp.Items[0]
	article := res.Snippet.Title
	return article, err
}
