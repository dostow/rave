package worker

import (
	"fmt"
	// "reflect"
	"strconv"
)

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
