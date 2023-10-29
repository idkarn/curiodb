package common

import (
	"fmt"
)

var ColumnsTypeEnum = [3]string{
	"number",
	"string",
	"bool",
}

var Store DatabaseStore

func GetRowById(tid TableIdType, id RowIdType) (Row[ColumnIdType], error) {
	if tid >= TableIdType(len(Store.Tables)) {
		return Row[ColumnIdType]{}, fmt.Errorf(ResponseStrings["T1"])
	}
	if id >= RowIdType(len(Store.Tables[tid].Rows)) {
		return Row[ColumnIdType]{}, fmt.Errorf(ResponseStrings["R1"])
	}

	return Store.Tables[tid].Rows[id], nil
}

func FindColumnByName(tid TableIdType, name string) (ColumnIdType, error) {
	for idx, col := range Store.TablesMetaData[tid].Columns {
		if col.Name == name {
			return ColumnIdType(idx), nil
		}
	}
	return 0, fmt.Errorf(ResponseStrings["C2"])
}

func AddNewRow(tid TableIdType, cols map[ColumnIdType]interface{}) (RowIdType, error) {
	var newRow Row[ColumnIdType]
	newRow.Id = RowIdType(len(Store.Tables[tid].Rows))
	newRow.Columns = make(map[ColumnIdType]interface{})

	// checking for type & assigning values to the row
	for cid, val := range cols {
		switch Store.TablesMetaData[tid].Columns[cid].Type {
		case 0:
			switch val.(type) {
			case int:
				newRow.Columns[cid] = float64(val.(int))
			case float64:
				newRow.Columns[cid] = val.(float64)
			}
		case 1:
			newRow.Columns[cid] = val.(string)
		case 2:
			newRow.Columns[cid] = val.(bool)
		}
	}

	Store.Tables[tid].Rows = append(Store.Tables[tid].Rows, newRow)

	return newRow.Id, nil
}

func (table *TableMetaData) CreateNewColumn(name string, colType uint8) (ColumnIdType, error) {
	columns := table.Columns
	newColumn := TableColumn{
		Id:         ColumnIdType(len(columns)),
		Name:       name,
		Type:       colType,
		IsOptional: false,
	}
	(*table).Columns = append(columns, newColumn)

	return newColumn.Id, nil
}

func (t Table) UpdateRow(rid RowIdType, data map[ColumnIdType]interface{}) error {
	if rid >= RowIdType(len(t.Rows)) {
		return fmt.Errorf(ResponseStrings["R1"])
	}

	// TODO: add data type checking
	for key := range t.Rows[rid].Columns {
		t.Rows[rid].Columns[key] = data[key]
	}

	return nil
}

func (t *Table) DeleteRow(rid RowIdType) error {
	if rid >= RowIdType(len(t.Rows)) {
		return fmt.Errorf(ResponseStrings["R1"])
	}

	t.Rows = append(t.Rows[:rid], t.Rows[rid+1:]...)

	return nil
}

func (t Table) GetAllRows() []Row[ColumnIdType] {
	return t.Rows
}
