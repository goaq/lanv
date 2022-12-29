
# lan video service

使用go实现的提供局域网视频播放的服务
A service implemented by Go to provide LAN video playback
![1236173455430](https://user-images.githubusercontent.com/45125070/209532915-653cc16e-7427-49dc-b7c2-d47bf1d525a4.jpg)


## 编译和运行

## Compile and run

```go build . ```
```go run .```
页面使用原生js，不需要编译
The page uses native js and does not need to be compiled

https://github.com/goaq/lanv/releases/tag/release

## 配置文件

## Configuration file

直接使用```lanv```命令，将加载lanv.ini，```-ini```选项会加载指定的ini配置文件
Using the ```lanv``` command directly, the lanv.ini will be loaded, and the ```-ini``` option will load the specified ini configuration file

```
addr=0.0.0.0
port=9000

[paths]
Music=E:\musicdir
Movies1=g:\moviesdir1
Movies2=g:\moviesdir2
Root=H:\videos

[types]
mp4=1
mkv=1
flv=1
ts=0
avi=0
mp3=1

```
```paths```中配置名称和路径，```Root```作为名称将在根路径下展开目录
The name and path are configured in ```paths```, and ```Root``` as the name will expand the directory under the root path

