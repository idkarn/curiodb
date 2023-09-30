package common

type TableIdType uint8
type ColumnIdType uint32
type RowIdType uint64

type NewRow struct {
	Columns map[string]interface{} `json:"columns"`
	Table   TableIdType            `json:"table"`
}

type IDecodedJson interface {
	NewRow | GetRow | NewColumn | UpdateRowData | DeleteRowType | GetAllRows
}

type GetRow struct {
	Id    RowIdType   `json:"id"`
	Table TableIdType `json:"table"`
}

type GetAllRows struct {
	Table TableIdType `json:"table"`
}

type UpdateRowData struct {
	Table   TableIdType            `json:"table"`
	Id      RowIdType              `json:"id"`
	Colunms map[string]interface{} `json:"columns"`
}

type DeleteRowType struct {
	Table TableIdType `json:"table"`
	Id    RowIdType   `json:"id"`
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
