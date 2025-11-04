---
title: 通过阿里云oss实现文件中转
date: 2025-10-28T17:04:28+08:00
lastmod: 2025-11-04T10:38:22+08:00
draft: false
categories:
tags:
  - 阿里云
  - OSS
---
## 创建bucket

阿里云oss管理： https://oss.console.aliyun.com/bucket  
创建bucket：  
	名字：  bucket的名字  
	地域：中国香港  
	endpoint：  oss-cn-hongkong.aliyuncs.com  
	存储冗余：本地冗余  
进入bucket：  
	将 Bucket 的读写权限设置为 **“公共读”**（Public Read）  
		- 权限控制->阻止公共访问->改为【关闭】  
		- 权限控制->读写权限->Bucket ACL->改为【公共读】  
		- 权限控制->访问控制RAM->跳转到RAM用户管理【 https://ram.console.aliyun.com/overview 】-> 身份管理 -> 用户 -> 创建用户  
			- 输入用户名称  
			- 取消【控制台访问】  
			- 勾选【使用永久AccessKey访问】  
				- 复制AccessKey：  
				- 复制AccessKey Secret：  
		- 权限控制->Bucket授权策略->新增授权  
			- 整个Bucket  
			- 子账号：选择刚创建的子账号  
			- 授权操作：读写  
			- 访问方式：https  
	24小时删除：  
		- 数据管理->生命周期->创建规则  
			- 启动  
			- 配置到整个Bucket  
			- 生命周期管理规则：最后一次【修改时间】【1】天后，【数据删除】


## SDK操作oss

### 引入库

```python
import alibabacloud_oss_v2 as oss
```
文档： https://gosspublic.alicdn.com/sdk-doc/alibabacloud-oss-python-sdk-v2/latest/alibabacloud_oss_v2.html

### 创建client对象

```python
credentials_provider = oss.credentials.StaticCredentialsProvider(
            self.config['access_key_id'], 
            self.config['access_key_secret']
        ) # 也可以使用环境变量，此处选择了用静态secret
cfg = oss.config.load_default()
cfg.credentials_provider = credentials_provider
cfg.region = self.config['region'] # bucket地域对应的region
cfg.endpoint = self.config['endpoint'] # bucket的endpoint
return oss.Client(cfg)
```

香港区的region：cn-hongkong  
其它区域的： https://help.aliyun.com/zh/oss/user-guide/regions-and-endpoints  
region不对会导致访问失败。

### 小文件上传

```python
self.client.put_object_from_file(
            oss.PutObjectRequest(
                bucket=bucket_name,
                key=object_key,
                progress_callback=progress_callback
            ),
            file_path
        )
```

- bucket名称就是目的bucket
- object_key就是远端的文件名。可以用原来的 `object_key = os.path.basename(file_path)` ，也可以某个hash值都可以。
- progress_callback 可以为None，不为None时就是进度回调
- 虽然文档说是支持小于5G的文件，但实际上几十M就很有可能超时，所以大点的文件还是用下面的分块上传更好。

### 大文件分块上传

```python
        """上传大文件（分片上传）"""
        file_size = os.path.getsize(file_path)
        
        # 初始化分片上传
        init_result = self.client.initiate_multipart_upload(oss.InitiateMultipartUploadRequest(
            bucket=bucket_name,
            key=object_key
        ))
        upload_id = init_result.upload_id
        
        # 设置分片大小
        part_size = 1 * 1024 * 1024  # 1MB per part
        
        # 获取文件总大小
        data_size = file_size
        part_number = 1
        upload_parts = []
        
        # 打开文件并分片上传
        with open(file_path, 'rb') as f:
            for start in range(0, data_size, part_size):
                # 计算当前分片大小
                n = part_size
                if start + n > data_size:
                    n = data_size - start
                
                # 创建分片读取器
                reader = oss.io_utils.SectionReader(oss.io_utils.ReadAtReader(f), start, n)
                
                # 上传当前分片
                up_result = self.client.upload_part(oss.UploadPartRequest(
                    bucket=bucket_name,
                    key=object_key,
                    upload_id=upload_id,
                    part_number=part_number,
                    body=reader
                ))
                
                # 保存分片信息
                upload_parts.append(oss.UploadPart(part_number=part_number, etag=up_result.etag))
                
                # 如果提供了进度回调，计算并调用
                if progress_callback:
                    uploaded_size = min(start + n, data_size)
                    progress_callback(uploaded_size, data_size)
                
                part_number += 1
        
        # 按分片编号排序
        parts = sorted(upload_parts, key=lambda p: p.part_number)
        
        # 完成分片上传
        self.client.complete_multipart_upload(oss.CompleteMultipartUploadRequest(
            bucket=bucket_name,
            key=object_key,
            upload_id=upload_id,
            complete_multipart_upload=oss.CompleteMultipartUpload(parts=parts)
        ))
```

### 获取分享url
因为我们的bucket是公共读的，所以可以直接拼接得到url。

```python
"""生成公共URL"""
endpoint = self.config['endpoint'].replace('https://', '').replace('http://', '')
bucket_name = self.config['bucket']
return f"https://{bucket_name}.{endpoint}/{object_key}"
```
类似于： https://bucket名称.oss-cn-hongkong.aliyuncs.com/object_key









