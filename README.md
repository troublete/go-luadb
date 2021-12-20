# luadb
> a Lua extension for database access

## Introduction

This extension provides database access. It currently supports:

* PostgreSQL

## Usage

```lua
local db = require('luadb')
db.connect_postgres('postgresql://...') -- connect via connection string

ok = db.ping() -- returns ok if ping successful; or throws an error
rows = db.query('select * from ...') -- runs query and returns rows as table of tables; or throws an error
lastId, rowsAffected = db.exec('insert into ...') -- executes query, and returns state vars as integers; or throws an error
```

## Build

Assumes a Lua version is installed as static lib (`liblua.a`) (which is
standard at least for Lua 5.4) in `/usr/local/lib/`.

To build run:

```bash
make build
```

## DBMS

### PostgreSQL

* Supports mapping for most generic types 
	* `SMALLINT`, `INT`, `BIGINT` are returned as integer (which correlates to C double double)
	* `REAL`, `DOUBLE` are returned as number (which correlates to C double)
	* `CHAR`, `VARCHAR`, `TEXT` are returned as string
	* `DATE`, `TIME`, `TIMETZ`, `TIMESTAMP`,  `TIMESTAMPTZ` are returned a ISO8601 representing string
	* `BOOL` is returned as boolean
	* `BYTEA` is returned as "byte table"
* Supports mapping of `JSON` and `JSONB` fields to Lua tables
* Supports mapping of `NUMERIC` and `DECIMAL` to numbers
* All not mapable data is returned as "byte table"

**Todo**

- [ ] `hstore` support

## Contribute

It would be nice to grow this library to include also support for other
relational DBMS's. 

Contributions are welcome. 