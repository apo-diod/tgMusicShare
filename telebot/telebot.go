package telebot

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	telebotrpc "apo.share.song/RPC-interface/telebot"
	spotapi "apo.share.song/spotify-api"
	ytapi "apo.share.song/youtube-api"
	tbot "github.com/yanzay/tbot"
)

var logger *log.Logger
var botSrv *tbot.Server
var bot *tbot.Client

//TelegramServer ...
type TelegramServer struct {
	telebotrpc.TelebotServer
}

//InitLog ...
func InitLog(lg *log.Logger) {
	logger = lg
}

//InitBot ...
func InitBot(token string) error {
	botSrv = tbot.New(token)
	bot = botSrv.Client()
	botSrv.HandleMessage(".", func(m *tbot.Message) {
		if m.Text == "/start" {
			bot.SendMessage(m.Chat.ID, "Я могу найти превьюшку музыки(пока что, ну уж извините). Напишите название композиции(желательно с автором) и я дам вам послушать 30 секунд. Выгодная сделка, да?")
			return
		}
		if strings.HasPrefix(m.Text, "https://open.spotify.com/track/") || strings.Contains(m.Text, "https://link.tospotify.com/") {
			var name string
			var err error
			if strings.HasPrefix(m.Text, "https://open.spotify.com/track/") {
				name, err = spotapi.SearchForSongByURI(strings.Split(strings.TrimPrefix(m.Text, "https://open.spotify.com/track/"), "?")[0])
				if err != nil {
					bot.SendMessage(m.Chat.ID, err.Error())
					return
				}
			} else if strings.Contains(m.Text, "https://link.tospotify.com/") {
				query := strings.Split(m.Text, "/")
				name, err = spotapi.SearchForSongByURI(query[len(query)-1])
				if err != nil {
					bot.SendMessage(m.Chat.ID, err.Error())
					return
				}
			}
			ytlink, err := ytapi.SearchForSong(name)
			if err != nil {
				bot.SendMessage(m.Chat.ID, err.Error())
				return
			}
			bot.SendMessage(m.Chat.ID, ytlink)
		}
		if strings.HasPrefix(m.Text, "https://music.youtube.com/watch?v=") {
			query, err := ytapi.SongQueryByID(strings.TrimPrefix(m.Text, "https://music.youtube.com/watch?v="))
			if err != nil {
				bot.SendMessage(m.Chat.ID, err.Error())
				return
			}
			_, links, err := spotapi.SearchForSong(query)
			if err != nil {
				bot.SendMessage(m.Chat.ID, err.Error())
				return
			}
			bot.SendMessage(m.Chat.ID, links[1])
			return
		}
		logger.Println("Message from: " + m.Chat.ID)
		article, songURIs, err := spotapi.SearchForSong(m.Text)
		if err != nil {
			bot.SendMessage(m.Chat.ID, "error")
		}

		//Sending Preview from spotify
		if songURIs[0] != "" {
			f, err := os.OpenFile("temp/"+article+".mp3", os.O_CREATE|os.O_RDWR, 0755)
			if err != nil {
				logger.Println(err)
				return
			}
			defer f.Close()

			req, err := http.NewRequest("GET", songURIs[0], nil)
			if err != nil {
				logger.Println(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.Println(err)
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Println(err)
			}
			f.Write(b)
		}
		ytURI, err := ytapi.SearchForSong(article)
		if err != nil {
			ytURI = "Not found"
			logger.Println(err.Error())
		}
		bot.SendMessage(m.Chat.ID, article+"\nSpotify: "+songURIs[1]+"\nYouTube Music: "+ytURI)
		if songURIs[0] != "" {
			bot.SendAudioFile(m.Chat.ID, "temp/"+article+".mp3", tbot.OptCaption("Превью предоставлено @apo_share_bot"), tbot.OptTitle(article))
		}
	})
	botSrv.HandleInlineQuery(func(q *tbot.InlineQuery) {
		results := []tbot.InlineQueryResult{}
		article, songURIs, err := spotapi.SearchForSong(q.Query)
		if err == nil {
			res := tbot.InlineQueryResultArticle{}
			res.URL = songURIs[1]
			res.Title = "Share " + article + "!"
			res.ID = "1"
			res.HideURL = true
			message := tbot.InputTextMessageContent{MessageText: article + "\nSpotify: " + songURIs[1]}
			ytURI, err := ytapi.SearchForSong(article)
			if err != nil {
				ytURI = "Not found"
				logger.Println(err.Error())
			}
			message.MessageText += "\nYouTube Music: " + ytURI
			res.InputMessageContent = message
			res.Type = "article"
			results = append(results, res)
			if songURIs[0] != "" {
				resSound := tbot.InlineQueryResultAudio{}
				resSound.AudioURL = songURIs[0]
				resSound.Type = "audio"
				resSound.ID = "2"
				resSound.Title = "Share " + article + " with preview!"
				resSound.Caption = article + "\nSpotify: " + songURIs[1] + "\nYouTube Music: " + ytURI
				resSound.Performer = article
				results = append(results, resSound)
			}
		}
		bot.AnswerInlineQuery(q.ID, results)
	})
	var err error
	go func() {
		err = botSrv.Start()
	}()
	if err == nil {
		logger.Println("logged in as bot")
	} else {
		logger.Fatalln(err)
	}
	return err
}

//SendMessage ...
func SendMessage(uid string, text string) error {
	_, err := bot.SendMessage(uid, text)
	return err
}

//Send ...
func (ts TelegramServer) Send(ctx context.Context, in *telebotrpc.TelegramMessage) (*telebotrpc.TelegramSent, error) {
	err := SendMessage(strconv.Itoa(int(in.UID)), in.Message)
	tsent := &telebotrpc.TelegramSent{}
	if err != nil {
		tsent.Error = err.Error()
		tsent.Sent = false
	}
	tsent.Error = ""
	tsent.Sent = true
	return tsent, err
}

func (ts TelegramServer) mustEmbedUnimplementedTelebotServer() {
	return
}
