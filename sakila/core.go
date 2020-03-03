package sakila

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Mapper JSON file parser into struct
type Mapper interface {
	FromJSON(jsonData []byte)
}

// DBdata database entry credentials
type DBdata struct {
	User, Pass, Schema, Protocol, IP string
	Port                             uint
}

// NullString Wrapper type for sql.NullString
type NullString sql.NullString

// NullStringArray Wrapper type for []string where is used
// a sql.NullString value for Split into the slice
type NullStringArray []string

// NullJSONArray Wrapper struct to handle json types in mysql
// procedures
type NullJSONArray struct {
	Array *NullStringArray
	Keys  []string
}

// NullSet Wrapper type for []string where is used
// a sql.NullString value for Split into the slice
type NullSet []string

// NullInt32 Wrapper type for sql.NullInt32
type NullInt32 sql.NullInt32

// NullTime Wrapper type for sql.NullTime
type NullTime sql.NullTime

// Film Contains information about films
type Film struct {
	ID              uint16          `json:"id"`
	Title           string          `json:"title"`
	Categories      NullStringArray `json:"categories,omitempty"`
	Description     *NullString     `json:"description,omitempty"`
	Year            NullInt32       `json:"year,omitempty"`
	Lang            string          `json:"language"`
	OriginalLang    *NullString     `json:"original language,omitempty"`
	RentalDuration  float32         `json:"rental duration"`
	RentalPrice     float32         `json:"rental rate"`
	Lenght          NullInt32       `json:"length,omitempty"`
	ReplacementCost float32         `json:"replacement cost"`
	Rating          *NullString     `json:"rating,omitempty"`
	SpecialFeatures NullSet         `json:"special features,omitempty"`
	LastUpdate      NullTime        `json:"last update,omitempty"`
	Actors          NullStringArray `json:"actors,omitempty"`
}

// Category contains information about categories
type Category struct {
	ID    uint16 `json:"id"`
	Name  string `json:"name"`
	Count uint16 `json:"film count"`
}

// Actor contains information about Actors
type Actor struct {
	ID        uint16          `json:"id"`
	FirstName string          `json:"first name"`
	LastName  string          `json:"last name"`
	Films     NullStringArray `json:"films"`
}

// FilmPagination search options
type FilmPagination struct {
	Limit   uint16
	Offset  uint16
	OrderBy string
	Order   string
}

// FilmResult Array of film pointers
type FilmResult []*Film

// CategoryResult Array of category pointers
type CategoryResult []*Category

// ActorResult Array of actor pointers
type ActorResult []*Actor

// QueryFilm contains a filter, search options, miscellaneous information and results
type QueryFilm struct {
	Filter  Film
	Options FilmPagination
	Result  FilmResult
	Info    ResultInformation
}

// ResultInformation miscellaneous information given to manage some JS scripts
type ResultInformation struct {
	SingleTitle    bool  `json:"-"`
	TotalData      int16 `json:"total data"`
	ActualPage     int16 `json:"actual page"`
	RemainingPages int16 `json:"remaining pages"`
	TotalPages     int16 `json:"total pages"`
}

// Statement main connection to database
var Statement *sql.DB

// Render implements Renderer interface from go-chi/render
func (f Film) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Render implements Renderer interface from go-chi/render
func (f FilmResult) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Render implements Renderer interface from go-chi/render
func (ri ResultInformation) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Render implements Renderer interface from go-chi/render
func (c Category) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Render implements Renderer interface from go-chi/render
func (a Actor) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Scan implements Scanner interface for type NullString
func (ns *NullString) Scan(value interface{}) error {
	var tmp sql.NullString
	if err := tmp.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ns = NullString{tmp.String, false}
	} else {
		*ns = NullString{tmp.String, true}
	}
	return nil
}

// Scan implements Scanner interface for type NullStringArray
func (nsa *NullStringArray) Scan(value interface{}) error {
	var tmp sql.NullString
	if err := tmp.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) != nil {
		*nsa = NullStringArray(strings.Split(tmp.String, ","))
	}
	return nil
}

// Scan implements Scanner interface for type NullStringArray
func (nset *NullSet) Scan(value interface{}) error {
	var tmp sql.NullString
	if err := tmp.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) != nil {
		*nset = NullSet(strings.Split(tmp.String, ","))
	}
	return nil
}

// Scan implements Scanner interface for type NullInt32
func (ni32 *NullInt32) Scan(value interface{}) error {
	var tmp sql.NullInt32
	if err := tmp.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ni32 = NullInt32{tmp.Int32, false}
	} else {
		*ni32 = NullInt32{tmp.Int32, true}
	}
	return nil
}

// Scan implements Scanner interface for type NullTime
func (nt *NullTime) Scan(value interface{}) error {
	var tmp sql.NullTime
	if err := tmp.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*nt = NullTime{tmp.Time, false}
	} else {
		*nt = NullTime{tmp.Time, true}
	}
	return nil
}

