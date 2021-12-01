package main

import (
	"encoding/json"
	"fmt"
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

	data := &Data{}
	err = json.Unmarshal(f, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	old := calculateOldModel(data)
	new, maxkWh := calculateNewModel(data)

	symbol := ""
	percentChange := (new - old) / old * 100
	if percentChange > 0 {
		symbol = "+"
	}

	fmt.Printf("Usage: %s kWh, Peak hour: %.2f kWh\n", "382.71", maxkWh)
	fmt.Printf("Sum old model: kr %.2f\n", old)
	fmt.Printf("Sum new model: kr %.2f\n", new)
	fmt.Printf("Diff: kr %f %s%.2f %%\n", new-old, symbol, percentChange)

}

const (
	OldConstPriceNOK      = 115.0
	NewConstPriceTier0NOK = 130.0 // 0-2kW
	NewConstPriceTier1NOK = 190.0 // 2-5kW
	NewConstPriceTier2NOK = 280.0 // 5-10kW
	NewConstPriceTier3NOK = 375.0 // 10-15kW
	NewConstPriceTier4NOK = 470.0 // 15-20kW

	OldEnergy      = 44.80
	NewEnergyDay   = 37.35
	NewEnergyNight = 31.10 // also weekends
)

func calculateOldModel(data *Data) float64 {
	price := 0.0
	kWhCounted := 0.0

	for _, month := range data.Years[0].Months {
		for _, day := range month.Days {
			for _, hour := range day.Hours {
				price = price + (hour.Consumption.Value * OldEnergy)
				kWhCounted = kWhCounted + hour.Consumption.Value
			}
		}
	}
	price = price / 100
	price = price + OldConstPriceNOK
	return price
}

func calculateNewModel(data *Data) (float64, float64) {
	price := 0.0
	kWhCounted := 0.0
	maxkWh := 0.0

	tier := 0

	for _, month := range data.Years[0].Months {
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
					tier = 1
				} else if hour.Consumption.Value > 5 {
					tier = 2
				} else if hour.Consumption.Value > 10 {
					tier = 3
				} else if hour.Consumption.Value > 15 {
					tier = 4
				} else if hour.Consumption.Value > 20 {
					tier = 5
				}

				if hour.Consumption.Value > maxkWh {
					maxkWh = hour.Consumption.Value
				}

				h, _ := strconv.Atoi(hour.Hour)

				if isWeekend {
					price = price + (hour.Consumption.Value * NewEnergyNight)
					kWhCounted = kWhCounted + hour.Consumption.Value
				} else if h >= 22 || (h >= 0 && h <= 6) {
					price = price + (hour.Consumption.Value * NewEnergyNight)
					kWhCounted = kWhCounted + hour.Consumption.Value
				} else {
					price = price + (hour.Consumption.Value * NewEnergyDay)
					kWhCounted = kWhCounted + hour.Consumption.Value
				}
			}
		}
	}

	constPrice := NewConstPriceTier0NOK
	switch tier {
	case 1:
		constPrice = NewConstPriceTier1NOK
	case 2:
		constPrice = NewConstPriceTier2NOK
	case 3:
		constPrice = NewConstPriceTier3NOK
	case 4:
		constPrice = NewConstPriceTier4NOK
	}

	price = price / 100
	price = price + constPrice
	return price, maxkWh
}
