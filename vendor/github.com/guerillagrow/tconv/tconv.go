package tconv

import (
	"errors"
	"strconv"
	"time"
)

type T interface{}

func T2Str(g interface{}) string {
	v, _ := T2String(g)
	return v
}

func T2String(g interface{}) (string, error) {

	// !TODO: rework as type switch!!!!

	if gstr, ok := g.(string); ok {
		return gstr, nil
	} else if durstr, ok := g.(time.Duration); ok {
		return strconv.Itoa(int(durstr)), nil
	} else if gbyteslc, ok := g.([]byte); ok {
		return string(gbyteslc), nil
	} else if gint64, ok := g.(int64); ok {
		return strconv.FormatInt(gint64, 10), nil
	} else if gint, ok := g.(int); ok {
		return strconv.Itoa(gint), nil
	} else if guint, ok := g.(uint32); ok {
		return strconv.FormatInt(int64(guint), 10), nil
	} else if guint64, ok := g.(uint64); ok {
		return strconv.FormatInt(int64(guint64), 10), nil
	} else if gbool, ok := g.(bool); ok {
		if gbool == true {
			return "1", nil
		} else {
			return "0", nil
		}
	} else if gfloat32, okf32 := g.(float32); okf32 {
		return strconv.FormatFloat(float64(gfloat32), 'f', -1, 32), nil
	} else if gbytes, okbytes := interface{}(g).([]byte); okbytes {
		return string(gbytes), nil
	} else if gfloat64, okf64 := g.(float64); okf64 {
		return strconv.FormatFloat(gfloat64, 'f', -1, 64), nil
	} else {
		return "", errors.New("Couldn't convert generic value to string!")
	}
}

func T2StringSlice(g interface{}) ([]string, error) {
	if gt, ok := g.([]string); ok {
		return gt, nil
	} else {
		return []string{}, errors.New("Couldn't convert generic value to []string!")
	}
}

func T2IntSlice(g interface{}) ([]int, error) {
	if gt, ok := g.([]int); ok {
		return gt, nil
	} else {
		return []int{}, errors.New("Couldn't convert generic value to []int!")
	}
}

func T2Int32Slice(g interface{}) ([]int32, error) {
	if gt, ok := g.([]int32); ok {
		return gt, nil
	} else {
		return []int32{}, errors.New("Couldn't convert generic value to []int32!")
	}
}

func T2Int64Slice(g interface{}) ([]int64, error) {
	if gt, ok := g.([]int64); ok {
		return gt, nil
	} else {
		return []int64{}, errors.New("Couldn't convert generic value to []int64!")
	}
}

func T2Float32Slice(g interface{}) ([]float32, error) {
	if gt, ok := g.([]float32); ok {
		return gt, nil
	} else {
		return []float32{}, errors.New("Couldn't convert generic value to []float64!")
	}
}

func T2Float64Slice(g interface{}) ([]float64, error) {
	if gt, ok := g.([]float64); ok {
		return gt, nil
	} else {
		return []float64{}, errors.New("Couldn't convert generic value to []float64!")
	}
}

func T2GenericSlice(g interface{}) ([]interface{}, error) {
	if gt, ok := g.([]interface{}); ok {
		return gt, nil
	} else {
		return []interface{}{}, errors.New("Couldn't convert generic value to []interface{}!")
	}
}

func T2Bytes(g interface{}) ([]byte, error) {
	v, err := T2String(g)
	return []byte(v), err
}

/*func T2ByteSlice(g interface{}) ([]byte, error) {
	return T2Bytes(g)
}*/

func T2Int(g interface{}) (int64, error) {
	v, err := T2String(g)
	if err != nil {
		return int64(0), err
	}
	i, ierr := strconv.Atoi(v)
	if ierr != nil {
		return int64(0), ierr
	}
	return int64(i), nil
}

func T2UInt32(g interface{}) (uint32, error) {
	v, err := T2String(g)
	if err != nil {
		return uint32(0), err
	}
	f64, ferr := strconv.ParseUint(v, 0, 64)
	return uint32(f64), ferr
}

func T2UInt64(g interface{}) (uint64, error) {
	v, err := T2String(g)
	if err != nil {
		return uint64(0), err
	}
	f64, ferr := strconv.ParseUint(v, 0, 64)
	return f64, ferr
}

func T2Float(g interface{}) (float64, error) {
	v, err := T2String(g)
	if err != nil {
		return float64(0), err
	}
	f64, ferr := strconv.ParseFloat(v, 64)
	return f64, ferr
}

func T2Float64(g interface{}) (float64, error) {
	v, err := T2String(g)
	if err != nil {
		return float64(0), err
	}
	f64, ferr := strconv.ParseFloat(v, 64)
	return f64, ferr
}

func T2Bool(g interface{}) (bool, error) {
	v, err := T2String(g)
	if err != nil {
		return false, err
	}
	if v == "1" || v == "true" {
		return true, nil
	} else {
		return false, nil
	}
}
