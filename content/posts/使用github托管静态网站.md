---
categories:
tags:
draft: false
title: 使用github托管静态网站
date: 2025-10-09T17:27:55+08:00
updated: 2025-10-10T11:46:46+08:00
---
远端创建repo，带个README，此时main分支被创建；  
web操作：Settings -> Actions -> General -> 在 **Workflow permissions** 部分，确保切换到 `Read and write permissions` 并保存。  
web操作：Settings -> Environments -> githug-page -> Selected branches and tags【source】

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

git submodule update --init --recursive

git clone --recurse-submodules https://github.com/your-username/your-repo.git

1. **进入仓库设置**:
    
    - 在你的 GitHub 仓库页面，点击 Settings。
        
2. **进入 Pages 设置**:
    
    - 在左侧菜单中，点击 Pages。
        
3. **修改部署源 (Source)**:
    
    - 在 Build and deployment 部分，你会看到 Source 当前被设置为 **Deploy from a branch**。
        
    - 点击下拉菜单，将其更改为 **GitHub Actions**。

## todo：

action的配置。


Warning:  
Cache folder path is retrieved but doesn't exist on disk: /home/runner/go/pkg/mod  
46  
Primary key was not generated. Please check the log messages above for more errors or information  
47  
##[debug]Node Action run completed with exit code 0


