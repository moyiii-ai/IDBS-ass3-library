package main



import (

	"testing"

)



func TestBuildDatabase(t *testing.T) {

	lib := BuildLibrary()
	//If func BuildLibrary can run, tables are created at the same time

	if lib.login_state != 0 {
		t.Errorf("can't build database")
	}
	
}



func TestAddBook(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 
	
	err := lib.AddBook("databasebook", "Xue", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't add book")
	}

	for i := 1; i <= 40; i++ {
		err = lib.AddBook("1", "1", "1")
		if err != nil {
			t.Errorf("can't add book")
		}
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.AddBook("databasebook", "Xue", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't add book")
	}

}



func TestAddUser(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	err := lib.AddUser("DeepC")
	if err != nil {
		t.Errorf("can't add user")
	}

	err = lib.AddUser("Lego")
	if err != nil {
		t.Errorf("can't add user")
	}
	err = lib.AddUser("Lego")
	if err != nil {
		t.Errorf("can't add user")
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0

	err = lib.AddUser("Lego")
	if err != nil {
		t.Errorf("can't add user")
	}

}



func TestDeleteBook(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 
	deleteid := 1
	deletereason := "lost"
	err := lib.DeleteBook(deleteid, deletereason)
	if err != nil {
		t.Errorf("can't delete book")
	}
	err = lib.DeleteBook(deleteid, deletereason)
	if err != nil {
		t.Errorf("can't delete book")
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.DeleteBook(deleteid, deletereason)
	if err != nil {
		t.Errorf("can't delete book")
	}

}



func TestSearchBook(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	err := lib.SearchBook("databasebook", "Xue", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't search book")
	}

	err = lib.SearchBook("databasebook", "no", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't search book")
	}

	err = lib.SearchBook("databasebook", "no", "no")
	if err != nil {
		t.Errorf("can't search book")
	}

	err = lib.SearchBook("no", "no", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't search book")
	}

	err = lib.SearchBook("no", "Xue", "no")
	if err != nil {
		t.Errorf("can't search book")
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.SearchBook("databasebook", "Xue", "no")
	if err != nil {
		t.Errorf("can't search book")
	}
	
	err = lib.SearchBook("no", "Xue", "7-301-04815-7")
	if err != nil {
		t.Errorf("can't search book")
	}

	err = lib.SearchBook("no", "no", "no")
	if err != nil {
		t.Errorf("can't search book")
	}

}



func TestBorrowBook(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 
	borrowid := 2
	err := lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}
	

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	for i := 2; i <= 32; i++ {
		err = lib.BorrowBook(i)
		if err != nil {
			t.Errorf("can't borrow book")
		}
	}

	borrowid = 33
	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//too much borrowed
	
	lib.user.user_id = "DeepC"
	for i := 33; i <= 35; i++ {
		sqlstr, err := lib.db.Prepare("INSERT INTO record(user_id, book_id, borrow_data) values(?, ?, ?)")
		CheckErr(err)

		_, err = sqlstr.Exec(lib.user.user_id, i, "2018-07-21")
		CheckErr(err)

		sqlstr, err = lib.db.Prepare(`UPDATE book
			SET is_available = 2
			WHERE book_id = ?`)
		CheckErr(err)

		_, err = sqlstr.Exec(i)
		CheckErr(err)
	}

	borrowid = 36
	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//too much overdue

	lib.user.user_id = "Lego"
	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//success

	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//has been borrowed

	borrowid = 1
	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//has been deleted

	borrowid = 61
	err = lib.BorrowBook(borrowid)
	if err != nil {
		t.Errorf("can't borrow book")
	}//do not exist

}



func TestAskforDelay(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	recordid := 1
	err := lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//success

	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//success

	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//success

	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//too many times

	recordid = 31
	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//do not belong to you

	lib.user.user_id = "DeepC"
	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//overdue

	err = lib.ReturnBook(recordid)
	err = lib.AskforDelay(recordid)
	if err != nil {
		t.Errorf("can't ask for delay")
	}//has been returned

}



func TestQueryHistory(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	userid := "zzh"
	err := lib.QueryHistory(userid)
	if err != nil {
		t.Errorf("can't query history")
	}//success

	userid = "moyi"
	err = lib.QueryHistory(userid)
	if err != nil {
		t.Errorf("can't query history")
	}//success

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.QueryHistory(userid)
	if err != nil {
		t.Errorf("can't query history")
	}//not belong to you

	userid = "zzh"
	err = lib.QueryHistory(userid)
	if err != nil {
		t.Errorf("can't query history")
	}//success

}



func TestQueryBorrow(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	userid := "zzh"
	err := lib.QueryBorrow(userid)
	if err != nil {
		t.Errorf("can't query borrow")
	}//success

	userid = "moyi"
	err = lib.QueryBorrow(userid)
	if err != nil {
		t.Errorf("can't query borrow")
	}//success

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.QueryBorrow(userid)
	if err != nil {
		t.Errorf("can't query borrow")
	}//do not belong to you

	userid = "zzh"
	err = lib.QueryBorrow(userid)
	if err != nil {
		t.Errorf("can't query borrow")
	}//success

}



func TestQueryDeadline(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	returnid := 1
	err := lib.QueryDeadline(returnid)
	if err != nil {
		t.Errorf("can't query deadline")
	}// success

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.QueryDeadline(returnid)
	if err != nil {
		t.Errorf("can't query deadline")
	}//success

	returnid = 31
	err = lib.QueryDeadline(returnid)
	if err != nil {
		t.Errorf("can't query deadline")
	}//do not belong to you

	lib.user.user_id = "DeepC"
	err = lib.QueryDeadline(returnid)
	if err != nil {
		t.Errorf("can't query deadline")
	}//has been returned

}



func TestCheckOverdue(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	userid := "moyi"
	err := lib.CheckOverdue(userid)
	if err != nil {
		t.Errorf("can't check overdue")
	}//success

	userid = "zzh"
	err = lib.CheckOverdue(userid)
	if err != nil {
		t.Errorf("can't check overdue")
	}//success

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.CheckOverdue(userid)
	if err != nil {
		t.Errorf("can't check overdue")
	}//success

	userid = "moyi"
	err = lib.CheckOverdue(userid)
	if err != nil {
		t.Errorf("can't check overdue")
	}//do not belong to you

}



func TestReturnBook(t *testing.T) {

	lib := BuildLibrary()

	//an administrator account should have been created and login
	//just suppose user_id = "moyi", password = "123"
	lib.login_state = 1
	lib.user.user_id = "moyi"
	lib.user.password = "123"
	lib.user.is_administrator = 1 

	returnid := 1
	err := lib.ReturnBook(returnid)
	if err != nil {
		t.Errorf("can't return book")
	}

	//suppose there is also a ordinary user, user_id = "zzh", password = "233"
	lib.login_state = 1
	lib.user.user_id = "zzh"
	lib.user.password = "233"
	lib.user.is_administrator = 0 

	err = lib.ReturnBook(returnid)
	if err != nil {
		t.Errorf("can't return book")
	}//success

	err = lib.ReturnBook(returnid)
	if err != nil {
		t.Errorf("can't return book")
	}//has been returned

	returnid = 32
	err = lib.ReturnBook(returnid)
	if err != nil {
		t.Errorf("can't return book")
	}//do not belong to you

}