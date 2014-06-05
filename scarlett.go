/**
 * Author:        Tony.Shao
 * Email:         xiocode@gmail.com
 * Github:        github.com/xiocode
 * File:          scarlett.go
 * Description:   scarlett
 */

package scarlett

import (
	"database/sql"
	log "github.com/xiocode/glog"
)

type Scarlett struct {
	Conn *sql.DB
}

func NewScarlett(driver, source string) *Scarlett {
	conn, err := sql.Open(driver, source)
	if err != nil {
		log.Fatalln(err)
	}
	return &Scarlett{
		Conn: conn,
	}
}

func (s *Scarlett) Query(dst interface{}, binder Binder, query string, params ...interface{}) error {
	rows, err := s.Conn.Query(query, params...)
	if err != nil {
		return err
	}

	err = s.Scan(dst, binder, rows)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scarlett) Exec(query string, params ...interface{}) (sql.Result, error) {
	return s.Conn.Exec(query, params...)
}
