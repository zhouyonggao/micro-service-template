# 项目说明

## 一、目录说名
```
├── Dockerfile // 构建 docker 文件
├── Makefile // makefile 工具，定义了一些命令，如 make api等
├── README.md // 说明文档
├── api // 微服务grpc使用的proto文件以及根据它们所生成的文件
│	 ├── myerr // 公共的可响应的错误信息
│	 │	 └── my_err.proto
│	 └── order 
│	     └── v1 // 项目版本号
│	         └── server.proto // 具体的
├── cmd //启动命令相关
│	 └── service // 服务启动的目录
│	     ├── main.go // 启动文件
│	     ├── wire.go // wire维护依赖注入的定义文件
│	     └── wire_gen.go // 通过 wire 生成的文件
├── configs // 配置文件目录，通常开发环境使用，线上等环境会指定到其它目录
│	 └── config.yaml // 具体的配置文件
├── go.mod // go module 的依赖定义文件
├── go.sum // go module 版本锁定的版本信息
├── internal // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
│	 ├── biz // 业务逻辑的组装层，类似 DDD 的 domain 层，在 api、biz、data三层中，它是最底层，它依赖的一些外部交互，应该定义成 interface，由 data 层去实现注入
│	 │	 ├── acl // 调用其它 rpc 或外部服务时的防腐层目录，注意不要依赖于具体的实现类，应该定义 interface，由 data 层去实现
│	 │	 ├── provider_set.go // 声明依赖的实例工厂方法
│	 │	 ├── do // 领域的实体
│	 │	 │	 └── order.go // 示例：具体的实体文件，类名以 Do 结尾
│	 │	 ├── event // 领域事件相关，同样，不要依赖于具体的实现类，应该定义 interface，由 data 层去实现
│	 │	 │	 └── send.go // 具体的定义，或包装，接口名以 Event 结尾
│	 │	 ├── repository // 持久化等相关的定义，同样，不要依赖于具体的实现类，应该定义 interface 抽象，由 data 层去实现
│	 │	 │	 ├── order.go // 示例：具体的定义，或包装，接口以 Repo 结尾
│	 │	 │	 └── transactio.go // 封装了对事务的支持
│	 │	 ├── valueobject // 领域中的值对象
│	 │	 │	     ├── order_status.go // 示例：类名以 Obj 结尾
│	 │	 ├── order.go // 具体的领域服务，类名以 Biz 结尾
│	 ├── conf // 维护配置文件的 proto 和 生成的 go 文件，框架会自动将配置解析生成到 go 对应的对象中
│	 │	 ├── conf.pb.go // 根据 proto 文件生成的 go
│	 │	 └── conf.proto // 定义的 proto 文件
│	 ├── data // 数据层实现层，通过依赖导致后，它是在最上层，即其它层都不依赖于它，它依赖于所有层，对去各层的抽象进行实现
│	 │	 ├── convert // po 与 do 的相互转换目录
│	 │	 │	 └── order.go // 示例：类名以为 Conv 结尾
│	 │	 ├── ent // entgo 框架生成的目录
│	 │	 │	 ├── generate.go // go的generate 命令，它会根据 schema 文件侠下的定义文件生成 go 代码 
│	 │	 │	 ├── schema // 存放表结构的定义文件
│	 │	 │	 │	 ├── order.go // 示例文件：订单表结构定义
│	 │	 ├── eventimpl // 事件相关的实现
│	 │	 │	 ├── consumer.go // consumer Service 的实现，类名以 EventImpl 结尾
│	 │	 │	 ├── provider_set.go // 声明依赖的实例工厂方法
│	 │	 │	 └── send.go // 发送消息的实现
│	 │	 ├── repositoryimpl // 持久化相关的实现
│	 │	 │	 ├── common_test.go // 持久化相关单元测试共用方法
│	 │	 │	 ├── order.go // 示例代码：持久化 order 的实现，类名以 RepoImpl 结尾
│	 │	 │	 └── provider_set.go // 声明依赖的实例工厂方法
│	 │	 ├── provider_set.go // 声明依赖的实例工厂方法
│	 │	 ├── data.go // data 构建工厂方法
│	 │	 └── transaction_impl.go // 事务支持的实现
│	 ├── pkg // 内部共用的代码，可以理解为 helper 层 或 工具层
│	 │	 ├── helper_log.go // 封装了对日志的包装
│	 │	 └── metadatamanager
│	 │	     └── metadata_manager.go // 对grpc metadata 数据的管理支持，如获取metadata中的数据
│	 ├── server // grpc和consumer服务的创建和配置
│	 │	 ├── consumer_server.go // 生成并配置consumer 服务
│	 │	 ├── grpc_server.go // 生成并配置 grpc 服务
│	 │	 └── provider_set.go // 声明依赖的实例工厂方法
│	 └── service // 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑，它的复用性较差，是由于它更多的表现外部的数据需求上
│	     ├── order.go // 示例代码：订单的 service，类名以 Service 结尾
│	     ├── order_consumer.go // 示例代码：订单消费的 service
│	     └── provider_set.go // 声明依赖的实例工厂方法
├── sync_table_entgo.sh // 同步数据库某表到 ent/schema 目录下的脚本
└── third_party // api 依赖的第三方proto
```
## 二、分层架构说明：
![分层架构](https://i.328888.xyz/2023/02/16/mQ7Yx.png)

采用 kratos 推荐的分层模型，与DDD四层依赖倒置模型类似，只是换了名称，和少了用户接口层。
* **Service 层：** 与 DDD 的 Application 类似，由它来协调多个领域服务共同实现某功能，并转换为其功能外部所关心的数据结构
，所有从外部进入应用的入口全在此层，如GRPC的接入、消费消息的接入

* **Biz 层：** 领域层，由领域服务来协调多个实体的行为来实现某功能，此层是最底层，即不能直接用 data 和 service 层，对外部的交互需要定义成抽象接口，由 data 层来实现，
如持久化数据、发送 MQ 消息、访问其它领域服务等

* **Data 层：** 可理解为 DDD 分层模型的基础层，在此层实现 Biz 中定义的抽象接口

* **API 层：** 它是一个单独的 git 仓库，通过 git submodule 引入，存放了的grpc和error的proto的相关文件

## 三、生成 ent 模型代码的步骤

1. 使用命令 `make ent t=[表名]` 如提供了第二个"表名"参数，则会从数据库中同步此表结构到 data/ent/schema 目录中
2. 所有schema 中的整型(Uint6,Int32...) 统一使用 Int，避免项目到处理充满了类型转换，若使用sync_talbe_entgo工具时已经自动替换为 Int 了
3. 若领域模型中也有此模型中相同的字段，则需要保持相同的类型，避免开发时来回转换类型

## 四、如何定义枚举类型为值对象
枚举类型的属性通常有一些行为，如某值对应的中文名称是什么，或某些值的操作逻辑等(如某值是否符合某行为的判定、状态机模式)
* 命令规则：领域对象文件名 + 属性名 结尾；对象名为：实体对象名 + 属性名 + Obj 的驼峰结构，如 orderStatus.go，类名为：OrderStatusObj
* 值对象里面必须实现GetName、GetValue方法
* 值对象在 go 语言里为新定义类型，基础类型是 int

## 五、编写单元测试
有效的单元测试可以检查我们的代码能否按照预期进行，代码逻辑是否有问题，以此可以提升代码质量。项目中已集成相关库：
* 单元测试快速编写工具：https://github.com/smartystreets/goconvey
* 方法/函数 Mock 工具：https://github.com/agiledragon/gomonkey
* 数据库 mock 工具：https://github.com/DATA-DOG/go-sqlmock
* Redis mock 工具：https://github.com/alicebob/miniredis
命名方式见 开发规范中：1.5 单元测试
原则上单元测试覆盖率越高越好，核心业务应保证100%覆盖
### 如何测试数据库？
我们知道，数据库的实现在data/repositoryimpl目录中，因此在目录中加入了commom_test.go，用于封装测试 data 相关的公共函数，如 GetMockData 方法，用于试获取 mock 好的 data，
示例：
```go
// mock 连接 data
data, mock := GetMockData(t)

// 生成impl 的实例，并注入 data
impl := NewOrderRepoImpl(data)
c.Convey("正确更新时间", t, func() {
    //mock sql 语句，当 sql 语句、参数完全匹配时，会返回 正常的结果，如果没匹配上，则返回 error
    now := time.Now()
    expectSql := "UPDATE `order_copy_tmp` SET `update_time` = \\? WHERE `order_copy_tmp`.`id` = \\?"
    mock.ExpectExec(expectSql).WithArgs(now, 1).WillReturnResult(sqlmock.NewResult(0, 1))
    //调用被测试方法
    err := impl.UpdateTimeNow(context.Background(), 1, now)
    //断言返回的 error 是 nil
    c.So(err, c.ShouldBeNil)
})
```
建议有更新数据库的操作，必须100%覆盖率，尽可能保证操作数据的正确性，以免造成数据错误

### 如何测试 Redis 
同测试数据库一样，调用 commom_test.go中的 mock 方法。比数据库更好的点是可以存入模拟数据，示例：
```go
  // mock redis
rds, mockRds := GetMockRds(t)
d := &data.Data{
    Rdb: rds,
}

//模拟redis中的数据
err := mockRds.Set("test_1", "{\"id\":1,\"ProductName\":\"测试数据\"}")
if err != nil {
    t.Fatalf("mock redis 数据错误：%s", err)
}
impl := NewOrderRepoImpl(d)

c.Convey("返回正常获取到信息", t, func() {
    do, err := impl.FindByIDByRedis(context.Background(), 1)
    c.So(do.Id, c.ShouldEqual, 1)
    c.So(err, c.ShouldBeNil)
})

c.Convey("返回未取到信息", t, func() {
    do, err := impl.FindByIDByRedis(context.Background(), 2)
    c.So(do, c.ShouldBeNil)
    c.So(err.Error(), c.ShouldEqual, "redis: nil")
})
```

### 如何测试 MQ
目前没有类似 redis 的 mock 工具，因此直接使用gomonkey工具mock掉

### 测试时出现包的循环引用
将测试文件中的包后缀添加为_test，那么其它的包就都会视为外部包，就不会有循环引用问题

### Mac apple 芯片有兼容问题
* 在使用到了方法/函数等 mock 时需指定环境变量：`GOARCH=amd64 go test -v`
* 在测试 data 层时，又不能加 GOARCH=amd64

### 配置代码检查工具 golangci-lint
为了保证代码质量，尽量减少 bug
* mac 安装：brew install golangci-lint
* windows 安装：可以用 docker 或者 install 方式安装
* 将安装好的执行文件路径加入到 path 中，mac：`vi ~/.profile` 加入 ``
1. 配置githook
   * cd 项目目录
   * vi .git/hooks/pre-commit
   * 添加如下代码：
   ```shell
    #!/bin/sh
    echo "Start lint code..."
    if !(golangci-lint --version); then
      echo "请先安装golangci-lint 命令，https://golangci-lint.run/usage/install/"
      exit 1
    fi
    # 目录名后跟上...表示对该目录进行递归查找
    if !(golangci-lint run ./...); then
      echo "Lint fail!"
      exit 1
    fi
    exit 0
   ```
2. 配置 goland
![eUeRd.png](https://i.328888.xyz/2023/02/27/eUeRd.png)


## 五、常见问题
### 1. 如何新增一个 GRPC 接口
* 在 api层定义 GRPC 方法
* `make api` 生成 go 代码 
* 如果是新创建的 service 类，则需要做如下操作
  * 在 service 目录创建对应的类，并实现 api 中生成的接口 `Unimplemented***Server`
  * 在 service.go 中将上面的工厂方法写入wire.NewSet中
  * 在 grpc_server.go 中注册 Service：`xx.RegisterOrderServer(srv, **Service)`
  * 运行 `make wire` 生成依赖代码
* 如果只是在已存在的 Service 中增加方法，则只需要在对应的 service 中实现新增的方法即可
    
### 2. 如何发送消息
项目中已实现了 rocketMQ 的对接，只需要做如下操作即可：
* 在config.ymal 中配置  data.rocketMq 连接信息
* 在需要使用的地方注入SendEvent接口
* 之后调用对应方法发送即可，SendASync 表示异步发送，SendSync表示同步发送（不建议使用，性能低，会阻塞Producer，直到返回了才能响应下一个send）

### 3. 如何创建消费者
消费者对于服务来讲，也是属于外部请求进入服务器，因此需要在 service 的 consumerService 中去注册消息与 service 的调用关系
* 首先在 service 层中创建对应的 receive 处理方法，此方法的格式为 xxx(body string) error，body 为取到的消息体
  * receive 方法中如果 panic 或返回 err 不为 nil，则此条消息消费 ack 被标识为 retry。只有返回了 nil 才会把ack 标识为 success
* 最后在 consumer_service 文件中注册 topic、group、receive 方法的关系：`csr.RegisterSubscribe("主题名", "消费组名", orderConsumer.Test, 2)`
  * 通常消费都组名使用与主题相同的名称，一个消费组只能对应一个 topic，请不要对应多个 topic，以免混乱

### 4. make 支持的命令有哪些
* `make api` 生成 api 目录中 grpc 和 error 的 proto 相关文件（已加入了 validate、errors）
* `make wire` 生成依赖
* `make config` 生成 /internal/con中proto 相关文件
* `make errors` 生成 api 中 error 的 proto 相关文件
* `make generate` 整理 mod 包、更新 wire、执行 go generate 相关命令(目前只有 wire)
* `make all` 执行全部，包含 api、errors、validate、config、wire、generate
* `make ent`  生成 ent 模型

### 5. 如何编写 grpc 的自动验证
proto 中可以定义非常多的自动验证，详细使用规则见 https://github.com/bufbuild/protoc-gen-validate

### 6. 如何使用事务？
框架中已将事务封装进了 context 中，data 的 repo 实现层，和正常写法一样，biz 中只需要注入repository.TransactionRepo接口，调用 Exec 方法即可
, 如下示例所示，将repository.TransactionRepo注入到本类中的 trans 属性中
```go
err := xx.trans.Exec(context.Background(), func(ctx context.Context) error {
    //这里的 ctx 不能自己生成，请使用使用参数中的 ctx，否则此执行将不会在事务中，因为 Exec 方法派生了一个有事务属性的 ctx
    err := ob.repo.UpdateTimeNow(ctx, 1)
    if err != nil {
        return err //返回的 err 不是 nil，则事务会回滚
    }
    err = ob.repo.UpdateTimeNow(ctx, 2)
    if err != nil {
        return err
    }
	//如果方法中产生了 panic，也会回滚
    return nil
})
```

### 7. 如何新增命令行脚本
1. 在 service 中新增脚本执行的方法，命令名称等，参考 demo_cli.go
2. 将其添加到该目录中 provider_set 中，交由wire 管理依赖
3. 在 server/cli.go 添加依赖，并注册命令与 service 的调用关系

更多用法请看github.com/spf13/cobra，注意要按止项目的方式来编写，否则无法用依赖管理、和一些资源的使用

### 8. 如何获取metadata 中请求的用户 id
使用 `metadatamanager.GetUserId(ctx)`

### 9. 引入了哪些工具库？
1. https://github.com/duke-git/lancet/blob/main/README_zh-CN.md 基本所有的工具方法都能找到

### 10. 如何控制 consumer 和 sql 的日志输出
* consumer 日志级别：配置环境变量 ROCKETMQ_GO_LOG_LEVEL=error


# 代码规范
### 1.1代码格式

- 代码都必须用 gofmt 进行格式化。

- 运算符和操作数之间要留空格。

- 建议一行代码不超过 120 个字符，超过部分，请采用合适的换行方式换行。但也有些例外场景，例如 import 行、工具自动生成的代码、带 tag 的 struct 字段。

- 文件长度不能超过 800 行。

- 函数长度不能超过 80 行。

- import 规范

  - 代码都必须用 goimports 进行格式化（建议将代码 Go 代码编辑器设置为：保存时运行 goimports）。

    - 不要使用相对路径引入包，例如 import …/util/net 。

    - 包名称与导入路径的最后一个目录名不匹配时，或者多个相同包名冲突时，则必须使用导入别名。

    - ```go
      	// bad
        "github.com/dgrijalva/jwt-go/v4"
      
        //good
        jwt "github.com/dgrijalva/jwt-go/v4"
      ```

    - 导入的包建议进行分组，匿名包的引用使用一个新的分组，并对匿名包引用进行说明。

    - ```go
      import (
        // go 标准包
        "fmt"
        
        // 第三方包
          "github.com/jinzhu/gorm"
          "github.com/spf13/cobra"
          "github.com/spf13/viper"
        
        // 匿名包单独分组，并对匿名包引用进行说明
          // import mysql driver
          _ "github.com/jinzhu/gorm/dialects/mysql"
        
        // 内部包
          v1 "github.com/marmotedu/api/apiserver/v1"
          metav1 "github.com/marmotedu/apimachinery/pkg/meta/v1"
          "github.com/marmotedu/iam/pkg/cli/genericclioptions"
      )
      ```

### 1.2声明、初始化和定义

- 当函数中需要使用到多个变量时，可以在函数开始处使用 var 声明。在函数外部声明必须使用 var ，不要采用 := ，容易踩到变量的作用域的问题。

- ```go
  var (
    Width  int
    Height int
  )
  ```

- 在初始化结构引用时，请使用 &T{}代替 new(T)，以使其与结构体初始化一致。

- ```go
  // bad
  sptr := new(T)
  sptr.Name = "bar"
  
  // good
  sptr := &T{Name: "bar"}
  ```

- struct 声明和初始化格式采用多行，定义如下。

- ```go
  type User struct{
      Username  string
      Email     string
  }
  
  user := User{
    Username: "colin",
    Email: "colin404@foxmail.com",
  }
  ```

- 相似的声明放在一组，同样适用于常量、变量和类型声明。

- ```go
  // bad
  import "a"
  import "b"
  
  // good
  import (
    "a"
    "b"
  )
  ```

- 尽可能指定容器容量，以便为容器预先分配内存，例如：

- ```go
  v := make(map[int]string, 4)
  v := make([]string, 0, 4)
  ```

- 在顶层，使用标准 var 关键字。请勿指定类型，除非它与表达式的类型不同。

- ```go
  // bad
  var _s string = F()
  
  func F() string { return "A" }
  
  // good
  var _s = F()
  // 由于 F 已经明确了返回一个字符串类型，因此我们没有必要显式指定_s 的类型
  // 还是那种类型
  
  func F() string { return "A" }
  ```

- 对于未导出的顶层常量和变量，使用 _ 作为前缀。

- ```go
  // bad
  const (
    defaultHost = "127.0.0.1"
    defaultPort = 8080
  )
  
  // good
  const (
    _defaultHost = "127.0.0.1"
    _defaultPort = 8080
  )
  ```

- 嵌入式类型（例如 mutex）应位于结构体内的字段列表的顶部，并且必须有一个空行将嵌入式字段与常规字段分隔开。

- ```go
  // bad
  type Client struct {
    version int
    http.Client
  }
  
  // good
  type Client struct {
    http.Client
  
    version int
  }
  ```

### 1.3错误处理

- error作为函数的值返回，必须对error进行处理，或将返回值赋值给明确忽略。对于defer xx.Close()可以不用显式处理。

- ```go
  func load() error {
    // normal code
  }
  
  // bad
  load()
  
  // good
   _ = load()
  ```

- error作为函数的值返回且有多个返回值的时候，error必须是最后一个参数。

- ```go
  // bad
  func load() (error, int) {
    // normal code
  }
  
  // good
  func load() (int, error) {
    // normal code
  }
  ```

- 尽早进行错误处理，并尽早返回，减少嵌套。

- ```go
  // bad
  if err != nil {
    // error code
  } else {
    // normal code
  }
  
  // good
  if err != nil {
    // error handling
    return err
  }
  // normal code
  ```

- 如果需要在 if 之外使用函数调用的结果，则应采用下面的方式。

- ```go
  // bad
  if v, err := foo(); err != nil {
    // error handling
  }
  
  // good
  v, err := foo()
  if err != nil {
    // error handling
  }
  ```

- 错误要单独判断，不与其他逻辑组合判断。

- ```go
  // bad
  v, err := foo()
  if err != nil || v  == nil {
    // error handling
    return err
  }
  
  // good
  v, err := foo()
  if err != nil {
    // error handling
    return err
  }
  
  if v == nil {
    // error handling
    return errors.New("invalid value v")
  }
  ```

- 如果返回值需要初始化，则采用下面的方式

- ```go
  v, err := f()
  if err != nil {
      // error handling
      return // or continue.
  }
  // use v
  ```

- 错误描述建议

  - 告诉用户他们可以做什么，而不是告诉他们不能做什么。

  - 当声明一个需求时，用 must 而不是 should。例如，must be greater than 0、must match regex ‘[a-z]+’。

  - 当声明一个格式不对时，用 must not。例如，must not contain。

  - 当声明一个动作时用 may not。例如，may not be specified when otherField is empty、only name may be specified。

  - 引用文字字符串值时，请在单引号中指示文字。例如，ust not contain ‘…’。

  - 当引用另一个字段名称时，请在反引号中指定该名称。例如，must be greater than request。

  - 指定不等时，请使用单词而不是符号。例如，must be less than 256、must be greater than or equal to 0 (不要用 larger than、bigger than、more than、higher than)。

  - 指定数字范围时，请尽可能使用包含范围。

  - 建议 Go 1.13 以上，error 生成方式为 fmt.Errorf("module xxx: %w", err)。

  - 错误描述用小写字母开头，结尾不要加标点符号，例如：

  - ```go
    // bad
    errors.New("Redis connection failed")
    errors.New("redis connection failed.")
      
    // good
    errors.New("redis connection failed")
    ```

### 1.4 panic 处理

- 在业务逻辑处理中禁止使用 panic。
- 在 main 包中，只有当程序完全不可运行时使用 panic，例如无法打开文件、无法连接数据库导致程序无法正常运行。
- 在 main 包中，使用 log.Fatal 来记录错误，这样就可以由 log 来结束程序，或者将 panic 抛出的异常记录到日志文件中，方便排查问题。
- 可导出的接口一定不能有 panic。
- 包内建议采用 error 而不是 panic 来传递错误。

### 1.5 单元测试

- 单元测试文件名命名规范为 example_test.go。
- 每个重要的可导出函数都要编写测试用例。
- 因为单元测试文件内的函数都是不对外的，所以可导出的结构体、函数等可以不带注释。
- 如果存在 func (b *Bar) Foo ，单测函数可以为 func TestBar_Foo。

### 1.6 类型断言失败处理

type assertion 的单个返回值针对不正确的类型将产生 panic。请始终使用 “comma ok”的惯用法。

```go
// bad
t := n.(int)

// good
t, ok := n.(int)
if !ok {
  // error handling
}
// normal code
```



## 2. 命名规范

命名规范是代码规范中非常重要的一部分，一个统一的、短小的、精确的命名规范可以大大提高代码的可读性，也可以借此规避一些不必要的 Bug。

### 2.1 包命名

- 包名必须和目录名一致，尽量采取有意义、简短的包名，不要和标准库冲突。
- 包名全部小写，没有大写或下划线，使用多级目录来划分层级。
- 项目名可以通过中划线来连接多个单词。
- 包名以及包所在的目录名，不要使用复数，例如，是net/url，而不是net/urls。
- 不要用 common、util、shared 或者 lib 这类宽泛的、无意义的包名。
- 包名要简单明了，例如 net、time、log。

### 2.2 函数命名

- 函数名采用驼峰式，首字母根据访问控制决定使用大写或小写，例如：MixedCaps 或者 mixedCaps。
- 代码生成工具自动生成的代码 (如 xxxx.pb.go) 和为了对相关测试用例进行分组，而采用的下划线 (如 TestMyFunction_WhatIsBeingTested) 排除此规则。

### 2.3 文件命名

- 文件名要简短有意义。
- 文件名应小写，并使用下划线分割单词。

### 2.4 结构体命名

- 采用驼峰命名方式，首字母根据访问控制决定使用大写或小写，例如 MixedCaps 或者 mixedCaps。

- 结构体名不应该是动词，应该是名词，比如 Node、NodeSpec。

- 避免使用 Data、Info 这类无意义的结构体名。

- 结构体的声明和初始化应采用多行，例如：

- ```go
  // User 多行声明
  type User struct {
      Name  string
      Email string
  }
  
  // 多行初始化
  u := User{
      UserName: "colin",
      Email:    "colin404@foxmail.com",
  }
  ```

### 2.5 接口命名

- 接口命名的规则，基本和结构体命名规则保持一致：

  - ~~单个函数的接口名以 “er"”作为后缀（例如 Reader，Writer），有时候可能导致蹩脚的英文，但是没关系。~~

  - ~~两个函数的接口名以两个函数名命名，例如 ReadWriter。~~

  - 三个以上函数的接口名，类似于结构体名。

  - ```go
    // Seeking to an offset before the start of the file is an error.
    // Seeking to any positive offset is legal, but the behavior of subsequent
    // I/O operations on the underlying object is implementation-dependent.
    type Seeker interface {
        Seek(offset int64, whence int) (int64, error)
    }
      
    // ReadWriter is the interface that groups the basic Read and Write methods.
    type ReadWriter interface {
        Reader
        Writer
    }
    ```

### 2.6 变量命名

- 变量名必须遵循驼峰式，首字母根据访问控制决定使用大写或小写。
- 在相对简单（对象数量少、针对性强）的环境中，可以将一些名称由完整单词简写为单个字母，例如：
  - user 可以简写为 u；
  - userID 可以简写 uid。
- 特有名词时，需要遵循以下规则：
  - 如果变量为私有，且特有名词为首个单词，则使用小写，如 apiClient。
  - 其他情况都应当使用该名词原有的写法，如 APIClient、repoID、UserID。

下面列举了一些常见的特有名词。

```go
// A GonicMapper that contains a list of common initialisms taken from golang/lint
var LintGonicMapper = GonicMapper{
    "API":   true,
    "ASCII": true,
    "CPU":   true,
    "CSS":   true,
    "DNS":   true,
    "EOF":   true,
    "GUID":  true,
    "HTML":  true,
    "HTTP":  true,
    "HTTPS": true,
    "ID":    true,
    "IP":    true,
    "JSON":  true,
    "LHS":   true,
    "QPS":   true,
    "RAM":   true,
    "RHS":   true,
    "RPC":   true,
    "SLA":   true,
    "SMTP":  true,
    "SSH":   true,
    "TLS":   true,
    "TTL":   true,
    "UI":    true,
    "UID":   true,
    "UUID":  true,
    "URI":   true,
    "URL":   true,
    "UTF8":  true,
    "VM":    true,
    "XML":   true,
    "XSRF":  true,
    "XSS":   true,
}
```

- 若变量类型为 bool 类型，则名称应以 Has，Is，Can 或 Allow 开头，例如：

```go
var hasConflict bool
var isExist bool
var canManage bool
var allowGitHook bool
```

- 局部变量应当尽可能短小，比如使用 buf 指代 buffer，使用 i 指代 index。
- 代码生成工具自动生成的代码可排除此规则 (如 xxx.pb.go 里面的 Id)

### 2.7 常量命名

- 常量名必须遵循驼峰式，首字母根据访问控制决定使用大写或小写。
- 如果是枚举类型的常量，需要先创建相应类型：

```go
// Code defines an error code type.
type Code int

// Internal errors.
const (
    // ErrUnknown - 0: An unknown error occurred.
    ErrUnknown Code = iota
    // ErrFatal - 1: An fatal error occurred.
    ErrFatal
)
```



### 2.8 Error 的命名

- Error 类型应该写成 FooError 的形式。

```go
type ExitError struct {
  // ....
}
```

- Error 变量写成 ErrFoo 的形式。

```go
var ErrFormat = errors.New("unknown format")
```



## 3. 注释规范

- 每个可导出的名字都要有注释，该注释对导出的变量、函数、结构体、接口等进行简要介绍。
- 全部使用单行注释，禁止使用多行注释。
- 和代码的规范一样，单行注释不要过长，禁止超过 120 字符，超过的请使用换行展示，尽量保持格式优雅。
- 注释必须是完整的句子，以需要注释的内容作为开头，句点作为结尾，格式为 // 名称 描述. 。例如

```go
// bad
// logs the flags in the flagset.
func PrintFlags(flags *pflag.FlagSet) {
  // normal code
}

// good
// PrintFlags logs the flags in the flagset.
func PrintFlags(flags *pflag.FlagSet) {
  // normal code
}
```

- 所有注释掉的代码在提交 code review 前都应该被删除，否则应该说明为什么不删除，并给出后续处理建议。
- 在多段注释之间可以使用空行分隔加以区分，如下所示：

```go
// Package superman implements methods for saving the world.
//
// Experience has shown that a small number of procedures can prove
// helpful when attempting to save the world.
package superman
```

### 3.1 包注释

- 每个包都有且仅有一个包级别的注释。
- 包注释统一用 // 进行注释，格式为 // Package 包名 包描述 ，例如

```go
// Package genericclioptions contains flags which can be added to you command, bound, completed, and produce
// useful helper functions.
package genericclioptions
```

### 3.2 变量 / 常量注释

- 每个可导出的变量 / 常量都必须有注释说明，格式为// 变量名 变量描述，例如：

```go
// ErrSigningMethod defines invalid signing method error.
var ErrSigningMethod = errors.New("Invalid signing method")
```

- 出现大块常量或变量定义时，可在前面注释一个总的说明，然后在每一行常量的前一行或末尾详细注释该常量的定义，例如：

```go
// Code must start with 1xxxxx.    
const (                         
    // ErrSuccess - 200: OK.          
    ErrSuccess int = iota + 100001    
                                                   
    // ErrUnknown - 500: Internal server error.    
    ErrUnknown    

    // ErrBind - 400: Error occurred while binding the request body to the struct.    
    ErrBind    
                                                  
    // ErrValidation - 400: Validation failed.    
    ErrValidation 
)
```

### 3.3 结构体注释

- 每个需要导出的结构体或者接口都必须有注释说明，格式为 // 结构体名 结构体描述.。
- 结构体内的可导出成员变量名，如果意义不明确，必须要给出注释，放在成员变量的前一行或同一行的末尾。例如：

```go
// User represents a user restful resource. It is also used as gorm model.
type User struct {
    // Standard object's metadata.
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Nickname string `json:"nickname" gorm:"column:nickname"`
    Password string `json:"password" gorm:"column:password"`
    Email    string `json:"email" gorm:"column:email"`
    Phone    string `json:"phone" gorm:"column:phone"`
    IsAdmin  int    `json:"isAdmin,omitempty" gorm:"column:isAdmin"`
}
```

### 3.4 方法注释

- 每个需要导出的函数或者方法都必须有注释，格式为// 函数名 函数描述.，例如：

```go
// BeforeUpdate run before update database record.
func (p *Policy) BeforeUpdate() (err error) {
  // normal code
  return nil
}
```

### 3.5 类型注释

- 每个需要导出的类型定义和类型别名都必须有注释说明，格式为 // 类型名 类型描述. ，例如：

```go
// Code defines an error code type.
type Code int
```

## 4. 类型

### 4.1 字符串

- 空字符串判断。

```go
// bad
if s == "" {
    // normal code
}

// good
if len(s) == 0 {
    // normal code
}
```

- []byte/string 相等比较。

```go
// bad
var s1 []byte
var s2 []byte
...
bytes.Equal(s1, s2) == 0
bytes.Equal(s1, s2) != 0

// good
var s1 []byte
var s2 []byte
...
bytes.Compare(s1, s2) == 0
bytes.Compare(s1, s2) != 0
```

- 复杂字符串使用 raw 字符串避免字符转义。

```go
// bad
regexp.MustCompile("\\.")

// good
regexp.MustCompile(`\.`)
```

### 4.2 切片

- 空 slice 判断。

```go
// bad
if len(slice) = 0 {
    // normal code
}

// good
if slice != nil && len(slice) == 0 {
    // normal code
}
```

上面判断同样适用于 map、channel。

- 声明 slice。

```go
// bad
s := []string{}
s := make([]string, 0)

// good
var s []string
```

- slice 复制。

```go
// bad
var b1, b2 []byte
for i, v := range b1 {
   b2[i] = v
}
for i := range b1 {
   b2[i] = b1[i]
}

// good
copy(b2, b1)
```

- slice 新增。

```go
// bad
var a, b []int
for _, v := range a {
    b = append(b, v)
}

// good
var a, b []int
b = append(b, a...)
```

### 4.3 结构体

- struct 初始化。

struct 以多行格式初始化。

```go
type user struct {
  Id   int64
  Name string
}

u1 := user{100, "Colin"}

u2 := user{
    Id:   200,
    Name: "Lex",
}
```

## 5. 控制结构

### 5.1 if

- if 接受初始化语句，约定如下方式建立局部变量。

```go
if err := loadConfig(); err != nil {
  // error handling
  return err
}
```

- if 对于 bool 类型的变量，应直接进行真假判断。

```go
var isAllow bool
if isAllow {
  // normal code
}
```

### 5.2 for

- 采用短声明建立局部变量。

```go
sum := 0
for i := 0; i < 10; i++ {
    sum += 1
}
```

- 不要在 for 循环里面使用 defer，defer 只有在函数退出时才会执行。

```go
// bad
for file := range files {
  fd, err := os.Open(file)
  if err != nil {
    return err
  }
  defer fd.Close()
  // normal code
}

// good
for file := range files {
  func() {
    fd, err := os.Open(file)
    if err != nil {
      return err
    }
    defer fd.Close()
    // normal code
  }()
}
```

### 5.3 range

- 如果只需要第一项（key），就丢弃第二个。

```go
for key := range keys {
// normal code
}
```

- 如果只需要第二项，则把第一项置为下划线。

```go
sum := 0
for _, value := range array {
    sum += value
}
```

### 5.4 switch

- 必须要有 default。

```go
switch os := runtime.GOOS; os {
    case "linux":
        fmt.Println("Linux.")
    case "darwin":
        fmt.Println("OS X.")
    default:
        fmt.Printf("%s.\n", os)
}
```

### 5.5 goto

- 业务代码禁止使用 goto 。
- 框架或其他底层源码尽量不用。

## 6. 函数

- 传入变量和返回变量以小写字母开头。
- 函数参数个数不能超过 5 个。
- 函数分组与顺序
  - 函数应按粗略的调用顺序排序。
  - 同一文件中的函数应按接收者分组。
- 尽量采用值传递，而非指针传递。
- 传入参数是 map、slice、chan、interface ，不要传递指针。

### 6.1函数参数

- 如果函数返回相同类型的两个或三个参数，或者如果从上下文中不清楚结果的含义，使用命名返回，其他情况不建议使用命名返回，例如：

```go
func coordinate() (x, y float64, err error) {
  // normal code
}
```

- 传入变量和返回变量都以小写字母开头。
- 尽量用值传递，非指针传递。
- 参数数量均不能超过 5 个。
- 多返回值最多返回三个，超过三个请使用 struct。

### 6.2 defer

- 当存在资源创建时，应紧跟 defer 释放资源 (可以大胆使用 defer，defer 在 Go1.14 版本中，性能大幅提升，defer 的性能损耗即使在性能敏感型的业务中，也可以忽略)。
- 先判断是否错误，再 defer 释放资源，例如：

```go
rep, err := http.Get(url)
if err != nil {
    return err
}

defer resp.Body.Close()
```

### 6.3 方法的接收器

- 推荐以类名第一个英文首字母的小写作为接收器的命名。
- 接收器的命名在函数超过 20 行的时候不要用单字符。
- 接收器的命名不能采用 me、this、self 这类易混淆名称。

### 6.4 嵌套

- 嵌套深度不能超过 4 层。

### 6.5 变量命名

- 变量声明尽量放在变量第一次使用的前面，遵循就近原则。
- 如果魔法数字出现超过两次，则禁止使用，改用一个常量代替，例如：

### 7. 最佳实践

- 尽量少用全局变量，而是通过参数传递，使每个函数都是“无状态”的。这样可以减少耦合，也方便分工和单元测试。
- 在编译时验证接口的符合性，例如：
- 服务器处理请求时，应该创建一个 context，保存该请求的相关信息（如 requestID），并在函数调用链中传递。

### 7.1 性能

- string 表示的是不可变的字符串变量，对 string 的修改是比较重的操作，基本上都需要重新申请内存。所以，如果没有特殊需要，需要修改时多使用 []byte。
- 优先使用 strconv 而不是 fmt。

### 7.2 注意事项
- append 要小心自动分配内存，append 返回的可能是新分配的地址。
- 如果要直接修改 map 的 value 值，则 value 只能是指针，否则要覆盖原来的值。
- map 在并发中需要加锁。
- 编译过程无法检查 interface{} 的转换，只能在运行时检查，小心引起 panic。
