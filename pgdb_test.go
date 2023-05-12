package healthman

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDbHealthier_Name(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err)

	healthier := NewDBHealthier(db)
	assert.Equal(t, PostgresDependencyName, healthier.Name())
}

func TestDbHealthier_GetHealthy(t *testing.T) {
	t.Run("no error on ping", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.Nil(t, err)
		mock.ExpectPing().WillReturnError(nil)
		var (
			healthier = NewDBHealthier(db)
			healthy   = healthier.CheckHealth()
		)
		assert.Equal(t, true, healthy)
	})
	t.Run("error on ping", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.Nil(t, err)
		mock.ExpectPing().WillReturnError(errors.New("some ping error"))
		var (
			healthier = NewDBHealthier(db)
			healthy   = healthier.CheckHealth()
		)
		assert.Equal(t, false, healthy)
	})
	t.Run("error due to the connection close", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		assert.Nil(t, err)

		mock.ExpectClose().WillReturnError(nil)
		if err = db.Close(); err != nil {
			t.Error("Error while closing the database")
		}

		var (
			healthier = NewDBHealthier(db)
			healthy   = healthier.CheckHealth()
		)
		assert.Equal(t, false, healthy)
	})
}
