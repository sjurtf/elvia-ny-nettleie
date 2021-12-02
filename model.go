package main

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

func (r *Result) GetPeakPower() float64 {
	s := 0.0
	for _, year := range r.Years {
		s = year.GetPeakPower()
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

func (y *Year) GetPeakPower() float64 {
	s := 0.0
	for _, month := range y.Months {
		s = month.GetPeakPower()
	}
	return s
}

func (m *Month) GetUsage() float64 {
	return m.Usage
}

func (m *Month) GetCost() float64 {
	return m.TotalCost
}

func (m *Month) GetPeakPower() float64 {
	return m.MaxKWh
}
