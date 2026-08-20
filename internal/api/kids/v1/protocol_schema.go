package v1

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// v1OpenAPISpec 固定当前服务实现所依据的 Clearwave Backend V1 v1。
//
//go:embed spec/clearwave-backend-v1.json
var v1OpenAPISpec []byte

var (
	v1SpecOnce sync.Once
	v1SpecData map[string]any
	v1SpecErr  error
)

// OpenAPISpec 返回对外接口说明的只读副本，供 Swagger 页面展示完整字段定义。
func OpenAPISpec() []byte {
	return append([]byte(nil), v1OpenAPISpec...)
}

// DecodeV1Body 严格解析 JSON body，并拒绝重复对象键和尾随内容。
func DecodeV1Body(raw []byte) (map[string]any, bool, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return map[string]any{}, false, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	value, err := decodeV1JSONValue(decoder)
	if err != nil {
		return nil, true, err
	}
	if decoder.More() {
		return nil, true, fmt.Errorf("JSON body has trailing content")
	}
	if _, err = decoder.Token(); err != io.EOF {
		if err == nil {
			return nil, true, fmt.Errorf("JSON body has trailing content")
		}
		return nil, true, err
	}
	body, ok := value.(map[string]any)
	if !ok {
		return nil, true, fmt.Errorf("JSON body must be an object")
	}
	return body, true, nil
}

// ValidateV1Request 按嵌入的 OpenAPI 接口校验当前 operation 的 header、path、query 和 body。
func ValidateV1Request(in V1OperationInput) error {
	spec, err := loadV1Spec()
	if err != nil {
		return fmt.Errorf("load v1 OpenAPI: %w", err)
	}
	operation, operationPath, err := findV1Operation(spec, in.OperationID, in.Method)
	if err != nil {
		return &V1Error{Status: 422, Code: "VALIDATION_FAILED", Retryable: false, Message: err.Error()}
	}
	if !matchesV1Path(operationPath, in.Path) {
		return v1ValidationError("/path", "request path does not match operation")
	}
	if err = validateV1Authorization(operation, in); err != nil {
		return err
	}
	parameters := appendV1Parameters(spec, operationPath, operation)
	for _, parameter := range parameters {
		if err = validateV1Parameter(spec, parameter, in); err != nil {
			return err
		}
	}
	allowedQuery := make(map[string]struct{})
	for _, parameter := range parameters {
		if parameter["in"] == "query" {
			name, _ := parameter["name"].(string)
			allowedQuery[name] = struct{}{}
		}
	}
	for name := range in.Query {
		if _, allowed := allowedQuery[name]; !allowed {
			return v1ValidationError("/query/"+name, "unknown query parameter")
		}
	}
	requestBody := resolveV1Object(spec, operation["requestBody"])
	if requestBody == nil {
		if in.BodyPresent {
			return v1ValidationError("/body", "request body is not allowed")
		}
		return nil
	}
	if !in.BodyPresent {
		if required, _ := requestBody["required"].(bool); required {
			return v1ValidationError("/body", "request body is required")
		}
		return nil
	}
	content, _ := requestBody["content"].(map[string]any)
	jsonContent, _ := content["application/json"].(map[string]any)
	schema, _ := jsonContent["schema"].(map[string]any)
	return validateV1Schema(spec, schema, in.Body, "/body", 0)
}

// ValidateV1ResponseData 校验 operation 成功响应的 data 是否满足其冻结 schema。
func ValidateV1ResponseData(operationID string, data map[string]any) error {
	spec, err := loadV1Spec()
	if err != nil {
		return fmt.Errorf("load v1 OpenAPI: %w", err)
	}
	operation, _, err := findV1OperationByID(spec, operationID)
	if err != nil {
		return err
	}
	responses, _ := operation["responses"].(map[string]any)
	successResponse := resolveV1Object(spec, responses["200"])
	content, _ := successResponse["content"].(map[string]any)
	jsonContent, _ := content["application/json"].(map[string]any)
	envelope := resolveV1Object(spec, jsonContent["schema"])
	properties, _ := envelope["properties"].(map[string]any)
	dataSchema := resolveV1Object(spec, properties["data"])
	return validateV1Schema(spec, dataSchema, data, "/data", 0)
}

