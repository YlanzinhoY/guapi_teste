package repository

import (
	"database/sql"
	"log"

	db "github.com/ylanzinhoy/guapi_teste/sql"

	_ "github.com/lib/pq"
)

type DatabaseRepository struct {
	ConnString string
	dbManager  *db.Queries
}

func NewDatabaseRepository(connString string) *DatabaseRepository {
	return &DatabaseRepository{
		ConnString: connString,
	}
}

func (s *DatabaseRepository) DatabaseConn() *sql.DB {
	connString := s.ConnString
	dbConn, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return dbConn
}

func (s *DatabaseRepository) DbHandler() *db.Queries {

	s.dbManager = db.New(s.DatabaseConn())
	return s.dbManager
}
