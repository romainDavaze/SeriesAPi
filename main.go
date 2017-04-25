package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

type Series struct {
	Id	int	`json:"id"`
	Type	string	`json:"type"`
	Seasons	int	`json:"seasons"`
	Author	string	`json:"author"`
	Title	string	`json:"title"`
}

type SoManySeries []Series

var seriesObj SoManySeries

func main() {

	seriesObj = InitSeries()

	router := mux.NewRouter().StrictSlash(true)
	//router.PathPrefix("/series")

	router.HandleFunc("/series", SoManySeriesHandler).Methods("GET","POST")
	router.HandleFunc("/series/{id}", SeriesHandler).Methods("GET","PUT","DELETE")


	log.Fatal(http.ListenAndServe(":8080", router))
}

func SoManySeriesHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == "GET" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(seriesObj)
	} else if method == "POST" {
		if CheckAuthentication(r) {
			decoder := json.NewDecoder(r.Body)
			var series Series
			err := decoder.Decode(&series)
			if err != nil {
				log.Fatal(err)
			}
			series.Id = len(seriesObj)+1
			seriesObj = append(seriesObj,series)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}
}

func SeriesHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	method := r.Method
	switch method {
	case "GET":
		series  := Series{}
		for _, s := range seriesObj {
			if strconv.Itoa(s.Id) == vars["id"] {
				series = s
			}
		}

		if (Series{}) == series {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(series)
		}

	case "PUT":
		if CheckAuthentication(r){
			decoder := json.NewDecoder(r.Body)
			var series Series
			err := decoder.Decode(&series)
			if err != nil {
				log.Fatal(err)
			}

			seriesID  := -1
			for i, s := range seriesObj {
				if strconv.Itoa(s.Id) == vars["id"] {
					seriesID = i
				}
			}

			if seriesID == -1 {
				series.Id,_ = strconv.Atoi(vars["id"])
				seriesObj = append(seriesObj,series)
				w.WriteHeader(http.StatusCreated)
			} else {
				series.Id = seriesObj[seriesID].Id
				seriesObj[seriesID] = series
				w.WriteHeader(http.StatusCreated)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}


	case "DELETE": w.WriteHeader(http.StatusNotImplemented)
	}

}

func InitSeries() SoManySeries {
 return SoManySeries{
	 Series{Id: 1, Title: "Game Of Thrones", Type: "Heroic fantasy", Author: "Georges RR Martin", Seasons: 6},
	 Series{Id: 2, Title: "Breaking Bad", Type: "Thriller", Author: "Vince Gilligan", Seasons: 5},
	 Series{Id: 3, Title: "Daredevil", Type: "Superheroes", Author: "Drew Goddard", Seasons: 2},
 }
}

func CheckAuthentication(r *http.Request) bool{

	header := r.Header
	authorization := header.Get("Authorization")

	return authorization == "toto"
}