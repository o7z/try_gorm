[gorm文档]:https://gorm.io/zh_CN/docs/index.html
[高级主题 > 自定义Logger]:https://gorm.io/zh_CN/docs/logger.html
[CRUD 接口 > 更新]:https://gorm.io/zh_CN/docs/update.html
[CRUD 接口 > 删除]:https://gorm.io/zh_CN/docs/delete.html
[关联 > Belong To]:https://gorm.io/zh_CN/docs/belongs_to.html
[关联 > Has One]:https://gorm.io/zh_CN/docs/has_one.html
[关联 > Has Many]:https://gorm.io/zh_CN/docs/has_many.html
[关联 > Has Many]:https://gorm.io/zh_CN/docs/has_many.html
[关联 > Many To Many]:https://gorm.io/zh_CN/docs/many_to_many.html
# 尝试记录

😏洒家正在一边阅读gorm官方的[gorm文档]并一边跟随文档中的示例尝试写一些代码

## get gorm
```
go get -u github.com/jinzhu/gorm
```

## 开始实验gorm

---

### 1.连接数据库
在本地mysql数据库中进行实验  
首先创建一个database：try_gorm  
#### package ./main
```go
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	ID         uint
	Name       string
	CreateTime time.Time
}

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		fmt.Println("connect success.")
	}
}
```

---

### 2.创建一个user表
#### package ./model > type User struct
```go
type User struct {
	gorm.Model
	ID         uint
	Name       string
	CreateTime time.Time
}
```
#### package ./bussniess > var db
```go
var db *gorm.DB
```
#### package ./bussniess > func Init(_db *gorm.DB) error
```go
func Init(_db *gorm.DB) error {
	db = _db
	return nil
}
```
#### package ./bussniess > func Test() error
```go
func Test() error {
	db.CreateTable(&model.User{})
	return nil
}
```
#### package ./main > func main()
```go
func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		if err := bussiness.Init(db); err != nil {
			fmt.Printf("bussiness.Init(db)->%s\n\r", err)
		} else if err := bussiness.Test(); err != nil {
			fmt.Printf("bussiness.Test()->%s\n\r", err)
		}
	}
}
```
执行没有发现输出错误，检查数据库中却没有找到应该被创建的user表  
发现了 [高级主题 > 自定义Logger] 这一页  
#### package ./bussiness > func Init(_db *gorm.DB) error
```go
func Init(_db *gorm.DB) error {
	db = _db
	db.LogMode(true) //日志还是要关心一下的
	return nil
}
```
再执行，这下看到Error了：  
```
Error 1060: Duplicate column name 'id' 

CREATE TABLE `users` (`id` int unsigned AUTO_INCREMENT,`created_at` timestamp NULL,`updated_at` timestamp NULL,`deleted_at` timestamp NULL,`id` int unsigned,`name` varchar(255),`create_time` timestamp NULL , PRIMARY KEY (`id`)) 
```
重复的列id!? 在日志显示的建表语句中发现了两个id!? why?  
发现在gorm.Model中已经定给出了4个字段：  
#### package jinzhu/gorm
```go
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
```
所以  
#### package ./main > type User struct
```go
type User struct {
	// gorm.Model
	ID         uint
	Name       string
	CreateTime time.Time
}
```
再执行未报错，检查数据库，发现了表users!?  
自动变复数，厉害了，但是洒家不乐意。  
#### package ./bussiness > func Init(_db *gorm.DB) error
```go
func Init(_db *gorm.DB) error {
	db = _db
	db.LogMode(true)
	db.SingularTable(true) // 单数表名
	return nil
}
```
手动删表再执行，检查数据库，发现了表user，并且id自动给上了自增主键。
就很棒，但是洒家不乐意，偏要搞UUID。
于是  
#### package ./model > type User struct
```go
type User struct {
	// ID         uint
	ID         string //改为字符串类型
	Name       string
	CreateTime time.Time
}
```
执行又报错了，主键太长，因为没指定长度和类型的字符串默认VARCHAR(255)。  
于是  
#### package ./main > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"` //指定长度50
	Name       string
	CreateTime time.Time
}
```
执行没问题，检查数据库，列id自动给上了主键，类型VARCHAR(50)，就很棒。  
这个表是不是太简陋了呢？
于是
#### package ./model > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	Gender     uint8 //加一个性别字段
	CreateTime time.Time
}
```
怎么改呢？删表再创建？不太好吧。
gorm提供了AutoMigrate方法，翻译过来是自动迁移。
#### package ./bussniess > func Test() error
```go
func Test() error {
	db.AutoMigrate(&model.User{})
	return nil
}
```
执行没报错
```
ALTER TABLE `user` ADD `gender` tinyint unsigned;  
```
注意了，日志显示，这个操作是给表添加了一个字段  
引发了两个需要思考的问题:

