package main

import (
	"encoding/json"
	"errors"
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

// list of drop rates by level; converted to culmulative rates per level
// level is first array position + 1, tier is second array position + 1
var dropRatesConv [9][5]float64

// results counter
type rollResult struct {
	Champion    string
	Cost        int
	Appearances int
	Percentage  string
}

var rollResults = make(map[string]rollResult)

// data structure for champions.json file
type Champion struct {
	Champion   string
	ChampionId string
	Cost       int
}

// provides the rates to the DropRates global variable
func assignDropRates() error {
	file, err := os.Open("./data/dropRates.json")

	dropRatesOK := true

	if err == nil {
		bv, err := ioutil.ReadAll(file)
		if err != nil {
			dropRatesOK = false
		} else {
			json.Unmarshal(bv, &dropRates)
		}
	} else { // create file if it does not exist
		dropRatesOK = false
	}
	defer file.Close()

	if !dropRatesOK {
		dropRates[0] = [5]float64{1.00, 0.00, 0.00, 0.00, 0.00}
		dropRates[1] = [5]float64{1.00, 0.00, 0.00, 0.00, 0.00}
		dropRates[2] = [5]float64{0.75, 0.25, 0.00, 0.00, 0.00}
		dropRates[3] = [5]float64{0.55, 0.30, 0.15, 0.00, 0.00}
		dropRates[4] = [5]float64{0.45, 0.33, 0.20, 0.02, 0.00}
		dropRates[5] = [5]float64{0.25, 0.40, 0.30, 0.05, 0.00}
		dropRates[6] = [5]float64{0.19, 0.30, 0.35, 0.15, 0.01}
		dropRates[7] = [5]float64{0.15, 0.20, 0.35, 0.25, 0.05}
		dropRates[8] = [5]float64{0.10, 0.15, 0.30, 0.30, 0.15}

		file, err := json.Marshal(dropRates)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
			return errors.New("could not retrieve droprates")
		}
		_ = ioutil.WriteFile("./data/dropRates.json", file, 0644)
	}

	// convert drop rates to cumulative percentages
	dropRatesConvOK := true
	cFile, cErr := os.Open("./data/dropRatesConv.json")

	if cErr == nil {
		bv, err := ioutil.ReadAll(cFile)
		if err != nil {
			dropRatesConvOK = false
		} else {
			json.Unmarshal(bv, &dropRatesConv)
		}
	} else {
		dropRatesConvOK = false
	}
	defer cFile.Close()

	if !dropRatesConvOK {
		dropRatesConv = dropRates
		for i, s := range dropRatesConv {
			for i2 := range s {
				if i2 != 0 {
					dropRatesConv[i][i2] = math.Round((dropRatesConv[i][i2]+dropRatesConv[i][i2-1])*100) / 100
				}
			}
		}

		cFile, err := json.Marshal(dropRatesConv)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
			return errors.New("could not retrieve droprates")
		}
		_ = ioutil.WriteFile("./data/dropRatesConv.json", cFile, 0644)
	}
	return nil
}

// returns a list of champion metadata
func openChampionsData() ([]Champion, error) {
	championsJSON, err := os.Open("./data/champions.json")

	if err != nil {
		fmt.Println(err)
	}

	defer championsJSON.Close()

	byteValue, err := ioutil.ReadAll(championsJSON)
	if err != nil {
		return nil, errors.New("could not retrieve champion data")
	}

	var champions []Champion

	json.Unmarshal(byteValue, &champions)

	return champions, nil
}

// returns a list of champions by cost and the number in pool.
func getchampionsByCost(championsData []Champion) {
	file, err := os.Open("./data/championsByCost.json")
	champOK := true

	if err == nil {
		bv, err := ioutil.ReadAll(file)
		if err != nil {
			champOK = false
		} else {
			json.Unmarshal(bv, &championsByCost)
		}
	} else {
		champOK = false
	}
	defer file.Close()

	if !champOK {
		for _, s := range championsData {
			n := s.Cost - 1
			championsByCost[n] = append(championsByCost[n], s.Champion)
		}

		nfile, err := json.Marshal(championsByCost)
		if err != nil {
			// don't need to handle error
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		_ = ioutil.WriteFile("./data/championsByCost.json", nfile, 0644)
	}
}

// randomly select 5 champions that are rolled
func roll(level int, rolls int) {
	rand.Seed(time.Now().UnixNano()) // make rand non-deterministic
	// roll the number of times required
	for r := 0; r < rolls; r++ {
		// each roll consists of 5 champion selections
		for n := 0; n < 5; n++ {
			// randomly select a cost
			var cost int
			rCost := rand.Float64()
			for i, s := range dropRatesConv[level-1] {
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
				rollResults[champion] = rollResult{champion, cost + 1, c.Appearances + 1, ""}
			} else {
				rollResults[champion] = rollResult{champion, cost + 1, 1, ""}
			}
		}
	}

	// calculate the percentage of appearance for each champion
	for i, s := range rollResults {
		p := (float32(s.Appearances) * 100) / (float32(rolls) * 5)
		s.Percentage = fmt.Sprintf("%.2f", p) + "%"
		rollResults[i] = s
	}
}

func rollLevel(w http.ResponseWriter, r *http.Request) {
	// sanitize input variables
	levelStr := r.URL.Query().Get("level") // returns empty string if not found
	level, err := strconv.Atoi(strings.Trim(levelStr, " "))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Bad Request: invalid level"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	rollNumsStr := r.URL.Query().Get("rolls") // number of rolls to perform
	rollNums, err := strconv.Atoi(strings.Trim(rollNumsStr, " "))

	if err != nil || rollNums > 200 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Bad Request: Provide a valid number of rolls (number less than 200)"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}

	// perform rolls
	err = assignDropRates()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Internal Server Error: " + err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}
	champData, err := openChampionsData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Internal Server Error: " + err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Printf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
		return
	}
	getchampionsByCost(champData)
	roll(level, rollNums)
	var rData []rollResult
	for _, s := range rollResults {
		rData = append(rData, s)
	}

	// send data back
	jsonResp, err := json.Marshal(rData)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error"))
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}

	// reset the results at the end of each roll
	rollResults = make(map[string]rollResult)
}

func getDropRates(w http.ResponseWriter, r *http.Request) {
	assignDropRates()

	type rates struct {
		Level                                           int
		OneCost, TwoCost, ThreeCost, FourCost, FiveCost string
	}

	var resp []rates

	for i, s := range dropRates {
		resp = append(resp, rates{
			i + 1,
			fmt.Sprint(int(s[0]*100)) + "%",
			fmt.Sprint(int(s[1]*100)) + "%",
			fmt.Sprint(int(s[2]*100)) + "%",
			fmt.Sprint(int(s[3]*100)) + "%",
			fmt.Sprint(int(s[4]*100)) + "%",
		})
	}

	// send data back
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error"))
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}

}

func main() {
	http.HandleFunc("/api/roll", rollLevel)
	http.HandleFunc("/api/droprates", getDropRates)
	http.HandleFunc("/", getDropRates)

	log.Fatal(http.ListenAndServe(":3080", nil))
}
