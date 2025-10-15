---
categories:
  - 技术
tags:
  - hugo
  - github
draft: false
title: 使用github托管静态网站
date: 2025-10-09T17:25:55+08:00
lastmod: 2025-10-15T12:13:47+08:00
---
## 1 web操作

web操作：创建repo，名字为`meizigz.github.io`。带个README，此时main分支被创建。  
web操作：Settings -> Actions -> General -> 在页面最下方 **Workflow permissions** 部分，确保切换到 `Read and write permissions` 并保存。  
web操作：Settings -> Pages -> Build and deployment -> Source设置为 **GitHub Actions**  
web操作：Settings -> Environments -> githug-page -> Selected branches and tags -> 添加**source**  

## 2 本地操作

```sh
git add .  
git commit -m "init"  
git branch -M source  
git remote add origin https://github.com/meizigz/meizigz.github.io.git  
# 创建 .gitignore  
git push -u origin source # -u 参数，建立追踪关系。
```

## 3 克隆库

```
git clone --recurse-submodules https://github.com/your-username/your-repo.git
```
通过参数指定，将主题一并clone下来。
## 4 常用git操作

```sh
git commit --allow-empty -m "chore: Trigger CI"  # 空提交，触发CI
git submodule update --init --recursive  # 单独克隆或更新主题
```

## 5 action脚本

```yaml
# .github/workflows/deploy.yml
name: Build and Deploy Hugo Site to Pages

on:
  push:
    branches:
      - source
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: true

env:
  # 建议锁定到你本地测试通过的具体版本
  HUGO_VERSION: 0.151.0 # 例如，或者你正在使用的版本
  GO_VERSION: 1.25.2 # 例如

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout source code
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Setup Pages
        id: pages
        uses: actions/configure-pages@v5

      - name: Convert Obsidian links for Hugo
        run: go run ./converter/main.go

      - name: Install Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: ${{ env.HUGO_VERSION }}
          extended: true

      - name: Setup Hugo cache
        # 注意：这里我们只定义了路径，实际的恢复/保存由 actions/cache 完成
        uses: actions/cache@v4
        with:
          path: /tmp/hugo_cache
          # key 包含 commit hash，确保每次 push 都有机会创建新缓存
          key: ${{ runner.os }}-hugocache-${{ github.sha }}
          # 如果当前 commit 没有缓存，则使用最近一次的缓存
          restore-keys: |
            ${{ runner.os }}-hugocache-

      - name: Build Hugo site
        run: |
          hugo --gc --minify \
            --contentDir temp_hugo_content \
            --baseURL "${{ steps.pages.outputs.base_url }}/" \
            --cacheDir /tmp/hugo_cache # 修复2：告诉 Hugo 使用缓存目录

      # 新增的调试步骤：打印 public 目录内容
      - name: List build output
        run: ls -R public

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```