>__如果没有表的时候AutoMigrate会自动建表吗？__  
>尝试之后并浅浅阅读了一下源码，得出一个结论：会

>__那当我在结构体中去掉Gender的时候，gorm会不会删除表中对应的这一列呢？__  
>↓

#### package ./model > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	// Gender     uint8
	CreateTime time.Time
}
```
执行结果，啥反应没有，日志没有，数据库中表结构也没有任何变化。  
可以理解，毕竟数据无价，那么当我确定我需要删除一列时该怎么做呢？  
在[gorm文档]中愣是没有找到，于是四处咨询了一下。  
#### package ./bussniess > func Test() error
```go
func Test() error {
	db.Model(&model.User{}).DropColumn("gender") //删除某一列
	return nil
}
```
OK，删掉了  
最后还是把gender字段上，方便后面条件查询  
#### package ./model > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	Gender     uint8
	CreateTime time.Time
}
```

这块理解了个大概，那么
#### package ./bussiness > func Init(_db *gorm.DB) error
```go
func Init(_db *gorm.DB) error {
	db = _db
	db.LogMode(true)
	db.SingularTable(true)
	db.AutoMigrate(&model.User{}) // 自动迁移user表，放到bussines的初始化方法中
	return nil
}
```
#### package ./bussniess > func Test() error
```go
func Test() error {
	// db.AutoMigrate(&model.User{}) 
	return nil
}
```

---

### 3.添加一些user数据
封装一个简单直接的创建UUID的方法  
这里我引入了 github.com/satori/go.uuid
#### package ./bussiness > func NewID() string
```go
// 创建UUID 返回一个字符串
func NewID() string {
	if newUUID, err := uuid.NewV4(); err != nil {
		fmt.Printf("uuid.NewV4()->%s\n\r", err)
		panic(err)
	} else {
		return newUUID.String()
	}
}
```
#### package ./bussniess > func Test() error
```go
func Test() error {
	db.Create(&User{
		ID:         NewID(),
		Name:       "Tom",
		Gender:     1,
		CreateTime: time.Now()})
	db.Create(&User{
		ID:         NewID(),
		Name:       "Jerry",
		Gender:     2,
		CreateTime: time.Now()})
	return nil
}
```
执行结果，数据插入成功。
>__可以批量插入吗？__  
>我目前使用的版本并不支持，听说gorm/v2将会开始支持批量插入  

顺手循环插入了100条记录，以供后边测试

---

