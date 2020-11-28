# 问卷调查系统

## 目录结构

```ini
Questionnaire                   # 项目名称
├── BE                          # 后端
│   ├── build                   # 构建目录
│   │   ├── docker              # 容器目录
│   │   └── sqls                # 数据库目录
│   ├── cmd                     # 入口目录
│   ├── configs                 # 配置目录
│   ├── go.mod                  # Go Modules 简略信息
│   ├── go.sum                  # Go Modules 详细信息
│   └── internal                # 内部目录
│       ├── common              # 公共目录
│       ├── controllers         # 控制器目录
│       ├── models              # 模型目录
│       └── routers             # 路由目录
├── FE                          # 前端
│   ├── babel.config.js         # Babel配置
│   ├── config                  # 编译配置目录
│   │   ├── dev.js              # 开发模式配置
│   │   ├── index.js            # 默认公共配置
│   │   └── prod.js             # 生产模式配置
│   ├── dist                    # 项目打包目录
│   ├── images                  # 图片目录
│   ├── package.json            # Node.js manifest
│   ├── project.config.json     # 项目配置
│   ├── src                     # 源码目录
│   │   ├── app.config.ts       # 全局配置
│   │   ├── app.scss            # 全局 SCSS
│   │   ├── app.ts              # 入口组件
│   │   ├── components          # 公共组件
│   │   ├── index.html          # 首页
│   │   └── pages               # 页面
│   └── yarn.lock               # Yarn manifest
└── README.md                   # 项目必读
```

## 功能

## 前端

### Mac

### Windows

### 运行项目

+ 开发环境:

```bash
export NODE_ENV="development"
yarn dev:weapp
```

+ 生成环境:

```bash
export NODE_ENV="production"
yarn dev:weapp
```

### ECharts

```bash
git clone https://github.com/ecomfe/echarts-for-weixin.git
```

## 后端

### 构建发布

+ 切换目录：

```bash
cd Questionnaire/BE
```

+ 构建发布：

```bash
./build.sh
```

### Linux

+ 将输出目录拷贝到服务器上：

```bash
scp -r Questionnaire/BE/output REMOTE:/opt/
```

+ 运行数据库：

```bash
cd /opt/output/build/docker/
docker-compose up -d
```

+ 创建数据库：

```bash
docker exec -it mysql bash
mysql -u root -p123456
```

```sql
-- 创建数据库
CREATE DATABASE questionnaire CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

+ 运行项目：

```bash
cd /opt/output/
./questionnaire
```

***
