package ctypes

type Cell struct {
	CompSize int    `json:"component_size"`
	Color    string `json:"color"`
}

type Field struct {
	Data [][]Cell `json:"data"`
}
