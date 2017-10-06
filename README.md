# Armory

The armory backend is responsible for storing and parsing d2s characters.

## MongoDB

### Start

```bash
$ mongod -dbpath=/path/to/mongodata
```

## Flags
These are flags that can be overridden when running the binary.

| Name     | Default                  |
|----------|--------------------------|
| listen   | `127.0.0.1:8090`         |
| db.url   | `127.0.0.1`              |
| db.name  | `armory`                 |
| d2s.path | `/home/slash/characters` |
