# mass
mass是一个通过ssh批量执行命令的工具

# 使用方法
var fIP = flag.String("ip", "", "ip list file")
var fCommand = flag.String("cmd", "", "command file")
var fPassword = flag.String("p", "", "password file")
var fConcurrency = flag.Int("c", 10, "concurrency number")

var fHead = flag.Int("head", -1, "head")
var fTail = flag.Int("tail", -1, "tail")

mass 
-ip   要执行命令的服务器ip列表文件  
-cmd  写有执行命令的文件
-p    密码文件 如果运行mass的服务器和执行命令的服务器没有做ssh免密，需要将密码写入文件。
-c    并发数量 默认为10
-head 从$head服务器开始执行
-tail 执行到$tail服务器
