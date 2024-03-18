package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

type TextArray[T ~string] pgtype.FlatArray[T]

// Scan implements the [database/sql.Scanner] interface.
func (s *TextArray[T]) Scan(src any) error {
	var typedArray []string
	err := pgtype.NewMap().SQLScanner(&typedArray).Scan(src)
	if err != nil {
		return err
	}

	(*s) = make(TextArray[T], len(typedArray))
	for i, value := range typedArray {
		(*s)[i] = T(value)
	}

	return nil
}

type numberField interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~int | ~uint
}

type numberTypeField interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | int | uint
}

var _ sql.Scanner = (*NumberArray[int8])(nil)

type NumberArray[F numberField] pgtype.FlatArray[F]

// Scan implements the [database/sql.Scanner] interface.
func (a *NumberArray[F]) Scan(src any) (err error) {
	var (
		mapper  func()
		scanner sql.Scanner
	)

	//nolint: exhaustive
	// only defined types
	switch reflect.TypeOf(*a).Elem().Kind() {
	case reflect.Int8:
		mapper, scanner = castedScan[int8](a)
	case reflect.Uint8:
		mapper, scanner = castedScan[uint8](a)
	case reflect.Int16:
		mapper, scanner = castedScan[int16](a)
	case reflect.Uint16:
		mapper, scanner = castedScan[uint16](a)
	case reflect.Int32:
		mapper, scanner = castedScan[int32](a)
	case reflect.Uint32:
		mapper, scanner = castedScan[uint32](a)
	case reflect.Int64:
		mapper, scanner = castedScan[int64](a)
	case reflect.Uint64:
		mapper, scanner = castedScan[uint64](a)
	case reflect.Int:
		mapper, scanner = castedScan[int](a)
	case reflect.Uint:
		mapper, scanner = castedScan[uint](a)
	}

	if err = scanner.Scan(src); err != nil {
		return err
	}
	mapper()

	return nil
}

func castedScan[T numberTypeField, F numberField](a *NumberArray[F]) (func(), sql.Scanner) {
	var typedArray []T

	mapper := func() {
		(*a) = make(NumberArray[F], len(typedArray))
		for i, value := range typedArray {
			(*a)[i] = F(value)
		}
	}
	scanner := pgtype.NewMap().SQLScanner(&typedArray)

	return mapper, scanner
}

type Map[V any] map[string]V

// Scan implements the [database/sql.Scanner] interface.
func (m *Map[V]) Scan(src any) error {
	if src == nil {
		return nil
	}
	var arr []byte
	// create valid payload
	src = hex.AppendEncode([]byte("\\x"), src.([]byte))
	err := pgtype.NewMap().SQLScanner(&arr).Scan(src)
	if err != nil || len(arr) == 0 {
		return err
	}

	return json.Unmarshal(arr, &m)
}

// Value implements the [database/sql/driver.Valuer] interface.
func (m Map[V]) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
