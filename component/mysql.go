// Copyright (c) Airy Author. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package component

import (
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

type Mysql struct {
	client *gorm.DB
	BaseComponent
}

func NewMysql() *Mysql {
	m := new(Mysql)
	m.SetType(constants.MysqlComponent)
	return m
}

func (m *Mysql) Init(config *config.AiryConfig) {
	dsn := config.Mysql.Username + ":" + config.Mysql.Password + "@tcp(" + config.Mysql.Path + ")/" + config.Mysql.Dbname + "?" + config.Mysql.Config
	newLogger := glogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold:             config.Mysql.SlowThreshold,             // Slow SQL threshold
			LogLevel:                  logLevel(config.Mysql.LogLevel),        // Log level
			IgnoreRecordNotFoundError: config.Mysql.IgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
			Colorful:                  config.Mysql.Colorful,                  // Disable color
		})

	ormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: config.Mysql.NamingStrategySingularTable,
		},
	})
	if err != nil {
		logger.Errorf("mysql init failed", err)
		panic(fmt.Sprintf("mysql init failed,err=%s", err.Error()))
	}

	sqlDB, _ := ormDB.DB()
	sqlDB.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	m.client = ormDB
	logger.Infof("mysql component init success")
}

func logLevel(logLevel string) glogger.LogLevel {
	var LogLevel glogger.LogLevel
	switch logLevel {
	case "silent", "Silent":
		LogLevel = glogger.Silent
	case "error", "Error":
		LogLevel = glogger.Error
	case "warn", "Warn":
		LogLevel = glogger.Warn
	case "info", "Info":
		LogLevel = glogger.Info
	}
	return LogLevel
}

//NEXT CRUD
