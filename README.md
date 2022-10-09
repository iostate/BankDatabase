# BankDatabase Backend

## Go Code Review Comments

A list of code review comments of the Go code is found at the link below. These are to be followed like code conventions.

[Go Code Review Comments] (<https://github.com/golang/go/wiki/CodeReviewComments#gofmt>)

# SQLC

SQLC allows code to be generated for SQL queries.
CRUD operations can be done through SQLC. 
[SQLC 1.4.0 Documentation] (<https://github.com/kyleconroy/sqlc/tree/v1.4.0>)

## Explanation of Files

## db/migration

Contains the migration files: Migrate down = Drop all tables. Migrate up = create tables again.

## db/query

Contains the SQLC code that will be generated to Golang code.

## db/sqlc

Contains the code generated from sqlc.
The code is generated using the `make sqlc` command, which is a Makefile command.

## /bank_dbdiagram.sql

A SQL diagram that I built on dbdiagram.io

## go.mod

Must initialize a go mod in order for Queries to be recognized.

## sqlc.yaml

Configuration file for SQLC. Dictate the location of the files for: queries, the database, and the models.

## Resources

<https://pkg.go.dev/database/sql#DB.QueryContext>
<https://gorm.io/docs/create.html#Create-Record>

## Video Notes

## Chapter 3 Migrations

(golang-migrate package) [https://github.com/golang-migrate/migrate]

options (flags):
  create [-ext E]
        Create a set of timestamped up/down migrations titled NAME 
  up
  down 
  goto v  Migrate schemas to version v

## Unit Test Notes

## 5 (<https://pkg.go.dev/database/sql#example-DB.BeginTx>)

1. When generating SQL code using SQLC, SQLC will generate any int64 it sees as a sql.sqlNullInt64.
2. When running tests, equal amount of transfers, and entries are created.
3. 24 Accounts being made.
4. Makefiles "./..." tests all Unit Tests
5. Unit Test functions must begin with a "Test"

## Closures

The callback function won't know exactly what type it should return.

## Database Transactions Notes Video #6

<https://www.youtube.com/watch?v=gBh__1eFwVI&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE>

DB Transaction file: ./db/sqlc/store.go
DBTX is a result of composition. It holds the DBTX object (implements DBTX interface).

- The New function adheres to ATOM, aka ATOMICITY, and errors out if the steps aren't carried out fully.
If the DB transaction error's out, it will abort the operation (<https://youtu.be/gBh__1eFwVI?t=342>).

## Channels

Since the go routine is inside of (another) go routine or function, there's no guarantee that if
a condition fails that it will stop the entire test. Create a channel that will take the error or results
back to the main routine. Once the go routine that sends the err and results to main exits, iterate
over the error and results (see store_test.go - TransferTx()).

## Goroutines

All code in Goroutines must be finished before the program can exit. The main function exits, The
rest of code stops running.

## Video 7 Deadlock

Prevent Deadlock by using the ROLLBACK feature

- Deadlock will occur when reaching into 2 tables at once
- Values won't update in SQL unless there is a FOR UPDATE when initially describing the tables
- `ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");`

SELECT but with FOR NO KEY UPDATE does ___?

Answer: Allows two update statements at once?

### Deadlock Example Using Terminal and 2 SQL Terminal Sessions
Two select statements being called at the same time means that
you get two of the same values twice. Updating them both at 
the same time will result in a logic bug. This can be demonstrated
using two terminal tabs.

Start Docker psql:
`docker exec -it postgresql12 psql -U root -d secret`

Execute two SELECT statements simultaneously with FOR UPDATE. One of the queries will be frozen while the other is beiing updated:

1st Window:
With FOR UPDATE:
```SQL
BEGIN;
SELECT * FROM accounts where id = 1 FOR UPDATE;
ROLLBACK;
```


- Postgresql debugging statement shiws exactly where the error is occurring in the SQL statement
- Context to pass values
- fmt.Sprintf("d", i + 1)

## Context with Values
- Go routines
  - Context values can be passed into goroutines

Example of how to pass a value into context.
```go
ctx := context.WithValue(context.Background(), txKey, txName)
```

### Main logic file & test file

Keep the logic in the main file and test in another file. Such as `store.go` and `store_test.go`

### gin.Context

Map entire request bodies using functions such as gin.Context.ShouldBindUri,
ShouldBindJSON, etc.

## Error Track

Running tool: /usr/local/go/bin/go test -timeout 30s -run ^TestTransferTx$ github.com/iostate/BankDatabase/db/sqlc

Wow, I am so good at debugging. I ran into an issue and traced the stack trace down to **/go/src/database/sql/sql.go:1288
which presented `db.mu.Lock()`. I've never used this library for go before. I went on Google and queried
"panic: runtime error: invalid memory address or nil pointer dereference sql lock connection"
and found an SO question. I traced the stack trace that the OP posted, and found the same line of code `db.mu.Lock()`.
I knew it had something to do with the lock. Being great at debugging is an invaluable skill to have especially in critical situations
such as an outage or a server going down.
