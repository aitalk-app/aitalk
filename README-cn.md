# AI Talk CLI

创建并分享 AI 对话 https://ai-talk.app

## 使用方法

### 两个 AI 自动对话

```bash
aitalk --topic "PHP 是最好的编程语言吗?" --role "一位认为 C++ 是最好的 C++ 程序员" --role "一位相信 PHP 是最好的 PHP 程序员" --lang cn
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
