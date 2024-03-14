package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/gaomqq/frame/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type mysqlConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
}

func InitMysql(serviceName string) error {
	type Val struct {
		Mysql mysqlConfig `yaml:"mysql" json:"mysql"`
	}
	mysqlConfigVal := Val{}
	content, err := config.GetConfig("DEFAULT_GROUP", serviceName)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(content), &mysqlConfigVal)
	if err != nil {
		fmt.Println("**********errr")
		return err
	}
	configM := mysqlConfigVal.Mysql
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		configM.Username,
		configM.Password,
		configM.Host,
		configM.Port,
		configM.Database,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return err
}

func WithTX(txFc func(tx *gorm.DB) error) {
	var err error
	tx := DB.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
