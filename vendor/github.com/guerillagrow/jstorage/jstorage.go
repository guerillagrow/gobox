package jstorage

// A Config package which is safe for concurrent access in your applications
// You can read some JSON config and access the Data like: conf.Get("path/to/my/key")
// !TODO: Mabye improve to sync.RWMutex

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/guerillagrow/tconv"
)

var ErrEmptySelector error = errors.New("Empty selector!")
var ErrKeyNotFound error = errors.New("Key was not found!")
var ErrInvalidType error = errors.New("Found element has an invalid type!")

func NewStorage() *Storage {
	c := Storage{}
	c.storage = make(map[string]interface{})
	c.Mux = &sync.Mutex{}
	return &c
}

type Storage struct {
	Status       int
	File         string
	FileChecksum string
	storage      map[string]interface{}
	Mux          *sync.Mutex
}

func (self *Storage) check() {
	if self.Mux == nil {
		log.Fatalln("JStorage Mutex was not initialized! Allways use" +
			" NewStorage() to initialize a JStorage instance.")
	}
}

func (self *Storage) FlushStorage() {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	self.storage = make(map[string]interface{})
}

func (self *Storage) LoadJSON(bjson []byte) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	jsonerr := json.Unmarshal(bjson, &self.storage)
	return jsonerr
}

func (self *Storage) SaveFile(filepath string) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	self.File = filepath
	data, jsonerr := json.Marshal(self.storage)
	if jsonerr != nil {
		return jsonerr
	}
	f, ferr := os.Create(filepath)
	if ferr != nil {
		return ferr
	}
	defer f.Close()
	_, err := f.Write(data)
	return err
}

func (self *Storage) SaveFilePretty(filepath string) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	self.File = filepath
	data, jsonerr := json.MarshalIndent(self.storage, "", "    ")
	if jsonerr != nil {
		return jsonerr
	}
	f, ferr := os.Create(filepath)
	if ferr != nil {
		return ferr
	}
	defer f.Close()
	_, err := f.Write(data)
	return err
}

func (self *Storage) LoadFile(filepath string) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	self.File = filepath
	buf, err := ioutil.ReadFile(self.File)
	if err != nil {
		return err
	}
	if self.FileChecksum != "" && Sha1Sum(buf) == self.FileChecksum {
		return nil
	}
	jsonerr := json.Unmarshal(buf, &self.storage)
	return jsonerr
}

func (self *Storage) Set(k string, v interface{}) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	err := self.set(k, v)
	return err
}

func (self *Storage) SetInt64(k string, v interface{}) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	vc, verr := tconv.T2UInt64(v)
	if verr != nil {
		return verr
	}
	err := self.set(k, vc)
	return err
}

func (self *Storage) SetInt32(k string, v interface{}) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	vc, verr := tconv.T2Int32(v)
	if verr != nil {
		return verr
	}

	err := self.set(k, vc)
	return err
}

func (self *Storage) SetFloat64(k string, v interface{}) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	vc, verr := tconv.T2Float64(v)
	if verr != nil {
		return verr
	}

	err := self.set(k, vc)
	return err
}

func (self *Storage) SetFloat32(k string, v interface{}) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	vc, verr := tconv.T2Float32(v)
	if verr != nil {
		return verr
	}

	err := self.set(k, vc)
	return err
}

func (self *Storage) SetString(k string, v interface{}) error {
	return self.Set(k, v)
}

func (self *Storage) set(k string, v interface{}) error {

	keyparts := strings.Split(strings.Trim(k, "/"), "/")
	klen := len(keyparts)
	if klen < 1 {
		return ErrEmptySelector
	}
	var s map[string]interface{} = self.storage

	if len(keyparts) == 1 {
		self.storage[keyparts[0]] = v
		return nil
	}

	for incr, kval := range keyparts {
		_, smok := s[kval]
		if incr+1 < klen {
			if !smok {
				s[kval] = make(map[string]interface{})
			}
			smi, smiok := s[kval].(map[string]interface{})
			if !smiok {
				s[kval] = make(map[string]interface{})
			}
			s = smi
		} else if incr+1 == klen {
			s[kval] = v
			return nil
		}
	}
	return errors.New("Could not set value by selector!")
}

func (self *Storage) Get(k string) (interface{}, error) {
	self.check()
	self.Mux.Lock()

	res, err := self.get(k)
	self.Mux.Unlock()
	return res, err
}

func (self *Storage) get(k string) (interface{}, error) {

	if k == "/" {
		return self.storage, nil
	}
	keyparts := strings.Split(strings.Trim(k, "/"), "/")
	klen := len(keyparts)
	if klen < 1 {
		return "", ErrEmptySelector
	}
	var s map[string]interface{} = self.storage

	for incr, kval := range keyparts {
		tobj, tobjFound := s[kval]
		if !tobjFound { //}|| tobj == nil {
			return "", ErrKeyNotFound
		}

		if incr+1 == klen && tobjFound && tobj != nil {
			return tobj, nil
		}
		sub, ismap := tobj.(map[string]interface{})

		if ismap == true && incr+1 < klen {
			s = sub
			continue
		}
	}

	return "", ErrKeyNotFound

}

func (self *Storage) GetInt(k string) (int, error) {
	r, err := self.Get(k)
	if err != nil {
		return int(0), err
	}
	v, verr := tconv.T2Int(r)
	return int(v), verr
}

func (self *Storage) GetInt64(k string) (int64, error) {
	r, err := self.Get(k)
	if err != nil {
		return int64(0), err
	}
	v, verr := tconv.T2Int(r)
	return v, verr
}

func (self *Storage) GetInt32(k string) (int32, error) {
	r, err := self.Get(k)
	if err != nil {
		return int32(0), err
	}
	v, verr := tconv.T2Int(r)
	return int32(v), verr
}

