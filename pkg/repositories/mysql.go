package repositories

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type MysqlAnalyze struct {
	db *sql.DB
}

func (d *MysqlAnalyze) InsertAnalyze(title string, body string) (id string, err error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return
	}
	id = uid.String()

	stmt, err := d.db.Prepare("INSERT INTO analyze VALUES(?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, title, body)
	if err != nil {
		return
	}
	return
}

func (d *MysqlAnalyze) GetAnalyze(id string) (status AnalyzeStatus, result string, err error) {
	panic("implement me")
}

func (d *MysqlAnalyze) AcquireAnalyze() (id string, title string, body string) {
	panic("implement me")
}

func (d *MysqlAnalyze) FinishAnalyze(id string, result string) error {
	panic("implement me")
}

func (d *MysqlAnalyze) name() {
	db := d.db
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)
}
