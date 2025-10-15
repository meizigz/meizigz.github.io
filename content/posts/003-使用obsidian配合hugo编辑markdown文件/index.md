---
draft: false
date: 2025-10-09T17:26:46+08:00
updated: 2025-10-10T11:46:32+08:00
title: 使用obsidian配合hugo编辑markdown文件
categories:
  - 技术
tags:
  - obsidian
  - hugo
lastmod: 2025-10-15T12:11:58+08:00
---
## vault

仓库根目录`content`  
新建笔记默认文件夹：当前笔记所在文件夹  
附件默认存放路径：当前笔记所在文件夹  
编辑器：关闭拼写检查  

## templater插件

新建templates文件夹放下面的模板文件。

```
<%*
// 1. 弹出输入框获取原始文件名
const rawFileName = await tp.system.prompt("请输入文档名称 (例如: 001-my-post)");
if (!rawFileName) {
  return; // 如果用户取消，则中止
}

// 2. 处理文件名以获取标题
const prefixRegex = /^[a-zA-Z0-9]+-/;
let processedTitle = rawFileName;
if (prefixRegex.test(rawFileName)) {
  processedTitle = rawFileName.replace(prefixRegex, "");
}

// 3. 设置 Hugo frontmatter 的标题
tR += `---
title: "${processedTitle}"
date: ${tp.date.now("YYYY-MM-DDTHH:mm:ssZ")}
lastmod: ${tp.date.now("YYYY-MM-DDTHH:mm:ssZ")}
draft: true
categories:
tags:
---



`;

// 4. 定义新文件的完整路径
const folderPath = `posts/${rawFileName}`;
const newFilePath = `${folderPath}/index.md`;

// 5. 确保文件夹存在 (Templater 没有直接创建文件夹的 API，但这通常不是问题)
//    我们将通过重命名和移动文件来实现
await tp.file.move(`${folderPath}/index`);

%>
```

Alt+N可以选择模板创建文章，每篇文章都是一个目录。
## linter 插件

保存时格式化文件
### yaml

插入YAML属性【打开】：增加以下字段，每行一个。  
```
	draft: true  
	title:  
	categories:  
	tags:  
```

YAML时间戳【打开】  
	- 创建日期键名`date`；创建日期数据源【YAML】  
	- 修改日期键名`updated`；修改日期数据源【系统】  
	- 格式：`YYYY-MM-DDTHH:mm:ssZ`  
YAML标题【打开】，键名就用title  

### 内容  

不同内容间换行，确保两个空格【打开】

## 转换脚本

将文章中的ob链接转为hugo支持的链接。除了图片支持以内容嵌入的形式链接，其它链接不能`!`开头。

更新.gitignore，忽略转换目录：
```
temp_hugo_content
/public/
/resources/
.hugo_build.lock
```

## 本地预览

```sh
go run converter\main.go  
hugo server --contentDir temp_hugo_content -D --disableFastRender -p 4000
```

## 示例

文章链接：[[posts/002-使用github托管静态网站/index]]  
普通附件（转换脚本）：[[main.zip]]  
嵌入图片：  
![[dianchi.jpg]]

