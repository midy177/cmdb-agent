#/bin/bash
if [[ $EUID -ne 0 ]]; then
   echo -e "\e[91m必须以root用户运行😒"
   exit 1
fi

if command -v cmdb-agent > /dev/null 2>&1;then
echo -e '\033[32m已经安装cmdb-agent\033[0m'
cmdb-agent service status
else
echo -e '\033[31m没有安装ycm-deploy-client,开始安装...\033[0m'
wget https://image.yeastar.com/tools/cmdb-agent -O /usr/local/bin/cmdb-agent
chmod +x /usr/local/bin/cmdb-agent
echo -e '\033[32mcmdb-agent安装成功\033[0m'
fi

if command -v trz > /dev/null 2>&1;then
echo -e '\033[32m已经安装trz\033[0m'
else
echo -e '\033[31m没有安装trz,开始安装...\033[0m'
wget https://image.yeastar.com/tools/trz -O /usr/local/bin/trz
chmod +x /usr/local/bin/trz
echo -e '\033[32mtrz安装成功\033[0m'
fi

if command -v tsz > /dev/null 2>&1;then
echo -e '\033[32m已经安装tsz\033[0m'
else
echo -e '\033[31m没有安装tsz,开始安装...\033[0m'
wget https://image.yeastar.com/tools/tsz -O /usr/local/bin/tsz
chmod +x /usr/local/bin/tsz
echo -e '\033[32mtsz安装成功\033[0m'
fi

if command -v acp > /dev/null 2>&1;then
echo -e '\033[32m已经安装acp\033[0m'
acp --setup
else
echo -e '\033[31m没有安装acp,开始安装...\033[0m'
wget https://image.yeastar.com/tools/acp -O /usr/local/bin/acp
chmod +x /usr/local/bin/acp
echo -e '\033[32macp安装成功\033[0m'
acp --setup
fi
