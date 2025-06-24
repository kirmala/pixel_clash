package ctypes

type Game struct {
	Size           int `json:"size"`
	FieldSize      int `json:"field_size"`
	Time           int `json:"time"`
	ThresholdSqare int `json:"threshold_square"`
	Cooldown int `json:"cooldown"`
}
