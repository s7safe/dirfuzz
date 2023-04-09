package fuzz

import (
    "bytes"
    "regexp"
    "strconv"
    "strings"
    "unicode"

    "github.com/gookit/color"
)

// 定义一个结构体用于存储过滤规则
type Filter struct {
    StatusCode []int         // 状态码过滤规则
    WordSize   []int         // 响应大小过滤规则
    WordRegexp []*regexp.Regexp // 正则表达式过滤规则
    WordList   []string      // 关键字过滤规则
    IgnoreWord []string      // 忽略关键字过滤规则
}

// 新建一个过滤器对象
func NewFilter() *Filter {
    return &Filter{}
}

// 解析状态码过滤规则
func (f *Filter) ParseStatus(status string) error {
    if status == "" {
        return nil
    }

    status = strings.ReplaceAll(status, " ", "") // 去除空格

    for _, s := range strings.Split(status, ",") {
        if strings.Contains(s, "-") { // 处理范围过滤规则，如：200-299
            parts := strings.Split(s, "-")
            if len(parts) != 2 {
                return ErrInvalidFilter
            }

            start, err := strconv.Atoi(parts[0])
            if err != nil {
                return ErrInvalidFilter
            }

            end, err := strconv.Atoi(parts[1])
            if err != nil {
                return ErrInvalidFilter
            }

            for i := start; i <= end; i++ {
                f.StatusCode = append(f.StatusCode, i)
            }
        } else { // 处理单个状态码过滤规则，如：404
            code, err := strconv.Atoi(s)
            if err != nil {
                return ErrInvalidFilter
            }

            f.StatusCode = append(f.StatusCode, code)
        }
    }

    return nil
}

// 解析响应大小过滤规则
func (f *Filter) ParseSize(size string) error {
    if size == "" {
        return nil
    }

    size = strings.ReplaceAll(size, " ", "") // 去除空格

    for _, s := range strings.Split(size, ",") {
        if strings.Contains(s, "-") { // 处理范围过滤规则，如：0-500
            parts := strings.Split(s, "-")
            if len(parts) != 2 {
                return ErrInvalidFilter
            }

            start, err := strconv.Atoi(parts[0])
            if err != nil {
                return ErrInvalidFilter
            }

            end, err := strconv.Atoi(parts[1])
            if err != nil {
                return ErrInvalidFilter
            }

            f.WordSize = append(f.WordSize, start, end)
        } else { // 处理单个响应大小过滤规则，如：1000
            size, err := strconv.Atoi(s)
            if err != nil {
                return ErrInvalidFilter
            }

            f.WordSize = append(f.WordSize, size)
        }
    }

    return nil
}

// 解析正则表达式过滤规则
func (f *Filter) ParseRegexp(regexpStr string) error {
    if regexpStr == "" {
        return nil
    }

    for _, s := range strings.Split(regexpStr, ",") {
        // 处理正则表达式过滤规则，如：(?i)admin|password
        if r, err := regexp.Compile(s); err == nil {
            f.WordRegexp = append(f.WordRegexp, r)
        } else {
            return ErrInvalidFilter
        }
    }

    return nil
}

// 解析关键字过滤规则
func (f *Filter) ParseWordList(wordList string) error {
    if wordList == "" {
        return nil
    }

    f.WordList = splitString(wordList)

    return nil
}


// 判断响应是否符合过滤规则
func (f *Filter) FilterResponse(status int, size int, body []byte) bool {
    // 判断状态码是否符合规则
    if len(f.StatusCode) > 0 && !contains(f.StatusCode, status) {
        return false
    }

    // 判断响应大小是否符合规则
    if len(f.WordSize) > 0 && !sizeInRange(size, f.WordSize) {
        return false
    }

    // 判断正则表达式是否符合规则
    if len(f.WordRegexp) > 0 {
        matched := false
        for _, r := range f.WordRegexp {
            if r.Match(body) {
                matched = true
                break
            }
        }
        if !matched {
            return false
        }
    }

    // 判断关键字是否符合规则
    if len(f.WordList) > 0 {
        found := false
        for _, word := range f.WordList {
            if bytes.Contains(body, []byte(word)) {
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }

    // 判断忽略关键字是否符合规则
    if len(f.IgnoreWord) > 0 {
        for _, word := range f.IgnoreWord {
            if bytes.Contains(body, []byte(word)) {
                return false
            }
        }
    }

    return true
}

// 判断一个值是否在一个整数数组中
func contains(arr []int, val int) bool {
    for _, v := range arr {
        if v == val {
            return true
        }
    }
    return false
}

// 判断响应大小是否在一个范围内
func sizeInRange(size int, rangeArr []int) bool {
    for i := 0; i < len(rangeArr); i += 2 {
        start := rangeArr[i]
        end := rangeArr[i+1]

        if start <= size && size <= end {
            return true
        }
    }
    return false
}

// 将字符串按照逗号分割成切片
func splitString(s string) []string {
    var ret []string
    for _, v := range strings.Split(s, ",") {
        if v != "" {
            ret = append(ret, v)
        }
    }
    return ret
}

// Filter 过滤器
type Filter struct {
    blacklist map[string]bool // 黑名单
    whitelist map[string]bool // 白名单
}

// NewFilter 创建一个过滤器
func NewFilter(blacklist, whitelist []string) *Filter {
    ret := &Filter{
        blacklist: make(map[string]bool),
        whitelist: make(map[string]bool),
    }

    // 将黑名单中的元素添加到 map 中
    for _, v := range blacklist {
        ret.blacklist[v] = true
    }

    // 将白名单中的元素添加到 map 中
    for _, v := range whitelist {
        ret.whitelist[v] = true
    }

    return ret
}

// Match 判断 s 是否符合过滤条件
func (f *Filter) Match(s string) bool {
    // 如果存在白名单，则只允许白名单中的元素通过
    if len(f.whitelist) > 0 {
        return f.whitelist[s]
    }

    // 如果存在黑名单，则只禁止黑名单中的元素通过
    if len(f.blacklist) > 0 {
        return !f.blacklist[s]
    }

    // 如果不存在黑名单或白名单，则允许所有元素通过
    return true
}
