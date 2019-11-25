## exchange_data gin+vue 快速部署获取各大数字货币交易所行情数据
![项目展示](https://github.com/GitTsewell/exchange_data/blob/master/depth.gif)
## 项目说明
```
    1.go在服务器端挂起一个WS客户端,利用goroutine建立各大交易所(okex,火币,币安,bitmex)的ws连接,获取实时行情数据
储存在redis中,量化机器人直接在redis中读取.提供一个行情管理后台,管理币种以及ws客户端的重启,支持分布式
    2.exchange_api 是go写的ws客户端和后台api接口, exchange_vue 是后台页面
```
## 使用说明
### GO
```
    golang api server 基于go.mod 如果golang版本低于1.11 请自行升级golang版本
    exchange_api/config/redis.go db配置
    cd exchange_api/cmd && go build ws_client.go  编译ws客户端 使用进程管理工具启动(比如supervisor)
    cd exchange_api/cmd && go build web_server.go  编译admin_api客户端
```
### vue
```
    exchange_vue/src/http.js 接口host配置
    cd exchange_vue
    yarn
    yarn run dev
    访问 localhost:8002
```
请愉快食用~


