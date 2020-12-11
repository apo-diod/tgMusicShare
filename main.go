package main

import (
	"log"

	"apo.share.song/logging"
	spotapi "apo.share.song/spotify-api"
	"apo.share.song/telebot"
	ytapi "apo.share.song/youtube-api"
)

func main() {
	log.Println("Setting loggers")
	//Log creations
	logging.InitLog("log.log")
	//databaseLog := log.New(logging.StandardLog, "[DB]", log.LstdFlags)
	//rpcLog := log.New(logging.StandardLog, "[gRPC]", log.LstdFlags)
	spotifyLog := log.New(logging.StandardLog, "[SpotifyApi]", log.LstdFlags)
	telebotLog := log.New(logging.StandardLog, "[TelegramBot]", log.LstdFlags)
	youtubeLog := log.New(logging.StandardLog, "[YouTube]", log.LstdFlags)
	mainLog := log.New(logging.StandardLog, "[Main]", log.LstdFlags)

	mainLog.Println("Services initialisation")

	telebot.InitLog(telebotLog)
	telebot.InitBot("###")
	spotapi.InitLog(spotifyLog)
	spotapi.InitSpotifyAPI()
	ytapi.InitLog(youtubeLog)
	err := ytapi.InitYoutubeAPI()
	if err != nil {
		mainLog.Println("Couldn't set youtube API")
	}
	ch := make(chan bool, 1)
	<-ch

}
