package scraper

type PageOptions struct {
	StepData StepData `json:"step_data"`
}

type StepData struct {
	List []Category `json:"list"`
}

type Category struct {
	Services []Service `json:"services"`
}

type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
