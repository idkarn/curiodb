package common

import "os"

type TableIdType uint8
type ColumnIdType uint32
type RowIdType uint64
type FilterType map[string][]string

type NewRow struct {
	Columns map[string]interface{} `json:"columns"`
	Table   TableIdType            `json:"table"`
}

type IDecodedJson interface {
	NewRow | GetRow | NewColumn | UpdateRowData | DeleteRowType
}

type filter struct {
	Filter FilterType `json:"filter"`
}

type GetRow struct {
	Table TableIdType `json:"table"`
	filter
}

type UpdateRowData struct {
	Table   TableIdType            `json:"table"`
	Colunms map[string]interface{} `json:"columns"`
	filter
}

type DeleteRowType struct {
	Table TableIdType `json:"table"`
	filter
}

type NewColumn struct {
	Name  string      `json:"name"`
	Table TableIdType `json:"table"`
	Type  string      `json:"type"`
}

type TableColumn struct {
	Id         ColumnIdType `json:"id"`
	Name       string       `json:"name"`
	Type       uint8        `json:"type"`
	IsOptional bool         `json:"isoptional"`
}

type Row[T ColumnIdType | string] struct {
	Id      RowIdType         `json:"id"`
	Columns map[T]interface{} `json:"columns"`
}

type TableMetaData struct {
	Name    string
	Columns []TableColumn
}

type Table struct {
	Id   TableIdType
	Rows []Row[ColumnIdType]
}

type DatabaseStore struct {
	Tables         []Table
	TablesMetaData []TableMetaData
}

type File struct {
	Path     string
	Desc     *os.File
	Content  []byte
	IsOpened bool
}
