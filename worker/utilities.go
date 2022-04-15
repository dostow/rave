package worker

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	// "reflect"
	"strconv"

	"github.com/apex/log"
	"github.com/sethgrid/pester"
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
	obj := reflect.ValueOf(o)
	fieldCount := 0
	obj = obj.Elem()
	fieldCount = obj.NumField()
	for i := 0; i < fieldCount; i++ {
		f := obj.Field(i)
		c := obj.Type().Field(i).Tag
		g := c.Get("ignore")
		if g == "-" {
			continue
		}
		val := f.Interface()
		switch f.Kind() {
		case reflect.String:
			r := gjson.Get(data, val.(string))
			if r.IsObject() {
				m := r.Value().(map[string]interface{})
				for k, v := range m {
					f.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
				}
			} else {
				if len(r.String()) > 0 {
					f.Set(reflect.ValueOf(r.String()))
				}
			}
		case reflect.Ptr:
			if !f.IsNil() {
				parseStructFields(data, val)
			}
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

func postRequest(client *pester.Client, config Config, data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	contentType := "application/json"
	if v, ok := config.Callback.Headers["content-type"].(string); ok {
		contentType = v
	}
	_, err = client.Post(config.Callback.URL, contentType, strings.NewReader(string(d)))
	if err != nil {
		log.WithField("callback", config.Callback).WithError(err).Error(client.LogString())
	} else {
		log.WithField("callback", config.Callback).Debug(client.LogString())
	}
	return err
}
