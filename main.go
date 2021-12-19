package main

/*
#cgo LDFLAGS: -L/usr/local/lib -llua
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <stdbool.h>
#include <stdlib.h>

typedef const char go_str;

extern void go_connect_postgres(lua_State *L, go_str *conn);
extern int go_ping(lua_State *L);
extern int go_query(lua_State *L, go_str *query);
extern int go_exec(lua_State *L, go_str *query);

static int ping(lua_State *L) {
	return go_ping(L);
}

static int query(lua_State *L) {
	go_str *query = luaL_checkstring(L, 1);
	return go_query(L, query);
}

static int exec(lua_State *L) {
	go_str *query = luaL_checkstring(L, 1);
	return go_exec(L, query);
}


static int connect_postgres(lua_State *L) {
	go_str *conn = luaL_checkstring(L, 1);
	go_connect_postgres(L, conn);
	return 0;
}

static int register_addon(lua_State *L) {
	lua_newtable(L);
	lua_pushstring(L, "connect_postgres");
	lua_pushcfunction(L, connect_postgres);
	lua_settable(L, -3);
	lua_pushstring(L, "ping");
	lua_pushcfunction(L, ping);
	lua_settable(L, -3);
	lua_pushstring(L, "query");
	lua_pushcfunction(L, query);
	lua_settable(L, -3);
	lua_pushstring(L, "exec");
	lua_pushcfunction(L, exec);
	lua_settable(L, -3);
	return 1;
}
*/
import "C"

import (
	"database/sql"
	"errors"
	"time"
	"unsafe"

	_ "github.com/lib/pq"
)

var database *sql.DB

func main() {}

func error_handling(L *C.lua_State, err error) {
	if err != nil {
		msg := C.CString(err.Error())
		defer C.free(unsafe.Pointer(msg))

		C.lua_pushstring(L, msg)
		C.lua_error(L)
	}
}

//export go_connect_postgres
func go_connect_postgres(L *C.lua_State, connectionString *C.go_str) {
	conn := C.GoString(connectionString)
	sql, err := sql.Open("postgres", conn)
	error_handling(L, err)
	database = sql
}

//export go_ping
func go_ping(L *C.lua_State) C.int {
	if database == nil {
		error_handling(L, errors.New("no database connected"))
	}

	err := database.Ping()
	error_handling(L, err)

	C.lua_pushboolean(L, C.int(1))
	return C.int(1)
}

//export go_query
func go_query(L *C.lua_State, query *C.go_str) C.int {
	if database == nil {
		error_handling(L, errors.New("no database connected"))
	}

	rows, err := database.Query(C.GoString(query))
	error_handling(L, err)

	if rows == nil {
		C.lua_pushnil(L)
		return C.int(1)
	}
	defer rows.Close()

	var columns []string
	columns, err = rows.Columns()
	error_handling(L, err)

	var columnTypes []*sql.ColumnType
	columnTypes, err = rows.ColumnTypes()
	error_handling(L, err)

	numberColumns := len(columns)
	i := 1

	C.lua_createtable(L, 0, 0)
	for rows.Next() {
		C.lua_pushinteger(L, C.longlong(i))
		C.lua_createtable(L, 0, 0)

		row := make([]interface{}, numberColumns)
		for i := range row {
			row[i] = &row[i]
		}

		err := rows.Scan(row...)
		error_handling(L, err)

		for i := range row {
			column := C.CString(columns[i])
			defer C.free(unsafe.Pointer(column))
			C.lua_pushstring(L, column)

			C.lua_createtable(L, 0, 0) // column value table
			ct := C.CString("type")
			defer C.free(unsafe.Pointer(ct))
			C.lua_pushstring(L, ct)

			columnType := C.CString(columnTypes[i].DatabaseTypeName())
			defer C.free(unsafe.Pointer(columnType))
			C.lua_pushstring(L, columnType)

			C.lua_settable(L, -3) // set type

			v := C.CString("value")
			defer C.free(unsafe.Pointer(v))
			C.lua_pushstring(L, v)

			value := row[i]
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
				C.lua_createtable(L, 0, 0)

				for idx, b := range v {
					C.lua_pushinteger(L, C.longlong(idx))
					C.lua_pushinteger(L, C.longlong(b))
					C.lua_settable(L, -3)
				}
			}

			C.lua_settable(L, -3) // set value
			C.lua_settable(L, -3) // set row
		}
		C.lua_settable(L, -3)
		i = i + 1
	}
	return C.int(1)
}

//export go_exec
func go_exec(L *C.lua_State, query *C.go_str) C.int {
	if database == nil {
		error_handling(L, errors.New("no database connected"))
	}

	result, err := database.Exec(C.GoString(query))
	error_handling(L, err)

	if result == nil {
		C.lua_pushnil(L)
		return 1
	}

	var lastId int64
	lastId, err = result.LastInsertId()
	error_handling(L, err)

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	error_handling(L, err)

	C.lua_pushinteger(L, C.longlong(lastId))
	C.lua_pushinteger(L, C.longlong(rowsAffected))
	return C.int(2)
}

//export luaopen_luadb
func luaopen_luadb(L *C.lua_State) C.int {
	return C.register_addon(L)
}
