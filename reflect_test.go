package golang

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReflection(t *testing.T) {
	printTestTitle("Test User Struct")
	user := User{
		"jack",
		20,
	}
	userType := reflect.TypeOf(user)
	fmt.Println(userType)        // golang.User
	fmt.Println(userType.Name()) // User
	fmt.Println(userType.Kind()) // struct
	for i := 0; i < userType.NumField(); i++ {
		fieldValue := userType.Field(i)
		fmt.Println(fieldValue.Name, fieldValue.Type, fieldValue.Tag)
	}
	/*使用refVal := reflect.ValueOf(val) 为变量创建一个reflect.Value实例。
	如果要使用反射来修改值，则必须获取指向变量的指针 refPtrVal := reflect.ValueOf(&val)，
	如果不这样做，则只能使用反射读取值，但不能修改它。*/
	userValue := reflect.ValueOf(&user).Elem()
	userValue.FieldByName("Name").SetString("rose")
	fmt.Println(user)
	newUser := reflect.New(userType) // return a pointer
	newUser.Elem().FieldByName("Name").SetString("rose")
	newUser.Elem().FieldByName("Age").SetInt(10)
	fmt.Println(user, newUser)
	// acceptUser(newUser) 会报错
	if u, ok := newUser.Interface().(User); ok {
		acceptUser(u)
	}

	printTestTitle("Test Map")
	mapInstance := make(map[int]string)
	mapType := reflect.TypeOf(mapInstance)
	fmt.Println(mapType)        // map[int]string
	fmt.Println(mapType.Name()) //
	fmt.Println(mapType.Kind()) // map
	fmt.Println(mapType.Elem()) // string

	reflectMap := reflect.MakeMap(mapType)
	rk := reflect.ValueOf(1)
	rv := reflect.ValueOf("a")
	reflectMap.SetMapIndex(rk, rv)
	if m, ok := reflectMap.Interface().(map[int]string); ok {
		acceptMap(m)
	}

	printTestTitle("Test Function")
	funcType := reflect.TypeOf(testMakeFunc)
	funcValue := reflect.ValueOf(testMakeFunc)
	fmt.Println(funcType, funcValue)
	newFunc := reflect.MakeFunc(funcType, func(args []reflect.Value) (results []reflect.Value) {
		start := time.Now()
		out := funcValue.Call(args)
		end := time.Now()
		fmt.Println(end.Sub(start))
		return out
	})
	newFunc.Call([]reflect.Value{reflect.ValueOf(10)})

	printTestTitle("Test Make Struct")
	structValue := reflect.ValueOf(testMakeStruct(1, true, "hello world"))
	fmt.Println(structValue)
}

func acceptUser(user User) {
	fmt.Println(user)
}

func acceptMap(mapInstance map[int]string) {
	fmt.Println(mapInstance)
}

func printTestTitle(title string) {
	var underline string
	for i := 0; i < len(title); i++ {
		underline += "="
	}
	fmt.Println(underline)
	fmt.Println(title)
	fmt.Println(underline)
}

func testMakeFunc(count int) {
	sum := 0
	for i := 0; i < count; i++ {
		sum += 1
	}
	fmt.Println(sum)
}

// 动态生成struct
func testMakeStruct(args ...interface{}) interface{} {
	var structList []reflect.StructField
	for index, value := range args {
		argType := reflect.TypeOf(value)
		item := reflect.StructField{
			Name: fmt.Sprintf("Item%d", index),
			Type: argType,
		}
		structList = append(structList, item)
	}
	structType := reflect.StructOf(structList)
	structValue := reflect.New(structType)
	return structValue.Interface()
}
