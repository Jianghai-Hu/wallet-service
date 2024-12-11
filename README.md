# wallet-service
A take-home assignment from crypto.com.

Can read the <b>Tech Design</b> doc of this demo first:

https://drive.google.com/file/d/1rYJyj5thYjp5blvta3IW9eJhntUN7SIE/view?usp=drive_link
## 1. Before You Run:
### Install PostgreSQL in Local Env(if you haven't done so)
`brew install postgresql`

`brew services start postgresql`

`brew services stop postgresql`
### Create a Database or Using Exist One
Modify the database config in `internal/config/config.go`:
```
var DefaultDBConfig = DBConfig{
    Host:     "localhost",
    Port:     "5432",
    Username: "jianghai",
    Password: "123456",
    DBName:   "demo_db",
}
```
### Create Table in the Database
Can copy the sql in `sql/postgreSQL/create_table.sql`
## 2. Run the Demo:
<b>Check Makefile for Available Cmd</b>

`sudo chmod 755 ./scripts/*`

`make help`

<b>Recommend Usage</b>

`make clean`: cleanup compiled file and logs

`make setup`: download golang dependencies

`make mod/tidy`: go mod tidy

`make build/all`: build all `main.go` under `cmd`

`make run`: run the demo

`make test`: run all the unit test

`make lint`: lint code

<b>Find log files here:</b>

`log/`
## 3. Note:
You can find the curl of provided APIs in <b>Tech Design</b> doc.

Following APIs haven't been implemented in this demo yet:
1. create user wallet
2. get user balance
3. get user transaction history

You may need to manually insert/check the value through `pgAdmin` or `psql` command line.

:bow: :bow: :bow:

And due to tight timeline, this demo may have unhandled edge cases/bugs.

:bow: :bow: :bow: