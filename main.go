package main

import (
	"bytes"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
)

var (
	//go:embed config.json
	configBytes []byte
	//go:embed assets
	staticAssets embed.FS
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	From     string `json:"from"`
	Security string `json:"security"`
}

func attachInline(
	e *email.Email, data []byte, filename string, contentType string,
) error {
	attachment, err := e.Attach(
		bytes.NewReader(data), filename, contentType,
	)

	if err != nil {
		return err
	}

	attachment.HTMLRelated = true

	return nil
}

func sendMail(
	emailAddress string, recipientName string,
	frontImagePath string, backImagePath string,
) {
	var config Config

	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalln("failed to unmarshal config: " + err.Error())
	}

	log.Println("from: " + config.From)
	log.Println("email address: " + emailAddress)
	log.Println("recipient: " + recipientName)
	log.Println("front image path: " + frontImagePath)
	log.Println("back image path: " + backImagePath)

	frontImage, err := os.ReadFile(frontImagePath)
	if err != nil {
		log.Fatalln("failed to read " + frontImagePath)
	}

	backImage, err := os.ReadFile(backImagePath)
	if err != nil {
		log.Fatalln("failed to read " + backImagePath)
	}

	htmlBytes, err := staticAssets.ReadFile("assets/showPostcard.html")
	if err != nil {
		log.Fatalln("failed to read showPostcard.html")
	}

	html := string(htmlBytes)
	html = strings.ReplaceAll(
		html, "{{name}}", template.HTMLEscapeString(recipientName),
	)

	// emailAddresses := strings.Split(emailAddress, "^")
	// for i := range len(emailAddresses) {
	// 	emailAddresses[i] = strings.TrimSpace(emailAddresses[i])
	// }

	e := email.NewEmail()
	e.From = config.From
	e.To = []string{emailAddress}
	e.Subject = "Postcard from " + recipientName
	e.HTML = []byte(html)

	staticAssetEntries, err := staticAssets.ReadDir("assets")
	if err != nil {
		log.Fatalln("failed to read static assets")
	}

	for _, entry := range staticAssetEntries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".jpg") {
			continue
		}

		data, err := staticAssets.ReadFile("assets/" + filename)
		if err != nil {
			log.Fatalf("failed to read %s\n", filename)
		}

		attachInline(e, data, filename, "image/jpeg")
	}

	attachInline(e, frontImage, "front.jpg", "image/jpeg")
	attachInline(e, backImage, "back.jpg", "image/jpeg")

	auth := smtp.PlainAuth(
		"", config.User, config.Pass, config.Host,
	)

	if config.Security == "plain" {
		err = e.Send(
			fmt.Sprintf("%s:%d", config.Host, config.Port), auth,
		)
	} else if config.Security == "starttls" {
		err = e.SendWithStartTLS(
			fmt.Sprintf("%s:%d", config.Host, config.Port), auth,
			&tls.Config{ServerName: config.Host},
		)
	} else {
		// default to tls
		err = e.SendWithTLS(
			fmt.Sprintf("%s:%d", config.Host, config.Port), auth,
			&tls.Config{ServerName: config.Host},
		)
	}

	if err != nil {
		log.Fatalln("failed to send email: " + err.Error())
	}

	log.Println("sent!")
}

func main() {
	if len(os.Args) < 8 {
		log.Fatalln("needs more arguments")
	}

	// gamePath := os.Args[1]
	// website := os.Args[2]
	// websiteEmail := os.Args[3]
	emailAddress := os.Args[4]
	recipientName := os.Args[5]
	frontImagePath := os.Args[6]
	backImagePath := os.Args[7]
	// language := os.Args[8]

	sendMail(emailAddress, recipientName, frontImagePath, backImagePath)
}
