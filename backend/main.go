package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// number of champions per level
var championPool = [5]int{29, 22, 18, 12, 10}

// list of champions and number in pool by its cost
type ChampionsByCost = [5][]string

// list of drop rates by level
// level is first array position + 1, tier is second array position + 1
type DropRates = [9][5]float64

var baseDR DropRates = DropRates{
	[5]float64{1.00, 0.00, 0.00, 0.00, 0.00},
	[5]float64{1.00, 0.00, 0.00, 0.00, 0.00},
	[5]float64{0.75, 0.25, 0.00, 0.00, 0.00},
	[5]float64{0.55, 0.30, 0.15, 0.00, 0.00},
	[5]float64{0.45, 0.33, 0.20, 0.02, 0.00},
	[5]float64{0.25, 0.40, 0.30, 0.05, 0.00},
	[5]float64{0.19, 0.30, 0.35, 0.15, 0.01},
	[5]float64{0.15, 0.20, 0.35, 0.25, 0.05},
	[5]float64{0.10, 0.15, 0.30, 0.30, 0.15},
}

// results counter
type rollResult struct {
	Champion    string
	Cost        int
	Appearances int
	Percentage  string
}

// data structure for champions.json file
type Champion struct {
	Name       string
	ChampionId string
	Cost       int
	Traits     []string
}

// convert base drop rates
func convertDropRates() DropRates {
	convertedDropRates := baseDR

	for i, s := range convertedDropRates {
		for i2 := range s {
			if i2 != 0 {
				convertedDropRates[i][i2] = math.Round((convertedDropRates[i][i2]+convertedDropRates[i][i2-1])*100) / 100
			}
		}
	}

	return convertedDropRates
}

func openDropRatesConv() DropRates {
	f, err := os.OpenFile("./data/dropRatesConv.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return convertDropRates()
	}
	defer f.Close()

	var dr DropRates
	// The only reason this call would error is if there is malformed JSON, but we are going to overwrite that JSON.
	if err := json.NewDecoder(f).Decode(&dr); err != nil || len(dr) == 0 {
		// If there was an error decoding or there is no data, then we overwrite the data with our own.
		// If there is an error when writing to the file, explicitly ignore it as we've done all we can here, and we can just proceed with the default data set
		// Though one may wish to omit a warning if this were to happen using a logging package set to trace/debug levels
		cdr := convertDropRates()
		_ = json.NewEncoder(f).Encode(cdr)
		return cdr
	}

	return dr
}

// returns a list of champion metadata
func openChampionsData() ([]Champion, error) {
	f, err := os.Open("./data/champions.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var champions []Champion
	return champions, json.NewDecoder(f).Decode(&champions)
}

func convertChampionsByCost() (ChampionsByCost, error) {
	var c ChampionsByCost
	champions, err := openChampionsData()
	if err != nil {
		return ChampionsByCost{}, err
	}
	for _, s := range champions {
		n := s.Cost - 1
		c[n] = append(c[n], s.Name)
	}

	return c, nil
}

// returns a list of champions by cost and the number in pool.
func openChampionsByCost() (ChampionsByCost, error) {
	f, err := os.OpenFile("./data/championsByCost.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return convertChampionsByCost()
	}
	defer f.Close()

	var c ChampionsByCost
	// The only reason this call would error is if there is malformed JSON, but we are going to overwrite that JSON.
	if err := json.NewDecoder(f).Decode(&c); err != nil || len(c) == 0 {
		// If there was an error decoding or there is no data, then we overwrite the data with our own.
		// If there is an error when writing to the file, explicitly ignore it as we've done all we can here, and we can just proceed with the default data set
		// Though one may wish to omit a warning if this were to happen using a logging package set to trace/debug levels
		cCosts, cErr := convertChampionsByCost()
		_ = json.NewEncoder(f).Encode(cCosts)
		return cCosts, cErr
	}

	return c, err
}

// randomly select 5 champions that are rolled
func roll(level int, rolls int) (*map[string]rollResult, error) {

	var results = make(map[string]rollResult)

	rand.Seed(time.Now().UnixNano()) // make rand non-deterministic
	// get meta data
	dropRates := openDropRatesConv()
	championsByCost, err := openChampionsByCost()
	if err != nil {
		return nil, err
	}

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

			if c, ok := results[champion]; ok {
				results[champion] = rollResult{champion, cost + 1, c.Appearances + 1, ""}
			} else {
				results[champion] = rollResult{champion, cost + 1, 1, ""}
			}
		}
	}

	return &results, nil
}

func rollLevel(w http.ResponseWriter, r *http.Request) {
	// sanitize input variables
	levelStr := r.URL.Query().Get("level") // returns empty string if not found
	level, err := strconv.Atoi(strings.Trim(levelStr, " "))
	if err != nil {
		JSONError(w,
			map[string]string{"message": "Bad Request: invalid level"},
			http.StatusBadRequest)
		return
	}

	rollNumsStr := r.URL.Query().Get("rolls") // number of rolls to perform
	rollNums, err := strconv.Atoi(strings.Trim(rollNumsStr, " "))
	if err != nil || rollNums > 200 {
		JSONError(w,
			map[string]string{"message": "Bad Request: Provide a valid number of rolls (number less than 200)"},
			http.StatusBadRequest)
		return
	}

	// perform rolls
	results, err := roll(level, rollNums)
	if err != nil {
		JSONError(w,
			map[string]string{"message": "Bad Request: Provide a valid number of rolls (number less than 200)"},
			http.StatusBadRequest)
		return
	}

	// calculate appearance percentages and prepare results
	var resp []rollResult
	for _, s := range *results {
		s.Percentage = fmt.Sprintf("%.2f%%", (float32(s.Appearances)*100)/(float32(rollNums)*5))
		resp = append(resp, s)
	}

	// success; return results
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getDropRates(w http.ResponseWriter, r *http.Request) {

	type rates struct {
		Level                                           int
		OneCost, TwoCost, ThreeCost, FourCost, FiveCost string
	}

	var resp []rates

	for i, s := range baseDR {
		resp = append(resp, rates{
			i + 1,
			fmt.Sprintf("%d%%", int(s[0]*100)),
			fmt.Sprintf("%d%%", int(s[1]*100)),
			fmt.Sprintf("%d%%", int(s[2]*100)),
			fmt.Sprintf("%d%%", int(s[3]*100)),
			fmt.Sprintf("%d%%", int(s[4]*100)),
		})
	}

	// success; return results
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func JSONError(w http.ResponseWriter, message interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func main() {
	http.HandleFunc("/api/roll", rollLevel)
	http.HandleFunc("/api/droprates", getDropRates)
	http.HandleFunc("/", getDropRates)

	log.Fatal(http.ListenAndServe(":3080", nil))
}
