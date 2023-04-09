# AI Talk CLI

创建并分享 AI 对话 https://ai-talk.app

## 安装
可以从 GitHub 的[发布页面](https://github.com/aitalk-app/aitalk/releases)下载`aitalk`。或者，可以使用 Go 在系统上安装`aitalk`：

```bash
go install github.com/aitalk-app/aitalk@latest
```
将 OpenAI API 密钥设置为环境变量

```bash
export OPENAI_API_KEY=<your_api_key>

# 如果需要，也可以设置一个自定义的API Host
export OPENAI_API_HOST=<your_custom_api_host_url>
```

## 使用方法

### 两个 AI 自动对话

```bash
aitalk --topic "甜豆腐脑还是咸豆腐脑" --role "喜欢甜豆腐脑" --role "喜欢咸豆腐脑" --lang cn
```

### 和AI 交互式创建对话

```bash
aitalk --topic "AI 会取代人类吗?"
```

## 管理对话 (可选)

如果想在 https://ai-talk.app 上管理对话，需要进行身份验证。运行以下命令并在浏览器中打开显示的 URL：

```
aitalk auth
```

如果跳过此步骤，可以稍后运行上述命令，并且所有以前创建的对话将自动绑定到登录的用户。
