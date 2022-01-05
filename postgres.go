package main

/*
#cgo LDFLAGS: -L/usr/local/lib -llua
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <stdlib.h>
*/
import "C"
import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
	"unsafe"
)

func map_postgres_row(L *C.lua_State, col string, value interface{}, columnType *sql.ColumnType) {
	column := C.CString(col)
	defer C.free(unsafe.Pointer(column))
	C.lua_pushstring(L, column) // set column name

	// map column value
	if v, ok := value.(int64); ok == true {
		C.lua_pushinteger(L, C.longlong(v))
	}

	if v, ok := value.(float64); ok == true {
		C.lua_pushnumber(L, C.double(v))
	}

	if v, ok := value.(string); ok == true {
		vp := C.CString(v)
		defer C.free(unsafe.Pointer(vp))
		C.lua_pushstring(L, vp)
	}

	if v, ok := value.(time.Time); ok == true {
		vp := C.CString(v.Format(time.RFC3339))
		defer C.free(unsafe.Pointer(vp))
		C.lua_pushstring(L, vp)
	}

	if v, ok := value.(bool); ok == true {
		if v == true {
			C.lua_pushboolean(L, 1)
		} else {
			C.lua_pushboolean(L, 0)
		}
	}

	if v, ok := value.([]byte); ok == true {
		switch {
		case columnType.DatabaseTypeName() == "JSONB":
			fallthrough
		case columnType.DatabaseTypeName() == "JSON":
			data := make(map[string]interface{})
			err := json.Unmarshal(v, &data)
			error_handling(L, err)

			go_to_lua(L, data)
		case columnType.DatabaseTypeName() == "NUMERIC":
			fallthrough
		case columnType.DatabaseTypeName() == "DECIMAL":
			val, err := strconv.ParseFloat(string(v), 8)
			error_handling(L, err)

			C.lua_pushnumber(L, C.double(val))
		default:
			C.lua_createtable(L, 0, 0)

			for idx, b := range v {
				C.lua_pushinteger(L, C.longlong(int64(idx+1)))
				C.lua_pushinteger(L, C.longlong(b))
				C.lua_settable(L, -3)
			}

		}
	}

	if value == nil {
		C.lua_pushnil(L)
	}
}
