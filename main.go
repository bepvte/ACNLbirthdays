package main

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/dghubble/oauth1"
	"github.com/lucasjones/go-twitter/twitter"
	"time"
	"net/http"
	"log"
	"io"
)

const dateStyle = `January 2`
type Character struct{
	Name string `selector:"td:nth-child(1) b a"`
	Phrase string `selector:"td:nth-child(6) i"`
	Personality string `selector:"td:nth-child(3) a"`
	PicQuote string
	ImageUrl string `selector:"img" attr:"data-src"`
}
func main() {
	b, err := ioutil.ReadFile("acnlwikiscraper/results.json")
	if err != nil {
		panic(err)
	}
	characterMap := make(map[string][]Character)
	if err := json.Unmarshal(b, &characterMap); err != nil {
		panic(err)
	}


	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	for _,y := range characterMap[time.Now().Format(dateStyle)]{
		respe, err := http.Get(y.ImageUrl)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.OpenFile("image.png", os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(f, respe.Body)
		if err != nil {
			log.Fatal(err)
		}
		respe.Body.Close()

		m, _, err := client.Media.UploadFile("image.png")
		if err != nil {
			log.Fatal(err)
		}
		os.Remove("image.png")

		_, resp, err := client.Statuses.Update(fmt.Sprintf("Today is %v's birthday! Their personality is %v and their favorite phrase is %v. \n\"%v\"", y.Name, y.Personality, y.Phrase, y.PicQuote), &twitter.StatusUpdateParams{MediaIds: []int64{m.MediaID}})
		fmt.Println(resp, err)
	}
}

