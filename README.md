# BankDatabase

# Explanation of Files

## db/migration 
Contains the migration files: Migrate down = Drop all tables. Migrate up = create tables again. 
# db/query
Contains the SQLC code that will be generated to Golang code. 

# db/sqlc 
Contains the code generated from sqlc. 
The code is generated using the `make sqlc` command, which is a Makefile command. 

# /bank_dbdiagram.sql
A SQL diagram that I built on dbdiagram.io 

# go.mod 
Must initialize a go mod in order for Queries to be recognized. 

# sqlc.yaml 
Configuration file for SQLC. Dictate the location of the files for: queries, the database, and the models.