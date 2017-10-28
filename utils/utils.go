package utils

import (
	"crypto/rand"
	r "math/rand"
	"strconv"
	"strings"
	"time"
)

// 随机字符串
func RandomCreateBytes(n int, alphabets ...byte) []byte {

	const alphanum = "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}

	return bytes
}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 {
		panic("end is wrong")
	}

	if end >= length {
		end = length
	}

	return string(rs[start:end])
}

// 将 k 转为1000
func ConKtoInt(num string) (float64, error) {
	var _num float64
	var err error
	if strings.HasSuffix(num, "k") {
		n := strings.TrimSuffix(num, "k")
		_num, err = strconv.ParseFloat(n, 10)
	} else if strings.HasSuffix(num, "K") {
		n := strings.TrimSuffix(num, "K")
		_num, err = strconv.ParseFloat(n, 10)
	}

	return _num * 1000, err
}

var months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

func Month(m string) int {
	return months[m]
}
