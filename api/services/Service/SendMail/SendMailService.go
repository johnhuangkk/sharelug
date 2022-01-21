package SendMail

import (
	"api/services/util/curl"
	"api/services/util/log"
	"api/services/util/tools"
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/spf13/viper"
)

type SetMassage struct {
	Username   string
	Subject    string
	Title      string
	Body       string
	Content    string
	Code       string
	Link       string
	Subcontent string
}

type MailRequest struct {
	Name  string `json:"name"`
	From  string `json:"from"`
	To    string `json:"to"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Html  string `json:"html"`
}

type MailResponse struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

func (meg SetMassage) ToSendMail(to string) error  {
	if err := meg.SendAwsMail(to); err != nil {
		log.Error("send Aws mail Error", err)
		if err := meg.SendGmail(to); err != nil {
			log.Error("send Gmail mail Error", err)
		}
	}
	return nil
}

func (meg SetMassage) ToMail(to string) error  {
	if err := meg.SendAwsMail(to); err != nil {
		log.Error("send Aws mail Error", err)
		if err := meg.SendGmail(to); err != nil {
			log.Error("send Gmail mail Error", err)
			return err
		}
	}
	return nil
}



func (meg SetMassage) SendAwsMail(to string) error {
	ENV := viper.GetString("ENV")

	t, err := template.ParseFiles("views/mail/mail.html")
	if err != nil {
		log.Error("Get Template Error", err)
		return fmt.Errorf("系統錯誤")
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, meg); err != nil {
		log.Error("Execute Template Error", err)
		return fmt.Errorf("系統錯誤")
	}
	body := tpl.String()
	from := viper.GetString("MAIL_SERVER.FROM")
	smtpHost := viper.GetString("MAIL_SERVER.SMTP_HOST")

	mail := MailRequest{}
	if ENV != "prod" {
		mail.Name = "Check'Ne (DEV)"
	} else {
		mail.Name = "Check'Ne"
	}
	mail.From = from
	mail.To = to
	mail.Title = meg.Subject
	mail.Html = body

	resp, err := curl.PostJson(smtpHost, mail)
	if err != nil {
		log.Error("SMTP Error", err)
		return err
	}
	log.Debug("smtp Response", resp)
	if resp == nil {
		return fmt.Errorf("mail error")
	}
	dd := MailResponse{}
	_ = tools.JsonDecode(resp, &dd)
	log.Debug("smtp Response", dd)
	return nil
}

func (meg SetMassage) SendMail(to string) error {
	ENV := viper.GetString("ENV")

	t, err := template.ParseFiles("views/mail/mail.html")
	if err != nil {
		log.Error("Get Template Error", err)
		return fmt.Errorf("系統錯誤")
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, meg); err != nil {
		log.Error("Execute Template Error", err)
		return fmt.Errorf("系統錯誤")
	}
	body := tpl.String()

	from := viper.GetString("MAIL_SERVER.FROM")
	smtpHost := viper.GetString("MAIL_SERVER.SMTP_HOST")

	mail := MailRequest{}
	if ENV != "prod" {
		mail.Name = "Check'Ne (DEV)"
	} else {
		mail.Name = "Check'Ne"
	}
	mail.From = from
	mail.To = to
	mail.Title = meg.Subject
	mail.Html = body

	resp, err := curl.PostJson(smtpHost, mail)
	if err != nil {
		log.Error("SMTP Error", err)
		return err
	}
	log.Debug("smtp Response", resp)
	if resp == nil {
		return fmt.Errorf("mail error")
	}
	dd := MailResponse{}
	_ = tools.JsonDecode(resp, &dd)
	log.Debug("smtp Response", dd)
	return nil
}

func (meg SetMassage) SendGmail(to string) error {
	from := "email@checkne.com"
	password := ">71&T&R9hW"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("views/mail/mail.html")
	var body bytes.Buffer		
	Name := ""
	if viper.GetString("ENV") != "prod" {
		Name = "Check'Ne (DEV)"
	} else {
		Name = "Check'Ne"
	}

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject:%s\r\nFrom:%s<%s>\r\n%s\r\n", meg.Subject, Name, from, mimeHeaders)))
	if err := t.Execute(&body, meg); err != nil {
		log.Error("Execute Error", err)
		return err
	}

	// Sending email.
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body.Bytes()); err != nil {
		log.Error("SendMail error", err)
		return err
	}
	return nil
}
