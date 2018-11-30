package db

import (
	"database/sql/driver"
	"time"

	"github.com/jinzhu/gorm"
)

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
	MOVIE_STATE_DOWNLOADING = MovieState("downloading")
	MOVIE_STATE_ERROR       = MovieState("error")
	MOVIE_STATE_FINISHED    = MovieState("finished")
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
	FILE_TYPE_MP4   = MoveFiletype("mp4")
	FILE_TYPE_WEBM  = MoveFiletype("webm")
	FILE_TYPE_OTHER = MoveFiletype("other")
)

type DbMovie struct {
	ID        string
	Name      string
	Link      string
	State     MovieState
	Filetype  MoveFiletype
	Selected  bool
	Error     string
	Progress  float64
	CreatedAt time.Time
}

func (dbm *DbMovie) Create(tx *gorm.DB) error {
	dbm.CreatedAt = time.Now()
	return tx.Create(&dbm).Error
}

func (dbm *DbMovie) Update(tx *gorm.DB) error {
	return tx.Model(&DbMovie{}).Save(dbm).Error
}

func (dbm *DbMovie) Delete(tx *gorm.DB) error {
	return tx.Delete(&dbm).Error
}

func GetMovieByUuid(db *gorm.DB, uuid string) (*DbMovie, error) {
	dm := &DbMovie{}
	err := db.Model(&DbMovie{}).Where("id = ?", uuid).Find(&dm).Error
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func GetMovieBySelected(db *gorm.DB) (*DbMovie, error) {
	dm := &DbMovie{}
	err := db.Model(&DbMovie{}).Where("selected = ?", true).Find(&dm).Error
	if err != nil {
		return nil, err
	}
	return dm, nil
}

func GetMoviesByPage(db *gorm.DB, limit int, fromDateAt *time.Time) ([]DbMovie, error) {
	movies := []DbMovie{}
	var err error
	if fromDateAt != nil {
		err = db.Model(&DbMovie{}).Where("created_at < ? ", *fromDateAt).Order("created_at DESC").Limit(limit).Find(&movies).Error
	} else {
		err = db.Model(&DbMovie{}).Order("created_at DESC").Limit(limit).Find(&movies).Error
	}
	return movies, err
}

func GetMovies(db *gorm.DB) ([]DbMovie, error) {
	movies := []DbMovie{}
	err := db.Model(&DbMovie{}).Find(&movies).Error
	return movies, err
}
