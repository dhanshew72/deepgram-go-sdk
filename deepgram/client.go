package deepgram

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// GetClient creates a new websocket client connection
func GetClient(header http.Header, u url.URL) (*websocket.Conn, *http.Response, error) {
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("error creating a new client connection: %s", err.Error())
		if resp != nil {
			log.Printf("handshake failed with status %s", resp.Status)
		} else {
			log.Printf("handshake failed with no response")
		}
		return retry(header, u)
	}
	return c, resp, nil
}

// retry creates a new websocket client connection with a timed retry logic
func retry(header http.Header, u url.URL) (*websocket.Conn, *http.Response, error) {
	log.Printf("Attemping to open a new client connection with a retry...")
	retryCount := 0
	for retryCount < 4 {
		// Every 5 seconds the transcriber will attempt to create a new websocket
		time.Sleep(time.Second * 5)
		c, res, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			log.Printf("Failed to to open an new client connection, retrying in 5 seconds...")
			retryCount++
		} else {
			return c, res, err
		}
	}
	return nil, nil, errors.New("failed to open a new client connection after 3 retries")
}
