package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// list of champions and number in pool by its cost
var championsByCost [5][]string

// number of champions per level
var championPool = [5]int{29, 22, 18, 12, 10}

// list of drop rates by level
// level is first array position + 1, tier is second array position + 1
var dropRates [9][5]float64

// results counter
type rollResult struct {
	Champion    string
	Cost        int
	Appearances int
}

var rollResults = make(map[string]rollResult)

// data structure for champions.json file
type Champion struct {
	Champion   string
	ChampionId string
	Cost       int
}

// provides the rates to the DropRates global variable
func assignDropRates() {
	file, err := os.Open("./data/dropRates.json")

	if err == nil {
		bv, _ := ioutil.ReadAll(file)
		json.Unmarshal(bv, &dropRates)
	} else { // create file if it does not exist
		dropRates[0] = [5]float64{1.00, 0.00, 0.00, 0.00, 0.00}
		dropRates[1] = [5]float64{1.00, 0.00, 0.00, 0.00, 0.00}
		dropRates[2] = [5]float64{0.75, 0.25, 0.00, 0.00, 0.00}
		dropRates[3] = [5]float64{0.55, 0.30, 0.15, 0.00, 0.00}
		dropRates[4] = [5]float64{0.45, 0.33, 0.20, 0.02, 0.00}
		dropRates[5] = [5]float64{0.25, 0.40, 0.30, 0.05, 0.00}
		dropRates[6] = [5]float64{0.19, 0.30, 0.35, 0.15, 0.01}
		dropRates[7] = [5]float64{0.15, 0.20, 0.35, 0.25, 0.05}
		dropRates[8] = [5]float64{0.10, 0.15, 0.30, 0.30, 0.15}

		for i, s := range dropRates {
			for i2 := range s {
				if i2 != 0 {
					dropRates[i][i2] = math.Round((dropRates[i][i2]+dropRates[i][i2-1])*100) / 100
				}
			}
		}

		file, _ := json.Marshal(dropRates)
		_ = ioutil.WriteFile("./data/dropRates.json", file, 0644)
	}
	defer file.Close()
}

// returns a list of champion metadata
func openChampionsData() []Champion {
	championsJSON, err := os.Open("./data/champions.json")

	if err != nil {
		fmt.Println(err)
	}

	defer championsJSON.Close()

	byteValue, _ := ioutil.ReadAll(championsJSON)

	var champions []Champion

	json.Unmarshal(byteValue, &champions)

	return champions
}

// returns a list of champions by cost and the number in pool.
func getchampionsByCost() {
	file, err := os.Open("./data/championsByCost.json")
	if err == nil {
		bv, _ := ioutil.ReadAll(file)
		json.Unmarshal(bv, &championsByCost)
	} else {
		for _, s := range openChampionsData() {
			n := s.Cost - 1
			championsByCost[n] = append(championsByCost[n], s.Champion)
		}

		file, _ := json.Marshal(championsByCost)
		_ = ioutil.WriteFile("./data/championsByCost.json", file, 0644)
	}
	defer file.Close()

}

// randomly select 5 champions that are rolled
func roll(level int, rolls int) {
	fmt.Println(dropRates)
	rand.Seed(time.Now().UnixNano()) // make rand non-deterministic
	// roll the number of times required
	for r := 0; r < rolls; r++ {
		// each roll consists of 5 champion selections
		for n := 0; n < 5; n++ {
			// randomly select a cost
			var cost int
			rCost := rand.Float64()
			for i, s := range dropRates[level-1] {
				if rCost <= s {
					cost = i
					break
				}
			}
			// randomly select a champion based on cost
			var champion string
			rChampion := rand.Intn((championPool[cost]*len(championsByCost[cost]))-1) + 1
			for i, s := range championsByCost[cost] {
				if rChampion <= championPool[cost]*(i+1) {
					champion = s
					break
				}
			}

			if c, ok := rollResults[champion]; ok {
				rollResults[champion] = rollResult{champion, cost + 1, c.Appearances + 1}
			} else {
				rollResults[champion] = rollResult{champion, cost + 1, 1}
			}
		}
	}
}

func rollLevel(w http.ResponseWriter, r *http.Request) {
	// sanitize input variables
	levelStr := r.URL.Query().Get("level") // returns empty string if not found
	level, err := strconv.Atoi(strings.Trim(levelStr, " "))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Incomplete request: provide a valid level."))
		log.Fatalf("Error: invalid level.")
		return
	}

	rollNumsStr := r.URL.Query().Get("rolls") // number of rolls to perform
	rollNums, err := strconv.Atoi(strings.Trim(rollNumsStr, " "))

	if err != nil || rollNums > 200 {
		w.WriteHeader(400)
		w.Write([]byte("Incomplete request: provide a valid number of rolls (number less than 200)."))
		log.Fatalf("Error: invalid number of rolls.")
		return
	}

	// perform rolls
	assignDropRates()
	getchampionsByCost()
	roll(level, rollNums)
	var rData []rollResult
	for _, s := range rollResults {
		rData = append(rData, s)
	}

	// send data back
	// w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(rData)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	// reset the results at the end of each roll
	rollResults = make(map[string]rollResult)

	fmt.Println(r.URL.Path, level, rollNums)
	fmt.Println(rData)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	fmt.Println(r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/api/roll", rollLevel)
	http.HandleFunc("/", handleRequest)

	log.Fatal(http.ListenAndServe(":3080", nil))
}
