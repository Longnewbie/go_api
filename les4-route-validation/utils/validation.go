package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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
					base := camelToSnake(part[:idx])
					index := part[idx:]
					parts[i] = base + index
				} else {
					parts[i] = camelToSnake(part)
				}
			}

			filePath := strings.Join(parts, ".")

			switch e.Tag() {
			case "gt":
				errors[filePath] = fmt.Sprintf("%s phai lon hon %s", filePath, e.Param())
			case "lt":
				errors[filePath] = fmt.Sprintf("%s phai nho hon %s", filePath, e.Param())
			case "gte":
				errors[filePath] = fmt.Sprintf("%s phai lon hon hoac bang %s", filePath, e.Param())
			case "lte":
				errors[filePath] = fmt.Sprintf("%s phai nho hon hoac bang %s", filePath, e.Param())
			case "uuid":
				errors[filePath] = fmt.Sprintf("%s phai la uuid hop le", filePath)
			case "slug":
				errors[filePath] = fmt.Sprintf("%s chi duoc chua chu thuong, so, dau gach ngang hoac dau cham", filePath)
			case "min":
				errors[filePath] = fmt.Sprintf("%s phai nhieu hon %s", filePath, e.Param())
			case "max":
				errors[filePath] = fmt.Sprintf("%s phai it hon %s", filePath, e.Param())
			case "min_int":
				errors[filePath] = fmt.Sprintf("%s phai co gia tri lon hon %s", filePath, e.Param())
			case "max_int":
				errors[filePath] = fmt.Sprintf("%s phai co gia tri be hon %s", filePath, e.Param())
			case "oneof":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[filePath] = fmt.Sprintf("%s phai la mot trong cac gia tri: %s", filePath, allowedValue)
			case "required":
				errors[filePath] = fmt.Sprintf("%s la bat buoc", filePath)
			case "search":
				errors[filePath] = fmt.Sprintf("%s chi duoc phep chu thuong, in hoa, so va khoang trang", filePath)
			case "email":
				errors[filePath] = fmt.Sprintf("%s phai dung dinh dang email", filePath)
			case "datetime":
				errors[filePath] = fmt.Sprintf("%s phai dung dinh dang YYY-MM-DD", filePath)
			case "file_ext":
				allowedValue := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[filePath] = fmt.Sprintf("%s chi cho phep nhung file co extension: %s", filePath, allowedValue)
			}
		}

		return gin.H{"error": errors}
	}

	return gin.H{"error": "Invalid" + err.Error()}
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)

	if !ok {
		return fmt.Errorf("Failed to get validator engine")
	}

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseInt(minStr, 10, 64)

		if err != nil {
			return false
		}

		return fl.Field().Int() >= minVal
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		maxVal, err := strconv.ParseInt(maxStr, 10, 64)

		if err != nil {
			return false
		}

		return fl.Field().Int() <= maxVal
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		// fileName = "https://example.com/image.jpg"
		// => filepath = ".jpg"
		fileName := fl.Field().String()

		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")

		for _, allowed := range allowedExt {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}

		return false
	})

	return nil
}
