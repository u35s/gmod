# gcmd
gmod本身实现了一套消息模块gcmd,提供了消息的定义及序列化,反序列化

gcmd消息tcp协议格式如下
```
----------------------------
| cmd | param | len | data |
----------------------------
```
cmd,param各一个字节,其中cmd的前两位分别标识消息是否压缩和消息是否连续,len是data的长度,此定义理论上支持无限长度消息.

gcmd下的processer实现了对gcmd的编码解码,并且实现了gnet的Processer的接口,因此你可以利用gcmd和gnet轻松的创建cs,ss服务器通信框架(examples/gnet下有简单的列子)
