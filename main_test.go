package main

import (
	"github.com/sjurtf/elvia-ny-nettleie/elvia"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestadata(consumptionData map[string]float64) *elvia.Data {
	var months []elvia.Month
	var years []elvia.Year
	var days []elvia.Day
	var hours []elvia.Hour

	for k, f := range consumptionData {
		h := elvia.Hour{
			Hour:        k,
			Id:          k,
			Consumption: elvia.Consumption{Value: f},
		}
		hours = append(hours, h)
	}

	d := elvia.Day{
		Day:   "01",
		Hours: hours,
	}
	days = append(days, d)

	m := elvia.Month{
		Month: "10",
		Days:  days,
	}
	months = append(months, m)

	y := elvia.Year{
		Year:        "2021",
		Months:      months,
		Consumption: elvia.Consumption{},
	}

	years = append(years, y)
	return &elvia.Data{Years: years}
}

func TestMonth_GetCost_Old(t *testing.T) {
	consumption := make(map[string]float64)
	consumption["10"] = 1.5
	consumption["22"] = 3

	data := getTestadata(consumption)
	r := CalculateOldModel(data)

	expected := elvia.OldConstPriceNOK + (elvia.OldEnergy * 1.5 / 100) + (elvia.OldEnergy * 3 / 100)
	assert.Equal(t, expected, r.GetCost())
}

func TestMonth_GetCost_New(t *testing.T) {
	consumption := make(map[string]float64)
	consumption["10"] = 1.5
	consumption["22"] = 3

	data := getTestadata(consumption)
	r := CalculateNewModel(data)

	expected := elvia.NewConstPriceTier1NOK + (elvia.NewEnergyDay * 1.5 / 100) + (elvia.NewEnergyNight * 3 / 100)
	assert.Equal(t, expected, r.GetCost())
}
