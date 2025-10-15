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

// --- 配置 ---
var (
	sourceContentDir = "content"
	destContentDir   = "temp_hugo_content"
	// 新增：需要忽略的顶级目录列表
	ignoreDirs             = []string{"templates", "scripts"}
	imageExtensions        = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".tiff"}
	downloadableExtensions = []string{".zip", ".pdf", ".rar", ".docx", ".xlsx", ".tar.gz"}
)

// --- 正则表达式 ---
var (
	wikilinkRegex  = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	attachmentRegex = regexp.MustCompile(`!\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	codeBlockRegex = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")
	dimensionRegex = regexp.MustCompile(`^\d+(x\d+)?$`)
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

		// 如果是目录，检查它是否在我们的忽略列表中
		if info.IsDir() {
			for _, dirToIgnore := range ignoreDirs {
				if info.Name() == dirToIgnore {
					fmt.Printf("⏭️  Skipping ignored directory: %s\n", path)
					return filepath.SkipDir // 告诉 Walk 不要进入此目录
				}
			}
		}

		// 跳过隐藏文件和目录的逻辑
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
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

func processMarkdownFile(srcPath, destPath string) error {
	contentBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	content := string(contentBytes)

	placeholders := make(map[string]string)
	i := 0
	content = codeBlockRegex.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := fmt.Sprintf("___CODEBLOCK_PLACEHOLDER_%d___", i)
		placeholders[placeholder] = match
		i++
		return placeholder
	})

	content = convertAttachments(content)
	content = convertWikilinks(content)

	for placeholder, originalBlock := range placeholders {
		content = strings.ReplaceAll(content, placeholder, originalBlock)
	}

	return os.WriteFile(destPath, []byte(content), 0644)
}

func generateImageHTML(src, alt, dimStr string) string {
	dimStr = strings.TrimSpace(dimStr)
	width, height := "", ""

	if strings.Contains(dimStr, "x") {
		parts := strings.Split(dimStr, "x")
		if len(parts) == 2 {
			width = strings.TrimSpace(parts[0])
			height = strings.TrimSpace(parts[1])
		}
	} else if dimStr != "" {
		width = dimStr
	}

	tag := fmt.Sprintf(`<img src="%s" alt="%s"`, src, alt)
	if width != "" {
		tag += fmt.Sprintf(` width="%s"`, width)
	}
	if height != "" {
		tag += fmt.Sprintf(` height="%s"`, height)
	}
	tag += ">"
	return tag
}

func convertWikilinks(content string) string {
	return wikilinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := wikilinkRegex.FindStringSubmatch(match)
		filename := strings.TrimSpace(parts[1])
		lowerFilename := strings.ToLower(filename)

		for _, ext := range imageExtensions {
			if strings.HasSuffix(lowerFilename, ext) {
				altText := filename
				dimStr := ""
				isDimension := false

				if len(parts) > 2 && parts[2] != "" {
					secondPart := strings.TrimSpace(parts[2])
					if dimensionRegex.MatchString(secondPart) {
						isDimension = true
						dimStr = secondPart
					} else {
						altText = secondPart
					}
				}
				if isDimension {
					return generateImageHTML(filename, altText, dimStr)
				} else {
					return fmt.Sprintf("![%s](%s)", altText, filename)
				}
			}
		}

		displayText := filename
		if len(parts) > 2 && parts[2] != "" {
			displayText = strings.TrimSpace(parts[2])
		}

		for _, ext := range downloadableExtensions {
			if strings.HasSuffix(lowerFilename, ext) {
				return fmt.Sprintf("[%s](%s)", displayText, filename)
			}
		}

		path := strings.ReplaceAll(filename, " ", "-")
		if strings.HasSuffix(path, ".md") {
			path = strings.TrimSuffix(path, ".md")
		}

		return fmt.Sprintf(`[%s]({{< ref "%s" >}})`, displayText, path)
	})
}

func convertAttachments(content string) string {
	return attachmentRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := attachmentRegex.FindStringSubmatch(match)
		filename := strings.TrimSpace(parts[1])
		altText := filename

		if len(parts) > 2 && parts[2] != "" {
			dimStr := strings.TrimSpace(parts[2])
			return generateImageHTML(filename, altText, dimStr)
		}

		return fmt.Sprintf("![%s](%s)", altText, filename)
	})
}

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