[官方文档]:https://gorm.io/docs/
[中文文档]:http://gorm.book.jasperxu.com/
[中文文档的模型部分]:http://gorm.book.jasperxu.com/models.html
# try_gorm

##### 一个兴趣使然的码贼，突然想搞一下Golang的orm
洒家英文水平有限，[官方文档]看起来十分吃力，所以正在一边阅读gorm官方的[中文文档]并一边跟随文档中的示例尝试写一些代码

## get gorm
```
go get -u github.com/jinzhu/gorm
```
## 开始实验gorm
在本地mysql数据库中进行实验  
首先创建一个database：try_gorm  
#### main.js
```
package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	ID         uint
	Name       string
	CreateTime *time.Time
}

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		db.CreateTable(&User{})
	}
}
```
执行没有发现输出错误，检查数据库中却没有找到应该被创建的user表  
仔细阅读了一下文档没有任何线索  
于是硬着头皮看了一下[官方文档]，发现了 Advanced Topics > Logger 这一页  
#### main.js > func main()
```
func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		db.LogMode(true) //关心一下日志
		db.CreateTable(&User{})
	}
}
```
再执行，这下看到Error了：重复的列id!? 在日志显示的建表语句中发现了两个id!? why?  
最终发现在gorm.Model中已经定给出了4个字段：  
```
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
```
所以
#### main.js > type User struct
```
type User struct {
	ID         uint
	Name       string
	CreateTime *time.Time
}
```
再执行，检查数据库，发现了表users!?  
自动变复数，厉害了，但是洒家不乐意。  
#### main.js > func main()
```
func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		"root", "123", "127.0.0.1", 3306, "try_gorm")
	if db, err := gorm.Open("mysql", connStr); err != nil {
		fmt.Printf("gorm.Open(`mysql`,`%s`)->%s\n\r", connStr, err)
	} else {
		defer db.Close()
		db.LogMode(true)
		db.SingularTable(true) //不要自作主张了
		db.CreateTable(&User{})
	}
}
```
手动删表再执行，检查数据库，发现了表user，并且id自动给上了自增主键。
就很棒，但是洒家不乐意，偏要搞UUID。
于是
#### main.js > type User struct
```
type User struct {
	ID         string
	Name       string
	CreateTime *time.Time
}
```
执行又报错了，也是哈，没指定长度让元芳怎么看。
于是
#### main.js > type User struct
```
type User struct {
	ID         string `gorm:"size:50"`
	Name       string
	CreateTime *time.Time
}
```
执行没问题，检查数据库，列id自动给上了主键，类型VARCHAR(50)，就很棒。  

--------

接下来就是洒家最关心的课题了，gorm如何处理表与表的关系  
在这之前，要先要填饱肚子
