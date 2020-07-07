package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"bytes"
	"math"
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

var lastName = []string{
    "赵", "钱", "孙", "李", "周", "吴", "郑", "王", "冯", "陈", "褚", "卫", "蒋",
    "沈", "韩", "杨", "朱", "秦", "尤", "许", "何", "吕", "施", "张", "孔", "曹", "严", "华", "金", "魏",
    "陶", "姜", "戚", "谢", "邹", "喻", "柏", "水", "窦", "章", "云", "苏", "潘", "葛", "奚", "范", "彭",
    "郎", "鲁", "韦", "昌", "马", "苗", "凤", "花", "方", "任", "袁", "柳", "鲍", "史", "唐", "费", "薛",
    "雷", "贺", "倪", "汤", "滕", "殷", "罗", "毕", "郝", "安", "常", "傅", "卞", "齐", "元", "顾", "孟",
    "平", "黄", "穆", "萧", "尹", "姚", "邵", "湛", "汪", "祁", "毛", "狄", "米", "伏", "成", "戴", "谈",
    "宋", "茅", "庞", "熊", "纪", "舒", "屈", "项", "祝", "董", "梁", "杜", "阮", "蓝", "闵", "季", "贾",
    "路", "娄", "江", "童", "颜", "郭", "梅", "盛", "林", "钟", "徐", "邱", "骆", "高", "夏", "蔡", "田",
    "樊", "胡", "凌", "霍", "虞", "万", "支", "柯", "管", "卢", "莫", "柯", "房", "裘", "缪", "解", "应",
    "宗", "丁", "宣", "邓", "单", "杭", "洪", "包", "诸", "左", "石", "崔", "吉", "龚", "程", "嵇", "邢",
    "裴", "陆", "荣", "翁", "荀", "于", "惠", "甄", "曲", "封", "储", "仲", "伊", "宁", "仇", "甘", "武",
    "符", "刘", "景", "詹", "龙", "叶", "幸", "司", "黎", "溥", "印", "怀", "蒲", "邰", "从", "索", "赖",
    "卓", "屠", "池", "乔", "胥", "闻", "莘", "党", "翟", "谭", "贡", "劳", "逄", "姬", "申", "扶", "堵",
    "冉", "宰", "雍", "桑", "寿", "通", "燕", "浦", "尚", "农", "温", "别", "庄", "晏", "柴", "瞿", "阎",
    "连", "习", "容", "向", "古", "易", "廖", "庾", "终", "步", "都", "耿", "满", "弘", "匡", "国", "文",
    "寇", "广", "禄", "阙", "东", "欧", "利", "师", "巩", "聂", "关", "荆", "司马", "上官", "欧阳", "夏侯",
    "诸葛", "闻人", "东方", "赫连", "皇甫", "尉迟", "公羊", "澹台", "公冶", "宗政", "濮阳", "淳于", "单于",
    "太叔", "申屠", "公孙", "仲孙", "轩辕", "令狐", "徐离", "宇文", "长孙", "慕容", "司徒", "司空"}
