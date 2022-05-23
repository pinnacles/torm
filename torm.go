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
	TableName             string
	Fields                []string
	HasAutoIncrement      bool
	AutoIncrementColumns  []string
	HasAutoCreateTime     bool
	AutoCreateTimeColumns map[string]string
	HasAutoUpdateTime     bool
	AutoUpdateTimeColumns map[string]string
}

func (m tableMeta) IsAutoIncrement(col string) bool {
	if !m.HasAutoIncrement {
		return false
	}
	for _, key := range m.AutoIncrementColumns {
		if key == col {
			return true
		}
	}
	return false
}

func (m tableMeta) IsAutoCreateTime(col string) bool {
	if !m.HasAutoCreateTime {
		return false
	}
	for k := range m.AutoCreateTimeColumns {
		if k == col {
			return true
		}
	}
	return false
}

func (m tableMeta) IsAutoUpdateTime(col string) bool {
	if !m.HasAutoUpdateTime {
		return false
	}
	for k := range m.AutoUpdateTimeColumns {
		if k == col {
			return true
		}
	}
	return false
}

type Schema interface {
	TableName() string
}

func Register(s Schema) {
	rv := reflect.ValueOf(s)
	rt := rv.Type()

	fs := []string{}
	hasAutoIncrement := false
	autoIncrementColumns := []string{}
	hasAutoCreateTime := false
	autoCreateTimeColumns := map[string]string{}
	hasAutoUpdateTime := false
	autoUpdateTimeColumns := map[string]string{}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		col := field.Tag.Get("db")
		if col == "" {
			continue
		}
		fs = append(fs, col)

		fn := field.Tag.Get("torm")
		if fn == "" {
			continue
		}
		switch fn {
		case "autoIncrement":
			hasAutoIncrement = true
			autoIncrementColumns = append(autoIncrementColumns, col)
		case "autoCreateTime":
			hasAutoCreateTime = true
			autoCreateTimeColumns[col] = field.Name
		case "autoUpdateTime":
			hasAutoUpdateTime = true
			autoUpdateTimeColumns[col] = field.Name
		default:
		}
	}

	metas[s.TableName()] = &tableMeta{
		TableName:             s.TableName(),
		Fields:                fs,
		HasAutoIncrement:      hasAutoIncrement,
		AutoIncrementColumns:  autoIncrementColumns,
		HasAutoCreateTime:     hasAutoCreateTime,
		AutoCreateTimeColumns: autoCreateTimeColumns,
		HasAutoUpdateTime:     hasAutoUpdateTime,
		AutoUpdateTimeColumns: autoUpdateTimeColumns,
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