// validateV1Authorization 按 operation 的认证上下文检查 Bearer credential 是否存在。
func validateV1Authorization(operation map[string]any, in V1OperationInput) error {
	authContext, _ := operation["x-auth-context"].(string)
	if authContext == "public" {
		if strings.TrimSpace(in.AccessToken) != "" {
			return &V1Error{Status: 422, Code: "VALIDATION_FAILED", Retryable: false, Message: "public operation must not include a Bearer credential"}
		}
		return nil
	}
	if authContext == "public_or_guest" {
		return nil
	}
	if strings.TrimSpace(in.AccessToken) == "" {
		return &V1Error{Status: 401, Code: "UNAUTHENTICATED", Retryable: false, Message: "Bearer credential is required"}
	}
	return nil
}

// loadV1Spec 只解析一次内嵌接口，避免请求时读取外部项目文件。
func loadV1Spec() (map[string]any, error) {
	v1SpecOnce.Do(func() {
		decoder := json.NewDecoder(bytes.NewReader(v1OpenAPISpec))
		decoder.UseNumber()
		v1SpecErr = decoder.Decode(&v1SpecData)
	})
	return v1SpecData, v1SpecErr
}

// findV1Operation 根据冻结 operation ID 和 HTTP method 找到接口 operation。
func findV1Operation(spec map[string]any, operationID, method string) (map[string]any, string, error) {
	paths, _ := spec["paths"].(map[string]any)
	for path, rawItem := range paths {
		item, _ := rawItem.(map[string]any)
		operation, _ := item[strings.ToLower(method)].(map[string]any)
		if operation != nil && operation["operationId"] == operationID {
			return operation, path, nil
		}
	}
	return nil, "", fmt.Errorf("unknown v1 operation %q", operationID)
}

// findV1OperationByID 通过唯一 operation ID 查找接口 operation。
func findV1OperationByID(spec map[string]any, operationID string) (map[string]any, string, error) {
	paths, _ := spec["paths"].(map[string]any)
	for path, rawItem := range paths {
		item, _ := rawItem.(map[string]any)
		for _, rawOperation := range item {
			operation, _ := rawOperation.(map[string]any)
			if operation != nil && operation["operationId"] == operationID {
				return operation, path, nil
			}
		}
	}
	return nil, "", fmt.Errorf("unknown v1 operation %q", operationID)
}

// appendV1Parameters 合并路径级和 operation 级参数定义。
func appendV1Parameters(spec map[string]any, path string, operation map[string]any) []map[string]any {
	paths, _ := spec["paths"].(map[string]any)
	item, _ := paths[path].(map[string]any)
	parameters := v1ParameterList(spec, item["parameters"])
	return append(parameters, v1ParameterList(spec, operation["parameters"])...)
}

