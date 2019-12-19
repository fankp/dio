package db

import (
	"devops-integral/basic/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/util/log"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	lock   sync.Mutex
	inited bool
)

func Init() {
	// 开始进行初始化，释放锁
	lock.Lock()
	// 方法执行完成后释放锁
	defer lock.Unlock()
	if inited {
		log.Logf("数据库连接已经初始化完成，请勿重复初始化")
		return
	}
	if config.GetDbConfig() != nil && config.GetDbConfig().GetEnable() {
		// 进行数据库初始化
		log.Logf("开始初始化数据库连接")
		var err error
		dbConfig := config.GetDbConfig()
		db, err = gorm.Open(dbConfig.GetType(), fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
			dbConfig.GetUser(), dbConfig.GetPassword(), dbConfig.GetHost(), dbConfig.GetPort(), dbConfig.GetName()))
		if err != nil {
			log.Errorf("初始化数据库连接失败：", err)
			panic(err)
		}
		// 设置表名前缀
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return dbConfig.GetTablePrefix() + defaultTableName
		}
		db.LogMode(true)
		// 设置最大空闲连接数
		db.DB().SetMaxIdleConns(dbConfig.GetMaxIdleConnection())
		//设置最大打开连接数
		db.DB().SetMaxOpenConns(dbConfig.GetMaxOpenConnection())
		db.SingularTable(true)
		db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
		db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
		db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	}
	inited = true
}

func GetDb() *gorm.DB {
	return db
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		_ = scope.SetColumn("UpdatedOn", time.Now().Unix())
	}
}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
