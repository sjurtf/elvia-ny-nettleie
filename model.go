package main

type Data struct {
	Years []Year
}

type Year struct {
	Year   string
	Months []Month
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
