package api

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	. "github.com/idkarn/curio-db/pkg/common"
)

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

type RouteHandler struct {
	HandlerFn     func(http.ResponseWriter, *http.Request)
	RequestStruct reflect.Type
}

func SetupRouting(routes []Route) {
	for _, route := range routes {
		currentRoute := route
		http.HandleFunc(currentRoute.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != currentRoute.Method {
				http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			} else {
				currentRoute.Handler(w, r)
			}
		})
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Health is being checked")
	fmt.Fprint(w, "ok")
}

func NewRowHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[NewRow](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.Table >= TableIdType(len(Store.Tables)) {
		http.Error(w, ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	dataColumns := make(map[TableColumn]interface{})
	for key, val := range data.Columns {
		colIdx, err := FindColumnByName(data.Table, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		currentColumn := Store.TablesMetaData[data.Table].Columns[colIdx]
		dataColumns[currentColumn] = val
	}

	var newRowId RowIdType
	newRowId, err = AddNewRow(data.Table, dataColumns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, newRowId)

	log.Printf(fmt.Sprintf("%s\n", ResponseStrings["R0"]), newRowId)
}

func GetRowHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[GetRow](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tid := data.Table

	row, err := GetRowById(tid, data.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	json, err := EncodeJson(userRow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, json)

	log.Printf(fmt.Sprintf("%s\n", ResponseStrings["R0"]), data.Id)
}

func NewColumnHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[NewColumn](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Table >= TableIdType(len(Store.Tables)) {
		http.Error(w, ResponseStrings["T1"], http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, newColumnId)
}

func UpdateRowHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[UpdateRowData](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Table >= TableIdType(len(Store.Tables)) {
		http.Error(w, ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	dataColumns := make(map[ColumnIdType]interface{})
	for key, val := range data.Colunms {
		colIdx, err := FindColumnByName(data.Table, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dataColumns[colIdx] = val
	}

	if err := Store.Tables[data.Table].UpdateRow(data.Id, dataColumns); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, 0)
}

func DeleteRowHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[DeleteRowType](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if data.Table >= TableIdType(len(Store.Tables)) {
		http.Error(w, ResponseStrings["T1"], http.StatusBadRequest)
		return
	}

	if err := Store.Tables[data.Table].DeleteRow(data.Id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%q\n", Store.Tables[data.Table].Rows)

	fmt.Fprint(w, 0)
}

func GetAllRowsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := DecodeJson[GetAllRows](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tid := data.Table

	if tid >= TableIdType(len(Store.Tables)) {
		http.Error(w, ResponseStrings["T1"], http.StatusBadRequest)
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

	json, err := EncodeJson(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, json)
}
