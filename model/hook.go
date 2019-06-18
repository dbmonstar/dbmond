package model

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dbmonstar/dbmond/common"
)

// Hook table
type Hook struct {
	ID          int       `json:"id" xorm:"id int(11) pk not null autoincr"`
	Name        string    `form:"name" json:"name" xorm:"varchar(32) not null index(01)"`
	Hash        string    `form:"hash" json:"hash" xorm:"varchar(32) not null index(01)"`
	Instance    string    `form:"instance" json:"instance" xorm:"instance varchar(32) not null index(01) index(02)"`
	Level       string    `form:"level" json:"level" xorm:"varchar(10) not null"`
	Subject     string    `form:"subject" json:"subject" xorm:"varchar(64) not null"`
	Status      string    `form:"status" json:"status" xorm:"varchar(10) not null "`
	Description string    `form:"description" json:"description" xorm:"text not null "`
	StartsAt    time.Time `json:"starts_at" xorm:"timestamp not null index(02)"`
	EndsAt      time.Time `json:"ends_at" xorm:"timestamp null"`
	CreatedAt   time.Time `json:"created_at" xorm:"timestamp not null created"`
	UpdatedAt   time.Time `json:"updated_at" xorm:"timestamp not null updated"`
}

// Exist check exists
func (o *Hook) Exist() bool {
	boolean, _ := orm.Exist(o)
	return boolean
}

// GetFirst get first one
func (o *Hook) GetFirst() (Hook, error) {
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
func (o *Hook) GetList(sort ...string) ([]Hook, error) {
	var err error
	var arr []Hook
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
func (o *Hook) Insert() error {
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
func (o *Hook) Update(to *Hook) (int64, error) {
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
func (o *Hook) Delete() (int64, error) {
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
func (o *Hook) InsertCheck() error {
	var err error
	return err
}

// UpdateCheck validation check
func (o *Hook) UpdateCheck() error {
	var err error
	return err
}

// DeleteCheck validation check
func (o *Hook) DeleteCheck() error {
	var err error
	return err
}

// rewriteCols rewrite column value
func (o *Hook) rewriteCols() {
}
