package fuzz

import (
	"bufio"
	"os"
)

// Dictionary是一个包含字典列表的结构
type Dictionary struct {
	words []string
}

// NewDictionary创建一个新的字典结构
func NewDictionary(filename string) (*Dictionary, error) {
	// 打开字典文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取字典文件并存储在words切片中
	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	// 创建字典结构并返回
	return &Dictionary{words}, nil
}

// Has返回字典中是否存在指定的单词
func (d *Dictionary) Has(word string) bool {
	for _, w := range d.words {
		if w == word {
			return true
		}
	}
	return false
}
