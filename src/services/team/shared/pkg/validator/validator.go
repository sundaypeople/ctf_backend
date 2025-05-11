package validator

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator
func NewValidator() echo.Validator {
	v := validator.New()

	v.RegisterValidation("cidr", cidrValidation)
	v.RegisterValidation("ip", ipValidation)

	return &CustomValidator{validator: v}
}

// cidrValidation はCIDR形式を検証するカスタムバリデーション関数です
func cidrValidation(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	_, _, err := net.ParseCIDR(s)
	return err == nil
}
func ipValidation(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	// 単一のIPアドレスとして解析
	ip := net.ParseIP(s)
	return ip != nil

}

// カスタムヴァリデータを編集
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
			fieldName := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s isn't email format.", fieldName))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long.", fieldName, err.Param()))
			case "cidr":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid CIDR (e.g., 0.0.0.0/24)", fieldName))
			case "ip":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid IP (e.g., 0.0.0.0)", fieldName))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is fail validation", fieldName))
			}
		}
		return fmt.Errorf(strings.Join(errorMessages, ", "))
	}
	return nil
}
