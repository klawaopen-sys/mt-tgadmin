package app

import (
	"encoding/base64"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mitoteam/goapp"
)

func apiCheckAuth(r *goapp.ApiRequest) bool {
	if r.SessionGet("auth") == true {
		return true
	} else {
		r.SetErrorStatus("Auth Required")
		return false
	}
}

func Api_HealthCheck(r *goapp.ApiRequest) error {
	r.SetOutData("auth", apiCheckAuth(r))
	r.SetOkStatus("API works: " + App.AppName)

	return nil
}

func Api_Password(r *goapp.ApiRequest) error {
	if r.GetInData("password") == Settings.GuiPassword {
		// Set user as authenticated
		r.SessionSet("auth", true)

		r.Session().Save()
		r.SetOkStatus("You are authorized")
	} else {
		r.SetErrorStatus("Wrong password!")
	}
	return nil
}

func Api_Logout(r *goapp.ApiRequest) error {
	if !apiCheckAuth(r) {
		return nil
	}

	r.SessionClear()
	r.SetOkStatus("Good bye!")

	return nil
}

func Api_Say(r *goapp.ApiRequest) error {
	if !apiCheckAuth(r) {
		return nil
	}

	text := r.GetInData("message")
	text = PrepareTelegramHtml(text)

	photoURL := r.GetInData("photo_url")
	photoBase64 := r.GetInData("photo_base64")
	photoName := r.GetInData("photo_name")

	buttonText := r.GetInData("button_text")
	buttonURL := r.GetInData("button_url")

	var replyMarkup interface{} = nil

	if buttonText != "" && buttonURL != "" {
		replyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(buttonText, buttonURL),
			),
		)
	}

	if photoBase64 != "" {
		if strings.Contains(photoBase64, ",") {
			parts := strings.SplitN(photoBase64, ",", 2)
			photoBase64 = parts[1]
		}

		photoBytes, err := base64.StdEncoding.DecodeString(photoBase64)
		if err != nil {
			r.SetErrorStatus("Ошибка чтения фото: " + err.Error())
			return nil
		}

		if photoName == "" {
			photoName = "photo.jpg"
		}

		photo := tgbotapi.NewPhoto(Settings.BotChatID, tgbotapi.FileBytes{
			Name:  photoName,
			Bytes: photoBytes,
		})

		photo.Caption = text
		photo.ParseMode = "HTML"

		if reply_to := r.GetInDataInt("reply_to", 0); reply_to > 0 {
			photo.ReplyToMessageID = reply_to
		}

		if r.GetInDataInt("silent", 0) != 0 {
			photo.DisableNotification = true
		}

		if replyMarkup != nil {
			photo.ReplyMarkup = replyMarkup
		}

		if _, err := tgBot.Send(photo); err != nil {
			r.SetErrorStatus("Ошибка отправки фото: " + err.Error())
			return nil
		}

		return nil
	}

	if photoURL != "" {
		photo := tgbotapi.NewPhoto(Settings.BotChatID, tgbotapi.FileURL(photoURL))
		photo.Caption = text
		photo.ParseMode = "HTML"

		if reply_to := r.GetInDataInt("reply_to", 0); reply_to > 0 {
			photo.ReplyToMessageID = reply_to
		}

		if r.GetInDataInt("silent", 0) != 0 {
			photo.DisableNotification = true
		}

		if replyMarkup != nil {
			photo.ReplyMarkup = replyMarkup
		}

		if _, err := tgBot.Send(photo); err != nil {
			r.SetErrorStatus("Ошибка отправки фото по ссылке: " + err.Error())
			return nil
		}

		return nil
	}

	msg := tgbotapi.NewMessage(Settings.BotChatID, text)
	msg.ParseMode = "HTML"

	if reply_to := r.GetInDataInt("reply_to", 0); reply_to > 0 {
		msg.ReplyToMessageID = reply_to
	}

	if r.GetInDataInt("silent", 0) != 0 {
		msg.DisableNotification = true
	}

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	if _, err := tgBot.Send(msg); err != nil {
		r.SetErrorStatus("Ошибка отправки сообщения: " + err.Error())
		return nil
	}

	return nil
}

func Api_ListMessages(r *goapp.ApiRequest) error {
	if !apiCheckAuth(r) {
		return nil
	}

	updates_config := tgbotapi.NewUpdate(0)
	updates_config.Timeout = 1
	updates_config.Limit = 100
	updates_config.Offset = -100
	updates_config.AllowedUpdates = []string{"message"}

	updates_list, err := tgBot.GetUpdates(updates_config)
	if err != nil {
		r.SetErrorStatus(err.Error())
		return nil
	}

	list := make([]*apiMessage, 0, len(updates_list))

	//from end to beginning
	for i := len(updates_list) - 1; i >= 0; i-- {
		update := updates_list[i]

		if update.Message.Chat.ID != Settings.BotChatID {
			//skip messages from other channels or chats
			continue
		}

		m := &apiMessage{}
		m.Message = update.Message.Text
		m.MessageId = update.Message.MessageID
		m.User = update.Message.From.FirstName + " " + update.Message.From.LastName + " = @" + update.Message.From.UserName
		m.Date = time.Unix(int64(update.Message.Date), 0).Format("2006-01-02 15:04:05")

		list = append(list, m)
	}

	//log.Println("Updates count: ", len(list))

	r.SetOutData("list", list)
	return nil
}

func Api_AiPost(r *goapp.ApiRequest) error {

	if r.GetInData("api_key") != "klava_super_secret_ai_post_key_2026_telegram_autopost" {
		r.SetErrorStatus("Invalid API key")
		return nil
	}

	text := r.GetInData("message")
	text = PrepareTelegramHtml(text)

	photoURL := r.GetInData("photo_url")

	buttonText := r.GetInData("button_text")
	buttonURL := r.GetInData("button_url")

	var replyMarkup interface{} = nil

	if buttonText != "" && buttonURL != "" {
		replyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(buttonText, buttonURL),
			),
		)
	}

	if photoURL != "" {

		photo := tgbotapi.NewPhoto(Settings.BotChatID, tgbotapi.FileURL(photoURL))
		photo.Caption = text
		photo.ParseMode = "HTML"

		if replyMarkup != nil {
			photo.ReplyMarkup = replyMarkup
		}

		if _, err := tgBot.Send(photo); err != nil {
			r.SetErrorStatus("Error sending photo: " + err.Error())
			return nil
		}

		r.SetOkStatus("Photo sent")
		return nil
	}

	msg := tgbotapi.NewMessage(Settings.BotChatID, text)
	msg.ParseMode = "HTML"

	if replyMarkup != nil {
		msg.ReplyMarkup = replyMarkup
	}

	if _, err := tgBot.Send(msg); err != nil {
		r.SetErrorStatus("Error sending message: " + err.Error())
		return nil
	}

	r.SetOkStatus("Message sent")

	return nil
}
