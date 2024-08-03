package services

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OneNetService struct {
	mux *http.ServeMux
}

func NewOneNet() *OneNetService {
	return &OneNetService{
		mux: http.NewServeMux(),
	}
}

func (oneNet *OneNetService) Init() *http.ServeMux {
	oneNet.mux.HandleFunc("/accept", oneNet.accept)
	return oneNet.mux
}

func (oneNet *OneNetService) accept(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		//验证
		oneNet.Auth(w, r)
		return
	}

	if r.Method == http.MethodPost {
		oneNet.dataResolve(w, r)
	}
	_, _ = w.Write([]byte(""))
}

func (oneNet *OneNetService) Auth(w http.ResponseWriter, r *http.Request) {
	logrus.Debug(r.URL.RawQuery)
	queryMap, err := Parse(r.URL.RawQuery)
	if err != nil {
		logrus.Debug(err)
		return
	}
	logrus.Debug(r.URL.RawQuery, queryMap)
	if msg, ok := queryMap["msg"]; ok {
		_, _ = w.Write([]byte(msg.(string)))
	}
}
func (oneNet *OneNetService) dataResolve(w http.ResponseWriter, r *http.Request) {
	//logrus.Debug(r.MultipartForm.Value)
	decoder := json.NewDecoder(r.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)
	type OneNetMessage struct {
		Msg       string `json:"msg"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
		Time      int64  `json:"time"`
		Id        string `json:"id"`
	}
	var msg OneNetMessage
	logrus.Debug(decoder.Decode(&msg), fmt.Sprintf("%#v", msg))
}

func (oneNet *OneNetService) ResponseSuc(r http.ResponseWriter) {

}

func Parse(s string) (result map[string]interface{}, err error) {
	if s == "" {
		return nil, nil
	}
	result = make(map[string]interface{})
	parts := strings.Split(s, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return nil, err
		}

		for len(key) > 0 && key[0] == ' ' {
			key = key[1:]
		}

		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return nil, err
		}
		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first
		// build nested map
		if err = build(result, keys, value); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// build nested map.
func build(result map[string]interface{}, keys []string, value interface{}) error {
	var (
		length = len(keys)
		key    = strings.Trim(keys[0], "'\"")
	)
	if length == 1 {
		result[key] = value
		return nil
	}

	// The end is slice. like f[], f[a][]
	if keys[1] == "" && length == 2 {
		// TODO nested slice
		if key == "" {
			return nil
		}
		val, ok := result[key]
		if !ok {
			result[key] = []interface{}{value}
			return nil
		}
		children, ok := val.([]interface{})
		if !ok {
			return fmt.Errorf("expected type '[]interface{}' for key '%s', but got '%T'", key, val)
		}
		result[key] = append(children, value)
		return nil
	}
	// The end is slice + map. like v[][a]
	if keys[1] == "" && length > 2 && keys[2] != "" {
		val, ok := result[key]
		if !ok {
			result[key] = []interface{}{}
			val = result[key]
		}
		children, ok := val.([]interface{})
		if !ok {
			return fmt.Errorf(
				"expected type '[]interface{}' for key '%s', but got '%T'",
				key, val,
			)
		}
		if l := len(children); l > 0 {
			if child, ok := children[l-1].(map[string]interface{}); ok {
				if _, ok := child[keys[2]]; !ok {
					_ = build(child, keys[2:], value)
					return nil
				}
			}
		}
		child := map[string]interface{}{}
		_ = build(child, keys[2:], value)
		result[key] = append(children, child)
		return nil
	}

	// map, like v[a], v[a][b]
	val, ok := result[key]
	if !ok {
		result[key] = map[string]interface{}{}
		val = result[key]
	}
	children, ok := val.(map[string]interface{})
	if !ok {
		return fmt.Errorf(
			"expected type 'map[string]interface{}' for key '%s', but got '%T'",
			key, val,
		)
	}
	if err := build(children, keys[1:], value); err != nil {
		return err
	}
	return nil
}
