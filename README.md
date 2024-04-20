# Using SQLBoiler

https://github.com/volatiletech/sqlboiler#sqlboiler

Introduction to SQLBoiler

## environment

- Start mysql in a container
- Installing dependncy
- Schema management with sql-migrate

## abstract 

For each code, you need to start a MySQL container in root directory. After that, perform migration in each code directory and run the code.

## internal/getting-started

Example code of [SQLBoiler - Getting Started (v3)](https://youtu.be/y5utRS9axfg)

## FAQ

- When you can't connect to MySQL.  You should not forget MySQL startup needs some time (10 or 30 sec). 
In this case, you should check the status of the mysql container.  (i.n. `docker ps`).  Also, check the health check status of the MySQL container. (i.n. `docker inspect --format "{{json .State.Health }}" database-db-1 | jq`).



