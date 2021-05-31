package api

import (
    "encoding/json"
    //"fmt"
    "log"
    "io/ioutil"
    "net/http"
    //"database/sql"
    "os"
	//"errors"
	"net/url"
	"strconv"

    _ "github.com/lib/pq"
)

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// HandleTelegramWebHook sends a message back to the chat by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {

	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

    content := ""

    log.Printf("Text : update.Message.Text")
    if (update.Message.Text == "?"){
        content = "Help :\n" +
                "\t? - For Help\n" +
                "\tNIM= - For Deadlines&ScheduleLink Link (Write your NIM after NIM== e.g. NIM=1900077)\n"
    } else if (update.Message.Text[0:4] == "NIM="){
        if len(update.Message.Text) > 11{
        } else {
            content = "Invalid NIM\n"
        }
    } else {
        content = "Unrecognized Command.\n" +
                "Help :\n" +
                "\t? - For Help\n" +
                "\tNIM= - For Deadlines&ScheduleLink Link (Write your NIM after NIM== e.g. NIM=1900077)\n"
    }

	// Send the content back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, content)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, reponse body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("content %s successfuly distributed to chat id %d", content, update.Message.Chat.Id)
	}
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("telegram_token") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}
