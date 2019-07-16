#!/usr/bin/env bash

EXEPATH="$(pwd)"

check_cmd_or_exit() {
  if ! command -v $1 > /dev/null;then
    red_text "未安装 $1"
    if [ $1 == "go" ];then
      echo "下载安装: https://studygolang.com/dl"
    fi
    if [ $1 == "protoc" ];then
      echo "下载安装: https://github.com/protocolbuffers/protobuf/releases"
    fi
    exit 1
  fi
}

red_text() {
  echo -e "\e[31m$*\e[0m"
}

green_text() {
  echo -e "\e[32m$*\e[0m"
}

act_compiler_proto() {
  for file in *; do
    if [ -d $file ]; then
    # echo "is dir"
      # if [ $file == *libs ];then
      #   continue
      # fi
      cd $file
      act_compiler_proto .
      cd ../
    else if [ $file == *.proto ];then
      pwd_path=`pwd -W`
      echo -e $pwd_path/$file "...\c"
      protoc -I=$GOPATH/src -I=. -I=$EXEPATH --proto_path=$EXEPATH --go_out=$EXEPATH $pwd_path/$file # TODO
      green_text " 完成"
      fi
    fi
  done
}

check_cmd_or_exit go
check_cmd_or_exit protoc
echo $EXEPATH

act_compiler_proto $EXEPATH

