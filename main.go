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

type Result struct {
	Years []Year
}
type Year struct {
	Months []Month
}

type Month struct {
	Usage         float64
	MaxKWh        float64
	UsageTierCost float64
	TotalCost     float64
}

func (r *Result) GetUsage() float64 {
	s := 0.0
	for _, year := range r.Years {
		s = year.GetUsage()
	}
	return s
}

func (r *Result) GetCost() float64 {
	s := 0.0
	for _, year := range r.Years {
		s = year.GetCost()
	}
	return s
}

func (y *Year) GetUsage() float64 {
	s := 0.0
	for _, month := range y.Months {
		s = month.GetUsage()
	}
	return s
}

func (y *Year) GetCost() float64 {
	s := 0.0
	for _, month := range y.Months {
		s = month.GetCost()
	}
	return s
}

func (m *Month) GetUsage() float64 {
	return m.Usage
}

func (m *Month) GetCost() float64 {
	return m.TotalCost
}

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

	oldResults := calculateOldModel(data)
	newResults := calculateNewModel(data)

	oldCost := oldResults.GetCost()
	newCost := newResults.GetCost()

	symbol := ""
	percentChange := (newCost - oldCost) / oldCost * 100
	if percentChange > 0 {
		symbol = "+"
	}

	fmt.Printf("Usage: %.2f kWh\n", totalConsumption)
	fmt.Printf("Sum old model: kr %.2f\n", oldCost)
	fmt.Printf("Sum new model: kr %.2f\n", newCost)
	fmt.Printf("Diff: kr %f %s%.2f %%\n", newCost-oldCost, symbol, percentChange)

}

func totalConsumption(data *elvia.Data) float64 {
	c := 0.0
	for _, year := range data.Years {
		c = c + year.Consumption.Value
	}
	return c
}

func calculateOldModel(data *elvia.Data) Result {
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

func calculateNewModel(data *elvia.Data) Result {
	var years []Year
	for _, year := range data.Years {
		var months []Month
		for _, month := range year.Months {
			price := 0.0
			kWhCounted := 0.0
			m := Month{}

			for _, day := range month.Days {

				date := fmt.Sprintf("%s-%s-%s", data.Years[0].Year, month.Month, day.Day)
				t, err := time.Parse("2006-01-02", date)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				isWeekend := false
				if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
					isWeekend = true
				}

				for _, hour := range day.Hours {
					if hour.Consumption.Value > 2 {
						m.UsageTierCost = elvia.NewConstPriceTier0NOK
					} else if hour.Consumption.Value > 5 {
						m.UsageTierCost = elvia.NewConstPriceTier1NOK
					} else if hour.Consumption.Value > 10 {
						m.UsageTierCost = elvia.NewConstPriceTier2NOK
					} else if hour.Consumption.Value > 15 {
						m.UsageTierCost = elvia.NewConstPriceTier3NOK
					} else if hour.Consumption.Value > 20 {
						m.UsageTierCost = elvia.NewConstPriceTier4NOK
					}

					if hour.Consumption.Value > m.MaxKWh {
						m.MaxKWh = hour.Consumption.Value
					}

					h, _ := strconv.Atoi(hour.Hour)

					if isWeekend {
						price = price + (hour.Consumption.Value * elvia.NewEnergyNight)
						kWhCounted = kWhCounted + hour.Consumption.Value
					} else if h >= 22 || (h >= 0 && h <= 6) {
						price = price + (hour.Consumption.Value * elvia.NewEnergyNight)
						kWhCounted = kWhCounted + hour.Consumption.Value
					} else {
						price = price + (hour.Consumption.Value * elvia.NewEnergyDay)
						kWhCounted = kWhCounted + hour.Consumption.Value
					}
				}
			}

			m.Usage = kWhCounted
			m.TotalCost = (price / 100) + m.UsageTierCost
			months = append(months, m)
		}

		y := Year{Months: months}
		years = append(years, y)
	}

	return Result{Years: years}
}
