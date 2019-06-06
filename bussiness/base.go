package bussiness

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/o7z/try_gorm/model"
	uuid "github.com/satori/go.uuid"
)

var db *gorm.DB

func Init(_db *gorm.DB) error {
	db = _db
	db.LogMode(true)
	db.SingularTable(true) //不要自作主张了
	// db.CreateTable(&model.User{})
	// db.DropTable(&model.User{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Article{})
	db.Model(&model.Article{}).AddForeignKey("create_user", "user(id)", "RESTRICT", "RESTRICT")
	// db.Model(&model.User{}).DropColumn("gender")
	return nil
}

func NewID() string {
	if newUUID, err := uuid.NewV4(); err != nil {
		fmt.Printf("uuid.NewV4()->%s\n\r", err)
		panic(err)
	} else {
		return newUUID.String()
	}
}

func Show(v interface{}) error {
	if vJsonBuf, err := json.MarshalIndent(v, "", "  "); err != nil {
		return fmt.Errorf("json.MarshalIndent(v, ``, `  `)->%s\n\r", err)
	} else {
		fmt.Println(string(vJsonBuf))
		return nil
	}
}

func Test() error {
	//TestCreateUsers
	//TestShowUsers()

	// db.Model(&user).Update("name", "王大锤")
	// db.Find(&users, "gender=?", 1)
	// db.Model(&users).Update("name", "懵逼")

	user := model.User{}
	db.First(&user)
	Show(user)
	// db.Create(model.Article{
	// 	ID:         NewID(),
	// 	Title:      "!",
	// 	Content:    "!",
	// 	CreateUser: user.ID,
	// 	CreateTime: time.Now()})

	// article := model.Article{}
	db.Model(&user).Related(&user.Articles, "Articles")

	Show(user)
	// Show(article)

	return nil
}

func TestCreateUsers() error {
	users := []model.User{}
	for i := 0; i < 100; i++ {
		users = append(users, model.User{
			ID:         NewID(),
			Name:       "test" + strconv.Itoa(i),
			Gender:     uint8(i%2) + 1,
			CreateTime: time.Now()})
	}
	// db.Create(&users) //批量插入目前版本还不行
	for _, user := range users {
		db.Create(&user)
	}
	return nil
}

func TestShowUsers() error {
	users := []model.User{}
	db.Find(&users)
	// db.Where("gender=?", 1).Find(&users)
	// db.Find(&users, "gender=?", 1)
	// db.Find(&users, &model.User{Gender: 1})
	// db.Find(&users, &model.User{Gender: 0})
	// gender := uint8(0)
	// db.Find(&users, &model.User{Gender: &gender})
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
