---
draft: false
date: 2025-10-09T17:26:46+08:00
updated: 2025-10-10T11:46:32+08:00
title: 使用obsidian配合hugo编辑markdown文件
categories:
  - 技术
  - hugo
tags:
  - obsidian
  - hugo
---
仓库根目录`content`  
新建笔记默认文件夹`posts`  
编辑器：关闭拼写检查  
## linter 插件配置

保存时格式化文件

yaml：  
	插入YAML属性【打开】：增加以下字段，每行一个。  
		`draft: true`  
	YAML时间戳【打开】  
	- 创建日期键名`date`；创建日期数据源【YAML】  
	- 修改日期键名`updated`；修改日期数据源【系统】  
	- 格式：`YYYY-MM-DDTHH:mm:ssZ`  
	YAML标题【打开】，键名就用title  
内容：  
	不同内容间换行，确保两个空格【打开】


## todo 

- 脚本转换链接

go run converter\main.go  
hugo server --contentDir temp_hugo_content -D --disableFastRender -p 4000

[[使用github托管静态网站]]