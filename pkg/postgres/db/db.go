package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var db *sql.DB

func Connect(port int, host, user, password, dbname string) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalln(err)
	}

	// Handle graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		shutdown()
	}()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
}

func shutdown() {
	fmt.Println("Shutting down...")
	// Close the database connection
	if db != nil {
		if err := db.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}
	os.Exit(0)
}

func InsertRecord(tbName string, columns []string, values []interface{}, outColumns []string, outValues ...interface{}) error {
	// Insert a record into the database
	if len(columns) == 0 || len(values) == 0 {
		return errors.New("columns should be given")
	}
	vals := make([]string, len(columns), len(columns))
	for i := 0; i < len(columns); i++ {
		vals[i] = fmt.Sprintf("$%d", i+1)
	}
	sqlStatement := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		tbName, strings.Join(columns, ","), strings.Join(vals, ","))
	var err error
	if len(outColumns) == 0 {
		_, err = db.Exec(sqlStatement, values...)
	} else {
		sqlStatement += " RETURNING " + strings.Join(outColumns, ", ")
		err = db.QueryRow(sqlStatement, values...).Scan(outValues...)
	}
	if err != nil {
		log.Println("Error inserting record:", err)
		return err
	}
	fmt.Println("record inserted successfully")
	return nil
}

type QueryConstruct struct {
	Columns         string
	WhereExpr       string
	WhereExprValues []interface{}
	SetExpr         string
	SetExprValues   []interface{}
	Offset          int
	Limit           int
}

func GetRecord(tbName string, q QueryConstruct) (*sql.Rows, error) {
	if q.Columns == "" {
		q.Columns = "*"
	}
	sqlStatement := fmt.Sprintf("SELECT %s FROM %s", q.Columns, tbName)
	if q.WhereExpr != "" {
		sqlStatement += fmt.Sprintf(" WHERE %s", q.WhereExpr)
	}
	if q.Offset > 0 {
		sqlStatement += fmt.Sprintf(" OFFSET %d", q.Offset)
	}
	if q.Limit > 0 {
		sqlStatement += fmt.Sprintf(" LIMIT %d", q.Limit)
	}
	res, err := db.Query(sqlStatement, q.WhereExprValues...)
	if err != nil {
		log.Println("Error getting record:", err)
		return nil, err
	}
	return res, nil
}

func UpdateRecord(tbName string, q QueryConstruct) error {
	sqlStatement := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tbName, q.SetExpr, q.WhereExpr)
	q.SetExprValues = append(q.SetExprValues, q.WhereExprValues...)
	_, err := db.Exec(sqlStatement, q.SetExprValues...)
	if err != nil {
		log.Println("Error updating record:", err)
		return err
	}
	fmt.Println("record updated successfully")
	return nil
}

func TransactUpdate(tbName string, qs ...QueryConstruct) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, q := range qs {
		sqlStatement := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			tbName, q.SetExpr, q.WhereExpr)
		q.SetExprValues = append(q.SetExprValues, q.WhereExprValues...)
		_, err = tx.Exec(sqlStatement, q.SetExprValues...)
		if err != nil {
			log.Println("Error updating record:", err)
			tx.Rollback()
			return err
		}
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
