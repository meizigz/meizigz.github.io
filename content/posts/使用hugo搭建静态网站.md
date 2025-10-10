---
categories:
tags:
draft: false
title: 使用hugo搭建静态网站
date: 2025-10-09T17:25:44+08:00
updated: 2025-10-10T11:46:48+08:00
---
## 环境准备

1. hugo： https://github.com/gohugoio/hugo/releases 。 选择对应平台的extended版本，以支持sass。
2. go： https://go.dev/dl/ 下载安装包后进行安装，终端命令行输入 `go version` 确认是否安装成功。
3. dart-sass： https://github.com/sass/dart-sass/releases 把sass转为css。

- hugo和dart-sass是zip包，解压后，把路径加入PATH环境变量。

## 创建项目

```sh
hugo new site my_site
cd my_site
```

## 配置主题

```sh
git init
git submodule add https://github.com/HEIGE-PCloud/DoIt.git themes/DoIt
```

修改`hugo.toml`如下（可参考主题repo的文档）：


```toml
baseURL = "http://example.org/"
# [en, zh-cn, fr, ...] 设置默认的语言
defaultContentLanguage = "zh-cn"
# 网站语言, 仅在这里 CN 大写
languageCode = "zh-CN"
# 是否包括中日韩文字
hasCJKLanguage = true
# 网站标题
title = "我的全新 Hugo 网站"

# 更改使用 Hugo 构建网站时使用的默认主题
theme = "DoIt"

[params]
  # DoIt 主题版本
  version = "0.2.X"

[menu]
  [[menu.main]]
    identifier = "posts"
    # 你可以在名称 (允许 HTML 格式) 之前添加其他信息, 例如图标
    pre = ""
    # 你可以在名称 (允许 HTML 格式) 之后添加其他信息, 例如图标
    post = ""
    name = "文章"
    url = "/posts/"
    # 当你将鼠标悬停在此菜单链接上时, 将显示的标题
    title = ""
    weight = 1
  [[menu.main]]
    identifier = "tags"
    pre = ""
    post = ""
    name = "标签"
    url = "/tags/"
    title = ""
    weight = 2
  [[menu.main]]
    identifier = "categories"
    pre = ""
    post = ""
    name = "分类"
    url = "/categories/"
    title = ""
    weight = 3

# Hugo 解析文档的配置
[markup]
  # 语法高亮设置 (https://gohugo.io/content-management/syntax-highlighting)
  [markup.highlight]
    # false 是必要的设置 (https://github.com/dillonzq/LoveIt/issues/158)
    noClasses = false
```

## 创建文章

```sh
hugo new posts/first_post.md
```
可以编辑一些内容。
## 本地预览

```sh
hugo server -D --disableFastRender # 构建（包含草稿）
```

