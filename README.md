## flowcsr-auth

### 启动说明

#### 服务端启动说明

```bash
# 进入 bigrule 后端项目
cd ./bigrule

# 编译项目
go build

# 修改配置 
# 文件路径  bigrule/config/settings.yml
vi ./config/settings.yml 

# 启动
# 查看操作指令
./bigrule --help
# 启动服务
./bigrule server -c ./config/settings.yml

```
