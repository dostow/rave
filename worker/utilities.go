package worker

import (
	"fmt"
	"reflect"

	// "reflect"
	"strconv"

	"github.com/tidwall/gjson"
)

// parseMapFields replaces map with parsed content
func parseMapFields(data string, o map[string]interface{}) {
	for k, v := range o {
		switch v.(type) {
		case string:
			r := gjson.Get(data, v.(string))
			// if r.Type()
			if r.IsObject() {
				o[k] = r.Value().(map[string]interface{})
			} else {
				o[k] = r.String()
			}
		case map[string]interface{}:
			parseMapFields(data, v.(map[string]interface{}))
		}
	}
}
func parseStructFields(data string, o interface{}) {
	v := reflect.ValueOf(o)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		val := v.Field(i).Interface()
		switch f.Kind() {
		case reflect.String:
		case reflect.Struct:
			parseStructFields(data, val)
		}
	}

}

// func StringFromPushData(path string, data *PushData) (interface{}, error) {
// 	keys := strings.Split(path, ".")
// 	var value interface{} = data
// 	var err error
// 	for _, key := range keys {
// 		if value, err = Get(key, value); err != nil {
// 			break
// 		}
// 	}
// 	if err == nil {
// 		return value, nil
// 	}
// 	return nil, err
// }

func Get(key string, s interface{}) (v interface{}, err error) {
	var (
		i  int64
		ok bool
	)
	switch s.(type) {
	case map[string]interface{}:
		if v, ok = s.(map[string]interface{})[key]; !ok {
			err = fmt.Errorf("Key not present. [Key:%s]", key)
		}
	case []interface{}:
		if i, err = strconv.ParseInt(key, 10, 64); err == nil {
			array := s.([]interface{})
			if int(i) < len(array) {
				v = array[i]
			} else {
				err = fmt.Errorf("Index out of bounds. [Index:%d] [Array:%v]", i, array)
			}
		}
		// case Signature:
		// 	r := reflect.ValueOf(s)
		// 	v = reflect.Indirect(r).FieldByName(key)
	}
	//fmt.Println("Value:", v, " Key:", key, "Error:", err)
	return v, err
}
