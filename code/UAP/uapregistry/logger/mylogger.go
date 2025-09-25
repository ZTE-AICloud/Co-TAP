package logger

import (
	"fmt"
	"os"

	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// Log level string representations (used in configuration files)
const (
	TraceStr    = "TRACE"
	DebugStr    = "DEBUG"
	InfoStr     = "INFO"
	WarnStr     = "WARNING"
	ErrorStr    = "ERROR"
	CriticalStr = "CRITICAL"
	OffStr      = "OFF"
)

const (
	moduleName = "uapregistry"
)

const (
	oddNumberErrMsg    = "Ignored key without a value."
	nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

type invalidPair struct {
	position   int
	key, value interface{}
}
type invalidPairs []invalidPair

type KVPair struct {
	Key   string
	Value interface{}
}

type MyLoggerInterface interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})

	Debug(params ...interface{})
	Info(params ...interface{})
	Error(params ...interface{})
	Warn(params ...interface{})

	Flush()
}

type MyLogger struct {
	hostName      string
	module        string
	transactionID string
	instanceID    string

	logger *logs.BeeLogger
}

type ServiceSpecificFields struct {
	ObjectID   string
	ObjectType string
	SFMap      map[string]string
}

func (log *MyLogger) Debugw(msg string, keysAndValues ...interface{}) {
	ssfStr := log.buildSSFStr(keysAndValues)
	prefixStr := log.buildLogPrefixStrWithSSF(ssfStr, DebugStr)
	log.logger.Debug(prefixStr + msg)
}

func (log *MyLogger) Infow(msg string, keysAndValues ...interface{}) {
	ssfStr := log.buildSSFStr(keysAndValues)
	prefixStr := log.buildLogPrefixStrWithSSF(ssfStr, InfoStr)
	log.logger.Info(prefixStr + msg)
}

func (log *MyLogger) Warnw(msg string, keysAndValues ...interface{}) {
	ssfStr := log.buildSSFStr(keysAndValues)
	prefixStr := log.buildLogPrefixStrWithSSF(ssfStr, WarnStr)
	log.logger.Warn(prefixStr + msg)
}
func (log *MyLogger) Errorw(msg string, keysAndValues ...interface{}) {
	ssfStr := log.buildSSFStr(keysAndValues)
	prefixStr := log.buildLogPrefixStrWithSSF(ssfStr, ErrorStr)
	log.logger.Error(prefixStr + msg)
}

func (log *MyLogger) Debugf(format string, params ...interface{}) {
	plog := log.buildLogPrefixStr(DebugStr)
	log.logger.Debug(plog+format, params...)
	//	log.logger.Debugf(format, params...)
}

func (log *MyLogger) Infof(format string, params ...interface{}) {
	plog := log.buildLogPrefixStr(InfoStr)
	log.logger.Info(plog+format, params...)
	//	log.logger.Infof(format, params...)
}

func (log *MyLogger) Warnf(format string, params ...interface{}) {
	plog := log.buildLogPrefixStr(WarnStr)
	log.logger.Warn(plog+format, params...)
	//	log.logger.Warnf(format, params...)
}

func (log *MyLogger) Errorf(format string, params ...interface{}) {
	plog := log.buildLogPrefixStr(ErrorStr)
	log.logger.Error(plog+format, params...)
	//	log.logger.Errorf(format, params...)
}

func (log *MyLogger) Debug(params ...interface{}) {
	plog := log.buildLogPrefixStr(DebugStr)
	log.logger.Debug(plog + fmt.Sprint(params...))
	//	log.logger.Debug(params...)
}

func (log *MyLogger) Info(params ...interface{}) {
	plog := log.buildLogPrefixStr(InfoStr)
	log.logger.Info(plog + fmt.Sprint(params...))
	//	log.logger.Info(params...)
}