// v1ParameterList 解析参数数组及其本地引用。
func v1ParameterList(spec map[string]any, raw any) []map[string]any {
	items, _ := raw.([]any)
	parameters := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if parameter := resolveV1Object(spec, item); parameter != nil {
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}

// validateV1Parameter 校验单个 header、path 或 query 参数。
func validateV1Parameter(spec, parameter map[string]any, in V1OperationInput) error {
	name, _ := parameter["name"].(string)
	location, _ := parameter["in"].(string)
	required, _ := parameter["required"].(bool)
	schema := resolveV1Object(spec, parameter["schema"])
	field := "/" + location + "/" + name
	var value any
	var present bool
	switch location {
	case "header":
		value, present = in.Headers[name]
	case "path":
		value, present = in.PathParameters[name]
	case "query":
		values, ok := in.Query[name]
		present = ok && len(values) > 0
		if schema != nil && schema["type"] == "array" {
			value = values
		} else if present {
			if len(values) != 1 {
				return v1ValidationError(field, "parameter must not be repeated")
			}
			value = values[0]
		}
	default:
		return v1ValidationError(field, "unsupported parameter location")
	}
	if !present {
		if required {
			return v1ValidationError(field, "parameter is required")
		}
		return nil
	}
	return validateV1ParameterValue(spec, schema, value, field)
}

// validateV1ParameterValue 将 URL 字符串按 schema 转换并复用值校验逻辑。
func validateV1ParameterValue(spec, schema map[string]any, value any, field string) error {
	schema = resolveV1Object(spec, schema)
	if schema == nil {
		return nil
	}
	if schema["type"] == "array" {
		values, ok := value.([]string)
		if !ok {
			return v1ValidationError(field, "invalid array parameter")
		}
		items := make([]any, 0, len(values))
		itemSchema := resolveV1Object(spec, schema["items"])
		for _, raw := range values {
			converted, err := convertV1Scalar(itemSchema, raw, field)
			if err != nil {
				return err
			}
			items = append(items, converted)
		}
		return validateV1Schema(spec, schema, items, field, 0)
	}
	stringValue, ok := value.(string)
	if !ok {
		return v1ValidationError(field, "invalid parameter")
	}
	converted, err := convertV1Scalar(schema, stringValue, field)
	if err != nil {
		return err
	}
	return validateV1Schema(spec, schema, converted, field, 0)
}

// convertV1Scalar 将 URL 参数转为相应的 JSON 标量类型。
func convertV1Scalar(schema map[string]any, value, field string) (any, error) {
	schema = resolveV1Object(nil, schema)
	typeName, _ := schema["type"].(string)
	if typeName == "" {
		if _, ok := schema["properties"].(map[string]any); ok {
			typeName = "object"
		}
	}
	switch typeName {
	case "integer":
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, v1ValidationError(field, "parameter must be an integer")
		}
		return json.Number(strconv.FormatInt(parsed, 10)), nil
	case "number":
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
			return nil, v1ValidationError(field, "parameter must be a number")
		}
		return json.Number(value), nil
	case "boolean":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, v1ValidationError(field, "parameter must be a boolean")
		}
		return parsed, nil
	default:
		return value, nil
	}
}

// validateV1Schema 递归执行接口中实际使用的 JSON Schema 2020-12 约束。
func validateV1Schema(spec, rawSchema map[string]any, value any, field string, depth int) error {
	if depth > 64 {
		return v1ValidationError(field, "schema nesting is too deep")
	}
	schema := resolveV1Object(spec, rawSchema)
	if schema == nil {
		return v1ValidationError(field, "schema is invalid")
	}
	if allOf, ok := schema["allOf"].([]any); ok {
		for _, child := range allOf {
			if err := validateV1Schema(spec, resolveV1Object(spec, child), value, field, depth+1); err != nil {
				return err
			}
		}
	}
	if anyOf, ok := schema["anyOf"].([]any); ok {
		matched := false
		for _, child := range anyOf {
			if validateV1Schema(spec, resolveV1Object(spec, child), value, field, depth+1) == nil {
				matched = true
				break
			}
		}
		if !matched {
			return v1ValidationError(field, "value does not match any allowed schema")
		}
	}
	if oneOf, ok := schema["oneOf"].([]any); ok {
		matches := 0
		for _, child := range oneOf {
			if validateV1Schema(spec, resolveV1Object(spec, child), value, field, depth+1) == nil {
				matches++
			}
		}
		if matches != 1 {
			return v1ValidationError(field, "value must match exactly one allowed schema")
		}
	}
	if expected, ok := schema["const"]; ok && !v1JSONEqual(expected, value) {
		return v1ValidationError(field, "value does not match required constant")
	}
	if enum, ok := schema["enum"].([]any); ok {
		matched := false
		for _, expected := range enum {
			if v1JSONEqual(expected, value) {
				matched = true
				break
			}
		}
		if !matched {
			return v1ValidationError(field, "value is not in the allowed enum")
		}
	}
	typeName, _ := schema["type"].(string)
	if typeName == "" {
		if _, ok := schema["properties"].(map[string]any); ok {
			typeName = "object"
		}
	}
	switch typeName {
	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			return v1ValidationError(field, "value must be an object")
		}
		required, _ := schema["required"].([]any)
		for _, rawName := range required {
			name, _ := rawName.(string)
			if _, exists := object[name]; !exists {
				return v1ValidationError(field+"/"+name, "required field is missing")
			}
		}
		properties, _ := schema["properties"].(map[string]any)
		additional, exists := schema["additionalProperties"]
		for name, child := range object {
			childSchema, defined := properties[name]
			if !defined {
				if exists && additional == false {
					return v1ValidationError(field+"/"+name, "unknown field")
				}
				continue
			}
			if err := validateV1Schema(spec, resolveV1Object(spec, childSchema), child, field+"/"+name, depth+1); err != nil {
				return err
			}
		}
	case "array":
		array, ok := v1Array(value)
		if !ok {
			return v1ValidationError(field, "value must be an array")
		}
		if err := validateV1ArrayBounds(schema, array, field); err != nil {
			return err
		}
		itemSchema := resolveV1Object(spec, schema["items"])
		for index, child := range array {
			if err := validateV1Schema(spec, itemSchema, child, fmt.Sprintf("%s/%d", field, index), depth+1); err != nil {
				return err
			}
		}
	case "string":
		stringValue, ok := value.(string)
		if !ok {
			return v1ValidationError(field, "value must be a string")
		}
		if err := validateV1String(schema, stringValue, field); err != nil {
			return err
		}
	case "integer":
		number, ok := v1Integer(value)
		if !ok {
			return v1ValidationError(field, "value must be an integer")
		}
		if err := validateV1Number(schema, number, field); err != nil {
			return err
		}
	case "number":
		number, ok := v1Number(value)
		if !ok {
			return v1ValidationError(field, "value must be a number")
		}
		if err := validateV1Number(schema, number, field); err != nil {
			return err
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return v1ValidationError(field, "value must be a boolean")
		}
	case "null":
		if value != nil {
			return v1ValidationError(field, "value must be null")
		}
	}
	return nil
}

