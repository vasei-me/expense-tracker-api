package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func (v *Validator) Validate(data interface{}) error {
	val := reflect.ValueOf(data)
	
	// اگر اشاره‌گر باشد، به مقدار اصلی برو
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	var errors []ValidationError
	
	// بررسی فیلدهای struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)
		tag := field.Tag.Get("validate")
		
		if tag == "" {
			continue
		}
		
		// تبدیل نام JSON به نام فیلد
		jsonName := field.Tag.Get("json")
		if jsonName == "" {
			jsonName = field.Name
		}
		jsonName = strings.Split(jsonName, ",")[0]
		
		// اعتبارسنجی بر اساس تگ
		if err := v.validateField(jsonName, fieldValue, tag); err != nil {
			errors = append(errors, *err)
		}
	}
	
	if len(errors) > 0 {
		return &ValidationErrors{Errors: errors}
	}
	
	return nil
}

func (v *Validator) validateField(fieldName string, fieldValue reflect.Value, tag string) *ValidationError {
	rules := strings.Split(tag, ",")
	
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		
		switch {
		case rule == "required":
			if isEmptyValue(fieldValue) {
				return &ValidationError{
					Field: fieldName,
					Error: fmt.Sprintf("%s is required", fieldName),
				}
			}
			
		case strings.HasPrefix(rule, "min="):
			minStr := strings.TrimPrefix(rule, "min=")
			min, err := strconv.Atoi(minStr)
			if err != nil {
				continue
			}
			
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.String() != "" && len(fieldValue.String()) < min {
					return &ValidationError{
						Field: fieldName,
						Error: fmt.Sprintf("%s must be at least %d characters", fieldName, min),
					}
				}
			}
			
		case strings.HasPrefix(rule, "max="):
			maxStr := strings.TrimPrefix(rule, "max=")
			max, err := strconv.Atoi(maxStr)
			if err != nil {
				continue
			}
			
			switch fieldValue.Kind() {
			case reflect.String:
				if fieldValue.String() != "" && len(fieldValue.String()) > max {
					return &ValidationError{
						Field: fieldName,
						Error: fmt.Sprintf("%s must be at most %d characters", fieldName, max),
					}
				}
			}
			
		case strings.HasPrefix(rule, "gt="):
			gtStr := strings.TrimPrefix(rule, "gt=")
			gt, err := strconv.ParseFloat(gtStr, 64)
			if err != nil {
				continue
			}
			
			switch fieldValue.Kind() {
			case reflect.Float64:
				if fieldValue.Float() <= gt {
					return &ValidationError{
						Field: fieldName,
						Error: fmt.Sprintf("%s must be greater than %g", fieldName, gt),
					}
				}
			}
			
		case rule == "email":
			if fieldValue.String() != "" {
				if !isValidEmail(fieldValue.String()) {
					return &ValidationError{
						Field: fieldName,
						Error: fmt.Sprintf("%s must be a valid email address", fieldName),
					}
				}
			}
			
		case strings.HasPrefix(rule, "datetime="):
			format := strings.TrimPrefix(rule, "datetime=")
			if fieldValue.String() != "" {
				_, err := time.Parse(format, fieldValue.String())
				if err != nil {
					return &ValidationError{
						Field: fieldName,
						Error: fmt.Sprintf("%s must be in format %s", fieldName, format),
					}
				}
			}
			
		case rule == "omitempty":
			// هیچ کاری نکن
		}
	}
	
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve *ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range ve.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Error))
	}
	return sb.String()
}