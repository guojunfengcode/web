package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Mysql *sql.DB

type UserInfo struct {
	username      string
	password      string
	register_time sql.NullString
	now_time      sql.NullString
}

func Query(id int) {
	var user UserInfo
	sqlstr := "select USERNAME, PASSWORD, REGISTER_TIME, LAST_LOGIN_TIME from LOGIN_INFO where id=?;"
	Mysql.QueryRow(sqlstr, id).Scan(&user.username, &user.password, &user.register_time, &user.now_time)
	log.Printf("respon %v\n", user)
}

func QueryMore(n int) {
	sqlstr := "select USERNAME, PASSWORD, REGISTER_TIME, LAST_LOGIN_TIME from LOGIN_INFO where id>?;"
	rows, err := Mysql.Query(sqlstr, n)
	if err != nil {
		log.Printf("exec failed, err %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user UserInfo
		user.register_time = sql.NullString{String: "", Valid: false}
		user.now_time = sql.NullString{String: "", Valid: false}
		err := rows.Scan(&user.username, &user.password, &user.register_time, &user.now_time)
		if err != nil {
			log.Printf("scan failed, err %v\n", err)
			return
		}
		log.Printf("respon %v\n", user)
	}

}

func Update() {
	sqlstr := "update LOGIN_INFO set REGISTER_TIME=str_to_date('%v', '%%Y-%%m-%%d %%H:%%i:%%s') where ID=%d;"
	//sqlstr := "update LOGIN_INFO set REGISTER_TIME=str_to_date('2006-01-02 15:04:05', '%Y-%m-%d %H:%i:%s') where ID=?;"
	var x string
	x = fmt.Sprintf(sqlstr, time.Now().Format("2006-01-02 15:04:05"), 1)

	ret, err := Mysql.Exec(x)
	log.Printf("datetimps:%v\n", time.Now())
	log.Printf("date:%v\n", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("update failed, err: %v\n", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		log.Printf("get id failed, err: %v\n", err)
		return
	}
	log.Printf("update data size: %v\n", n)
}

func UpdateLastTime(username string) {
	sqlstr := "update LOGIN_INFO set LAST_LOGIN_TIME=str_to_date('%v', '%%Y-%%m-%%d %%H:%%i:%%s') where USERNAME='%v';"
	smt := fmt.Sprintf(sqlstr, time.Now().Format("2006-01-02 15:04:05"), username)

	ret, err := Mysql.Exec(smt)
	if err != nil {
		log.Printf("update failed, err: %v\n", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		log.Printf("get id failed, err: %v\n", err)
		return
	}
	log.Printf("update data size: %v\n", n)
}

func IsExistUser(username string) int {
	var user string
	sqlstr := "select * from LOGIN_INFO where USERNAME='?'"
	err := Mysql.QueryRow(sqlstr, username).Scan(&user)
	if err != nil {
		return 0
	}
	return 1
}

func Select(username string, password string) int {

	sqlstr := "select USERNAME, PASSWORD from LOGIN_INFO"
	rows, err := Mysql.Query(sqlstr)
	if err != nil {
		log.Printf("exec failed, err %v\n", err)
		return -1
	}
	defer rows.Close()

	var SelectError int
	for rows.Next() {
		var user UserInfo
		err := rows.Scan(&user.username, &user.password)
		if err != nil {
			log.Printf("scan failed, err %v\n", err)
			return -1
		}
		if user.username == username {
			if user.password == password {
				log.Printf("%v login success\n", username)
				SelectError = 1
				return 1
			} else {
				SelectError = -2
				return SelectError
			}
		} else {
			SelectError = -1
		}
	}
	log.Printf("not %v info, error value:%v\n", username, SelectError)
	return SelectError

}

func Insert(name string, password string, basic string) {
	sqlstr := `insert into LOGIN_INFO(USERNAME, PASSWORD, REGISTER_TIME, LAST_LOGIN_TIME, BASIC) 
		value("%s", "%s", str_to_date('%v', '%%Y-%%m-%%d %%H:%%i:%%s'), str_to_date('%v', '%%Y-%%m-%%d %%H:%%i:%%s'), "%v")`
	smt := fmt.Sprintf(sqlstr, name, password, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), basic)
	ret, err := Mysql.Exec(smt)
	if err != nil {
		log.Printf("insert failed, err: %v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		log.Printf("get id failed, err: %v\n", err)
		return
	}
	log.Printf("insert %v data success\n", n)
}

func MysqlPing(username, password, hostip, port, database string) (err error) {
	str := username + ":" + password + "@tcp(" + hostip + ":" + port + ")/" + database
	//str := "guojunfeng:123456@tcp(192.168.18.128:3306)/login"
	Mysql, err = sql.Open("mysql", str)
	if err != nil {
		log.Printf("mysql login cmd invaild, err %v\n", err)
		return nil
	}

	err = Mysql.Ping()
	if err != nil {
		log.Printf("open mysql fail, err %v\n", err)
		return nil
	}
	log.Println("Successfully connected to the database")
	return
}

func MysqlInit() {
	//MysqlPing()
	Query(1)
	QueryMore(0)
	//Update()
}
