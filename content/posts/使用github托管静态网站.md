---
draft: true
title: 使用github托管静态网站
date: 2025-10-09T17:27:55+08:00
updated: 2025-10-09T22:38:54+08:00
---
远端创建repo，带个README，此时main分支被创建；  
web操作：Settings -> Actions -> General -> 在 **Workflow permissions** 部分，确保切换到 `Read and write permissions` 并保存。


本地操作：  
git add .  
git commit -m "init"  
git branch -M source  
git remote add origin https://github.com/meizigz/meizigz.github.io.git  
创建 .gitignore  
```text
abc
```
git push -u origin source # -u 参数，建立追踪关系。

## todo：

action的配置。
