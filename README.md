# mass
mass是一个通过sshpass批量执行命令的工具

## 使用方法
mass   
-port       服务器ssh端口
-ip       要执行命令的服务器ip列表文件    
-cmd      写有执行命令的文件  
-password        密码文件 如果运行mass的服务器和执行命令的服务器没有做ssh免密，需要将密码写入文件。  
-c        并发数量 默认为10  
-head     从$head服务器开始执行  
-tail     执行到$tail服务器  
