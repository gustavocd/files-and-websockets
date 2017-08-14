package handlers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	upgrader  = websocket.Upgrader{}
	broadcast = make(chan []byte)
	clients   = make(map[*websocket.Conn]bool)
)

type data struct {
	File   string `json:"file"`
	Filter string `json:"filter"`
}

// Upload ...
func Upload(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	checkErr(err)
	defer ws.Close()

	clients[ws] = true
	var info data
	for {
		err := ws.ReadJSON(&info)
		if err != nil {
			logrus.Errorf("Could not read as json due to %v", err.Error())
			return
		}

		input := info.File[strings.IndexByte(info.File, ',')+1:]

		data, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			logrus.Errorf("Could not decode string into base64 due to %v", err.Error())
			return
		}

		err = ioutil.WriteFile("image.jpg", data, 0655)
		if err != nil {
			logrus.Errorf("Could not %v", err.Error())
			delete(clients, ws)
			break
		}

		f, err := imgio.Open("./image.jpg")
		checkErr(err)

		inverted := effect.Invert(f)
		err = imgio.Save("image_changed", inverted, imgio.JPEG)
		checkErr(err)

		response, err := ioutil.ReadFile("./image_changed.jpg")
		checkErr(err)
		broadcast <- response
	}
}

// HandleFile ...
func HandleFile() {
	for {
		data := <-broadcast
		for client := range clients {
			err := client.WriteMessage(2, data)
			if err != nil {
				logrus.Fatal(err.Error())
				client.Close()
				delete(clients, client)
				return
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}
}
