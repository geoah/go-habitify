package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/geoah/go-habitify"
	"github.com/jinzhu/now"
)

//go:embed assets
var assets embed.FS

func main() {
	apiKey := os.Getenv("HABITIFY_API_KEY")
	if apiKey == "" {
		log.Fatal("missing HABITIFY_API_KEY env var")
	}

	client := habitify.New(apiKey)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(assets, "assets/calendar.html")
		if err != nil {
			log.Fatal(err)
		}
		habits, err := client.GetHabits()
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(w, habits); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		logs, err := client.GetHabitLogs(
			r.URL.Query().Get("habit_id"),
			now.BeginningOfYear(),
			now.EndOfMonth(),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		dataLogs := map[string]float64{}
		for _, log := range logs {
			dataLogs[fmt.Sprintf("%d", log.TargetDate.Unix())] += log.Value
		}
		res, err := json.Marshal(dataLogs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(res)
		w.Header().Set("Content-Type", "application/json")
	})

	http.HandleFunc("/logs/up", func(w http.ResponseWriter, r *http.Request) {
		targetDate, _ := time.Parse(time.RFC3339, r.URL.Query().Get("target_date"))
		_, err := client.AddHabitLogs(
			r.URL.Query().Get("habit_id"),
			targetDate,
			r.URL.Query().Get("unit"),
			"1",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
