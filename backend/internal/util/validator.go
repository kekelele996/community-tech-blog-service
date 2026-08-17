package util

import "github.com/go-playground/validator/v10"

// Validate 全局 validator/v10 实例
var Validate = validator.New()

// ValidateStruct 校验结构体，返回第一个错误文案
func ValidateStruct(obj interface{}) error {
	if err := Validate.Struct(obj); err != nil {
		return err
	}
	return nil
}
