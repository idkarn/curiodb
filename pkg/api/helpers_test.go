package api

import (
	"reflect"
	"testing"

	"github.com/idkarn/curio-db/pkg/common"
)

func newRow(id common.RowIdType, name string, age int) common.Row[string] {
	return common.Row[string]{
		Id: id,
		Columns: map[string]interface{}{
			"name": name,
			"age":  float64(age),
		},
	}
}

func config() {
	common.Config(common.DatabaseStore{
		Tables: []common.Table{{
			Id:   0,
			Rows: make([]common.Row[common.ColumnIdType], 0),
		}},
		TablesMetaData: []common.TableMetaData{
			{Name: "", Columns: []common.TableColumn{
				{Id: 0, Name: "name", Type: 1, IsOptional: false},
				{Id: 1, Name: "age", Type: 0, IsOptional: false},
			}},
		},
	})
	var defaultRows = [][2]any{
		{"none", 0},
		{"null", 100},
		{"noname", 42},
	}
	for _, cols := range defaultRows {
		common.AddNewRow(0, map[common.ColumnIdType]interface{}{
			0: cols[0],
			1: cols[1],
		})
	}
}

func TestSearchForRecordsEqual(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {"=none"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(0, "none", 0),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsStartsWith(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {"<no"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(0, "none", 0), newRow(2, "noname", 42),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsContains(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {".l"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(1, "null", 100),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsNot(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {"!null"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(0, "none", 0),
		newRow(2, "noname", 42),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsMultiple(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {"<n", ">me"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(2, "noname", 42),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsMultiField(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"name": {">e"},
		"age":  {"!42"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(0, "none", 0),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}

func TestSearchForRecordsById(t *testing.T) {
	config()
	rows, err := SearchForRecords(0, common.FilterType{
		"id": {"=1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []common.Row[string]{
		newRow(1, "null", 100),
	}
	if !reflect.DeepEqual(rows, expected) {
		t.Fatalf("expected: %+v, but returned %+v", expected, rows)
	}
}
