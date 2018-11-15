package db

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetMoviesPagination(t *testing.T) {
	dbtx, err := NewDB("/tmp/test.db")
	assert.Nil(t, err)
	ms := []DbMovie{}
	for i := 0; i < 10; i++ {
		m1 := DbMovie{}
		m1.Create(dbtx)
		time.Sleep(1 * time.Second)
		ms = append(ms, m1)
	}
	nms, _ := GetMoviesByPage(dbtx, 2, &ms[1].CreatedAt)
	log.Println(nms)
	nms, _ = GetMoviesByPage(dbtx, 2, &ms[9].CreatedAt)
	log.Println(nms)

	nms, _ = GetMoviesByPage(dbtx, 2, nil)
	log.Println(nms)
	nms, _ = GetMoviesByPage(dbtx, 2, &nms[1].CreatedAt)
	log.Println(nms)
}
