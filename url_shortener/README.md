# URL Shortener

_Inspired by [Gophercises](https://courses.calhoun.io/courses/cor_gophercises)_

HTTP Server that redirects to endpoint by mapping paths to endpoint urls

Mapping is given through json and/or yaml files and/or postgres db

## Getting Started

```bash
go run main.go [options]
```

### Setup

* Put mapping `path -> url` into json/yaml file (see [json_example](mapping.json), [yaml_example](mapping.yaml))
  Or into postgres db (columns: {path: VARCHAR, url: VARCHAR})
* If multiple mapping sources are specified (json, yaml, postgres) then mapping resolves in the following order:
  1. Postgres 2. JSON 3. Yaml
* Run program specifying paths to the json and/or yaml files using `-json` and `-yaml` flags respectively

  ```go run main.go -json=mapping.json -yaml=mapping.yaml -pg=table_name```
