package api

import (
	"errors"
	"strconv"
	"strings"

	"github.com/idkarn/curio-db/pkg/common"
)

const EqualOperator = '='
const NotOperator = '!'
const LessOperator = '<'
const GreaterOperator = '>'
const ContainOperator = '.'

func SearchForRecords(tid common.TableIdType, filter common.FilterType) ([]common.Row[string], error) {
	if tid >= common.TableIdType(len(common.Store.Tables)) {
		return nil, errors.New(common.ResponseStrings["T1"])
	}

	rows := []common.Row[string]{}
	columnsMeta := common.Store.TablesMetaData[tid].Columns
loop:
	for _, row := range common.Store.Tables[tid].Rows {
		var cols = make(map[string]interface{})
		for field, conds := range filter {
			for id, col := range row.Columns {
				colType := columnsMeta[id].Type
				if field == "id" {
					col = float64(row.Id)
					colType = 0
				} else if columnsMeta[id].Name != field {
					continue
				}

				for _, cond := range conds {
					op := cond[0]
					if !checkCondition(colType, col, op) {
						return nil, errors.New("wrong condition")
					}

					val, err := convert(cond[1:], colType)
					if err != nil {
						return nil, err
					}

					ok := processOperation(op, colType, col, val)
					if !ok {
						continue loop
					}
				}

				if field == "id" {
					break
				}
			}
		}

		for colid, val := range row.Columns {
			cols[columnsMeta[colid].Name] = val
		}

		rows = append(rows, common.Row[string]{
			Id:      row.Id,
			Columns: cols,
		})
	}
	return rows, nil
}

func convert(cond string, typeId uint8) (any, error) {
	var val any
	var err error
	if typeId == 0 {
		val, err = strconv.ParseFloat(cond, 64)
		if err != nil {
			return nil, errors.New("wrong number")
		}
	} else if typeId == 1 {
		val = cond
	} else if typeId == 2 {
		if cond == "false" {
			val = false
		} else if cond == "true" {
			val = true
		} else {
			return nil, errors.New("only true and false are allowed")
		}
	}
	return val, err
}

func checkCondition(typ uint8, col interface{}, op byte) bool {
	if op != EqualOperator && op != NotOperator && op != LessOperator && op != GreaterOperator && op != ContainOperator {
		return false
	}
	if typ == 0 {
		if op == ContainOperator {
			return false
		}
	}
	if typ == 2 {
		if op == LessOperator || op == GreaterOperator || op == ContainOperator {
			return false
		}
	}
	return true
}

func processOperation(op byte, typ uint8, a, b any) bool {
	result := false

	switch op {
	case EqualOperator:
		if a == b {
			result = true
		}
	case NotOperator:
		if a != b {
			result = true
		}
	case LessOperator:
		if typ == 0 && a.(float64) < b.(float64) {
			result = true
		} else if typ == 1 && strings.HasPrefix(a.(string), b.(string)) {
			result = true
		}
	case GreaterOperator:
		if typ == 0 && a.(float64) > b.(float64) {
			result = true
		} else if typ == 1 && strings.HasSuffix(a.(string), b.(string)) {
			result = true
		}
	case ContainOperator:
		if strings.Contains(a.(string), b.(string)) {
			result = true
		}
	}

	return result
}

// func filterParser(conds []string) ([]any, int) {
// 	parsed := []any{}
// 	for _, cond := range conds {
// 		for i, o := range OperatorEnum {
// 			if cond[0] == o {
// 				parsed = append(parsed, i)
// 				break
// 			}
// 		}
// 	}
// 	return nil, 0
// }