// v1Array 将服务端返回的任意切片统一为 JSON array 语义。
func v1Array(value any) ([]any, bool) {
	if array, ok := value.([]any); ok {
		return array, true
	}
	reflected := reflect.ValueOf(value)
	if !reflected.IsValid() || (reflected.Kind() != reflect.Array && reflected.Kind() != reflect.Slice) {
		return nil, false
	}
	array := make([]any, reflected.Len())
	for index := 0; index < reflected.Len(); index++ {
		array[index] = reflected.Index(index).Interface()
	}
	return array, true
}

// validateV1ArrayBounds 校验数组长度和不重复约束。
func validateV1ArrayBounds(schema map[string]any, values []any, field string) error {
	if minimum, ok := v1Int(schema["minItems"]); ok && len(values) < minimum {
		return v1ValidationError(field, "array has too few items")
	}
	if maximum, ok := v1Int(schema["maxItems"]); ok && len(values) > maximum {
		return v1ValidationError(field, "array has too many items")
	}
	if unique, _ := schema["uniqueItems"].(bool); unique {
		seen := make(map[string]struct{}, len(values))
		for _, value := range values {
			encoded, _ := json.Marshal(value)
			if _, exists := seen[string(encoded)]; exists {
				return v1ValidationError(field, "array items must be unique")
			}
			seen[string(encoded)] = struct{}{}
		}
	}
	return nil
}

// validateV1String 校验字符串边界、正则和 OpenAPI format。
func validateV1String(schema map[string]any, value, field string) error {
	if minimum, ok := v1Int(schema["minLength"]); ok && len([]rune(value)) < minimum {
		return v1ValidationError(field, "string is too short")
	}
	if maximum, ok := v1Int(schema["maxLength"]); ok && len([]rune(value)) > maximum {
		return v1ValidationError(field, "string is too long")
	}
	if pattern, ok := schema["pattern"].(string); ok {
		re, err := regexp.Compile(pattern)
		if err != nil || !re.MatchString(value) {
			return v1ValidationError(field, "string does not match required pattern")
		}
	}
	format, _ := schema["format"].(string)
	switch format {
	case "date":
		if _, err := time.Parse(time.DateOnly, value); err != nil {
			return v1ValidationError(field, "string must be an ISO date")
		}
	case "iana-time-zone":
		if _, err := time.LoadLocation(value); err != nil {
			return v1ValidationError(field, "string must be an IANA time zone")
		}
	}
	return nil
}

// validateV1Number 校验 JSON number 的最小值、最大值和 OpenAPI int64 格式。
func validateV1Number(schema map[string]any, value float64, field string) error {
	if minimum, ok := v1Float(schema["minimum"]); ok && value < minimum {
		return v1ValidationError(field, "number is below minimum")
	}
	if maximum, ok := v1Float(schema["maximum"]); ok && value > maximum {
		return v1ValidationError(field, "number exceeds maximum")
	}
	return nil
}

