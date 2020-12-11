package spotApi

import (
	"context"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var logger *log.Logger
var client spotify.Client

//InitLog ...
func InitLog(lg *log.Logger) {
	logger = lg
}

//InitSpotifyAPI ...
func InitSpotifyAPI() {
	config := &clientcredentials.Config{
		ClientID:     "###",
		ClientSecret: "###",
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client = spotify.Authenticator{}.NewClient(token)
}

//SearchForSong ...
func SearchForSong(query string) (string, []string, error) {
	result := []string{}
	logger.Println("Searching for " + query)
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return "", result, err
	}
	if results.Tracks != nil {
		if len(results.Tracks.Tracks) == 0 {
			return query, []string{"Not Found", "Not Found"}, nil
		}
		track := results.Tracks.Tracks[0]
		result = append(result, track.PreviewURL)
		result = append(result, "https://open.spotify.com/track/"+string(track.ID))
		logger.Println(track.ID)
		return track.Artists[0].Name + "-" + track.Name, result, err
	}
	return "", result, err
}

//SearchForSongByURI ...
func SearchForSongByURI(URI string) (string, error) {
	logger.Println("searching for song by URI: " + URI)
	track, err := client.GetTrack(spotify.ID(URI))
	if err != nil {
		return "", err
	}
	return track.Artists[0].Name + "-" + track.Name, err
}
