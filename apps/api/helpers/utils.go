package helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

/*
 * ObjectToJSON is a simple function to cast
 * any object/struct to a JSON string
 */
func ObjectToJSON(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)

	if err != nil {
		return "", err
	}
	return string(b), nil
}

/* JSONToObject is a function to cast
 * any JSON string to
 */
func JSONToObject(jsonStr string, obj interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return err
	}
	return nil
}

func ObjectToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	jsonStr, err := ObjectToJSON(obj)

	if err != nil {
		return nil, err
	}
	err = JSONToObject(jsonStr, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func HashPassword(password string) (string, error) {
	// Convert the password string to a byte slice
	pwd := []byte(password)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func ComparePassword(hashedPassword, password string) bool {
	// CompareHashAndPassword will return an error if the hash does not match the password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// If the error is nil, the password was a match
	return err == nil
}

func JsonEncode(w io.Writer, val interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(val)
}

func JsonDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

func FormDecode(r *http.Request, val interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return json.NewDecoder(strings.NewReader(r.Form.Encode())).Decode(val)
}

func WriteError(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(graphql.Result{
		Errors: []gqlerrors.FormattedError{
			{
				Message: message,
			},
		},
	})
}

func ValidString(s interface{}) bool {
	strType := reflect.TypeOf(s)
	if s != nil && s != "" && strType == reflect.TypeOf("") {
		return true
	}
	return false
}

func ValidFloat64(f interface{}) bool {
	float64Type := reflect.TypeOf(f)
	return f != nil && float64Type == reflect.TypeOf(float64(1.5))
}
func ValidFloat32(f interface{}) bool {
	float32Type := reflect.TypeOf(f)
	return f != nil && float32Type == reflect.TypeOf(float32(1.5))
}

func ValidInt16(i interface{}) bool {
	int16Type := reflect.TypeOf(i)

	if int16Type == reflect.TypeOf(int16(12)) {
		return true
	}

	if reflect.TypeOf(i) == reflect.TypeOf("") {
		try := StringToInt16(i.(string))

		if try == 0 && i != nil && i != "" && i != "0" {
			return false
		}
		int16Type := reflect.TypeOf(try)
		return i != nil && int16Type == reflect.TypeOf(int16(12))
	}

	return false
}
func ValidInt32(i interface{}) bool {
	int32Type := reflect.TypeOf(i)

	if int32Type == reflect.TypeOf(int32(12)) {
		return true
	}

	if reflect.TypeOf(i) == reflect.TypeOf("") {
		try := StringToInt32(i.(string))

		if try == 0 && i != nil && i != "" && i != "0" {
			return false
		}
		int32Type := reflect.TypeOf(try)
		return i != nil && int32Type == reflect.TypeOf(int32(12))
	}

	return false
}
func ValidInt64(i interface{}) bool {
	int64Type := reflect.TypeOf(i)

	if int64Type == reflect.TypeOf(int64(12)) {
		return true
	}

	if reflect.TypeOf(i) == reflect.TypeOf("") {
		try := StringToInt64(i.(string))

		if try == 0 && i != nil && i != "" && i != "0" {
			return false
		}
		int64Type := reflect.TypeOf(try)
		return i != nil && int64Type == reflect.TypeOf(int64(12))
	}

	return false
}

func ValidStringArray(sa interface{}) bool {
	stringArrayType := []string{"a", "b", "c"}

	if _, ok := sa.(primitive.A); ok {
		sa = PrimitiveAToArray(sa.(primitive.A))
	}
	if _, ok := sa.(string); ok {
		sa = StringToArray(sa.(string))
	}
	if _, ok := sa.([]string); ok {
		sa = sa.([]string)
	}
	if _, ok := sa.([]interface{}); ok {
		sa = StringArrayToStringSlice(sa.([]interface{}))
	}

	tmpSA := sa.([]string)
	return reflect.TypeOf(tmpSA) == reflect.TypeOf(stringArrayType)
}

func ValidMap(m interface{}) bool {
	o, ok := m.(map[string]interface{})
	return ok && len(o) > 0
}

func StringIncludes(p string, e string) bool {
	included := false

	if p == e || strings.Contains(p, e) {
		included = true
	}

	return included
}

func StringArrayIncludes(p []string, e string) bool {
	included := false

	for _, v := range p {
		if StringIncludes(v, e) {
			included = true
			break
		}
	}

	return included
}

func StringToArray(s string) []string {
	// Convert s like "[ACTIVE DRAFT]" to a []string as ["ACTIVE", "DRAFT"]
	s = s[1 : len(s)-1]
	return strings.Split(s, " ")
}

func StringArrayToString(p []string) string {
	// Convert p like ["ACTIVE", "DRAFT"] to a string as "[ACTIVE DRAFT]"
	return "[" + strings.Join(p, " ") + "]"
}

func StringArrayToStringSlice(p []interface{}) []string {
	var result []string
	for _, v := range p {
		if v != nil {
			result = append(result, v.(string))
		}
	}
	return result
}
func StringToMap(s string) map[string]interface{} {
	// Convert s like "{ACTIVE DRAFT}" to a map[string]interface{} as {"ACTIVE": "DRAFT"}
	s = s[1 : len(s)-1]
	result := make(map[string]interface{})
	pairs := strings.Split(s, " ")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

func CleanQueryString(query string) string {
	queryStr := strings.ReplaceAll(query, "\t", " ")
	queryParamsArr := strings.Split(strings.ReplaceAll(queryStr, "\n", ""), " ")

	var tmpQueryParamsArr []string
	for i := 0; i < len(queryParamsArr); i++ {
		if queryParamsArr[i] != "" {
			tmpQueryParamsArr = append(tmpQueryParamsArr, queryParamsArr[i])
		}
	}

	return strings.Join(tmpQueryParamsArr, " ")
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}
func StringToInt16(s string) int16 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int16(i)
}
func StringToInt32(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int32(i)
}
func StringToInt64(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return int64(i)
}
func ArrayToString(arr []string) string {
	return strings.Join(arr, ",")
}

func PrimitiveAToArray(p primitive.A) []string {
	var result []string
	for _, v := range p {
		if v != nil {
			result = append(result, v.(string))
		}
	}
	return result
}
