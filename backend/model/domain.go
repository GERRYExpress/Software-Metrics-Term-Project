package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type ID int
type Name string
type Description string
type Price float64
type PosterURL string
type ReleaseDate time.Time
type Username string
type Password string
type Date time.Time
type Count int
type Amount float64
type Hidden bool

func NewID(value int) ID                         { return ID(value) }
func NewName(value string) Name                  { return Name(value) }
func NewDescription(value string) Description    { return Description(value) }
func NewPrice(value float64) Price               { return Price(value) }
func NewPosterURL(value string) PosterURL        { return PosterURL(value) }
func NewReleaseDate(value time.Time) ReleaseDate { return ReleaseDate(value) }
func NewUsername(value string) Username          { return Username(value) }
func NewPassword(value string) Password          { return Password(value) }
func NewDate(value time.Time) Date               { return Date(value) }
func NewCount(value int) Count                   { return Count(value) }
func NewAmount(value float64) Amount             { return Amount(value) }
func NewHidden(value bool) Hidden                { return Hidden(value) }

func (d Date) MarshalJSON() ([]byte, error) {
	return time.Time(d).MarshalJSON()
}

func (d *Date) Scan(value any) error {
	parsed, err := scanTime(value)
	if err != nil {
		return err
	}
	*d = Date(parsed)
	return nil
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d ReleaseDate) MarshalJSON() ([]byte, error) {
	return time.Time(d).MarshalJSON()
}

func (d *ReleaseDate) Scan(value any) error {
	parsed, err := scanTime(value)
	if err != nil {
		return err
	}
	*d = ReleaseDate(parsed)
	return nil
}

func (d ReleaseDate) Value() (driver.Value, error) {
	return time.Time(d), nil
}

func (d *Date) UnmarshalJSON(value []byte) error {
	var parsed time.Time
	if err := json.Unmarshal(value, &parsed); err != nil {
		return err
	}
	*d = Date(parsed)
	return nil
}

func scanTime(value any) (time.Time, error) {
	switch value := value.(type) {
	case time.Time:
		return value, nil
	case nil:
		return time.Time{}, nil
	case string:
		return time.Parse(time.RFC3339, value)
	case []byte:
		return time.Parse(time.RFC3339, string(value))
	default:
		return time.Time{}, fmt.Errorf("cannot scan %T into time value", value)
	}
}

func NewGame(id int, name, description string, price float64, posterURL string) Game {
	return Game{
		ID:          NewID(id),
		Name:        NewName(name),
		Description: NewDescription(description),
		Price:       NewPrice(price),
		PosterUrl:   NewPosterURL(posterURL),
		Categories:  make(Categories, 0),
	}
}

func NewGameDetail(id int, name, description string, price float64, posterURL string, releaseDate time.Time) GameDetail {
	return GameDetail{
		ID:                  NewID(id),
		Name:                NewName(name),
		Description:         NewDescription(description),
		Price:               NewPrice(price),
		PosterUrl:           NewPosterURL(posterURL),
		ReleaseDate:         NewReleaseDate(releaseDate),
		Categories:          make(Categories, 0),
		HardwareRequirement: make([]Hardware, 0),
		Media:               make([]Media, 0),
	}
}

func NewCategory(id int, name string, isHidden bool) Category {
	return Category{ID: NewID(id), Name: NewName(name), IsHidden: NewHidden(isHidden)}
}

func NewAccount(id int, username string) Account {
	return Account{ID: NewID(id), Username: NewUsername(username)}
}

type HardwareData struct {
	ID             int
	OS             string
	Processor      string
	Memory         string
	GraphicsCard   string
	DirectXVersion string
	Network        string
	Storage        string
	Soundcard      string
	AdditionalNote string
	Type           string
}

func NewHardware(data HardwareData) Hardware {
	return Hardware{
		ID:             NewID(data.ID),
		OS:             NewDescription(data.OS),
		Processor:      NewDescription(data.Processor),
		Memory:         NewDescription(data.Memory),
		GraphicsCard:   NewDescription(data.GraphicsCard),
		DirectXVersion: NewDescription(data.DirectXVersion),
		Network:        NewDescription(data.Network),
		Storage:        NewDescription(data.Storage),
		Soundcard:      NewDescription(data.Soundcard),
		AdditionalNote: NewDescription(data.AdditionalNote),
		Type:           NewName(data.Type),
	}
}
