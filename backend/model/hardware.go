package model

type Hardware struct {
	ID             int    `json:"id"`
	OS             string `json:"operatingSystem"`
	Processor      string `json:"processor"`
	Memory         string `json:"memory"`
	GraphicsCard   string `json:"graphicsCard"`
	DirectXVersion string `json:"directXVersion"`
	Network        string `json:"network"`
	Storage        string `json:"storage"`
	Soundcard      string `json:"soundcard"`
	AdditionalNote string `json:"additionalNote"`
	Type           string `json:"type"`
}
