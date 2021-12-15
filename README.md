# luadb
> a Lua extension for database access

## Introduction

This extension provides database access. It currently supports:

* PostgreSQL

## Usage

```lua
local db = require('luadb')
db.connect_postgres('postgresql://...') -- connect via connection string

ok = db.ping() -- returns ok if ping successful; or throws on error
rows = db.query('select * from ...') -- runs query and returns rows as table of tables; or throws on error
lastId, rowsAffected = db.exec('insert into ...') -- executes query, and returns state vars as numbers; or throws on error
```

## Build

Assumes ANY Lua version is installed as static lib (`liblua.a`) (which is
standard at least for Lua 5.4) in `/usr/local/lib/`.

```bash
make build
```

## Contribute

It would be nice to grow this library to include also support for other
relational DBMS's. To extend the support to any `database/sql` driver it is
only necessary to add a connector function and a custom value interpreter
(i.e. for special values). (see `main.go` for more details).

Contributions are welcome. 