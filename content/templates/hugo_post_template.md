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