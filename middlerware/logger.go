package middlerware

import (
    "fmt"
    "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
    "os"
    "time"

	"../config"
	"../config/bean"
)

var logLevels = map[bean.LogLevel]logrus.Level{
	bean.Debug: logrus.DebugLevel,
	bean.Info:  logrus.InfoLevel,
	bean.Warn:  logrus.WarnLevel,
	bean.Error: logrus.ErrorLevel,
	bean.Fatal: logrus.FatalLevel,
	bean.Panic: logrus.PanicLevel,
}

var lg *logrus.Logger

func init() {
	var lh *lfshook.LfsHook
	var wm lfshook.WriterMap = make(lfshook.WriterMap, 6)
	var rs = make(map[logrus.Level]*rotatelogs.RotateLogs, 6)

	lc := config.AppConfig.Logger

	fmt.Println(lc)

	lf := new(logrus.JSONFormatter)
	lf.TimestampFormat = `2006-01-02 15:04:05`
	lf.DisableTimestamp = false

    lg = logrus.New()
    lg.SetFormatter(lf)
    lg.SetOutput(os.Stdout)

	if l, ok := logLevels[lc.StandLevel]; ok {
        lg.SetLevel(l)
	} else {
        lg.SetLevel(logrus.InfoLevel)
	}

	if len(lc.Files) > 0 {
		for _, f := range lc.Files {
			var wr *rotatelogs.RotateLogs
			var err error
			if f.LinkName != "" {
				wr, err = rotatelogs.New(
					f.FileNameFormat,
					rotatelogs.WithLinkName(f.LinkName),
					rotatelogs.WithRotationTime(lc.RotationTime*time.Second),
					rotatelogs.WithRotationCount(lc.RotationCount),
				)
			} else {
				wr, err = rotatelogs.New(
					f.FileNameFormat,
					rotatelogs.WithRotationTime(lc.RotationTime*time.Second),
					rotatelogs.WithRotationCount(lc.RotationCount),
				)
			}
			if err != nil {
				logrus.Errorf("[middleware]-`logger初始化异常`, error:`%v`", err)
				continue
			}
			for _, lv := range f.Level {
				if l, ok := logLevels[lv]; ok {
					rs[l] = wr
				}
			}
		}
	}
	if len(rs) > 0 {
		for l, w := range rs {
			wm[l] = w
		}
	}
	lh = lfshook.NewHook(wm, &logrus.JSONFormatter{})
    lg.AddHook(lh)
	Cont.Register("log", lg)
}
