package admin

import (
	"database/sql/driver"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&DbMovie{})
	return db, nil
}

type MovieState string

func (m MovieState) Value() (driver.Value, error) {
	return string(m), nil
}

func (m *MovieState) Scan(value interface{}) error {
	str, ok := value.([]uint8)
	if !ok {
		*m = ""
		return nil
	}
	*m = MovieState(str)
	return nil
}

const (
	MovieDownloading = MovieState("downloading")
	MovieError       = MovieState("error")
	MovieFinished    = MovieState("finished")
)

type MoveFiletype string

func (n MoveFiletype) Value() (driver.Value, error) {
	return string(n), nil
}

func (n *MoveFiletype) Scan(value interface{}) error {
	str, ok := value.([]uint8)
	if !ok {
		*n = ""
		return nil
	}
	*n = MoveFiletype(str)
	return nil
}

const (
	FiletypeMp4   = MoveFiletype("mp4")
	FiletypeWebm  = MoveFiletype("webm")
	FiletypeOther = MoveFiletype("other")
)

type DbMovie struct {
	Id       string
	Link     string
	State    MovieState
	Filetype MoveFiletype
	Error    string
	Progress float64
}

func (dbm *DbMovie) Create(tx *gorm.DB) error {
	return tx.Create(&dbm).Error
}

func (dbm *DbMovie) Update(tx *gorm.DB) error {
	return tx.Save(&dbm).Error
}

func (dbm *DbMovie) Delete(tx *gorm.DB) error {
	return tx.Delete(&dbm).Error
}