// MarshalJSON implements Marshaler interface for type NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return []byte(`null`), nil
}

// MarshalJSON implements Marshaler interface for type NullInt32
func (ni32 NullInt32) MarshalJSON() ([]byte, error) {
	if ni32.Valid {
		return json.Marshal(ni32.Int32)
	}
	return []byte(`null`), nil
}

// MarshalJSON implements Marshaler interface for type NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return []byte(`null`), nil
}

// UnmarshalJSON implements Unmarshaler interface for type NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if string(data) != `null` {
		ns.String = tmp
		ns.Valid = true
		return nil
	}
	ns.Valid = false
	return err
}

// UnmarshalJSON implements Unmarshaler interface for type NullStringArray
func (nsa *NullStringArray) UnmarshalJSON(data []byte) error {
	var tmp []string
	err := json.Unmarshal(data, &tmp)
	if string(data) != `null` {
		*nsa = tmp
		return nil
	}
	*nsa = []string{}
	return err
}

// UnmarshalJSON implements Unmarshaler interface for type NullSet
func (nset *NullSet) UnmarshalJSON(data []byte) error {
	var tmp []string
	err := json.Unmarshal(data, &tmp)
	if string(data) != `null` {
		*nset = tmp
		return nil
	}
	return err
}

// UnmarshalJSON implements Unmarshaler interface for type NullInt32
func (ni32 *NullInt32) UnmarshalJSON(data []byte) error {
	var tmp int32
	err := json.Unmarshal(data, &tmp)
	if string(data) != `null` {
		ni32.Int32 = tmp
		ni32.Valid = true
		return nil
	}
	ni32.Valid = false
	return err
}

// UnmarshalJSON implements Unmarshaler interface for type NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var tmp time.Time
	err := json.Unmarshal(data, &tmp)
	if string(data) != `null` {
		nt.Time = tmp
		nt.Valid = true
		return nil
	}
	nt.Valid = false
	return err
}

// GetResult returns json data as string from the results stored in the QueryFilm pointer
func (f *QueryFilm) GetResult() []byte {
	buffer := new(bytes.Buffer)
	var (
		data []byte
		err  error
	)

	switch {
	case len(f.Result) == 1:
		data, err = json.MarshalIndent(f.Result[0], "", "\t")
	case len(f.Result) > 1:
		data, err = json.MarshalIndent(f.Result, "", "\t")
	default:
		return nil
	}

	if err != nil {
		return nil
	}

	if _, err := buffer.Write(data); err != nil {
		return nil
	}

	return buffer.Bytes()
}

// ToString returns json data as string from the film
func (f Film) ToString() string {
	buffer := new(bytes.Buffer)
	var (
		data []byte
		err  error
	)

	data, err = json.MarshalIndent(f, "", "\t")

	if err != nil {
		return ""
	}

	if _, err := buffer.Write(data); err != nil {
		return ""
	}

	return buffer.String()
}

// FromJSON fills a film pointer using json data in bytes
func (f *Film) FromJSON(jsonData []byte) error {
	if err := json.Unmarshal(jsonData, f); err != nil {
		return err
	}
	return nil
}

// FromJSON fills a category pointer using json data in bytes
func (c *Category) FromJSON(jsonData []byte) error {
	if err := json.Unmarshal(jsonData, c); err != nil {
		return err
	}
	return nil
}

// FromJSON fills an actor pointer using json data in bytes
func (a *Actor) FromJSON(jsonData []byte) error {
	if err := json.Unmarshal(jsonData, a); err != nil {
		return err
	}
	return nil
}

// SplitValues function that allows to split a single string or an array of strings
// into a 2D array
func SplitValues(s string, arr ...string) (data [][]string) {
	for _, v := range arr {
		data = append(data, strings.Split(v, s))
	}
	return
}

// ToString is a replacement for String() override function since the struct has a
// nested string type
func (ns *NullString) ToString() string {
	if ns == nil || !ns.Valid {
		return "null"
	}
	return fmt.Sprintf("'%s'", ns.String)
}

func (nsa NullStringArray) String() string {
	if len(nsa) == 0 {
		return "null"
	}
	var array string
	for i, item := range nsa {
		array += fmt.Sprintf("\"%v\"", item)
		if i+1 < len(nsa) {
			array += ","
		}
	}
	return fmt.Sprintf("'[%s]'", array)
}

func (nset NullSet) String() string {
	if len(nset) == 0 {
		return "null"
	}
	return fmt.Sprintf("(\"%v\")", strings.Join(nset, ","))
}

func (ni32 NullInt32) String() string {
	if !ni32.Valid {
		return "null"
	}
	return fmt.Sprintf("%d", ni32.Int32)
}

func (nt NullTime) String() string {
	if !nt.Valid {
		return "null"
	}
	return fmt.Sprintf("'%v'", nt.Time)
}
