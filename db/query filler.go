package db

import (
	"github.com/rohan5564/sakila-inventario/sakila"

	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

/****************************************************************
**	Functions
****************************************************************/

func scanFilmResult(result *sakila.FilmResult, rows *sql.Rows) error {
	for rows.Next() {
		tmp := new(sakila.Film)
		if err := rows.Scan(
			&tmp.ID,
			&tmp.Title,
			&tmp.Categories,
			&tmp.Description,
			&tmp.Year,
			&tmp.Lang,
			&tmp.OriginalLang,
			&tmp.RentalDuration,
			&tmp.RentalPrice,
			&tmp.Lenght,
			&tmp.ReplacementCost,
			&tmp.Rating,
			&tmp.SpecialFeatures,
			&tmp.LastUpdate,
			&tmp.Actors); err != nil {
			return err
		}

		*result = append(*result, tmp)
	}
	return nil
}

// fromJSONArray returns a JSON formatted string and an error if occurs
func fromJSONArray(values [][]string, keys ...string) (string, error) {
	var str string
	lk := len(keys)
	lvs := len(values)
	if lk == 0 || lvs == 0 {
		return "null", errors.New("wrong parameters input")
	}
	lfv := len(values[0])

	str += "'["
	for i, value := range values {
		if lv := len(value); lfv != lv || lk != lv {
			return "null", fmt.Errorf("can't do the mapping with %d keys and %d subvalues in each values", lk, lv)
		}
		str += "{"
		for j, key := range keys {
			str += fmt.Sprintf("\"%s\":\"%s\"", key, value[j])
			if j+1 < lk {
				str += ","
			}
		}
		str += "}"
		if i+1 < lvs {
			str += ","
		}
	}
	str += "]'"
	return str, nil
}

func call(interfaces ...interface{}) (str string) {
	for i, ifc := range interfaces {
		switch ifc.(type) {
		case string:
			if val := ifc.(string); val == "" {
				str += "null"
			} else {
				str += "'" + ifc.(string) + "'"
			}
		case *sakila.NullString:
			ns := ifc.(*sakila.NullString)
			str += ns.ToString()

		case sakila.NullStringArray,
			sakila.NullTime:
			str += fmt.Sprintf("%s", ifc)

		case sakila.NullJSONArray:
			data := ifc.(sakila.NullJSONArray)
			names := sakila.SplitValues(" ", *data.Array...)
			actors, _ := fromJSONArray(names, data.Keys...)
			str += actors

		case sakila.NullSet,
			sakila.NullInt32:
			str += fmt.Sprintf("%s", ifc)

		default:
			if reflect.ValueOf(ifc).IsZero() {
				str += "null"
			} else {
				str += fmt.Sprintf("%v", ifc)
			}
		}
		if i+1 < len(interfaces) {
			str += ","
		}
	}
	return
}
