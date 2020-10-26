# Questionnaire

## 简介

## 前端

+ 创建项目:

```bash
taro init Questionnaire
```

+ 目录结构:

```bash
Questionnaire/
├── babel.config.js             # Babel配置
├── .eslintrc.js                # ESLint配置
├── config                      # 编译配置目录
│   ├── dev.js                  # 开发模式配置
│   ├── index.js                # 默认公共配置
│   └── prod.js                 # 生产模式配置
├── images                      # 图片目录
├── package.json                # Node.js manifest
├── dist                        # 打包目录
├── project.config.json         # 小程序项目配置
├── src                         # 源码目录
│   ├── app.config.js           # 全局配置
│   ├── app.css                 # 全局 CSS
│   ├── app.js                  # 入口组件
│   ├── index.html              # H5 入口 HTML
│   └── pages                   # 页面组件
│       └── index               # 首页
│           ├── index.config.js # 页面配置
│           ├── index.css       # 页面CSS
│           └── index.jsx       # 页面组件
└── yarn.lock                   # Yarn manifest
```

+ 创建页面:

```bash
taro create --name questionnaire
taro create --name answersheet
taro create --name analysis

taro create --name login
```

## ECharts

```bash
git clone https://github.com/ecomfe/echarts-for-weixin.git
```

```bash
cp -r echarts-for-weixin/ec-canvas Questionnaire/src/components
```

+ 图床: `https://imgchr.com/`
+ 图标: `https://88icon.com/`

***


+ Bug:
  + 在新建投票后，先创建问题，再更新标题，标题未更新
  + 选项模态框当选项过多时, 模态框无滚动条。