var firstName = []string{
    "伟", "刚", "勇", "毅", "俊", "峰", "强", "军", "平", "保", "东", "文", "辉", "力", "明", "永", "健", "世", "广", "志", "义",
    "兴", "良", "海", "山", "仁", "波", "宁", "贵", "福", "生", "龙", "元", "全", "国", "胜", "学", "祥", "才", "发", "武", "新",
    "利", "清", "飞", "彬", "富", "顺", "信", "子", "杰", "涛", "昌", "成", "康", "星", "光", "天", "达", "安", "岩", "中", "茂",
    "进", "林", "有", "坚", "和", "彪", "博", "诚", "先", "敬", "震", "振", "壮", "会", "思", "群", "豪", "心", "邦", "承", "乐",
    "绍", "功", "松", "善", "厚", "庆", "磊", "民", "友", "裕", "河", "哲", "江", "超", "浩", "亮", "政", "谦", "亨", "奇", "固",
    "之", "轮", "翰", "朗", "伯", "宏", "言", "若", "鸣", "朋", "斌", "梁", "栋", "维", "启", "克", "伦", "翔", "旭", "鹏", "泽",
    "晨", "辰", "士", "以", "建", "家", "致", "树", "炎", "德", "行", "时", "泰", "盛", "雄", "琛", "钧", "冠", "策", "腾", "楠",
    "榕", "风", "航", "弘", "秀", "娟", "英", "华", "慧", "巧", "美", "娜", "静", "淑", "惠", "珠", "翠", "雅", "芝", "玉", "萍",
    "红", "娥", "玲", "芬", "芳", "燕", "彩", "春", "菊", "兰", "凤", "洁", "梅", "琳", "素", "云", "莲", "真", "环", "雪", "荣",
    "爱", "妹", "霞", "香", "月", "莺", "媛", "艳", "瑞", "凡", "佳", "嘉", "琼", "勤", "珍", "贞", "莉", "桂", "娣", "叶", "璧",
    "璐", "娅", "琦", "晶", "妍", "茜", "秋", "珊", "莎", "锦", "黛", "青", "倩", "婷", "姣", "婉", "娴", "瑾", "颖", "露", "瑶",
    "怡", "婵", "雁", "蓓", "纨", "仪", "荷", "丹", "蓉", "眉", "君", "琴", "蕊", "薇", "菁", "梦", "岚", "苑", "婕", "馨", "瑗",
    "琰", "韵", "融", "园", "艺", "咏", "卿", "聪", "澜", "纯", "毓", "悦", "昭", "冰", "爽", "琬", "茗", "羽", "希", "欣", "飘",
    "育", "滢", "馥", "筠", "柔", "竹", "霭", "凝", "晓", "欢", "霄", "枫", "芸", "菲", "寒", "伊", "亚", "宜", "可", "姬", "舒",
    "影", "荔", "枝", "丽", "阳", "妮", "宝", "贝", "初", "程", "梵", "罡", "恒", "鸿", "桦", "骅", "剑", "娇", "纪", "宽", "苛",
    "灵", "玛", "媚", "琪", "晴", "容", "睿", "烁", "堂", "唯", "威", "韦", "雯", "苇", "萱", "阅", "彦", "宇", "雨", "洋", "忠",
    "宗", "曼", "紫", "逸", "贤", "蝶", "菡", "绿", "蓝", "儿", "翠", "烟", "小", "轩"}

var lastNameLen = len(lastName)
var firstNameLen = len(firstName)

func GetRandomName() string {
    rand.Seed(time.Now().UnixNano()) //设置随机数种子
    var first string                 //名
    for i := 0; i <= rand.Intn(1); i++ { //随机产生2位或者3位的名
        first = fmt.Sprint(firstName[rand.Intn(firstNameLen-1)])
    }
    //返回姓名
    return fmt.Sprintf("%s%s", fmt.Sprint(lastName[rand.Intn(lastNameLen-1)]), first)
}

func TimeStamp(str string)int64{
	date,_ := time.Parse("2006-01-02 15:04:05",str)
	return date.Unix()
}

/**
* @des 时间转换函数
* @param atime string 要转换的时间戳（秒）
* @return string
*/
func StrTime (atime int64) string{
    var byTime = []int64{365*24*60*60,24*60*60,60*60,60,1}
    var unit = []string{"年前","天前","小时前","分钟前","秒钟前"}
    now := time.Now().Unix()
    ct := now - atime
    if ct < 0{
        return "刚刚"
    }
    var res string
    for i := 0;i < len(byTime);i++{
        if ct < byTime[i]{
            continue
        }
        var temp = math.Floor(float64(ct / byTime[i]))
        ct = ct % byTime[i];
        if temp > 0 {
            var tempStr string
            tempStr = strconv.FormatFloat(temp,'f',-1,64)
            res = MergeString(tempStr,unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
        }
        break//我想要的形式是精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可（别问我怎么调整，这如果都不会我也是无语）
    }
    return res
}

/**
* @des 拼接字符串
* @param args ...string 要被拼接的字符串序列
* @return string
*/
func MergeString (args ...string) string {
    buffer := bytes.Buffer{}
    for i:=0; i<len(args); i++ {
        buffer.WriteString(args[i])
    }
    return buffer.String()
}

func ShowNum(str string)string{
	num,_ := strconv.Atoi(str);
	if num >= 100000000{
		return strconv.Itoa(num/100000000.0)+"亿"
	}
	if num >= 10000 {
		return strconv.Itoa(num/10000.0) +"万"		
	}
	if num >= 1000 {
		return strconv.Itoa(num/1000.0)+"千"
	}
	return strconv.Itoa(num)
}
