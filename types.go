package doljara

type Response struct {
	Room Room `json:"room"`
}

type Room struct {
	Code string `json:"code"`
	Raw  Raw    `json:"raw"`
}

type Raw struct {
	Rev  int      `json:"rev"`
	Data Audition `json:"data"`
}

type Audition struct {
	EastP      Producer `json:"eastP"`
	WestP      Producer `json:"westP"`
	SouthP     Producer `json:"southP"`
	NorthP     Producer `json:"northP"`
	Applicants []Idol   `json:"applicants"`
}

type Producer struct {
	Debut           bool   `json:"debut"`
	ScoutedIdols    []Idol `json:"scoutedIdols"`
	UnselectedIdols []Idol `json:"unselectedIdols"`
}

type Idol struct {
	ID        string `json:"id"`
	ImagePath string `json:"imagePath"`
	IsReached bool   `json:"isReached"`
	IsDrawn   bool   `json:"isDrawn"`
}
