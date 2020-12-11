package database

import (
	"log"
	"strconv"
	"strings"

	badger "github.com/dgraph-io/badger"
)

//TOKENYT ...
var TOKENYT = 1

//TOKENSF ...
var TOKENSF = 2

var logger *log.Logger
var db *badger.DB

//User ...
type User struct {
	UID          int32
	YoutubeToken string
	SpotifyToken string
}

//InitLog ...
func InitLog(lg *log.Logger) {
	logger = lg
}

//InitDB ...
func InitDB() error {
	db, err := badger.Open(badger.DefaultOptions("/badger"))
	if err != nil {
		logger.Println(err)
		return err
	}
	logger.Println(db.Opts().Dir)
	return nil
}

//AddUser ...
func AddUser(usr User) error {
	value := usr.tokensToString()
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(strconv.Itoa(int(usr.UID))), []byte(value))
		return err
	})
	return err
}

//UpdateUser ...
func UpdateUser(uid int32, tokenType int, token string) error {
	usr, err := GetUser(uid)
	if err != nil {
		return err
	}
	switch tokenType {
	case TOKENYT:
		usr.YoutubeToken = token
		break
	case TOKENSF:
		usr.SpotifyToken = token
		break
	}
	err = AddUser(usr)
	return err
}

//GetUser ...
func GetUser(uid int32) (User, error) {
	user := User{UID: uid}
	result := ""
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(int(uid))))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			result = string(val)
			return nil
		})
		return err
	})
	if err != nil {
		return user, err
	}
	tokenMap := parseUserResponse(result)
	for k, v := range tokenMap {
		switch strings.ToLower(k) {
		case "youtube":
			user.YoutubeToken = v
			break
		case "spotify":
			user.SpotifyToken = v
			break
		}
	}
	return user, nil
}

func parseUserResponse(response string) map[string]string {
	result := map[string]string{}
	temp := strings.Split(strings.TrimSpace(response), ",")
	for _, item := range temp {
		kvPair := strings.Split(item, ":")
		result[kvPair[0]] = kvPair[1]
	}
	return result
}

func (usr User) tokensToString() string {
	result := ""
	if usr.SpotifyToken != "" {
		result += "spotify:" + usr.SpotifyToken
	}
	if usr.YoutubeToken != "" {
		result += "youtube:" + usr.YoutubeToken
	}
	return result
}
