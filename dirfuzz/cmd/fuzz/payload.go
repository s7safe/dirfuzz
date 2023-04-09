package main

import (
    "bufio"
    "compress/gzip"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

const (
    // 文件的压缩格式
    GZIP_EXT = ".gz"
)

// 读取文件并将其返回为字符串切片。
func ReadPayloadFile(filename string) ([]string, error) {
    // 确定文件的绝对路径。
    absPath, err := filepath.Abs(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path for %s: %v", filename, err)
    }

    // 检查文件是否存在。
    _, err = os.Stat(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read payload file %s: %v", absPath, err)
    }

    // 读取文件内容。
    file, err := os.Open(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open payload file %s: %v", absPath, err)
    }

    defer file.Close()

    // 如果文件是Gzip压缩的，则使用Gzip解压缩。
    if strings.HasSuffix(strings.ToLower(absPath), GZIP_EXT) {
        gz, err := gzip.NewReader(file)
        if err != nil {
            return nil, fmt.Errorf("failed to create gzip reader for %s: %v", absPath, err)
        }
        defer gz.Close()

        reader := bufio.NewReader(gz)

        // 逐行读取文件内容。
        var lines []string
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                break
            }
            lines = append(lines, strings.TrimRight(line, "\r\n"))
        }

        return lines, nil
    }

    // 如果文件不是Gzip压缩的，则直接读取文件内容。
    content, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read payload file %s: %v", absPath, err)
    }

    return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}
