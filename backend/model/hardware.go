package model

type Hardware struct {
	ID             ID          `json:"id"`
	OS             Description `json:"operatingSystem"`
	Processor      Description `json:"processor"`
	Memory         Description `json:"memory"`
	GraphicsCard   Description `json:"graphicsCard"`
	DirectXVersion Description `json:"directXVersion"`
	Network        Description `json:"network"`
	Storage        Description `json:"storage"`
	Soundcard      Description `json:"soundcard"`
	AdditionalNote Description `json:"additionalNote"`
	Type           Name        `json:"type"`
}
