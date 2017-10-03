package main

import "strings"

type Filter struct {
	Property       string `json:"property"`
	ComparisonType string `json:"comparisontype"`
	Value          string `json:"value"`
}

type ComplexFilter struct {
	Filter           *Filter `json:"filter"`
	Junction         string  `json:"junction"`
	AdditionalFilter *Filter `json:"additionalFilter"`
}

type FilterRequest struct {
	Filters  []ComplexFilter `json:"filters"`
	Junction string          `json:"junction"`
}

func isOrRequest(filter *FilterRequest) bool {
	return strings.EqualFold("or", filter.Junction)
}

func isAndRequest(filter *FilterRequest) bool {
	return strings.EqualFold("and", filter.Junction)
}

func isOrFilter(filter ComplexFilter) bool {
	return strings.EqualFold("or", filter.Junction)
}

func isAndFilter(filter ComplexFilter) bool {
	return strings.EqualFold("and", filter.Junction)
}

func isEqualsComparision(filter Filter) bool {
	return strings.EqualFold("equals", filter.ComparisonType)
}

func isNotEqualsComparision(filter Filter) bool {
	return strings.EqualFold("notequals", filter.ComparisonType)
}

func isContainsComparision(filter Filter) bool {
	return strings.EqualFold("contains", filter.ComparisonType)
}
