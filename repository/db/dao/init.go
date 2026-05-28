package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	conf "github.com/YasinDoyle/e-mall/config"
	trackutil "github.com/YasinDoyle/e-mall/utils/track"
)

var (
	_db *gorm.DB
)

func InitMysql() {
	mConfig := conf.Config.MySql["default"]
	pathRead := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":", mConfig.DbPort, ")/", mConfig.DbName, "?charset=" + mConfig.Charset + "&parseTime=true"}, "")
	pathWrite := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":", mConfig.DbPort, ")/", mConfig.DbName, "?charset=" + mConfig.Charset + "&parseTime=true"}, "")

	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	_db = db
	_ = _db.Use(dbresolver.
		Register(dbresolver.Config{
			// `db2` 作为 sources，`db3`、`db4` 作为 replicas
			Sources:  []gorm.Dialector{mysql.Open(pathRead)},                         // 写操作
			Replicas: []gorm.Dialector{mysql.Open(pathWrite), mysql.Open(pathWrite)}, // 读操作
			Policy:   dbresolver.RandomPolicy{},                                      // sources/replicas 负载均衡策略
		}))
	_db = _db.Set("gorm:table_options", "charset=utf8mb4")
	err = migrate()
	if err != nil {
		panic(err)
	}
	registerTracingCallbacks(_db)
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}

const gormTraceSpanKey = "trace:span"

func registerTracingCallbacks(db *gorm.DB) {
	registerTraceCallback(db, "query", db.Callback().Query().Before("gorm:query"), db.Callback().Query().After("gorm:after_query"))
	registerTraceCallback(db, "create", db.Callback().Create().Before("gorm:create"), db.Callback().Create().After("gorm:after_create"))
	registerTraceCallback(db, "update", db.Callback().Update().Before("gorm:update"), db.Callback().Update().After("gorm:after_update"))
	registerTraceCallback(db, "delete", db.Callback().Delete().Before("gorm:delete"), db.Callback().Delete().After("gorm:after_delete"))
	registerTraceCallback(db, "row", db.Callback().Row().Before("gorm:row"), db.Callback().Row().After("gorm:row"))
	registerTraceCallback(db, "raw", db.Callback().Raw().Before("gorm:raw"), db.Callback().Raw().After("gorm:raw"))
}

type gormTraceRegistrar interface {
	Register(name string, fn func(*gorm.DB)) error
}

func registerTraceCallback(db *gorm.DB, operation string, before gormTraceRegistrar, after gormTraceRegistrar) {
	beforeName := fmt.Sprintf("trace:%s:before", operation)
	afterName := fmt.Sprintf("trace:%s:after", operation)

	_ = before.Register(beforeName, func(tx *gorm.DB) {
		ctx := tx.Statement.Context
		if ctx == nil || opentracing.SpanFromContext(ctx) == nil {
			return
		}

		spanName := fmt.Sprintf("db.%s.%s", operation, tx.Statement.Table)
		span, spanCtx := trackutil.WithSpan(ctx, spanName)
		tx.Statement.Context = spanCtx
		tx.InstanceSet(gormTraceSpanKey, span)
	})

	_ = after.Register(afterName, func(tx *gorm.DB) {
		spanValue, ok := tx.InstanceGet(gormTraceSpanKey)
		if !ok {
			return
		}

		span, ok := spanValue.(opentracing.Span)
		if !ok {
			return
		}

		if tx.Error != nil {
			ext.Error.Set(span, true)
			span.SetTag("error.message", tx.Error.Error())
		}
		span.Finish()
	})
}
