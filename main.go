package main

import (
	"encoding/json"
	"fmt"
	"github.com/sjurtf/elvia-ny-nettleie/elvia"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	filename := os.Args[1]
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := &elvia.Data{}
	err = json.Unmarshal(f, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	totalConsumption := totalConsumption(data)

	oldResults := CalculateOldModel(data)
	newResults := CalculateNewModel(data)

	oldCost := oldResults.GetCost()
	newCost := newResults.GetCost()

	symbol := ""
	percentChange := (newCost - oldCost) / oldCost * 100
	if percentChange > 0 {
		symbol = "+"
	}

	fmt.Printf("Usage: %.2f kWh - Peak hour %.2f kWh\n", totalConsumption, newResults.GetPeakPower())
	fmt.Printf("Sum old model: %.2f kr\n", oldCost)
	fmt.Printf("Sum new model: %.2f kr\n", newCost)
	fmt.Printf("Diff: %.2f kr %s%.2f %%\n", newCost-oldCost, symbol, percentChange)

}

func totalConsumption(data *elvia.Data) float64 {
	c := 0.0
	for _, year := range data.Years {
		c = c + year.Consumption.Value
	}
	return c
}

func CalculateOldModel(data *elvia.Data) Result {
	var years []Year

	for _, year := range data.Years {
		var months []Month
		for _, month := range year.Months {
			price := 0.0
			kWhCounted := 0.0
			m := Month{}

			for _, day := range month.Days {
				for _, hour := range day.Hours {
					price = price + (hour.Consumption.Value * elvia.OldEnergy)
					kWhCounted = kWhCounted + hour.Consumption.Value
				}
			}

			m.Usage = kWhCounted
			m.UsageTierCost = elvia.OldConstPriceNOK
			m.TotalCost = (price / 100) + m.UsageTierCost
			months = append(months, m)
		}

		y := Year{Months: months}
		years = append(years, y)
	}
	return Result{Years: years}
}

func CalculateNewModel(data *elvia.Data) Result {
	var years []Year
	for _, year := range data.Years {
		var months []Month
		for _, month := range year.Months {
			cost := 0.0
			kWh := 0.0
			m := Month{}

			for _, day := range month.Days {
				date := fmt.Sprintf("%s-%s-%s", year.Year, month.Month, day.Day)
				t, err := time.Parse("2006-01-02", date)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				isWeekend := false
				if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
					isWeekend = true
				}

				summer := false
				if t.Month() >= 4 && t.Month() <= 10 {
					summer = true
				}

				for _, hour := range day.Hours {
					if hour.Consumption.Value > m.MaxKWh {
						m.MaxKWh = hour.Consumption.Value
					}

					c, k := hourCost(isWeekend, summer, hour)

					cost = cost + c
					kWh = kWh + k
				}
			}
			m.Usage = kWh
			m.TotalCost = (cost / 100) + constPriceTier(m.MaxKWh)
			months = append(months, m)
		}

		y := Year{Months: months}
		years = append(years, y)
	}

	return Result{Years: years}
}

func hourCost(isWeekend bool, summer bool, hour elvia.Hour) (float64, float64) {
	h, _ := strconv.Atoi(hour.Hour)

	cost := 0.0
	kWh := 0.0
	if isWeekend {
		if summer {
			cost = hour.Consumption.Value * elvia.NewEnergyNight
		} else {
			cost = hour.Consumption.Value * elvia.NewEnergyNightWinter
		}
		kWh = hour.Consumption.Value

	} else if h >= 22 || (h >= 0 && h <= 6) {
		if summer {
			cost = hour.Consumption.Value * elvia.NewEnergyNight
		} else {
			cost = hour.Consumption.Value * elvia.NewEnergyNightWinter
		}
		kWh = hour.Consumption.Value

	} else {
		if summer {
			cost = hour.Consumption.Value * elvia.NewEnergyDay
		} else {
			cost = hour.Consumption.Value * elvia.NewEnergyDayWinter
		}
		kWh = hour.Consumption.Value
	}
	return cost, kWh
}

func constPriceTier(maxKWh float64) float64 {
	if maxKWh > 15 {
		return elvia.NewConstPriceTier4NOK
	} else if maxKWh > 10 {
		return elvia.NewConstPriceTier3NOK
	} else if maxKWh > 5 {
		return elvia.NewConstPriceTier2NOK
	} else if maxKWh > 2 {
		return elvia.NewConstPriceTier1NOK
	}
	return elvia.NewConstPriceTier0NOK
}
