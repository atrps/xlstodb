# xlstodb
Loading data from Excel files

## Install

### PostgreSQL

Run scripts from xlstodb/scripts for create table and procedure
```shell
1_Create_tbl.sql
2_Create_proc.sql
3_Fill_tbl.sql
```

Set params for database in xlstodb/configs/.env 

## Usage

### Sample

Simple loading from files xlstodb/examples
```shell
go run cmd/sample/main.go
```

### User interface

Run application
```shell
go run cmd/choice/choice.go
```

In your browser
```shell
http://localhost:8080
```
