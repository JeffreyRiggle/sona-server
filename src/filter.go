package main

type Filter struct {
	Property       string `json:"property"`
	ComparisonType string `json:"comparisontype"`
	Value          string `json:"value"`
}

type ComplexFilter struct {
	Filter           string `json:"filter"`
	Junction         string `json:"junction"`
	AdditionalFilter Filter `json:"additionalFilter"`
}

type FilterRequest struct {
	Filters  []ComplexFilter `json:"filters"`
	Junction string          `json:"junction"`
}
