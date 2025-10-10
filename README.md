# 海巴迪

## 带主题克隆

```sh
git clone --recurse-submodules https://github.com/meizigz/meizigz.github.io.git
```

## 分步克隆

```sh
git clone https://github.com/meizigz/meizigz.github.io.git
git submodule update --init --recursive
```

## 切换到文章分支

```sh
git switch source
```

## 空提交，触发action

```sh
git commit --allow-empty -m "chore: Trigger CI"
```

## 本地预览

```
# 切换到source分支
preview.bat
```