package entity

type Filter struct {
	Ingredients []*Ingredient `json:"ingredients"`
	TimeBetween string        `json:"timeBetween"`
	SortTime    string        `json:"sortTime"`
}
