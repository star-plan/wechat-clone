# wechat-clone

macOS 微信分身管理 CLI 工具，支持创建、管理、启动多个微信分身。

## 功能

- 创建指定数量的微信分身
- 列出已有分身
- 启动指定或所有分身
- 修复分身签名（微信更新后可能需要）
- 删除分身
- 环境检查（自动检测依赖和权限）

## 前提条件

- macOS 系统
- 已安装微信（`/Applications/WeChat.app`）
- 已安装 Xcode 或 Xcode Command Line Tools（提供 `codesign` 和 `PlistBuddy`）
- sudo 权限（复制和修改 App 需要）

运行 `wechat-clone doctor` 可以自动检查以上条件。

## 安装

### Homebrew（推荐）

```bash
brew install deali/tap/wechat-clone
```

### 从 GitHub Release 下载

前往 [Releases](https://github.com/deali/wechat-clone/releases) 页面下载对应架构的二进制文件。

```bash
# Apple Silicon (M1/M2/M3/M4)
curl -L -o wechat-clone https://github.com/deali/wechat-clone/releases/latest/download/wechat-clone-darwin-arm64
chmod +x wechat-clone
sudo mv wechat-clone /usr/local/bin/

# Intel Mac
curl -L -o wechat-clone https://github.com/deali/wechat-clone/releases/latest/download/wechat-clone-darwin-amd64
chmod +x wechat-clone
sudo mv wechat-clone /usr/local/bin/
```

### 从源码编译

```bash
git clone https://github.com/deali/wechat-clone.git
cd wechat-clone
go build -o wechat-clone .
sudo mv wechat-clone /usr/local/bin/
```

## 使用

### 交互式模式（推荐）

直接运行 `wechat-clone` 进入交互式 TUI 界面：

```bash
wechat-clone
```

使用方向键选择功能，Enter 确认，Esc 返回。

### 命令行模式

```bash
# 检查环境
wechat-clone doctor

# 创建 3 个分身
wechat-clone create 3

# 强制覆盖已存在的分身
wechat-clone create 3 --force

# 列出所有分身
wechat-clone list

# 查看分身启动指引（显示路径，可选在 Finder 中定位）
wechat-clone open

# 查看指定编号分身的启动指引
wechat-clone open 2

# 修复所有分身签名
wechat-clone repair

# 修复指定分身
wechat-clone repair 1

# 删除指定分身（会要求确认）
wechat-clone remove 2

# 强制删除所有分身
wechat-clone remove all --force
```

## 分身命名规则

| 项目 | 值 |
|------|-----|
| 分身路径 | `/Applications/WeChat Clone 1.app`、`WeChat Clone 2.app` ... |
| Bundle ID | `com.tencent.xinWeChat.clone1`、`.clone2` ... |

## 常见问题

### Q: 创建分身时提示"未找到微信应用"

确认微信已安装在 `/Applications/WeChat.app`。如果微信在其他位置，使用 `--source` 参数指定：

```bash
wechat-clone create 1 --source "/path/to/WeChat.app"
```

### Q: 提示需要 sudo 权限

创建和修复分身需要修改 `/Applications` 目录下的文件，系统会自动提示输入密码。

### Q: 微信更新后分身打不开了

运行修复命令：

```bash
wechat-clone repair
```

### Q: 分身之间数据会互相影响吗？

不会。每个分身有独立的 Bundle ID，微信会为每个分身创建独立的数据目录。

### Q: 怎么打开微信分身？

创建分身后，在 Finder 或启动台中找到 `WeChat Clone 1.app`，双击即可启动。也可以运行 `wechat-clone open` 查看路径并快速在 Finder 中定位。建议把常用的分身拖到程序坞固定。

### Q: 可以同时登录多个微信账号吗？

可以。每个分身是独立的 App 实例，可以同时登录不同的微信账号。

### Q: PlistBuddy 不可用

需要安装 Xcode 或 Xcode Command Line Tools：

```bash
xcode-select --install
```

## 风险说明

- 本工具仅通过复制 App、修改 Bundle ID、重新签名的方式实现分身，不修改原始微信
- 不注入 dylib、不实现防撤回、不绕过任何安全机制
- 微信官方可能在后续版本中限制分身登录，请自行评估风险
- 本工具不收集任何用户信息，不上传任何数据

## 卸载

### 删除分身

```bash
wechat-clone remove all --force
```

### 删除工具本身

```bash
# 如果通过 Homebrew 安装
brew uninstall wechat-clone

# 如果手动安装
sudo rm /usr/local/bin/wechat-clone
```

### 清理残留数据（可选）

分身的聊天数据存储在 `~/Library/Containers/` 下，删除分身后如需清理：

```bash
# 查看相关数据目录
ls ~/Library/Containers/ | grep xinWeChat
```

## 开发

```bash
# 克隆项目
git clone https://github.com/deali/wechat-clone.git
cd wechat-clone

# 安装依赖
go mod tidy

# 编译
go build -o wechat-clone .

# 运行测试
go test ./...
```

## License

[Apache License 2.0](LICENSE)
