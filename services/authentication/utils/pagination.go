package utils

import (
	"database/sql"
	"gorm.io/gorm"
)

type Pagination struct {
	PageIndex int
	PageSize  int
}

func (m *Pagination) GetPageIndex() int {
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

func (m *Pagination) GetOffSet() int {
	offset := (m.GetPageIndex() - 1) * m.GetPageSize()
	return offset
}

func Paginate(pageSize, pageIndex int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (pageIndex - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//returns (1,columns 2,every row 3, err)
func Rows2Map(rows *sql.Rows) ([]string, []map[string]interface{}, error) {
	var rowMap []map[string]interface{}
	cols, _ := rows.Columns()

	for rows.Next() {
		record := map[string]interface{}{}

		columns := make([]interface{}, len(cols))
		columnsPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnsPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnsPointers...); err != nil {
			return cols, rowMap, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnsPointers[i].(*interface{})
			m[colName] = *val
		}

		for k, v := range m {
			var val interface{}
			switch v.(type) {
			case []uint8:
				val = B2S(v.([]uint8))
			case nil:
				val = ""
			default:
				val = v
			}

			record[k] = val
		}
		rowMap = append(rowMap, record)
	}

	return cols, rowMap, nil
}

func B2S(bs []uint8) string {
	var ba []byte
	for _, b := range bs {
		ba = append(ba, b)
	}
	return string(ba)
}
