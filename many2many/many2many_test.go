package many2many

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type User struct {
	gorm.Model
	Languages []*Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	gorm.Model
	Name  string
	Users []*User `gorm:"many2many:user_languages;"`
}

func TestQueryAgentWithClues(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	defer sqlDB.Close()
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}))

	mock.ExpectQuery("SELECT * FROM `languages` WHERE `languages`.`deleted_at` IS NULL").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "zh").
			AddRow(2, "en"))

	mock.ExpectQuery("SELECT * FROM `user_languages` WHERE `user_languages`.`language_id` IN (?,?)").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "language_id"}).
			AddRow(1, 1).
			AddRow(1, 2).
			AddRow(2, 1))

	mock.ExpectQuery("SELECT * FROM `users` WHERE `users`.`id` IN (?,?) AND `users`.`deleted_at` IS NULL").
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "user1").
			AddRow(2, "user2"))

	var languages []Language
	err := gormDB.Model(&Language{}).Preload("Users").Find(&languages).Error
	if err != nil {
		t.Fatal(err)
	}
}
