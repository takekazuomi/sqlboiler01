#!/bin/sh

# for sqlboiler binary
go install github.com/volatiletech/sqlboiler/v4@v4.16.2
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.16.2

# for sql-migrate binary
go install github.com/rubenv/sql-migrate/sql-migrate@v1.7.0

# for staticcheck binary
go install honnef.co/go/tools/cmd/staticcheck@latest

# other dependencies
go get github.com/volatiletech/null/v8@latest
