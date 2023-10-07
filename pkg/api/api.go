package api

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/idkarn/curio-db/pkg/common"
	"github.com/idkarn/curio-db/pkg/middleware"
)

func SetupRouting(routes []middleware.Route) {
	for _, route := range routes {
		currentRoute := route
		http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			middleware.HandleWith(w, r, currentRoute)
		})
	}
}

func HealthHandler(ctx middleware.RequestContext) {
	log.Println("Health is being checked")
	ctx.Send("ok")
}

func NewRowHandler(ctx middleware.RequestContext) {
	var data NewRow
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	if data.Table >= TableIdType(len(Store.Tables)) {
		ctx.Error(ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	dataColumns := make(map[TableColumn]interface{})
	for key, val := range data.Columns {
		colIdx, err := FindColumnByName(data.Table, key)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		currentColumn := Store.TablesMetaData[data.Table].Columns[colIdx]
		dataColumns[currentColumn] = val
	}

	var newRowId RowIdType
	newRowId, err := AddNewRow(data.Table, dataColumns)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Send(newRowId)

	log.Printf(fmt.Sprintf("%s\n", ResponseStrings["R0"]), newRowId)
}

func GetRowHandler(ctx middleware.RequestContext) {
	var data GetRow
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	tid := data.Table

	row, err := GetRowById(tid, data.Id)
	if err != nil {
		ctx.Error(err.Error(), http.StatusNotFound)
		return
	}

	userRow := Row[string]{
		Id:      row.Id,
		Columns: make(map[string]interface{}),
	}

	for id, val := range row.Columns {
		name := Store.TablesMetaData[tid].Columns[id].Name
		userRow.Columns[name] = val
	}

	ctx.SendJSON(userRow)

	log.Printf(fmt.Sprintf("%s\n", ResponseStrings["R0"]), data.Id)
}

func NewColumnHandler(ctx middleware.RequestContext) {
	var data NewColumn
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	if data.Table >= TableIdType(len(Store.Tables)) {
		ctx.Error(ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	var colType uint8
	for key, val := range ColumnsTypeEnum {
		if val == data.Type {
			colType = uint8(key)
			break
		}
	}

	newColumnId, err := Store.TablesMetaData[data.Table].CreateNewColumn(data.Name, colType)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Send(newColumnId)
}

func UpdateRowHandler(ctx middleware.RequestContext) {
	var data UpdateRowData
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	if data.Table >= TableIdType(len(Store.Tables)) {
		ctx.Error(ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	dataColumns := make(map[ColumnIdType]interface{})
	for key, val := range data.Colunms {
		colIdx, err := FindColumnByName(data.Table, key)
		if err != nil {
			ctx.Error(err.Error(), http.StatusBadRequest)
			return
		}
		dataColumns[colIdx] = val
	}

	if err := Store.Tables[data.Table].UpdateRow(data.Id, dataColumns); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	ctx.Send(0)
}

func DeleteRowHandler(ctx middleware.RequestContext) {
	var data DeleteRowType
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	if data.Table >= TableIdType(len(Store.Tables)) {
		ctx.Error(ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	if err := Store.Tables[data.Table].DeleteRow(data.Id); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%q\n", Store.Tables[data.Table].Rows)

	ctx.Send(0)
}

func GetAllRowsHandler(ctx middleware.RequestContext) {
	var data GetAllRows
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	tid := data.Table

	if tid >= TableIdType(len(Store.Tables)) {
		ctx.Error(ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	rows := Store.Tables[tid].GetAllRows()

	var resp []Row[string]
	for _, row := range rows {
		newRow := Row[string]{
			Id:      row.Id,
			Columns: make(map[string]interface{}),
		}

		for id, val := range row.Columns {
			name := Store.TablesMetaData[tid].Columns[id].Name
			newRow.Columns[name] = val
		}

		resp = append(resp, newRow)
	}

	ctx.SendJSON(resp)
}
