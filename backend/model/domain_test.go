package model

import (
	"reflect"
	"testing"
	"time"
)

var testTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
var testTimeStr = "2023-01-01T12:00:00Z"
var testTimeBytes = []byte(`"` + testTimeStr + `"`)

func TestNewID(t *testing.T) {
	if got := NewID(42); got != ID(42) {
		t.Errorf("NewID() = %v, want 42", got)
	}
}

func TestNewName(t *testing.T) {
	if got := NewName("test"); got != Name("test") {
		t.Errorf("NewName() = %v, want 'test'", got)
	}
}

func TestNewDescription(t *testing.T) {
	if got := NewDescription("desc"); got != Description("desc") {
		t.Errorf("NewDescription() = %v, want 'desc'", got)
	}
}

func TestNewPrice(t *testing.T) {
	if got := NewPrice(19.99); got != Price(19.99) {
		t.Errorf("NewPrice() = %v, want 19.99", got)
	}
}

func TestNewPosterURL(t *testing.T) {
	if got := NewPosterURL("url"); got != PosterURL("url") {
		t.Errorf("NewPosterURL() = %v, want 'url'", got)
	}
}

func TestNewReleaseDate(t *testing.T) {
	if got := NewReleaseDate(testTime); !reflect.DeepEqual(got, ReleaseDate(testTime)) {
		t.Errorf("NewReleaseDate() = %v, want %v", got, ReleaseDate(testTime))
	}
}

func TestNewUsername(t *testing.T) {
	if got := NewUsername("user"); got != Username("user") {
		t.Errorf("NewUsername() = %v, want 'user'", got)
	}
}

func TestNewPassword(t *testing.T) {
	if got := NewPassword("pass"); got != Password("pass") {
		t.Errorf("NewPassword() = %v, want 'pass'", got)
	}
}

func TestNewDate(t *testing.T) {
	if got := NewDate(testTime); !reflect.DeepEqual(got, Date(testTime)) {
		t.Errorf("NewDate() = %v, want %v", got, Date(testTime))
	}
}

func TestNewCount(t *testing.T) {
	if got := NewCount(5); got != Count(5) {
		t.Errorf("NewCount() = %v, want 5", got)
	}
}

func TestNewAmount(t *testing.T) {
	if got := NewAmount(10.5); got != Amount(10.5) {
		t.Errorf("NewAmount() = %v, want 10.5", got)
	}
}

func TestNewHidden(t *testing.T) {
	if got := NewHidden(true); got != Hidden(true) {
		t.Errorf("NewHidden() = %v, want true", got)
	}
}

func TestNewOrderDate(t *testing.T) {
	if got := NewOrderDate(testTime); !reflect.DeepEqual(got, OrderDate(testTime)) {
		t.Errorf("NewOrderDate() = %v, want %v", got, OrderDate(testTime))
	}
}

func TestNewPaymentMethod(t *testing.T) {
	if got := NewPaymentMethod("CC"); got != PaymentMethod("CC") {
		t.Errorf("NewPaymentMethod() = %v, want 'CC'", got)
	}
}

func TestNewFinalPrice(t *testing.T) {
	if got := NewFinalPrice(99.99); got != FinalPrice(99.99) {
		t.Errorf("NewFinalPrice() = %v, want 99.99", got)
	}
}

func TestNewMediaURL(t *testing.T) {
	if got := NewMediaURL("img"); got != MediaURL("img") {
		t.Errorf("NewMediaURL() = %v, want 'img'", got)
	}
}

