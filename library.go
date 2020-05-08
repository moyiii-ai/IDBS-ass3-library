package main



import (

	"fmt"
	"time"
	// mysql connector

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"

)



const (
	User     = "user"//please fill your username there
	Password = "password"//please fill your password there
	DBName   = "ass3"
	timeLayout = "2006-01-02" 
)

var sum = [13]int {
	0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365, 
}

var monthday = [13]int {
	0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
} 

var emonthday = [13]int {
	0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31,
} 

type account struct {
	user_id string
	password string
	is_administrator int
}

type Library struct {
	db *sqlx.DB
	user account
	login_state int
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}



func (lib *Library) ConnectDB() {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	CheckErr(err)
	lib.db = db
}



// CreateTables created the tables in MySQL

func (lib *Library) CreateTables() error {
	
	sqlstr := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS user(
		user_id VARCHAR(32) NOT NULL PRIMARY KEY,
		password VARCHAR(32) DEFAULT "233",
		is_adminisrator INT DEFAULT false)`)

	_, err := lib.db.Exec(sqlstr)
	CheckErr(err)
	
	sqlstr = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS book(
		book_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(64),
		author VARCHAR(32),
		isbn VARCHAR(128),
		is_available INT DEFAULT 1,
		explanation VARCHAR(256) DEFAULT "")`)

	_, err = lib.db.Exec(sqlstr)
	CheckErr(err)

	sqlstr = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS record(
		record_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		user_id VARCHAR(32),
		book_id INT,
		borrow_data char(16),
		return_data char(16) DEFAULT "",
		deadline INT DEFAULT 90,
		delay_time INT DEFAULT 0,
		FOREIGN KEY(user_id) REFERENCES user(user_id),
		FOREIGN KEY(book_id) REFERENCES book(book_id))`)

	_, err = lib.db.Exec(sqlstr)
	CheckErr(err)

	return nil

}



func BuildLibrary() *Library {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
	CheckErr(err)

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", DBName))
	CheckErr(err)
	//drop for test

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", DBName))
	CheckErr(err)

	lib := &Library{
		db:  nil,
		user: account{
			user_id: "",
			password: "",
			is_administrator: 0,
		},
		login_state: 0,
	}

	lib.ConnectDB()
	lib.CreateTables()
	return lib
}



func (lib *Library) CheckBook(book_id int) int {
	
	sqlstr, err := lib.db.Prepare(`SELECT is_available FROM book WHERE book_id = ?`)
	CheckErr(err)

	defer sqlstr.Close()

	rows, err := sqlstr.Query(book_id)
	CheckErr(err)
	defer rows.Close()

	var ret int

	for rows.Next() {
		err := rows.Scan(&ret)
		CheckErr(err)
		return ret
	}

	return 4

}

/*	return 1 : the book is available
	return 2 : the book is not available because it is borrowed
	return 3 : the book is not available because is has been deleted
	return 4 : the book do not exist*/


func GetDay(data string) int {
	var ans, _year, _month, _day int
	ans = 0
	for i := 0; i < 4; i++ {
		_year = _year * 10 + int(data[i] - '0')
	}
	_month = int(data[5] - '0') * 10 + int(data[6] - '0')
	_day = int(data[8] - '0') * 10 + int(data[9] - '0')
	ans = sum[_month] + _day
	if (_year % 100 != 0 && _year % 4 == 0) || (_year % 400 == 0) {
		ans++
	}

	ans += _year * 365 +  _year / 4 - _year / 100 + _year / 400 + 1	

	return ans

}



func (lib *Library) CheckUser() int {
	
	sqlstr, err := lib.db.Prepare(`SELECT borrow_data, deadline
		FROM record
		WHERE user_id = ? AND return_data = ""`)
	CheckErr(err)

	defer sqlstr.Close()

	rows, err := sqlstr.Query(lib.user.user_id)
	CheckErr(err)

	defer rows.Close()

	var BorrowNum, OverdueNum, deadline int = 0, 0, 0
	borrow_data := "" 

	for rows.Next() {
		err := rows.Scan(&borrow_data, &deadline)
		CheckErr(err)
		
		BorrowNum++
		
		if GetDay(time.Now().Format(timeLayout)) - GetDay(borrow_data) > deadline {
			OverdueNum++
		}
	}

	if OverdueNum >= 3 {
		return 2
	}
	if BorrowNum >= 30 {
		return 3
	}
	return 1
}

/*	return 1 : the user can borrow books
	return 2 : the user cannot borrow books because overdue
	return 3 : the user cannot borrow books because he has borrowed too much*/



func (lib *Library) AddBook(title, author, ISBN string) error{
	
	if lib.user.is_administrator != 1 {
		fmt.Println("Error 2: You cannot AddBook because you are not administrator!")
		return nil
	}

	sqlstr, err := lib.db.Prepare(`INSERT INTO book(title, author, isbn, is_available) values(?, ?, ?, 1)`)
	CheckErr(err)

	_, err = sqlstr.Exec(title, author, ISBN)
	CheckErr(err)

	fmt.Printf("Book %s has been successfully added!\n", title)

	return nil

}



func (lib *Library) AddUser(user_id string) error{

	if lib.user.is_administrator != 1 {
		fmt.Println("Error 2: You cannot AddUser because you are not administrator!")
		return nil
	}

	rows, err := lib.db.Query(fmt.Sprintf("select * from user where user_id = '%s'", user_id))
	CheckErr(err)
	defer rows.Close()

	num := 0
	for rows.Next(){
		num++
	}

	if num > 0 {
		fmt.Printf("Error 5: User already exists!\n")
		return nil
	}

	sqlstr, err := lib.db.Prepare(`INSERT INTO user(user_id) values(?)`)
	CheckErr(err)

	_, err = sqlstr.Exec(user_id)
	CheckErr(err)

	fmt.Printf("User %s has been successfully created!\n", user_id)

	return nil

}



func (lib *Library) DeleteBook(book_id int, explanation string) error{

	if lib.user.is_administrator != 1 {
		fmt.Println("Error 2: You cannot DeleteBook because you are not administrator!")
		return nil
	}

	book_type := lib.CheckBook(book_id)
	if book_type == 3 {
		fmt.Println("Error 6: You cannot DeleteBook because it has been deleted!")
		return nil
	}

	sqlstr, err := lib.db.Prepare(`UPDATE book
		SET is_available = 3, explanation = ?
		WHERE book_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(explanation, book_id)
	CheckErr(err)

	sqlstr, err = lib.db.Prepare(`UPDATE record
		SET return_data = ?
		WHERE book_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(time.Now().Format(timeLayout), book_id)
	CheckErr(err)

	fmt.Printf("You have successfully deleted the book %d.\n", book_id)

	return nil

}


func (lib *Library) BorrowBook(book_id int) error{
	
	if lib.user.is_administrator == 1 {
		fmt.Println("Error 2: You cannot borrow book because you are administrator")
		return nil
	}

	user_type := lib.CheckUser()
	if user_type == 2 {
		fmt.Println("Error 8: You cannot borrow book because too much book overdue")
		return nil
	}
	if user_type == 3 {
		fmt.Println("Error 9: You cannot borrow book because too much book borrowed")
		return nil
	}

	book_type := lib.CheckBook(book_id) 
	if book_type == 2 {
		fmt.Println("Error 10: You cannot borrow this book because it has been borrowed")
		return nil
	} 
	if book_type == 3 {
		fmt.Println("Error 6: You cannot borrow this book because it has been deleted")
		return nil
	}
	if book_type == 4 {
		fmt.Println("Error 11: You cannot borrow this book because it do not exist")
		return nil
	}

	sqlstr, err := lib.db.Prepare("INSERT INTO record(user_id, book_id, borrow_data) values(?, ?, ?)")
	CheckErr(err)

	_, err = sqlstr.Exec(lib.user.user_id, book_id, time.Now().Format(timeLayout))
	CheckErr(err)

	sqlstr, err = lib.db.Prepare(`UPDATE book
		SET is_available = 2
		WHERE book_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(book_id)
	CheckErr(err)

	fmt.Printf("You have successfully borrowed %d.\n", book_id)

	return nil
}



