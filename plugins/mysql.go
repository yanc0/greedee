package plugins

import (
	"fmt"
	"github.com/yanc0/greedee/events"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type MySQLPlugin struct {
	Conn *gorm.DB
	MySQLPluginConfig *MySQLPluginConfig
}

type MySQLPluginConfig struct {
	Active bool `toml:"active"`
	Host string `toml:"host"`
	Port int `toml:"port"`
	User string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func NewMySQLPlugin(config *MySQLPluginConfig) *MySQLPlugin {
	return &MySQLPlugin{
		MySQLPluginConfig: config,
	}
}

func (mysql *MySQLPlugin) Name() string {
	return "MySQL"
}

func (mysql *MySQLPlugin) Init() error {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mysql.MySQLPluginConfig.User,
		mysql.MySQLPluginConfig.Password,
		mysql.MySQLPluginConfig.Host,
		mysql.MySQLPluginConfig.Port,
		mysql.MySQLPluginConfig.Database,
	)

	d, err := gorm.Open("mysql", conn)
	if err != nil {
		return err
	}
	mysql.Conn = d
	mysql.Conn.LogMode(false)

	mysql.migrate()

	return nil
}

func (mysql *MySQLPlugin) migrate() {
	mysql.Conn.AutoMigrate(&events.Event{})
}

func (mysql *MySQLPlugin) Send(event events.Event) error {
	if mysql.Conn == nil {
		log.Println("[WARN] MySQL is not initialized, retrying ...")
		err := mysql.Init()
		if err != nil {
			return err
		}
		return nil
	}

	db := mysql.Conn.Create(event)
	if db.Error != nil {
		log.Println("[WARN] MySQL", db.Error.Error(), "retrying ...")
		err := mysql.Init()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
