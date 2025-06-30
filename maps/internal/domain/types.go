package domain

type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type WorkingHours struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DaySchedule struct {
	WorkingHours []WorkingHours `json:"working_hours"`
}

type Schedule struct {
	Monday    DaySchedule `json:"Mon"`
	Tuesday   DaySchedule `json:"Tue"`
	Wednesday DaySchedule `json:"Wed"`
	Thursday  DaySchedule `json:"Thu"`
	Friday    DaySchedule `json:"Fri"`
	Saturday  DaySchedule `json:"Sat"`
	Sunday    DaySchedule `json:"Sun"`
}

type Place struct {
	Name     string   `json:"name"`
	Point    Point    `json:"point"`
	Id       string   `json:"id"`
	Schedule Schedule `json:"schedule"`
}

type ShopInfo struct {
	Info []Place `json:"info"`
	Shop string  `json:"shop"`
}

type MinTimeRoute struct {
	Points    []int  `json:"points"`
	Duration  int    `json:"duration"`
	Transport string `json:"transport"`
}
