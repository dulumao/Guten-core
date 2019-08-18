package validation

// zh-CN language messages
var zhCN = map[string]string{
	"_": "{field} 没有通过验证",
	// int
	"min": "{field} 的最小值是 %d",
	"max": "{field} 的最大值是 %d",
	// Length
	"minLength": "{field} 的最小长度是 %d",
	"maxLength": "{field} 的最大长度是 %d",
	// range
	"enum":  "{field} 值必须在下列枚举中 %v",
	"range": "{field} 值必须在此范围内 %d - %d",
	// required
	"required": "{field} 是必填项",
	// email
	"email": "{field}不是合法邮箱",
	// field compare
	"eqField":  "{field} 值必须等于该字段 %s",
	"neField":  "{field} 值不能等于该字段 %s",
	"ltField":  "{field} 值应小于该字段 %s",
	"lteField": "{field} 值应小于等于该字段 %s",
	"gtField":  "{field} 值应大于该字段 %s",
	"gteField": "{field} 值应大于等于该字段 %s",
}

var en = map[string]string{
	"_":       "{field} did not pass validate", // default message
	"_filter": "{field} data is invalid",       // data filter error
	// int value
	"min": "{field} min value is %d",
	"max": "{field} max value is %d",
	// type check
	"isInt":    "{field} value must be an integer",
	"isInts":   "{field} value must be an int slice",
	"isUint":   "{field} value must be an unsigned integer(>= 0)",
	"isString": "{field} value must be an string",
	// length
	"minLength": "{field} min length is %d",
	"maxLength": "{field} max length is %d",
	// string length. calc rune
	"stringLength": "{field} length must be in the range %d - %d",

	"isFile":  "{field} must be an uploaded file",
	"isImage": "{field} must be an uploaded image file",

	"enum":  "{field} value must be in the enum %v",
	"range": "{field} value must be in the range %d - %d",
	// required
	"required": "{field} is required",
	// field compare
	"eqField":  "{field} value must be equal the field %s",
	"neField":  "{field} value cannot be equal the field %s",
	"ltField":  "{field} value should be less than the field %s",
	"lteField": "{field} value should be less than or equal to field %s",
	"gtField":  "{field} value must be greater the field %s",
	"gteField": "{field} value should be greater or equal to field %s",
}
