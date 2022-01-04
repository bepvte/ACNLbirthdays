package main

import (
	"github.com/gocolly/colly"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"time"
	"strings"
)

//type Character struct{
//	Name string `selector:"big big"`
//	Phrase string `selector:"tr:nth-child(6) td.roundytr:nth-child(2)"`
//	Personality string `selector:"tr:nth-child(2) td:nth-child(2) a"`
//	PicQuote string `selector:"small i"`
//	Birthday string `selector:"tr:nth-child(5) td:nth-child(2) a:nth-child(1)"`
//}
type Character struct{
		Name string `selector:"td:nth-child(1) b a"`
		Phrase string `selector:"td:nth-child(6) i"`
		Personality string `selector:"td:nth-child(3) a"`
		PicQuote string
		Birthday string `json:"-"`
		ImageUrl string
	}
const dateStyle = `January 2`
func main() {
	c := colly.NewCollector(colly.Async(false), colly.MaxDepth(2))

	characterMap := make(map[string]*Character)
	///table.roundy tr tbody tr
	c.OnHTML("div#mw-content-text table.roundy:nth-child(3) tr tbody tr", func(e *colly.HTMLElement) {
		var char Character
		if err := e.Unmarshal(&char); err != nil {
			log.Fatal(err)
		}
		if char.Name != "" {
			characterMap[char.Name] = &char
		}
		char.Birthday = strings.TrimSpace(e.DOM.Find("td:nth-child(5)").Contents().First().Text())
		e.Request.Visit(e.ChildAttr("td:nth-child(1) a", "href"))
	})

	c.OnHTML(".portable-infobox", func(e *colly.HTMLElement) {
		name := e.DOM.Find("h2:nth-child(1)").Text()
		if _, ok := characterMap[name]; !ok {
			characterMap[name] = &Character{}
		}
		characterMap[name].PicQuote = e.DOM.Find("figcaption").Text()
		imageraw, _ := e.DOM.Find("figure a.image").Attr("href")
		characterMap[name].ImageUrl = strings.Split(imageraw, "/revision")[0]
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})
	c.Visit("http://animalcrossing.wikia.com/wiki/List_of_Villagers_in_New_Leaf")
	c.Wait()

	characterDateMap := make(map[string][]*Character)
	for x, y := range characterMap {
		if y.Phrase == "" {
			fmt.Println("Missing phrase for ", x)
		}
		if y.Personality == "" {
			fmt.Println("Missing personality for ", x)
		}
		if y.PicQuote == "" {
			fmt.Println("Missing picture quote for ", x)
		}
		if _, err := time.Parse(dateStyle, y.Birthday); err != nil {
			fmt.Println("Invalid date for ", x, " error is ", err)
		}
		characterDateMap[y.Birthday] = append(characterDateMap[y.Birthday], y)
	}
	b, err := json.MarshalIndent(characterDateMap, "", "	")
	if err != nil {
		panic(err)
	}


	ioutil.WriteFile("results.json", b, 0777)
}
