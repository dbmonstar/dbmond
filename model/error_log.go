package model

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dbmonstar/dbmond/common"
)

// ErrorLog table
type ErrorLog struct {
	ID        int64     `json:"id" xorm:"id bigint pk not null autoincr"`
	Hostname  string    `json:"hostname" xorm:"varchar(50) not null index(01)"`
	Timestamp float64   `json:"timestamp" xorm:"double not null index(01) index(02)"`
	LogedAt   string    `json:"loged_at" xorm:"varchar(50) not null"`
	Level     string    `json:"level" xorm:"varchar(20) not null index(03)"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at" xorm:"timestamp not null created"`
	UpdatedAt time.Time `json:"updated_at" xorm:"timestamp not null updated"`
}

// Exist check exists
func (o *ErrorLog) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *ErrorLog) GetFirst() (ErrorLog, error) {
	var err error

	ret := *o
	boolean, err := orm.Get(&ret)
	if err != nil {
		return ret, err
	}

	if !boolean {
		return ret, fmt.Errorf("no rows")
	}

	return ret, err
}

// GetList get rows
func (o *ErrorLog) GetList(sort ...string) ([]ErrorLog, error) {
	var err error
	var arr []ErrorLog
	var order string

	for i, s := range sort {
		if i > 0 {
			order += ","
		}
		order += s
	}
	err = orm.OrderBy(order).Find(&arr, o)
	common.Log.Info(reflect.TypeOf(o), len(arr), " rows selected")
	return arr, err
}

// Insert new row
func (o *ErrorLog) Insert() error {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	o.rewriteCols()

	if err = o.InsertCheck(); err != nil {
		return err
	}

	if affected, err = session.Insert(o); err != nil {
		return err
	}
	common.Log.Info(reflect.TypeOf(o), affected, "rows inserted!")

	return err
}

// Update update row (partitial column)
func (o *ErrorLog) Update(to *ErrorLog) (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	to.rewriteCols()

	if err = to.UpdateCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Update(to, o); err != nil {
		common.Log.Error(err)
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows updated!")
	return affected, err
}

// Delete delete row
func (o *ErrorLog) Delete() (int64, error) {
	var err error
	var affected int64

	session := orm.NewSession()
	defer session.Close()

	if err = o.DeleteCheck(); err != nil {
		return affected, err
	}

	if affected, err = session.Delete(o); err != nil {
		return affected, err
	}

	common.Log.Info(reflect.TypeOf(o), affected, "rows deleted!")
	return affected, err
}

// InsertCheck validation check
func (o *ErrorLog) InsertCheck() error {
	var err error
	return err
}

// UpdateCheck validation check
func (o *ErrorLog) UpdateCheck() error {
	var err error
	return err
}

// DeleteCheck validation check
func (o *ErrorLog) DeleteCheck() error {
	var err error
	return err
}

// rewriteCols rewrite column value
func (o *ErrorLog) rewriteCols() {
}
