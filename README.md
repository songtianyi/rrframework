#### 代码规范
1. 使用fmt.Errorf来创建一个error
2. 错误信息首字母大写, 当用户无换行的时候便于找到开头
3. 不需要导出的方法首字母小写
4. Create或者New 命名要保持一致

#### 依赖包
* go get gopkg.in/redis.v5  
 		 commit: b6bfe529a846fbb9a58c832ce71c61b6fde12c15
* go get github.com/golang/protobuf/proto  
         commit: 8ee79997227bf9b34611aee7946ae64735e6fd93
* go get github.com/samuel/go-zookeeper/zk
