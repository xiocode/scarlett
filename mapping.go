/**
 * Author:        Tony.Shao
 * Email:         xiocode@gmail.com
 * Github:        github.com/xiocode
 * File:          mapping.go
 * Description:   mapping
 */

package scarlett

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	STRUCT_TAG = "scarlett"
)

type FieldsAttribute struct {
	index int
}

type Record struct {
	columns []string
	fields  map[string]*FieldsAttribute
}

type Binder func(rows []map[string]interface{}, target interface{}) error

// CustomScanner binds a database column value to a Go type
type scanner struct {
	// After a row is scanned, Holder will contain the value from the database column.
	// Initialize the CustomScanner with the concrete Go type you wish the database
	// driver to scan the raw column into.
	rows []map[string]interface{}
	// Target typically holds a pointer to the target struct field to bind the Holder
	// value to.
	target interface{}
	// Binder is a custom function that converts the holder value to the target type
	// and sets target accordingly.  This function should return error if a problem
	// occurs converting the holder to the target.
	binder Binder
}

// Bind is called automatically by gorp after Scan()
func (s scanner) Bind() error {
	if s.binder == nil {
		return binder(s.rows, s.target)
	}
	return s.binder(s.rows, s.target)
}

func (s *Scarlett) Scan(dst interface{}, binder Binder, rows *sql.Rows) error {
	// make sure we always close rows
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	var holders []map[string]interface{}
	for rows.Next() {
		holder := make(map[string]interface{})
		var fields []interface{}
		for i := 0; i < len(columns); i++ {
			var field interface{}
			fields = append(fields, &field)
		}
		err := rows.Scan(fields...)
		if err != nil {
			return err
		}
		for i, column := range columns {
			holder[column] = fields[i]
		}
		holders = append(holders, holder)
	}

	scan := &scanner{
		rows:   holders,
		target: dst,
		binder: binder,
	}

	return scan.Bind()
}

func binder(rows []map[string]interface{}, target interface{}) error {
	record, err := analysis(target)
	if err != nil {
		return err
	}
	structVal := reflect.ValueOf(target).Elem()
	for _, filed := range record.columns {
		for _, targets := range rows {
			if len(targets) != len(record.columns) {
				return fmt.Errorf("scarlett.binder: mismatch in number of columns (%d) and targets (%s)",
					len(record.columns), len(targets))
			}
			if value, ok := targets[filed]; ok {
				fieldAddr := structVal.Field(record.fields[filed].index).Addr().Interface()
				fmt.Println(reflect.TypeOf(fieldAddr))
				fmt.Println(strreflect.Indirect(reflect.ValueOf(value)))
			} else {
				fmt.Print("NOTHING!")
			}
		}
	}
	return nil
}

func analysis(dst interface{}) (*Record, error) {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("scarlett called with non-pointer destination %v", dstType)
	}
	structType := dstType.Elem()
	if structType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("scarlett called with pointer to non-struct %v", dstType)
	}
	record := new(Record)
	record.fields = make(map[string]*FieldsAttribute)

	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)

		// skip non-exported fields
		if f.PkgPath != "" {
			continue
		}

		// examine the tag for metadata
		tag := strings.Split(f.Tag.Get(STRUCT_TAG), ",")

		// was this field marked for skipping?
		if len(tag) > 0 && tag[0] == "-" {
			continue
		}

		// default to the field name
		name := f.Name

		// the tag can override the field name
		if len(tag) > 0 && tag[0] != "" {
			name = tag[0]
		}

		if _, present := record.fields[name]; present {
			return nil, fmt.Errorf("scarlett found multiple fields for column %s", name)
		}

		record.columns = append(record.columns, name)
		record.fields[name] = &FieldsAttribute{
			index: i,
		}
	}

	return record, nil
}