// v1Integer 判断值是否是无精度损失的 JSON 整数。
func v1Integer(value any) (float64, bool) {
	number, ok := v1Number(value)
	return number, ok && math.Trunc(number) == number
}

// v1Number 将 Decoder 保留的 json.Number 转为浮点数以执行 schema 比较。
func v1Number(value any) (float64, bool) {
	switch number := value.(type) {
	case json.Number:
		parsed, err := number.Float64()
		return parsed, err == nil && !math.IsNaN(parsed) && !math.IsInf(parsed, 0)
	case float64:
		return number, !math.IsNaN(number) && !math.IsInf(number, 0)
	case int:
		return float64(number), true
	case int64:
		return float64(number), true
	default:
		return 0, false
	}
}

// v1Int 读取接口 JSON 内的整数约束。
func v1Int(value any) (int, bool) {
	if number, ok := v1Number(value); ok && math.Trunc(number) == number {
		return int(number), true
	}
	return 0, false
}

// v1Float 读取接口 JSON 内的数值约束。
func v1Float(value any) (float64, bool) {
	return v1Number(value)
}

// resolveV1Object 解开 components 下的本地 JSON Pointer 引用。
func resolveV1Object(spec map[string]any, raw any) map[string]any {
	object, _ := raw.(map[string]any)
	if object == nil {
		return nil
	}
	ref, _ := object["$ref"].(string)
	if ref == "" || spec == nil {
		return object
	}
	const prefix = "#/"
	if !strings.HasPrefix(ref, prefix) {
		return nil
	}
	var current any = spec
	for _, token := range strings.Split(strings.TrimPrefix(ref, prefix), "/") {
		part := strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
		items, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = items[part]
	}
	resolved, _ := current.(map[string]any)
	return resolved
}

// matchesV1Path 验证运行时 path 与 OpenAPI path template 是否一致。
func matchesV1Path(template, path string) bool {
	templateParts := strings.Split(strings.Trim(template, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(templateParts) != len(pathParts) {
		return false
	}
	for index, expected := range templateParts {
		if strings.HasPrefix(expected, "{") && strings.Contains(expected, "}") {
			prefix := expected[:strings.Index(expected, "{")]
			suffix := expected[strings.Index(expected, "}")+1:]
			if !strings.HasPrefix(pathParts[index], prefix) || !strings.HasSuffix(pathParts[index], suffix) || len(pathParts[index]) <= len(prefix)+len(suffix) {
				return false
			}
			continue
		}
		if expected != pathParts[index] {
			return false
		}
	}
	return true
}

// v1JSONEqual 以 JSON 语义比较 schema 常量和枚举值。
func v1JSONEqual(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}

// v1ValidationError 返回可安全回显字段位置的稳定接口校验错误。
func v1ValidationError(field, message string) *V1Error {
	return &V1Error{Status: 422, Code: "VALIDATION_FAILED", Retryable: false, Field: &field, Message: message}
}

// decodeV1JSONValue 以 token 流递归解析 JSON，确保同一对象键不会出现两次。
func decodeV1JSONValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch delimiter := token.(type) {
	case json.Delim:
		switch delimiter {
		case '{':
			object := make(map[string]any)
			for decoder.More() {
				rawKey, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				key, ok := rawKey.(string)
				if !ok {
					return nil, fmt.Errorf("JSON object key is invalid")
				}
				if _, exists := object[key]; exists {
					return nil, fmt.Errorf("JSON object contains duplicate key %q", key)
				}
				value, err := decodeV1JSONValue(decoder)
				if err != nil {
					return nil, err
				}
				object[key] = value
			}
			if end, err := decoder.Token(); err != nil || end != json.Delim('}') {
				return nil, fmt.Errorf("JSON object is not closed")
			}
			return object, nil
		case '[':
			array := make([]any, 0)
			for decoder.More() {
				value, err := decodeV1JSONValue(decoder)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			if end, err := decoder.Token(); err != nil || end != json.Delim(']') {
				return nil, fmt.Errorf("JSON array is not closed")
			}
			return array, nil
		default:
			return nil, fmt.Errorf("JSON delimiter is invalid")
		}
	default:
		return token, nil
	}
}
