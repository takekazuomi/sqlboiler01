#!/bin/sh

# for sqlboiler binary
go install github.com/aarondl/sqlboiler/v4@v4.19.5
go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-mysql@v4.19.5

# for sql-migrate binary
go install github.com/rubenv/sql-migrate/sql-migrate@v1.8.0

# for staticcheck binary
go install honnef.co/go/tools/cmd/staticcheck@latest

# other dependencies
go get github.com/aarondl/null/v8@latest
