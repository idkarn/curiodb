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

	dataColumns := make(map[ColumnIdType]interface{})
	for key, val := range data.Columns {
		colIdx, err := FindColumnByName(data.Table, key)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}
		dataColumns[colIdx] = val
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

	userRow, err := SearchForRecords(data.Table, data.Filter)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	ctx.SendJSON(userRow)

	// log.Printf(fmt.Sprintf("%s\n", ResponseStrings["R0"]), -1)
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

	rows, err := SearchForRecords(data.Table, data.Filter)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	var failed []RowIdType
	for _, row := range rows { // FIXME: unrevertable changes
		dataColumns := make(map[ColumnIdType]interface{})
		for key, val := range data.Colunms {
			colIdx, err := FindColumnByName(data.Table, key)
			if err != nil {
				ctx.Error(err.Error(), http.StatusBadRequest)
				return
			}
			dataColumns[colIdx] = val
		}

		if err := Store.Tables[data.Table].UpdateRow(row.Id, dataColumns); err != nil {
			failed = append(failed, row.Id)
		}
	}

	if len(failed) == 0 {
		ctx.SendJSON(map[string]any{
			"ok":     true,
			"failed": nil,
		})
	} else {
		ctx.SendJSON(map[string]any{
			"ok":     false,
			"failed": failed,
		})
	}
}

func DeleteRowHandler(ctx middleware.RequestContext) {
	var data DeleteRowType
	if err := ctx.Read(&data); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	rows, err := SearchForRecords(data.Table, data.Filter)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	var failed []RowIdType
	for _, row := range rows {
		if err := Store.Tables[data.Table].DeleteRow(row.Id); err != nil {
			failed = append(failed, row.Id)
		}
	}

	if len(failed) == 0 {
		ctx.SendJSON(map[string]any{
			"ok":     true,
			"failed": nil,
		})
	} else {
		ctx.SendJSON(map[string]any{
			"ok":     false,
			"failed": failed,
		})
	}
}
