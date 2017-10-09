package event

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanc0/greedee/events"
	"log"
)

type MySQLPlugin struct {
	Conn              *gorm.DB
	MySQLPluginConfig *MySQLPluginConfig
}

type MySQLPluginConfig struct {
	Active   bool   `toml:"active"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
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

	mysql.ProcessAll(event)

	return nil
}

func (mysql *MySQLPlugin) GetExpiredAndNotProcessed() ([]events.Event, error) {
	if mysql.Conn == nil {
		log.Println("[WARN] MySQL is not initialized, retrying ...")
		err := mysql.Init()
		if err != nil {
			return nil, err
		}
	}
	var events []events.Event
	mysql.Conn.Where("expires_at <= NOW() and processed = 0 and ttl <> 0 and status = 0").Find(&events)

	return events, nil
}

func (mysql *MySQLPlugin) Process(e events.Event, expired bool) error {
	if mysql.Conn == nil {
		log.Println("[WARN] MySQL is not initialized, retrying ...")
		err := mysql.Init()
		if err != nil {
			return err
		}
	}

	e.Processed = true
	e.Expired = expired
	mysql.Conn.Save(&e)

	return nil
}

func (mysql *MySQLPlugin) ProcessAll(e events.Event) error {
	if mysql.Conn == nil {
		log.Println("[WARN] MySQL is not initialized, retrying ...")
		err := mysql.Init()
		if err != nil {
			return err
		}
	}

	db := mysql.Conn.Exec("update events set processed = 1, expired = 0 "+
		"where name = ? and auth_user_source = ? and id <> ? and ttl > 0",
		e.Name, e.AuthUserSource, e.ID)
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
