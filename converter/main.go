// converter/main.go
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// --- 配置 (保持不变) ---
var (
	sourceContentDir       = "content"
	destContentDir         = "temp_hugo_content"
	downloadsURLPath       = "/downloads/"
	imageURLPath           = "/images/"
	downloadableExtensions = []string{".zip", ".pdf", ".rar", ".docx", ".xlsx", ".tar.gz"}
)

// --- 正则表达式升级 ---
var (
	// 匹配 wikilinks: [[link|text]]
	wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	// 匹配图片附件: ![[image.png]]
	attachmentRegex = regexp.MustCompile(`!\[\[([^\]]+)\]\]`)
	// 关键新增：匹配 Markdown/Hugo 的代码块。
	// 这会匹配 ```...``` 和 ~~~...~~~ 形式的代码块
	codeBlockRegex = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
)

func main() {
	fmt.Println("🚀 Starting Go conversion process...")

	if err := os.RemoveAll(destContentDir); err != nil {
		fmt.Printf("Error removing old destination directory: %v\n", err)
	}
	if err := os.MkdirAll(destContentDir, 0755); err != nil {
		fmt.Printf("Error creating destination directory: %v\n", err)
		os.Exit(1)
	}

	err := filepath.Walk(sourceContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				fmt.Printf("⏭️  Skipping hidden directory: %s\n", path)
				return filepath.SkipDir
			}
			fmt.Printf("⏭️  Skipping hidden file: %s\n", path)
			return nil
		}

		relativePath, err := filepath.Rel(sourceContentDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destContentDir, relativePath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return processMarkdownFile(path, destPath)
		} else {
			return copyFile(path, destPath)
		}
	})

	if err != nil {
		fmt.Printf("❌ Error during file processing: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Go conversion complete. Output in '%s'.\n", destContentDir)
}


// --- 核心逻辑大改动 ---
// processMarkdownFile 读取、转换并写入 Markdown 文件
func processMarkdownFile(srcPath, destPath string) error {
	contentBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	// --- 全新的转换策略 ---
	// 1. 找到所有代码块，并用唯一的占位符替换它们
	placeholders := make(map[string]string)
	i := 0
	content = codeBlockRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := fmt.Sprintf("___CODEBLOCK_PLACEHOLDER_%d___", i)
		placeholders[placeholder] = match
		i++
		return placeholder
	})

	// 2. 现在，对不包含代码块的内容进行 wikilink 和 attachment 转换
	content = convertWikilinks(content)
	content = convertAttachments(content)

	// 3. 将代码块恢复原状
	for placeholder, originalBlock := range placeholders {
		content = strings.ReplaceAll(content, placeholder, originalBlock)
	}
	// --- 策略结束 ---

	return os.WriteFile(destPath, []byte(content), 0644)
}


// convertWikilinks 智能转换 [[wikilinks]]
func convertWikilinks(content string) string {
	return wikilinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := wikilinkRegex.FindStringSubmatch(match)
		filename := strings.TrimSpace(parts[1])
		
		// 修正：确保即使没有 `|`，displayText 也有一个默认值
		displayText := filename
		if len(parts) > 2 && parts[2] != "" {
			displayText = strings.TrimSpace(parts[2])
		}

		// 检查是否是可下载文件
		for _, ext := range downloadableExtensions {
			if strings.HasSuffix(strings.ToLower(filename), ext) {
				url := downloadsURLPath + filename
				// 始终生成标准的 Markdown 链接
				return fmt.Sprintf("[%s](%s)", displayText, url)
			}
		}

		// 默认视为指向另一篇 Markdown 文章
		path := strings.ReplaceAll(filename, " ", "-")
		if strings.HasSuffix(path, ".md") {
			path = strings.TrimSuffix(path, ".md")
		}
		
		// 关键修复：始终使用 [displayText]({{< ref >}}) 结构来确保链接是可点击的
		return fmt.Sprintf(`[%s]({{< ref "%s" >}})`, displayText, path)
	})
}

// convertAttachments 转换 ![[attachments]]
func convertAttachments(content string) string {
	return attachmentRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := attachmentRegex.FindStringSubmatch(match)
		filename := strings.TrimSpace(parts[1])
		baseFilename := filepath.Base(filename)
		return fmt.Sprintf("![%s](%s%s)", baseFilename, imageURLPath, baseFilename)
	})
}

// copyFile 是一个简单的文件复制工具
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}