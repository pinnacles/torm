package torm

import (
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

var (
	metas = map[string]*tableMeta{}
)

type tableMeta struct {
	TableName string
	Fields    map[string]string
}

type schema interface {
	TableName() string
}

func Register(s schema) {
	rv := reflect.ValueOf(s)
	rt := rv.Type()

	fs := map[string]string{}
	for i := 0; i < rt.NumField(); i++ {
		fn := rt.Field(i).Tag.Get("db")
		if fn == "" {
			continue
		}
		fs[fn] = fn
	}

	metas[s.TableName()] = &tableMeta{
		TableName: s.TableName(),
		Fields:    fs,
	}
}

func VerboseLevel(level int) {
	switch level {
	case 0:
		logrus.SetLevel(logrus.PanicLevel)
	case 1:
		logrus.SetLevel(logrus.InfoLevel)
	case 2:
		logrus.SetLevel(logrus.DebugLevel)
	case 3:
		logrus.SetLevel(logrus.TraceLevel)
	default:
		panic("VerboseLevel is must be in range of 0 to 3")
	}
}
