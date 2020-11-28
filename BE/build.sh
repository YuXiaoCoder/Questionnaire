#!/bin/bash

SCRIPT=$(readlink -f $0)
CWD=$(dirname ${SCRIPT})

# 创建输入目录
mkdir -p ${CWD}/output

# 拷贝目录
cp -r ${CWD}/build ${CWD}/output/
cp -r ${CWD}/configs ${CWD}/output/

# 构建主程序
export GOOS="linux"
go build -o ${CWD}/output/questionnaire ${CWD}/cmd/main.go
