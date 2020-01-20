package main

import (
	"strings"
	"time"
	"fmt"
  "log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	r := gin.Default()
	r.GET("/schedule", func(c *gin.Context) {
		c.JSON(200, GetSchedule())
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

type ScheduleItem struct {
	Title string `json:"title"`
	Caster string `json:"caster"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	Cancelled bool `json:"cancelled"`
}


func GetSchedule() []ScheduleItem {
  // Request the HTML page.
  res, err := http.Get("https://www.bonjwa.de/programm")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
	}
	
	items := []ScheduleItem{}

  // Find the review items
  doc.Find(".stream-plan > table > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
		//content := s.Find("p")

		if (len(strings.TrimSpace(s.Text())) > 0) {
			title := strings.TrimSpace(s.Find("p:nth-child(even)").Text())
			caster := strings.TrimSpace(s.Find("p:nth-child(odd)").Text())
			startHour := s.AttrOr("data-hour-start", "-1")
			endHour := s.AttrOr("data-hour-end", "-1")
			date := s.AttrOr("data-date", "-1")
			
			classes := s.AttrOr("class", "");
			cancelled := strings.Contains(classes, "cancelled-streaming-slot")
			
			loc, _ := time.LoadLocation("Europe/Berlin")
			startDate, err := time.ParseInLocation("2006-1-02 15:04", fmt.Sprintf("%s %02s:00", date, startHour), loc)
			if err != nil { fmt.Println(err) }
			endDate,err := time.ParseInLocation("2006-1-02 15:04", fmt.Sprintf("%s %02s:00", date, endHour), loc)
			if err != nil { fmt.Println(err) }

			item := ScheduleItem{
				Title: title,
				Caster: caster,
				StartDate: startDate.UTC().Format(time.RFC3339),
				EndDate: endDate.UTC().Format(time.RFC3339),
				Cancelled: cancelled,
			}

			items = append(items, item)
		}

	})
	
	return items
}
