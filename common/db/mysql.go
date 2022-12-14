package db


import (
	"bigrule/common/global"
	blogger "bigrule/common/logger"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	level = map[string]logger.LogLevel{
		"info":   logger.Info,
		"warn":   logger.Warn,
		"error":  logger.Error,
	}
	//
)

func DBSetUp(addr, levelstr string) {
	//db, err := gorm.Open("mysql", Cmqconfig.Mysql.Master)
	db, err := gorm.Open(
		mysql.Open(
			addr),
			&gorm.Config{
				Logger:Default.LogMode(level[levelstr]),
				DisableForeignKeyConstraintWhenMigrating: true,
			},
		)
	if err != nil {
		panic("连接数据库失败")
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	global.DBMysql = db
}


var (
	Discard = New(log.New(ioutil.Discard, "", log.LstdFlags), GormConfig{})
	Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), GormConfig{
		SlowThreshold: 400 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      false,
	})
)

func New(writer Writer, config GormConfig) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &GormLogger{
		Writer:       writer,
		GormConfig:   config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

type GormConfig struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      logger.LogLevel
}

type GormLogger struct {
	Writer
	GormConfig
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	if level == logger.Info{
		newLogger.Colorful = true
	}
	if newLogger.Colorful {
		newLogger.infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		newLogger.warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		newLogger.errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		newLogger.traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		newLogger.traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		newLogger.traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}
	return &newLogger
}

func (g *GormLogger) Info(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel >= logger.Info {
		g.Printf(g.infoStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Warn(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		g.Printf(g.warnStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Error(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.Printf(g.errStr+message, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
			if rows == -1 {
				g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case g.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func (g *GormLogger) Printf(message string, data ...interface{}) {
	switch g.LogLevel {
	case logger.Error:
		blogger.Errorf(message, data...)
	case logger.Warn:
		blogger.Warnf(message, data...)
	case logger.Info:
		blogger.Infof(message, data...)
	}
}