func (lib *Library) ReturnBook(record_id int) error{
	if lib.user.is_administrator == 1 {
		fmt.Println("Error 2: You cannot return book because you are administrator")
		return nil
	}

	book_id := -1
	return_data := ""
	var user_id string

	rows, err := lib.db.Query(fmt.Sprintf("select user_id, book_id, return_data from record where record_id = %d", record_id))
	CheckErr(err)
	defer rows.Close()

	for rows.Next(){
		rows.Scan(&user_id, &book_id, &return_data)
	}

	if book_id == -1 {
		fmt.Println("Error 13: You cannot return this book because the record do not exist.")
		return nil
	}
	if user_id != lib.user.user_id {
		fmt.Println("Error 2: You cannot return this book because the record do not belong to you.")
		return nil
	}
	if return_data != "" {
		fmt.Println("Error 12: You cannot return this book because it has been deleted or returned.")
		return nil
	}


	sqlstr, err := lib.db.Prepare(`UPDATE record
		SET return_data = ?
		WHERE record_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(time.Now().Format(timeLayout), record_id)
	CheckErr(err)

	sqlstr, err = lib.db.Prepare(`UPDATE book
		SET is_available = 1
		WHERE book_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(book_id)
	CheckErr(err)

	fmt.Printf("You have successfully returned the book %d.\n", book_id)

	return nil
}



func (lib *Library) AskforDelay(record_id int) error{
	if lib.user.is_administrator == 1 {
		fmt.Println("Error 2: You cannot ask for delay because you are administrator")
		return nil
	}

	var book_id, delay_time  int = -1, 0
	var user_id, borrow_data, return_data string

	rows, err := lib.db.Query(fmt.Sprintf(`select user_id, book_id, borrow_data, return_data, delay_time 
		from record 
		where record_id = %d`, record_id))
	CheckErr(err)
	defer rows.Close()

	for rows.Next(){
		rows.Scan(&user_id, &book_id, &borrow_data, &return_data, &delay_time)
	}

	if user_id != lib.user.user_id {
		fmt.Println("Error 2: You cannot ask for delay because the record do not belong to you.")
		return nil
	}
	if return_data != "" {
		fmt.Println("Error 12: You cannot ask for delay because the book has been deleted or returned.")
		return nil
	}
	if delay_time >= 3 {
		fmt.Println("Error 14: You cannot ask for delay because you have asked for too many times")
		return nil
	}
	
	var expect_time int = delay_time * 30 + 90
	if GetDay(time.Now().Format(timeLayout)) - GetDay(borrow_data) > expect_time{
		fmt.Println("Error 15: You cannot ask for delay because deadline has been missed")
		return nil
	}

	sqlstr, err := lib.db.Prepare(`UPDATE record
		SET delay_time = delay_time + 1, deadline = deadline + 30
		WHERE record_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(record_id)
	CheckErr(err)

	fmt.Printf("You have successfully asked for delay record %d.\n", record_id)

	return nil
}



func (lib *Library) ChangePassword() error{
	var str, password string
	
	for ;; {
		fmt.Print("Please enter your new password, or input 'quit' to leave.\n> ")
		fmt.Scan(&str)

		if str == "quit" {
			fmt.Println("Change cancelled.")
			return nil
		}

		fmt.Print("Please repeat you password\n> ")
		fmt.Scan(&password)

		if str != password {
			fmt.Println("Error 7: password doesn't match!")
			continue
		}

		break
	}
	
	sqlstr, err := lib.db.Prepare(`UPDATE user 
		SET password = ?
		WHERE user_id = ?`)
	CheckErr(err)

	_, err = sqlstr.Exec(password, lib.user.user_id)
	CheckErr(err)
	
	fmt.Println("Change the possword successfully")

	lib.user.password = password

	return nil
}



func (lib *Library) QueryHistory(user_id string) error{
	if lib.user.is_administrator != 1 && user_id != lib.user.user_id {
		fmt.Println("Error 2: You can only query for your own history.")
		return nil
	}

	rows, err := lib.db.Query(fmt.Sprintf("select * from record where user_id = '%s'", user_id))
	CheckErr(err)
	defer rows.Close()
	
	var record_id, book_id, deadline, delay_time int
	var borrow_data, return_data string
	for rows.Next() {
		rows.Scan(&record_id, &user_id, &book_id, &borrow_data, &return_data, &deadline, &delay_time)
		fmt.Printf("record_id = %d, user_id = %s, book_id = %d, borrow_data = %s, return_data = %s, deadline = %d, delay_time = %d\n",
					record_id, user_id, book_id, borrow_data, return_data, deadline, delay_time)
	}

	return nil

}



func (lib *Library) QueryBorrow(user_id string) error{
	if lib.user.is_administrator != 1 && user_id != lib.user.user_id {
		fmt.Println("Error 2: You can only query for your own information.")
		return nil
	}

	rows, err := lib.db.Query(fmt.Sprintf(`select record_id, book_id, borrow_data
		from record 
		where user_id = '%s' and return_data = ""`, user_id))
	CheckErr(err)
	defer rows.Close()
	
	var record_id, book_id int
	var borrow_data string
	for rows.Next() {
		rows.Scan(&record_id, &book_id, &borrow_data)
		fmt.Printf("record_id = %d, book_id = %d, borrow_data = %s\n",
					record_id, book_id, borrow_data)
	}

	return nil

}



func (lib *Library) QueryDeadline(record_id int) error{

	var book_id, delay_time int = -1, 0
	var user_id, borrow_data, return_data string

	rows, err := lib.db.Query(fmt.Sprintf(`select user_id, book_id, borrow_data, return_data, delay_time 
		from record 
		where record_id = %d`, record_id))
	CheckErr(err)
	defer rows.Close()

	for rows.Next(){
		rows.Scan(&user_id, &book_id, &borrow_data, &return_data, &delay_time)
	}

	if lib.user.is_administrator != 1 && user_id != lib.user.user_id {
		fmt.Println("Error 2: You cannot query deadline because the record do not belong to you.")
		return nil
	}
	if return_data != "" {
		fmt.Println("Error 12: You cannot query deadline because the book has been deleted or returned.")
		return nil
	}
	
	var _year, _month, _day int
	for i := 0; i < 4; i++ {
		_year = _year * 10 + int(borrow_data[i] - '0')
	}
	_month = int(borrow_data[5] - '0') * 10 + int(borrow_data[6] - '0')
	_day = int(borrow_data[8] - '0') * 10 + int(borrow_data[9] - '0')

	_day += 90 + 30 * delay_time
	if (_year % 100 != 0 && _year % 4 == 0) || (_year % 400 == 0) {
		for ;_day > emonthday[_month]; {
			_day -= emonthday[_month]
			_month++
			if _month == 13 {
				_year++
				for ;_day > monthday[_month]; {
					_day -= monthday[_month]
					_month++
				}
			}
		}
	}else {
		for ;_day > monthday[_month]; {
			_day -= monthday[_month]
			_month++
			if _month == 13 {
				_year++
				if (_year % 100 != 0 && _year % 4 == 0) || (_year % 400 == 0) {
					for ;_day > emonthday[_month]; {
						_day -= emonthday[_month]
						_month++
					}
				}else {
					for ;_day > monthday[_month]; {
						_day -= monthday[_month]
						_month++
					}
				}
			}
		}
	}
	
	fmt.Printf("Deadline : %d-%d-%d\n", _year, _month, _day)

	return nil
}



func (lib *Library) CheckOverdue(user_id string) error{
	if lib.user.is_administrator != 1 && user_id != lib.user.user_id {
		fmt.Println("Error 2: You can only check your own information.")
		return nil
	}
	
	sqlstr, err := lib.db.Prepare(`SELECT record_id, book_id, borrow_data, deadline
		FROM record
		WHERE user_id = ? AND return_data = ""`)
	CheckErr(err)

	defer sqlstr.Close()

	rows, err := sqlstr.Query(user_id)
	CheckErr(err)

	defer rows.Close()

	var num, record_id, book_id, deadline int = 0, 0, 0, 0
	borrow_data := "" 

	for rows.Next() {
		err := rows.Scan(&record_id, &book_id, &borrow_data, &deadline)
		CheckErr(err)
		
		if GetDay(time.Now().Format(timeLayout)) - GetDay(borrow_data) > deadline {
			num++
			fmt.Printf("Overdue : record_id = %d, book_id = %d, deadline = %d, borrow_data = %s.\n",
						record_id, book_id, deadline, borrow_data)
		}
	}

	if num == 0 {
		fmt.Println("No book overdue.")
	}

	return nil

}



func (lib *Library) SearchBook(title, author, isbn string) error {
	book_id := 0 
	var is_available int
	var explanation string
	if title != "no" && author != "no" && isbn != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where title = '%s' and author = '%s' and isbn = '%s'`, title, author, isbn))
		CheckErr(err)
		defer rows.Close()
		
		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if title != "no" && author != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where title = '%s' and author = '%s'`, title, author))
		CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if title != "no" && isbn != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where title = '%s' and isbn = '%s'`, title, isbn))
		CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if author != "no" && isbn != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where author = '%s' and isbn = '%s'`, author, isbn))
		CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if title != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where title = '%s'`, title))
		CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if author != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where author = '%s'`, author))
		CheckErr(err)
		defer rows.Close()
		
		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}
	
	if isbn != "no" {
		rows, err := lib.db.Query(fmt.Sprintf(`select *
			from book
			where isbn = '%s'`, isbn))
		CheckErr(err)
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book_id, &title, &author, &isbn, &is_available, &explanation)
			CheckErr(err)
			fmt.Printf("book_id = %d, title = %s, author = %s, isbn = %s, is_available = %d, explanation = %s\n", 
						book_id, title, author, isbn, is_available, explanation)
		}
		return nil
	}

	return nil

}



func (lib *Library) Logout(){
	fmt.Printf("Bye! %s!\n", lib.user.user_id)
	lib.user = account{
		user_id: "",
		password: "",
		is_administrator: 0,
	}
	lib.login_state = 0
}



func (lib *Library) Solve() {

	var opt string
	for ;; {
		if lib.user.is_administrator == 1 {
			fmt.Print("administrator@library$ ")
		}else{
			fmt.Print("user@library# ")
		}

		fmt.Scan(&opt)
		if opt == "addbook" {
			var title, author, isbn string
			fmt.Println("Please enter title, author, isbn of the new book, spaced by blank")
			fmt.Scan(&title, &author, &isbn)
			lib.AddBook(title, author, isbn)
			continue
		}

		if opt == "adduser" {
			var user_id string
			fmt.Println("Please enter user_id for the new account")
			fmt.Scan(&user_id)
			lib.AddUser(user_id)
			continue
		}
		
		if opt == "deletebook" {
			var book_id int
			var explanation string
			fmt.Println("Please enter book_id and explanation for the deleted book")
			fmt.Scan(&book_id, &explanation)
			lib.DeleteBook(book_id, explanation)
			continue
		}
		
		if opt == "borrowbook" {
			var book_id int
			fmt.Println("Please enter book_id you want to borrow")
			fmt.Scan(&book_id)
			lib.BorrowBook(book_id)
			continue
		}
		
		if opt == "returnbook" {
			var record_id int
			fmt.Println("Please enter record_id you want to return")
			fmt.Scan(&record_id)
			lib.ReturnBook(record_id)
			continue
		}
		
		if opt == "askfordelay" {
			var record_id int
			fmt.Println("Please enter record_id you want to delay")
			fmt.Scan(&record_id)
			lib.AskforDelay(record_id)
			continue
		}
		
		if opt == "searchbook" {
			var title, author, isbn string
			fmt.Println("Please enter title, author, isbn to search, spaced by blank")
			fmt.Println("Please input no if you do not know details, for example : db no no")
			fmt.Scan(&title, &author, &isbn)
			lib.SearchBook(title, author, isbn)
			continue
		}
		
		if opt == "queryhistory" {
			var user_id string
			fmt.Println("Please enter user_id you want to query")
			fmt.Scan(&user_id)
			lib.QueryHistory(user_id)
			continue
		}

		if opt == "queryborrow" {
			var user_id string
			fmt.Println("Please enter user_id you want to query")
			fmt.Scan(&user_id)
			lib.QueryBorrow(user_id)
			continue
		}
		
		if opt == "querydeadline" {
			var record_id int
			fmt.Println("Please enter record_id you want to query")
			fmt.Scan(&record_id)
			lib.QueryDeadline(record_id)
			continue
		}
		
		if opt == "checkoverdue" {
			var user_id string
			fmt.Println("Please enter user_id you want to check")
			fmt.Scan(&user_id)
			lib.CheckOverdue(user_id)
			continue
		}

		if opt == "changepassword" {
			lib.ChangePassword()
			continue
		}

		if opt == "logout" {
			lib.Logout()
			return
		}
		
		//if len(opt) != 0 {
			fmt.Println("Error 1: Invalid Input")
		//}
	}
}



func (lib *Library) Enroll(){
	var password, repeat_password, str string

	for ;; {
		fmt.Print("Please enter the Password of database to prove you have the authority, or input 'quit' to leave.\n> ")
		fmt.Scan(&password)
		if password == "quit" {
			fmt.Println("Enroll cancelled.")
			return
		}

		if password == Password {
			fmt.Println("Welcome to enroll.")
			break
		}

		fmt.Println("Error 4: Wrong Password!")
	}

	for ;; {
		fmt.Print("Please input your user_id, or input 'quit' to leave.\n> ")
		fmt.Scan(&str)
		if str == "quit" {
			fmt.Println("Enroll cancelled.")
			return
		}

		rows, err := lib.db.Query(fmt.Sprintf("select * from user where user_id = '%s'", str))
		CheckErr(err)
		defer rows.Close()

		num := 0
		for rows.Next(){
			num++
		}

		if num > 0 {
			fmt.Printf("Error 5: User already exists!\n")
			continue
		}

		fmt.Print("Please input your password\n> ")
		fmt.Scan(&password)

		fmt.Print("Please repeat you password\n> ")
		fmt.Scan(&repeat_password)

		if password != repeat_password {
			fmt.Println("Error 7: password doesn't match!")
			continue
		}

		break

	}

	sqlstr, err := lib.db.Prepare("insert into user values(?, ?, 1)")
	CheckErr(err)

	_, err = sqlstr.Exec(str, password)
	CheckErr(err)	
	
	fmt.Printf("Administrator account '%s' has been created successfully!\n", str)

}



func (lib *Library) Login(){
	var correct_password, password, user_id, str string
	var is_administrator int

	for ;; {
		fmt.Printf("Please enter your user_id, or input 'quit' to leave.\n> ")

		fmt.Scan(&str)
		if str == "quit" {
			fmt.Println("Login cancelled.")
			return
		}

		rows, err := lib.db.Query(fmt.Sprintf("select * from user where user_id = '%s'", str))
		CheckErr(err)
		defer rows.Close()
		
		num := 0
		for rows.Next() {
			num++
			rows.Scan(&user_id, &correct_password, &is_administrator)
		}

		if num == 0 {
			fmt.Println("Error 3: No such user!")
			continue
		}

		fmt.Print("Please enter your password.\n> ")
		fmt.Scan(&password)
		if correct_password != password {
			fmt.Println("Error 4: Wrong password!")
		}else {
			fmt.Println("Login successfully!")
			break
		}
	}

	lib.user = account{
		user_id: user_id,
		password: password,
		is_administrator: is_administrator,
	}

	lib.login_state = 1
}



func main() {
	lib := BuildLibrary()

	fmt.Println("Welcome to the Library Management System!")
	fmt.Printf("Please Input 'login' to sign in, 'enroll' to sign up an admin account or 'exit' to exit.\n>")

	var str string

	for ;; {
		fmt.Scan(&str)

		if str == "exit"{
			fmt.Println("Hope to see you again!")
			break
		}
		
		if str == "login"{
			lib.Login()
			
			if lib.login_state == 1 {
				lib.Solve()
			}

			fmt.Print("> ")
			continue
		}

		if str == "enroll"{
			lib.Enroll()
			fmt.Print("> ")
			continue
		}
			
		fmt.Println("Error 1: Invalid input.")
		fmt.Println("Please Input 'login' to sign in, 'enroll' to sign up an admin account or 'exit' to exit.")

		fmt.Print("> ")
	}
}