### 4.查询user数据
废话说太多，休息下，直接代码
#### package ./bussniess > func Show(v interface{}) error
```go
func Show(v interface{}) error {
	if vJsonBuf, err := json.MarshalIndent(v, "", "  "); err != nil {
		return fmt.Errorf("json.MarshalIndent(v, ``, `  `)->%s\n\r", err)
	} else {
		fmt.Println(string(vJsonBuf))
		return nil
	}
}
```
#### package ./bussniess > func Test() error
```go
func Test() error {
	users := []model.User{}
	db.Find(&users)
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
```
查询OK，试试条件查询
#### package ./bussniess > func Test() error
```go
func Test() error {
	users := []model.User{}
	db.Where("gender=?", 1).Find(&users)
	// db.Find(&users, "gender=?", 1)          // 也可以这么写
	// db.Find(&users, &model.User{Gender: 1}) // 也可以这么写，然而：这么可能会有问题：往下看
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
```
如果结构体字段是值类型且为默认值时，gorm会认为这不是一个条件  
#### package ./bussniess > func Test() error
```go
func Test() error {
	users := []model.User{}
	db.Find(&users, &model.User{Gender: 0}) 
	// db.Find(&users, &model.User{}) // 两句话是一样，因 Gender uint8 默认为0
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
```
走一个，查询到了所有数据
#### package ./bussniess > func Test() error
```go
func Test() error {
	users := []model.User{}
	db.Find(&users, &model.User{Gender: 0})
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
```
走一个，查询到gender为0的所有数据  
结果不一样，所以这里尽量还是使用字符串条件的写法  
如果非要以拿个对象当条件，那就  
#### package ./model > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	Gender     *uint8
	CreateTime time.Time
}
```
#### package ./bussniess > func Test() error
```go
func Test() error {
	users := []model.User{}
	gender := uint8(0)
	db.Find(&users, &model.User{Gender: &gender})
	if err := Show(users); err != nil {
		return fmt.Errorf("Show(users)->%s", err)
	}
	return nil
}
```

---

### 5.修改、删除user数据
查看官方文档中的[CRUD 接口 > 更新]和[CRUD 接口 > 删除]  
这个很简单，也没有遇到什么问题  
非要说该注意什么，也在批量更新的时候  
```go
users := []model.User{}
db.Where("gender=?", 1).Find(&users)
//不要天真的以为可以像这样直接以users更新数据
db.Model(&users).Update("name", "懵逼")
```
这样写没有报错，进而导致所有人清一色“懵逼”

---

### 6.关系型数据库怎么少得了关联
这是洒家最关心的课题  
阅读了文档中  
[关联 > Belong To]  
[关联 > Has One]  
[关联 > Has Many]  
[关联 > Many To Many]  

🤣：  
I’m not sure why this part of the translation is missing.  
My English is not proficient, it's a bit tough for me to understand.  

Belong To 和 Has One 在我脑译过来都是 One To One，  
而且我并不关心这种情况，我认为一对一就应该搞成一个表  

>为什么不是 One To One / One To Many / Many To Many ？  

首先创建一个article表 
#### package ./model > type Article struct
```go
type Article struct {
	ID         string `gorm:"size:50"`
	Title      string
	Content    string `gorm:"type:blob"`
	CreateUser string `gorm:"size:50"`
	CreateTime time.Time
}
```
#### package ./bussiness > func Init(_db *gorm.DB) error
```go
func Init(_db *gorm.DB) error {
	db = _db
	db.LogMode(true)
	db.SingularTable(true)
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Article{}) // 添加一个自动迁移
	db.Model(&model.Article{}).AddForeignKey("create_user", "user(id)", "RESTRICT", "RESTRICT") // 并添加一个外键 
	return nil
}
```
给User结构体添加一个[]Article类型的字段，并标记gorm.foreignkey为CreateUser
#### package ./model > type User struct
```go
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	CreateTime time.Time
	Gender     uint8
	Articles   []Article `gorm:"foreignkey:CreateUser"` // 这里标记一个外键名称，数据库中是create_user，试了一下这里都可以CreateUser和create_user都可以
}
```
这么来试一下，获取第一条user，以这个user的id作为一条article的create_user
#### package ./bussniess > func Test() error
```go
func Test() error {
	user := model.User{}
	db.First(&user)
	db.Create(model.Article{
		ID:         NewID(),
		Title:      "hello",
		Content:    "HELLO!",
		CreateUser: user.ID,
		CreateTime: time.Now()})
	return nil
}
```
然后，查
#### package ./bussniess > func Test() error
```go
func Test() error {
	user := model.User{}
	db.First(&user)
	// article := model.Article{}
	db.Model(&user).Related(&user.Articles, "Articles")
	if err := Show(user); err != nil {
		return fmt.Errorf("Show(user)->%s", err)
	}
	return nil
}
```
输出正常
放假了，回家陪闺女~

