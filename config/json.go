package rrconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type JsonConfig struct {
	m     map[string]interface{}
	cache map[string]interface{}
	mu    sync.RWMutex
}

func LoadJsonConfigFromFile(path string) (*JsonConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadJsonConfigFromBytes(b)
}


func LoadJsonConfigFromBytes(b [] byte) (*JsonConfig, error) {
	var jm map[string]interface{}
	if err := json.Unmarshal(b, &jm); err != nil {
		return nil, err
	}
	s := &JsonConfig{
		m:     jm,
		cache: make(map[string]interface{}),
	}
	return s, nil
}

// Get("a.b.c")
func (s *JsonConfig) Get(key string) (interface{}, error) {
	s.mu.RLock()
	if v, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return v, nil
	}
	s.mu.RUnlock()
	nodes := strings.Split(key, ".")
	m := s.m
	for i := 0; i < len(nodes); i++ {
		if v, ok := m[nodes[i]]; ok {
			// exist
			if vv, okk := v.(map[string]interface{}); okk {
				// not end
				m = vv
			} else {
				s.mu.Lock()
				s.cache[key] = v
				s.mu.Unlock()
				return v, nil
			}
		} else {
			return nil, fmt.Errorf("no value for key %s", key)
		}
	}
	return m, nil
}

func (s *JsonConfig) GetStringSlice(key string) ([]string, error) {
	f, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	if _, ok := f.([]interface{}); !ok {
		return nil, fmt.Errorf("value for key %s is not slice", key)
	}
	sf := f.([]interface{})
	ss := make([]string, len(sf))
	for i, v := range sf {
		if vv, ok := v.(string); ok {
			ss[i] = vv
		}else{
			return nil, fmt.Errorf("%s[%d] is not a string", key, i)
		}
	}
	return ss, nil
}

func (s *JsonConfig) GetString(key string)(string, error) {
	f, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	if _, ok := f.(string); !ok {
		return nil, fmt.Errorf("value for key %s is not string", key)
	}
	return f.(string), nil
}
