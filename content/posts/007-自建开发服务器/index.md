---
title: "自建开发服务器"
date: 2025-12-05T14:53:29+08:00
lastmod: 2025-12-05T16:55:02+08:00
draft: true
categories:
tags:
---
购买服务器后安装系统，记录root密码。

## 流量转发

### 服务端

https://github.com/Alvin9999/Sing-Box-Plus 一键搭建Sing-Box多协议脚本，无需域名，开箱即用 18 个节点

```
mkdir airport
cd airport
wget -O sing-box-plus.sh https://raw.githubusercontent.com/Alvin9999/Sing-Box-Plus/main/sing-box-plus.sh
chmod +x sing-box-plus.sh
bash sing-box-plus.sh
# 选择1 安装部署。得到18个节点。
```

### 客户端

到 https://github.com/2dust/v2rayN/releases 下载客户端。  
启动后选择【配置项】--> 【从剪贴板导入分享链接】，然后就可以把生成的18个节点导入了。  
选中所有配置 --> 右键【测试真连接延迟】，可以看到配置是否连通。  
也可以测试哪种方法速度最快。【我这儿是TUIC最快】  
代理方式选择【自动配置系统代理】 V4-白名单

https://github.com/Loyalsoldier/v2ray-rules-dat?tab=readme-ov-file 可以去找最新的geoip和geosite文件，放到bin目录里面。【路由就用软件自带的即可】

【设置】--> 【以管理员身份重启】  
【设置】-->【v2rayN设置】 --> 开机启动【On】/ 自动更新Geo文件 【720小时】









