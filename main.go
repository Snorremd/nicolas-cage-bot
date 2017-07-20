package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

var (
	netflixURL = "https://www.netflix.com/title/"
)

// Movie contains movie info
type Movie struct {
	Title       string `json:"show_title"`
	ReleaseYear string `json:"release_year"`
	ID          int    `json:"show_id"`
	Rating      string `json:"rating"`
	PosterURL   string `json:"poster"`
	Summary     string `json:"summary"`
}

func fetchMovies() ([]Movie, error) {
	// Create request for Nicols Cage movies from netflix roulette
	res, err := http.Get("http://netflixroulette.net/api/api.php?actor=Nicolas%20Cage")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Try to read the contents
	moviesRaw, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Try to decode JSON
	var movies []Movie
	err = json.Unmarshal(moviesRaw, &movies)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return movies, nil
}

func main() {

	api := slack.New(os.Getenv("SLACK_TOKEN"))

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			if strings.Contains(ev.Text, "movie") || strings.Contains(ev.Text, "film") {
				movies, err := fetchMovies()
				if err != nil {
					log.Println(err)
				} else {
					movie := movies[rand.Intn(len(movies))]
					params := slack.PostMessageParameters{}
					attachment := slack.Attachment{
						Title:     movie.Title,
						TitleLink: netflixURL + strconv.Itoa(movie.ID),
						Pretext:   "Here is something to watch!",
						Text:      movie.Summary,
						ImageURL:  movie.PosterURL,
						Footer:    "Results driven by https://netflixroulette.net",
					}
					params.Attachments = []slack.Attachment{attachment}
					rtm.PostMessage(ev.Channel, "", params)
				}

			}

		case *slack.PresenceChangeEvent:
			log.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			log.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())
			return

		case *slack.InvalidAuthEvent:
			log.Fatalln("Invalid credentials")

		default:
		}
	}
}
