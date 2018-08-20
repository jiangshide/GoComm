package utils

import (
	"fmt"
	"time"
	"regexp"
	"math/rand"
	"strings"
	"strconv"
	"unicode/utf8"
	"crypto/md5"
)

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func Password(len int, pwdO string) (pwd string, salt string) {
	defaultLen := 4
	if len < defaultLen {
		len = defaultLen
	}
	salt = GetRandomString(len)
	defaultPwd := "jsd123"
	if pwdO != "" {
		defaultPwd = pwdO
	}
	pwd = Md5([]byte(defaultPwd + salt))
	return pwd, salt
}

func Md5(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SizeFormat(size float64) string {
	units := []string{"Byte", "KB", "MB", "GB", "TB"}
	n := 0
	for size > 1024 {
		size /= 1024
		n += 1
	}

	return fmt.Sprintf("%.2f %s", size, units[n])
}

func IsEmail(b []byte) bool {
	return emailPattern.Match(b)
}

// 生成32位MD5
// func MD5(text string) string{
//    ctx := md5.New()
//    ctx.Write([]byte(text))
//    return hex.EncodeToString(ctx.Sum(nil))
// }

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}


/**
 * 字符串首字母转化为大写
 */
func StrFirstToUpper(s string) string {
	if len(s) == 0 {
		return s
	}else if len(s) == 1{
		return strings.ToUpper(s)
	}else{
		return strings.ToUpper(string(s[0]))+s[1:len(s)]
	}
}

/**
 * 字符串首字母转化为小写
 */
func StrFirstToLower(s string) string {
	if len(s) == 0 {
		return s
	}else if len(s) == 1{
		return strings.ToLower(s)
	}else{
		return strings.ToLower(string(s[0]))+s[1:len(s)]
	}
}

func Unicode(rs string) string {
	json := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

func ContainsStr(str string,strArr[] string) bool {
	if len(str) > 0 && len(strArr) > 0{
		for _,v := range strArr{
			if v == str{
				return true
			}
		}
	}
	return false
}

// 过滤 emoji 表情
func FilterEmoji(str string) string {
	newStr := ""
	for _, value := range str {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newStr += string(value)
		}
	}
	return newStr
}

func ContainsNum(str string) bool{
	var reg = regexp.MustCompile("[0-9]+?")
	return reg.MatchString(str)
}

func StringExtractNum(str string ) string{
	var reg = regexp.MustCompile("[0-9]")
	arrStr := reg.FindAllStringSubmatch(str,-1)
	result := ""
	for _,v := range arrStr{
		result += v[0]
	}
	return result
}
