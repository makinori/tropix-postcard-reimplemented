package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

//go:embed config.json
var configBytes []byte

//go:embed html/showPostcard.html
var showPostcard string

//go:embed html/postcardBorder_01.jpg
var postcardBorder_01 []byte

//go:embed html/postcardBorder_02.jpg
var postcardBorder_02 []byte

//go:embed html/postcardBorder_03.jpg
var postcardBorder_03 []byte

//go:embed html/postcardBorder_04.jpg
var postcardBorder_04 []byte

//go:embed html/postcardBorder_06.jpg
var postcardBorder_06 []byte

//go:embed html/postcardBorder_07.jpg
var postcardBorder_07 []byte

//go:embed html/postcardBorder_08.jpg
var postcardBorder_08 []byte

//go:embed html/postcardBorder_09.jpg
var postcardBorder_09 []byte

//go:embed html/tropixBamboo.jpg
var tropixBamboo []byte

//go:embed html/tropixTitle.jpg
var tropixTitle []byte

type Config struct {
	Email struct {
		Host string
		Port string
		User string
		Pass string
		From string
	}
}

func getConfig() Config {
	var config Config
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func sendMail(emailTo string, name string, frontImagePath string, backImagePath string) {
	config := getConfig()

	fmt.Println("Email: " + emailTo)
	fmt.Println("Name: " + name)
	fmt.Println("Front image path: " + frontImagePath)
	fmt.Println("Back image path: " + backImagePath)

	e := email.NewEmail()
	e.From = config.Email.From
	e.To = []string{emailTo}
	e.Subject = "Postcard from " + name
	e.HTML = []byte(strings.ReplaceAll(showPostcard, "{{name}}", name))

	e.Attach(bytes.NewReader(postcardBorder_01), "postcardBorder_01.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_02), "postcardBorder_02.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_03), "postcardBorder_03.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_04), "postcardBorder_04.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_06), "postcardBorder_06.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_07), "postcardBorder_07.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_08), "postcardBorder_08.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_09), "postcardBorder_09.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(postcardBorder_09), "postcardBorder_09.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(tropixBamboo), "tropixBamboo.jpg", "image/jpeg")
	e.Attach(bytes.NewReader(tropixTitle), "tropixTitle.jpg", "image/jpeg")

	front, err := os.ReadFile(frontImagePath)
	if err != nil {
		panic(err)
	}
	e.Attach(bytes.NewReader(front), "front.jpg", "image/jpeg")

	back, err := os.ReadFile(backImagePath)
	if err != nil {
		panic(err)
	}
	e.Attach(bytes.NewReader(back), "back.jpg", "image/jpeg")

	e.Send(
		config.Email.Host+":"+config.Email.Port,
		smtp.PlainAuth("", config.Email.User, config.Email.Pass, config.Email.Host),
	)
}

func main() {
	if len(os.Args) < 8 {
		fmt.Println("Needs more arguments")
		os.Exit(1)
	}

	// game_path := os.Args[1]
	// website := os.Args[2]
	// website_email := os.Args[3]
	fromEmail := os.Args[4]
	toName := os.Args[5]
	frontImagePath := os.Args[6]
	backImagePath := os.Args[7]
	// language := os.Args[8]

	sendMail(fromEmail, toName, frontImagePath, backImagePath)
}