func (log *MyLogger) Warn(params ...interface{}) {
	plog := log.buildLogPrefixStr(WarnStr)
	log.logger.Warn(plog + fmt.Sprint(params...))
	//	log.logger.Warn(params...)
}

func (log *MyLogger) Error(params ...interface{}) {
	plog := log.buildLogPrefixStr(ErrorStr)
	log.logger.Error(plog + fmt.Sprint(params...))
	//	log.logger.Error(params...)
}

func (log *MyLogger) Flush() {
	log.logger.Flush()
}

func (log *MyLogger) buildLogPrefixStr(level string) string {
	return log.buildLogPrefixStrWithSSF("", level)
}

func (log *MyLogger) buildLogPrefixStrWithSSF(ssfStr, level string) string {
	var (
		unfixedHeaderFields string
	)

	fixedHeaderFields := fmt.Sprintf("%s\t%s\t%s\t", level, log.hostName, log.module)
	if ssfStr != "" {
		unfixedHeaderFields = fmt.Sprintf("TransactionID=%s\tInstanceID=%s\t[%s]\t", log.transactionID, log.instanceID, ssfStr)
	} else {
		unfixedHeaderFields = fmt.Sprintf("TransactionID=%s\tInstanceID=%s\t", log.transactionID, log.instanceID)

	}

	return fixedHeaderFields + unfixedHeaderFields
}

func (log *MyLogger) sweetenKVPair(argstemp []interface{}, fieldstemp []KVPair, invalidtemp invalidPairs) (fields []KVPair, invalid invalidPairs) {
	for i := 0; i < len(argstemp); {
		if f, ok := argstemp[i].(KVPair); ok {
			fieldstemp = append(fieldstemp, f)
			i++
			continue
		}

		if i == len(argstemp)-1 {
			log.logger.Warn(oddNumberErrMsg, "ignored", argstemp[i])
			break
		}

		key, val := argstemp[i], argstemp[i+1]
		if keyStr, ok := key.(string); !ok {
			if cap(invalidtemp) == 0 {
				invalidtemp = make(invalidPairs, 0, len(argstemp)/2)
			}
			invalidtemp = append(invalidtemp, invalidPair{i, key, val})
		} else {
			fieldstemp = append(fieldstemp, KVPair{keyStr, val})
		}
		i += 2
	}
	return fieldstemp, invalidtemp
}

func (log *MyLogger) sweetenKVPairs(args []interface{}) []KVPair {
	if len(args) == 0 {
		return nil
	}

	fields := make([]KVPair, 0, len(args))
	var invalid invalidPairs

	fields, invalid = log.sweetenKVPair(args, fields, invalid)

	if len(invalid) > 0 {
		log.logger.Warn(nonStringKeyErrMsg, KVPair{"invalid", invalid})
	}
	return fields
}

func (log *MyLogger) buildSSFStr(keysAndValues []interface{}) string {
	kvps := log.sweetenKVPairs(keysAndValues)

	var (
		ssfStr string
		//		haveObjectID   bool
		//		haveObjectType bool
	)
	for _, kvp := range kvps {
		switch kvp.Key {
		case "ObjectID":
			if kvp.Value.(string) != "" {
				//				haveObjectID = true
				strTmp := fmt.Sprint(kvp.Key, "=", kvp.Value)
				ssfStr = ssfStr + strTmp + ","
			}
		case "ObjectType":
			if kvp.Value.(string) != "" {
				//				haveObjectType = true
				strTmp := fmt.Sprint(kvp.Key, "=", kvp.Value)
				ssfStr = ssfStr + strTmp + ","
			}
		default:
			strTmp := fmt.Sprint(kvp.Key, "=", kvp.Value)
			ssfStr = ssfStr + strTmp + ","
		}
	}

	return strings.TrimSuffix(ssfStr, ",")
}

func getHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to get hostname:%v", err)
		return "null"
	}

	return hostName
}

func getInstanceID() string {
	return strconv.Itoa(os.Getpid())
}
