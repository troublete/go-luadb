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

func struct_to_table(L *C.lua_State, data interface{}) {
	if d, ok := data.(map[string]interface{}); ok == true {
		if len(d) > 0 {
			C.lua_createtable(L, 0, 0)

			for k, v := range d {
				key := C.CString(k)
				defer C.free(unsafe.Pointer(key))

				C.lua_pushstring(L, key)

				if val, ok := v.(string); ok == true {
					value := C.CString(val)
					defer C.free(unsafe.Pointer(value))

					C.lua_pushstring(L, value)
				}

				if val, ok := v.(int64); ok == true {
					C.lua_pushinteger(L, C.longlong(val))
				}

				if val, ok := v.(float64); ok == true {
					C.lua_pushnumber(L, C.double(val))
				}

				if val, ok := v.(bool); ok == true {
					if val == true {
						C.lua_pushboolean(L, 1)
					} else {
						C.lua_pushboolean(L, 0)
					}
				}

				if val, ok := v.(map[string]interface{}); ok == true {
					struct_to_table(L, val)
				}

				if val, ok := v.(map[int]interface{}); ok == true {
					struct_to_table(L, val)
				}

				if v == nil {
					C.lua_pushnil(L)
				}

				C.lua_settable(L, -3)
			}
		} else {
			C.lua_pushnil(L)
		}
	} else {
		C.lua_pushnil(L)
	}
}
