package utils

import (
	"encoding/xml"
	"errors"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhcn "github.com/go-playground/validator/v10/translations/zh"
	"github.com/hashicorp/go-version"
	"go.uber.org/zap"
)

const defaultConn = "default"

// X is a convenient alias for a map[string]interface{}.
type X map[string]interface{}

// CDATA XML CDATA section which is defined as blocks of text that are not parsed by the parser, but are otherwise recognized as markup.
type CDATA string

// MarshalXML 将接收器编码为零个或多个 XML 元素。
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

// Date格式化本地时间日期
//返回根据给定格式字符串使用给定的 int64 时间戳格式化的字符串。
//默认格式为：2006-01-02 15:04:05。
func Date(timestamp int64, layout ...string) string {
	l := "2006-01-02 15:04:05"

	if len(layout) != 0 {
		l = layout[0]
	}

	date := time.Unix(timestamp, 0).Local().Format(l)

	return date
}

// StrToTime 将英文文本日期时间描述解析为 Unix 时间戳。
// 默认格式为： 2006-01-02 15:04:05.
func StrToTime(datetime string, layout ...string) int64 {
	l := "2006-01-02 15:04:05"

	if len(layout) != 0 {
		l = layout[0]
	}

	t, err := time.ParseInLocation(l, datetime, time.Local)

	// mismatch layout
	if err != nil {
		zap.L().Error("parse layout mismatch", zap.Error(err))

		return 0
	}

	return t.Unix()
}

// WeekAround 返回当前周的星期一和星期日的日期
func WeekAround(t time.Time) (monday, sunday string) {
	weekday := t.Local().Weekday()

	// monday
	offset := int(time.Monday - weekday)

	if offset > 0 {
		offset = -6
	}

	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	monday = today.AddDate(0, 0, offset).Format("20060102")

	// sunday
	offset = int(time.Sunday - weekday)

	if offset < 0 {
		offset += 7
	}

	sunday = today.AddDate(0, 0, offset).Format("20060102")

	return
}

// IP2Long 将包含 (IPv4) Internet 协议点地址的字符串转换为 uint32 整数。
func IP2Long(ip string) uint32 {
	ipv4 := net.ParseIP(ip).To4()

	if ipv4 == nil {
		return 0
	}

	return uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3])
}

// Long2IP将 uint32 整数地址转换为 (IPv4) Internet 标准点分格式的字符串。
func Long2IP(ip uint32) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// Validator validator can be used for Gin
type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validator) ValidateStruct(obj interface{}) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}

	if err := v.validator.Struct(obj); err != nil {
		e, ok := err.(validator.ValidationErrors)

		if !ok {
			return err
		}

		errM := e.Translate(v.translator)
		msgs := make([]string, 0, len(errM))

		for _, v := range errM {
			msgs = append(msgs, v)
		}

		return errors.New(strings.Join(msgs, ";"))
	}

	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v10
func (v *Validator) Engine() interface{} {
	return v.validator
}

// NewValidator returns a new validator.
// Used for Gin: binding.Validator = yiigo.NewValidator()
func NewValidator() *Validator {
	locale := zh.New()
	uniTrans := ut.New(locale)

	validate := validator.New()
	validate.SetTagName("valid")

	translator, _ := uniTrans.GetTranslator("zh")

	zhcn.RegisterDefaultTranslations(validate, translator)

	return &Validator{
		validator:  validate,
		translator: translator,
	}
}

// VersionCompare 比较语义版本范围，支持: >, >=, =, !=, <, <=, | (or), & (and)
// eg: 1.0.0, =1.0.0, >2.0.0, >=1.0.0&<2.0.0, <2.0.0|>3.0.0, !=4.0.4
func VersionCompare(rangeVer, curVer string) bool {
	if rangeVer == "" || curVer == "" {
		return true
	}

	semVer, err := version.NewVersion(curVer)

	// invalid semantic version
	if err != nil {
		zap.L().Warn("invalid semantic version", zap.Error(err), zap.String("range_version", rangeVer), zap.String("cur_version", curVer))

		return true
	}

	orVers := strings.Split(rangeVer, "|")

	for _, ver := range orVers {
		andVers := strings.Split(ver, "&")

		constraints, err := version.NewConstraint(strings.Join(andVers, ","))

		if err != nil {
			zap.L().Error("version compared error", zap.Error(err), zap.String("range_version", rangeVer), zap.String("cur_version", curVer))

			return true
		}

		if constraints.Check(semVer) {
			return true
		}
	}

	return false
}
