# Takara

An experimental finance management application for Personal Use and to Demystify Financial welbeing of the family

## Scope

A general tracker for daily expense, tracking assets and liability, and analysing spending pattern for future welfare

## Funtional Requirement
- Visualize my daily expenses and liability
- Analyse trends in expenses and assets
- Plan for saving for the future
- Update everyday at the end of the day
- Multiple Users at home

## Non Functional Requirement
- Easy of use
- Runnable @ home

## Stack

#### Frontend
- React - Already worked with react and its easy to get started
- Typescript - I like Typescript 💗
- React complier - IDK why not (This is experimental anyway)
- antd + antv - less styling more learning

#### Backend
- Go - I like Go and no need to over complicate this
- Gin - Framework to reduce boiler plate
- modernc.org/sqlite + database/sql - No CGO requirement and works with database/sql module from Go
- sqlx - Helps convert rows into typed structs datas 

#### Database
SQLite - simple and easy to use

### Current 
- Writing CRUD for Table types (whether to write it or just write custom handlers for each table)
  - Implemented CRUD with reflection and something like this:

    ```go
    type crudModel struct {
        Table   string
        Columns []string
        Model   any
    }

    var models = map[string]crudModel{
        "user": {
            Table: "user",
            Columns: []string{
                "name",
                "email",
                "password",
                "active",
            },
            Model: schema.User{},
        },
    }
    ```
    - This is absolutely horrible because we have to modify this array when ever i add a new field and it becomes unmanageable quickly
    - Started learning about go generics
    - Decided to write CRUD for tables because there might just be three tables
    

### References
- [`Why Sqlite With Go` - oneuptime.com](https://oneuptime.com/blog/post/2026-02-02-sqlite-go/view#why-sqlite-with-go)
- [`Quick Start for Gin`- gin-gonic.com](https://gin-gonic.com/en/docs/quickstart/)
- [`Go Database SQL Guide` - go-database-sql.org](http://go-database-sql.org/)
- [`Guide to SQLX` - jmoiron.github.io/sqlx/](https://jmoiron.github.io/sqlx/)
- [`SQL Injection` - go.dev](https://go.dev/doc/database/sql-injection)
- [`Mastering Go Generics`- medium.com](https://medium.com/hprog99/mastering-generics-in-go-a-comprehensive-guide-4d05ec4b12b)
- [`Are breadcrumbs still fresh for ux` - medium.com](https://medium.com/madison-ave-collective/are-breadcrumbs-still-fresh-for-ux-6e72b0f96e9b)