func (self *Storage) GetFloat(k string) (float64, error) {
	r, err := self.Get(k)
	if err != nil {
		return float64(0), err
	}
	v, verr := tconv.T2Float(r)
	return float64(v), verr
}

func (self *Storage) GetFloat64(k string) (float64, error) {
	r, err := self.Get(k)
	if err != nil {
		return float64(0), err
	}
	v, verr := tconv.T2Float(r)
	return v, verr
}

func (self *Storage) GetFloat32(k string) (float32, error) {
	r, err := self.Get(k)
	if err != nil {
		return float32(0), err
	}
	v, verr := tconv.T2Float(r)
	return float32(v), verr
}

func (self *Storage) GetString(k string) (string, error) {
	r, err := self.Get(k)
	if err != nil {
		return string(""), err
	}
	v, verr := tconv.T2String(r)
	return v, verr
}

func (self *Storage) GetStringSlice(k string) ([]string, error) {
	r, err := self.Get(k)
	if err != nil {
		return []string{}, err
	}
	v, verr := tconv.T2StringSlice(r)
	return v, verr
}

func (self *Storage) AppendStringSlice(k string, v string) error {
	_, terr := self.Get(k)
	if terr == ErrKeyNotFound {
		self.Set(k, []string{})
	}
	tvalGeneric, _ := self.Get(k)
	tslc, cerr := tconv.T2StringSlice(tvalGeneric)
	if cerr != nil {
		return cerr
	}
	tslc = append(tslc, v)
	self.Set(k, tslc)
	return nil
}

func (self *Storage) GetInt64Slice(k string) ([]int64, error) {
	r, err := self.Get(k)
	if err != nil {
		return []int64{}, err
	}
	v, verr := tconv.T2Int64Slice(r)
	return v, verr
}

func (self *Storage) AppendInt64Slice(k string, v int64) error {
	_, terr := self.Get(k)
	if terr == ErrKeyNotFound {
		self.Set(k, []int64{})
	}
	tvalGeneric, _ := self.Get(k)
	tslc, cerr := tconv.T2Int64Slice(tvalGeneric)
	if cerr != nil {
		return cerr
	}
	tslc = append(tslc, v)
	self.Set(k, tslc)
	return nil
}

func (self *Storage) GetFloat64Slice(k string) ([]float64, error) {
	r, err := self.Get(k)
	if err != nil {
		return []float64{}, err
	}
	v, verr := tconv.T2Float64Slice(r)
	return v, verr
}

func (self *Storage) AppendFloat64Slice(k string, v float64) error {
	_, terr := self.Get(k)
	if terr == ErrKeyNotFound {
		self.Set(k, []float64{})
	}
	tvalGeneric, _ := self.Get(k)
	tslc, cerr := tconv.T2Float64Slice(tvalGeneric)
	if cerr != nil {
		return cerr
	}
	tslc = append(tslc, v)
	self.Set(k, tslc)
	return nil
}

func (self *Storage) GetGenericSlice(k string) ([]interface{}, error) {
	r, err := self.Get(k)
	if err != nil {
		return []interface{}{}, err
	}
	v, verr := tconv.T2GenericSlice(r)
	return v, verr
}

func (self *Storage) GetBool(k string) (bool, error) {
	r, err := self.Get(k)
	if err != nil {
		return false, err
	}
	v, verr := tconv.T2Bool(r)
	return v, verr
}
func (self *Storage) GetMap(k string) (map[string]interface{}, error) {
	r, err := self.Get(k)
	if err != nil {
		return make(map[string]interface{}, 0), err
	}
	ifc, ifcOk := r.(map[string]interface{})
	if !ifcOk {
		return make(map[string]interface{}, 0),
			errors.New("Could not convert interface to map!")
	}

	return ifc, nil
}
func (self *Storage) GetStringMap(k string) (map[string]string, error) {
	r, err := self.Get(k)
	if err != nil {
		return make(map[string]string, 0), err
	}
	ifc, ifcOk := r.(map[string]interface{})
	if !ifcOk {
		return make(map[string]string, 0),
			errors.New("Could not convert interface to map!")
	}

	ms := make(map[string]string)

	for k, v := range ifc {
		va, vok := v.(string)
		if vok {
			ms[k] = va
		} else {
			return make(map[string]string, 0),
				errors.New("Could not convert interface to map!")
		}
	}

	return ms, nil
}

func (self *Storage) GetSlice(k string) ([]interface{}, error) {
	r, err := self.Get(k)
	if err != nil {
		return make([]interface{}, 0), err
	}
	ifc, ifcOk := r.([]interface{})
	if !ifcOk {
		return make([]interface{}, 0),
			errors.New("Could not convert interface to slice!")
	}

	return ifc, nil
}

func (self *Storage) Delete(k string) error {
	self.check()
	self.Mux.Lock()
	defer self.Mux.Unlock()

	err := self.delete(k)
	return err
}
func (self *Storage) delete(k string) error {

	keyparts := strings.Split(strings.Trim(k, "/"), "/")
	klen := len(keyparts)
	if klen < 1 {
		return ErrEmptySelector
	}
	var s map[string]interface{} = self.storage

	for incr, kval := range keyparts {
		tobj, tobjFound := s[kval]
		if !tobjFound {
			return ErrKeyNotFound
		}

		if incr+1 == klen && tobjFound {
			delete(s, kval)
			return nil
		}
		sub, ismap := tobj.(map[string]interface{})

		if ismap == true && incr+1 < klen {
			s = sub
			continue
		}
	}

	return ErrKeyNotFound
}

func Sha1Sum(b []byte) string {
	h := sha1.New()
	_, werr := h.Write(b)
	if werr != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
