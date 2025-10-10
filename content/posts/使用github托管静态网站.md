---
draft: false
title: 使用github托管静态网站
date: 2025-10-09T17:27:55+08:00
updated: 2025-10-10T10:14:19+08:00
---
远端创建repo，带个README，此时main分支被创建；  
web操作：Settings -> Actions -> General -> 在 **Workflow permissions** 部分，确保切换到 `Read and write permissions` 并保存。  
web操作：Settings -> Environments -> githug-page -> Deployment branches and tags 【No restriction】

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


git commit --allow-empty -m "chore: Trigger CI after service recovery"

## todo：

action的配置。


Warning:  
Cache folder path is retrieved but doesn't exist on disk: /home/runner/go/pkg/mod  
46  
Primary key was not generated. Please check the log messages above for more errors or information  
47  
##[debug]Node Action run completed with exit code 0