func TestHidden_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    Hidden
		wantErr bool
	}{
		{"bool true", true, true, false},
		{"int64 1", int64(1), true, false},
		{"int64 0", int64(0), false, false},
		{"string true", "true", true, false},
		{"string 1", "1", true, false},
		{"byte true", []byte("true"), true, false},
		{"invalid type", 12.34, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h Hidden
			err := h.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hidden.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && h != tt.want {
				t.Errorf("Hidden.Scan() got = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestHidden_scanString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Hidden
		wantErr bool
	}{
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"invalid", "not-a-bool", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h Hidden
			if err := h.scanString(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Hidden.scanString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && h != tt.want {
				t.Errorf("Hidden.scanString() got = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestHidden_Value(t *testing.T) {
	h := Hidden(true)
	val, err := h.Value()
	if err != nil || val != true {
		t.Errorf("Hidden.Value() = %v, %v; want true, nil", val, err)
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	d := Date(testTime)
	got, err := d.MarshalJSON()
	if err != nil || !reflect.DeepEqual(got, testTimeBytes) {
		t.Errorf("Date.MarshalJSON() = %s, want %s", got, testTimeBytes)
	}
}

func TestDate_Scan(t *testing.T) {
	var d Date
	if err := d.Scan(testTimeStr); err != nil || time.Time(d) != testTime {
		t.Errorf("Date.Scan() failed, got %v", time.Time(d))
	}
}

func TestDate_Value(t *testing.T) {
	d := Date(testTime)
	val, err := d.Value()
	if err != nil || val != testTime {
		t.Errorf("Date.Value() failed")
	}
}

func TestOrderDate_MarshalJSON(t *testing.T) {
	d := OrderDate(testTime)
	got, _ := d.MarshalJSON()
	if string(got) != string(testTimeBytes) {
		t.Errorf("OrderDate.MarshalJSON() = %s", got)
	}
}

func TestOrderDate_Scan(t *testing.T) {
	var d OrderDate
	_ = d.Scan(testTime)
	if time.Time(d) != testTime {
		t.Errorf("OrderDate.Scan() failed")
	}
}

func TestOrderDate_Value(t *testing.T) {
	d := OrderDate(testTime)
	val, _ := d.Value()
	if val != testTime {
		t.Errorf("OrderDate.Value() failed")
	}
}

func TestReleaseDate_MarshalJSON(t *testing.T) {
	d := ReleaseDate(testTime)
	got, _ := d.MarshalJSON()
	if string(got) != string(testTimeBytes) {
		t.Errorf("ReleaseDate.MarshalJSON() = %s", got)
	}
}

func TestReleaseDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		value   []byte
		want    time.Time
		wantErr bool
	}{
		{"RFC3339", testTimeBytes, testTime, false},
		{"YYYY-MM-DD", []byte(`"2023-01-01"`), time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"invalid", []byte(`"bad-date"`), time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d ReleaseDate
			if err := d.UnmarshalJSON(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("ReleaseDate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !time.Time(d).Equal(tt.want) {
				t.Errorf("ReleaseDate.UnmarshalJSON() got = %v, want %v", time.Time(d), tt.want)
			}
		})
	}
}

func TestReleaseDate_Scan(t *testing.T) {
	var d ReleaseDate
	_ = d.Scan(testTime)
	if time.Time(d) != testTime {
		t.Errorf("ReleaseDate.Scan() failed")
	}
}

func TestReleaseDate_Value(t *testing.T) {
	d := ReleaseDate(testTime)
	val, _ := d.Value()
	if val != testTime {
		t.Errorf("ReleaseDate.Value() failed")
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	var d Date
	err := d.UnmarshalJSON(testTimeBytes)
	if err != nil || !time.Time(d).Equal(testTime) {
		t.Errorf("Date.UnmarshalJSON() failed")
	}
}

func Test_scanTime(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    time.Time
		wantErr bool
	}{
		{"time.Time", testTime, testTime, false},
		{"string", testTimeStr, testTime, false},
		{"[]byte", []byte(testTimeStr), testTime, false},
		{"nil", nil, time.Time{}, false},
		{"invalid type", 123, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scanTime(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("scanTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	got := NewGame(1, "Name", "Desc", 9.99, "url")
	if got.ID != 1 || got.Name != "Name" || got.Description != "Desc" || got.Price != 9.99 || got.PosterUrl != "url" {
		t.Errorf("NewGame() mapped fields incorrectly: %v", got)
	}
	if got.Categories == nil {
		t.Errorf("NewGame() did not initialize Categories slice")
	}
}

func TestNewGameDetail(t *testing.T) {
	got := NewGameDetail(1, "Name", "Desc", 9.99, "url", testTime)
	if got.ID != 1 || time.Time(got.ReleaseDate) != testTime {
		t.Errorf("NewGameDetail() mapped fields incorrectly")
	}
	if got.Categories == nil || got.HardwareRequirement == nil || got.Media == nil {
		t.Errorf("NewGameDetail() did not initialize slices")
	}
}

func TestNewCategory(t *testing.T) {
	got := NewCategory(1, "RPG", true)
	if got.ID != 1 || got.Name != "RPG" || got.IsHidden != true {
		t.Errorf("NewCategory() failed")
	}
}

func TestNewAccount(t *testing.T) {
	got := NewAccount(1, "user")
	if got.ID != 1 || got.Username != "user" {
		t.Errorf("NewAccount() failed")
	}
}

func TestNewHardware(t *testing.T) {
	data := HardwareData{ID: 1, OS: "Win 10", Type: "Minimum"}
	got := NewHardware(data)
	if got.ID != 1 || got.OS != "Win 10" || got.Type != "Minimum" {
		t.Errorf("NewHardware() failed mapping fields")
	}
}

func TestNewMedia(t *testing.T) {
	t.Run("video type mapping", func(t *testing.T) {
		got := NewMedia(1, "https://youtube.com/watch?v=123")
		if got.Type != "video" {
			t.Errorf("NewMedia() expected type 'video', got %v", got.Type)
		}
	})
	t.Run("image type mapping", func(t *testing.T) {
		got := NewMedia(2, "https://example.com/image.png")
		if got.Type != "image" {
			t.Errorf("NewMedia() expected type 'image', got %v", got.Type)
		}
	})
}

func TestNewOrder(t *testing.T) {
	got := NewOrder(1, 2, testTime, "Cash", 50.0)
	if got.ID != 1 || got.UserID != 2 || time.Time(got.OrderDate) != testTime || got.PaymentMethod != "Cash" || got.FinalPrice != 50.0 {
		t.Errorf("NewOrder() failed mapping fields")
	}
	if got.Games == nil {
		t.Errorf("NewOrder() did not initialize Games slice")
	}
}
