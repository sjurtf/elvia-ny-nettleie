package elvia

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

type Data struct {
	Years []Year
}

type Year struct {
	Year        string
	Months      []Month
	Consumption Consumption
}

type Month struct {
	Month string
	Days  []Day
}

type Day struct {
	Day   string
	Hours []Hour
}

type Hour struct {
	Hour        string
	Id          string
	Consumption Consumption
}

type Consumption struct {
	Value float64
}
