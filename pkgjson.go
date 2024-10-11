package pkgjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("无法读取文件: %v", err)
	}

	if err := json.Unmarshal(bytes, &pj.Data); err != nil {
		return fmt.Errorf("无法解析 JSON: %v", err)
	}

	return nil
}

// Modify 修改指定键的值
func (pj *PackageJSON) Modify(key string, value interface{}) error {
	pj.mu.Lock()
	defer pj.mu.Unlock()

	if _, exists := pj.Data.Get(key); !exists {
		return fmt.Errorf("键 '%s' 不存在", key)
	}

	// pj.Data[key] = value
	pj.Data.Set(key, value)
	return nil
}

// Add 添加一个新的键值对
func (pj *PackageJSON) Add(key string, value interface{}) error {
	pj.mu.Lock()
	defer pj.mu.Unlock()

	if _, exists := pj.Data.Get(key); exists {
		return fmt.Errorf("键 '%s' 已存在", key)
	}

	pj.Data.Set(key, value)
	return nil
}

// Delete 删除指定键
func (pj *PackageJSON) Delete(key string) error {
	pj.mu.Lock()
	defer pj.mu.Unlock()

	if _, exists := pj.Data.Get(key); !exists {
		return fmt.Errorf("键 '%s' 不存在", key)
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

	if err := ioutil.WriteFile(pj.FilePath, bytes, 0644); err != nil {
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
