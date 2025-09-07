# Sensitive-lexicon (中文敏感词库)

![Commit Activity](https://img.shields.io/github/commit-activity/y/Konsheng/Sensitive-lexicon)
![License: MIT](https://img.shields.io/github/license/Konsheng/Sensitive-lexicon)
![GitHub stars](https://img.shields.io/github/stars/Konsheng/Sensitive-lexicon)

> **一个持续更新的中文敏感词库，帮助开发者和内容审核者快速识别并过滤不当文本。**

## 目录

* [简介](#简介)
* [功能特点](#功能特点)
* [目录结构](#目录结构)
* [快速开始](#快速开始)
  * [集成到项目](#集成到项目)
  * [贡献词汇](#贡献词汇)
* [注意事项](#注意事项)
* [开源许可](#开源许可)
* [项目支持与致谢](#项目支持与致谢)

## 简介

Sensitive‑lexicon 提供了一份广泛覆盖政治、色情、暴力等敏感领域的词汇列表，方便快速嵌入任何文本审核流程，并通过社区协作保持长期更新。

## 功能特点

* **广泛覆盖**：涵盖数万条词汇，覆盖主流敏感领域。
* **持续更新**：根据社会语境变化定期更新，保持时效性与准确性。
* **易于集成**：纯文本格式，可在任意语言/框架中直接引用。
* **社区驱动**：欢迎 Issue / PR，携手打造更完整的词库。

## 目录结构
```
Sensitive-lexicon/
├── ThirdPartyCompatibleFormats/        # 用于第三方格式
├── Organized/                          # 已经进行整理的词库
├── Vocabulary/                         # 词汇库
├── LICENSE                             # 许可证
└── README.md                           # 项目说明
```

## 快速开始

### 集成到项目

1. 克隆或下载本仓库。
2. 在您的代码中读取 `词库中的 .txt 文件`（或您需要的分支文件）。
3. 根据业务场景，选择合适的匹配算法（如 DFA、Trie、正则表达式等）进行过滤。

```bash
# 示例：使用 Git 克隆
git clone https://github.com/Konsheng/Sensitive-lexicon.git
```

### 贡献词汇

* **Pull Request**：在 `Vocabulary/` 词汇库目录新增或修改词条，并提交 PR。
* **Issue**：如果不确定具体实现，欢迎通过 Issue 提出建议或讨论。

> **提示**：PR 请附上来源或用例，便于维护者审核。

## 注意事项

* 使用时请遵守当地法律法规及平台政策。
* 敏感词定义受文化/地域/语境影响，实际应用中请结合业务需求自行评估与调整。

## 开源许可

本项目采用 **MIT License**，在保留版权与许可声明的前提下，可自由使用、修改与分发。

## 项目支持与致谢

* **中国数字时代** ([https://chinadigitaltimes.net](https://chinadigitaltimes.net))
* **中国农业科学院信息化办公室**

感谢所有贡献者的关注与支持！

## Star History
<a href="https://star-history.com/#konsheng/Sensitive-lexicon&Date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=konsheng/Sensitive-lexicon&type=Date&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=konsheng/Sensitive-lexicon&type=Date" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=konsheng/Sensitive-lexicon&type=Date" />
  </picture>
</a>
