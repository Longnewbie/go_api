package validation

import (
	"fmt"
	"strings"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("Failed to get validator engine")
	}

	RegisterCustomValidation(v)
	return nil
}

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			root := strings.Split(e.Namespace(), ".")[0]
			rawPath := strings.TrimPrefix(e.Namespace(), root+".")
			parts := strings.Split(rawPath, ".")

			for i, part := range parts {
				if strings.Contains(part, "[") {
					idx := strings.Index(part, "[")
					base := utils.CamelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = utils.CamelToSnake(part)
				}
			}

			filePath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[filePath] = fmt.Sprintf("%s phải lớn hơn %s", filePath, e.Param())
			case "lt":
				errors[filePath] = fmt.Sprintf("%s phải nhỏ hơn %s", filePath, e.Param())
			case "gte":
				errors[filePath] = fmt.Sprintf("%s phải lớn hơn hoặc bằng %s", filePath, e.Param())
			case "lte":
				errors[filePath] = fmt.Sprintf("%s phải nhỏ hơn hoặc bằng %s", filePath, e.Param())
			case "uuid":
				errors[filePath] = fmt.Sprintf("%s phải là uuid hợp lệ", filePath)
			case "slug":
				errors[filePath] = fmt.Sprintf("%s chỉ được chứa chữ thường, số, dấu gạch gang hoặc dấu chấm", filePath)
			case "min":
				errors[filePath] = fmt.Sprintf("%s phải nhiều hơn %s", filePath, e.Param())
			case "max":
				errors[filePath] = fmt.Sprintf("%s phải ít hơn %s", filePath, e.Param())
			case "min_int":
				errors[filePath] = fmt.Sprintf("%s phải có giá trị lớn hơn %s", filePath, e.Param())
			case "max_int":
				errors[filePath] = fmt.Sprintf("%s phải có giá trị bé hơn %s", filePath, e.Param())
			case "oneof":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[filePath] = fmt.Sprintf("%s phải là một trong các giá trị: %s", filePath, allowedValue)
			case "required":
				errors[filePath] = fmt.Sprintf("%s là bắt buộc", filePath)
			case "search":
				errors[filePath] = fmt.Sprintf("%s chỉ được phép chữ thường, in hoa, số và khoảng trắng", filePath)
			case "email":
				errors[filePath] = fmt.Sprintf("%s phải đúng định dạng email", filePath)
			case "email_advanced":
				errors[filePath] = fmt.Sprintf("%s này trong danh sách bị cấm", filePath)
			case "password_strong":
				errors[filePath] = fmt.Sprintf("%s phải ít nhất 8 kí tự bao gồm: chữ thường, chữ in hoa, số và kí tự đặc biệt", filePath)
			case "datetime":
				errors[filePath] = fmt.Sprintf("%s phải đúng định dạng YYY-MM-DD", filePath)
			case "file_ext":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[filePath] = fmt.Sprintf("%s chi cho phep nhung file co extension: %s", filePath, allowedValue)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{
		"error":   "Yêu cầu không hợp lệ",
		"details": err.Error()}
}
