package pkgjson

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/iancoleman/orderedmap"
)

// PackageJSON 结构体用于存储 package.json 的内容
type PackageJSON struct {
	FilePath string
	Data     *orderedmap.OrderedMap
	mu       sync.RWMutex
}

// NewPackageJSON 创建一个新的 PackageJSON 实例并读取文件内容
func NewPackageJSON(filePath string) (*PackageJSON, error) {
	pj := &PackageJSON{
		FilePath: filePath,
		Data:     orderedmap.New(),
	}

	if err := pj.Read(); err != nil {
		return nil, err
	}

	return pj, nil
}

// Read 读取 package.json 文件内容到 Data 字段
func (pj *PackageJSON) Read() error {
	pj.mu.Lock()
	defer pj.mu.Unlock()

	file, err := os.Open(pj.FilePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("无法读取文件: %v", err)
	}

	if err := json.Unmarshal(bytes, &pj.Data); err != nil {
		return fmt.Errorf("无法解析 JSON: %v", err)
	}

	return nil
}

// Update 修改指定键的值
func (pj *PackageJSON) Update(key string, value interface{}) error {
	pj.mu.Lock()
	defer pj.mu.Unlock()
	pj.Data.Set(key, value)
	return nil
}

// Delete 删除指定键
func (pj *PackageJSON) Delete(key string) error {
	pj.mu.Lock()
	defer pj.mu.Unlock()

	if _, exists := pj.Data.Get(key); !exists {
		// return fmt.Errorf("键 '%s' 不存在", key)
		return nil
	}

	// delete(pj.Data, key)
	pj.Data.Delete(key)
	return nil
}

// Save 将 Data 字段的内容写回到 package.json 文件
func (pj *PackageJSON) Save() error {
	pj.mu.RLock()
	defer pj.mu.RUnlock()

	bytes, err := json.MarshalIndent(pj.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("无法序列化 JSON: %v", err)
	}

	// 处理 unicode 转义字符
	bytes = []byte(convertUnicodeEscape(string(bytes)))

	if err := os.WriteFile(pj.FilePath, bytes, 0644); err != nil {
		return fmt.Errorf("无法写入文件: %v", err)
	}

	return nil
}

// Print 打印 package.json 的内容
func (pj *PackageJSON) Print() error {
	pj.mu.RLock()
	defer pj.mu.RUnlock()

	bytes, err := json.MarshalIndent(pj.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("无法序列化 JSON: %v", err)
	}

	fmt.Println(string(bytes))
	return nil
}

func convertUnicodeEscape(s string) string {
	// Handle both \u0026 (4-digit) and \U0001F600 (8-digit) style escapes
	// First handle 4-digit escapes (\uXXXX)
	s = regexp.MustCompile(`\\u([0-9A-Fa-f]{4})`).ReplaceAllStringFunc(s, func(match string) string {
		code, _ := strconv.ParseInt(match[2:], 16, 32)
		return string(rune(code))
	})

	// Then handle 8-digit escapes (\UXXXXXXXX)
	s = regexp.MustCompile(`\\U([0-9A-Fa-f]{8})`).ReplaceAllStringFunc(s, func(match string) string {
		code, _ := strconv.ParseInt(match[2:], 16, 32)
		return string(rune(code))
	})

	return s
}
