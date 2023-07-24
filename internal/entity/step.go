package entity

type Step struct {
	StepNumber  uint8  `json:"stepNumber"`
	Description string `json:"description"`
	Time        int    `json:"time"`
}
