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
	"unsafe"
)

// go_to_lua converts a arbitrary go value into a Lua value pushed onto the
// stack, it's main purpose here is to convert a generic unmarshaled JSON
// value into a Lua value
func go_to_lua(L *C.lua_State, data interface{}) {
	if d, ok := data.(map[string]interface{}); ok == true {
		if len(d) > 0 {
			C.lua_createtable(L, 0, 0)

			for k, v := range d {
				key := C.CString(k)
				defer C.free(unsafe.Pointer(key))

				C.lua_pushstring(L, key)
				go_to_lua(L, v)
				C.lua_settable(L, -3)
			}
		} else {
			C.lua_pushnil(L)
		}
	} else if d, ok := data.(map[int]interface{}); ok == true {
		if len(d) > 0 {
			C.lua_createtable(L, 0, 0)

			for k, v := range d {
				C.lua_pushinteger(L, C.longlong(k))
				go_to_lua(L, v)
				C.lua_settable(L, -3)
			}
		} else {
			C.lua_pushnil(L)
		}
	} else if val, ok := data.(string); ok == true {
		value := C.CString(val)
		defer C.free(unsafe.Pointer(value))

		C.lua_pushstring(L, value)
	} else if val, ok := data.(int64); ok == true {
		C.lua_pushinteger(L, C.longlong(val))
	} else if val, ok := data.(float64); ok == true {
		C.lua_pushnumber(L, C.double(val))
	} else if val, ok := data.(bool); ok == true {
		if val == true {
			C.lua_pushboolean(L, 1)
		} else {
			C.lua_pushboolean(L, 0)
		}
	} else if val, ok := data.(map[string]interface{}); ok == true {
		go_to_lua(L, val)
	} else if val, ok := data.(map[int]interface{}); ok == true {
		go_to_lua(L, val)
	} else if data == nil {
		C.lua_pushnil(L)
	} else {
		C.lua_pushnil(L)
	}
}
