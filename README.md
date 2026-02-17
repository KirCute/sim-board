# sim-board

通用桌游模拟器，一个基于 WebSocket 的轻量级多人在线桌面游戏平台。它不包含任何游戏规则引擎，玩家可以自由进行任意类型的桌面游戏，只需自行遵守规则。适用于熟人之间线上聚会，提供了 UNO、扑克、骰子和筹码四种基础牌具。

![演示](https://upyun.kircute.top/sim_board.png)

### 快速开始
- 直接编译：

  1. 编译[前端](https://github.com/KirCute/sim-board-frontend)。
  2. 将编译好的前端放在`./public/dist`。
  3. 编译后端，运行。

- 使用 Docker：

  1. 构建镜像：

     ```bash
     docker build -t sim-board .
     ```

  2. 启动容器：

     ```bash
     docker run -d --name sim-board -p 6700:6700 sim-board
     ```

### 添加自定义牌具

1. 在`deck`下新建 package，在其中添加：

   1. 实现以下接口的结构体：

      ```go
      type Deck interface {
      	Type() string  // 返回字符串常量
      	Name() string  // 返回添加牌具后牌具实例的名称，用于区分一场游戏中类型相同用途不同的牌具
      	RestLen() int  // 牌具中的剩余牌数，返回-1视为牌具中有无限张牌
      	MaxLen() int  // 牌具的总牌数，返回-1视为牌具中有无限张牌
      	Return(card sim_board.Card)  // 将一张已经发出的牌放回牌堆中
      	Draw(count int) []sim_board.Card  // 发出 count 张牌
      }
      ```
   
   2. 一个参数结构体，成员为牌具的构造参数，成员需要具有以下标签：
   
      - `json`：（必须）参数序列化后的 key。
      - `label`：（必须）前端看到的参数名。
      - `type`：（必须）参数类型，有效值有`string`、`int`、`bool`。
      - `min`：（可选）仅适用于`int`参数，最小值。
      - `max`：（可选）仅适用于`int`参数，最大值。
      - `default`：（可选）默认值。
   
   3. 牌具结构体的构造函数`Create`，形参仅有参数结构体的指针，返回牌具结构体的指针。
   
   4. `GetHTML`方法，用于将`sim_board.Card`类型的牌转换为其展示在前端的 HTML 字符串。
   
   5. ```go
      const Name = "牌具结构体Type函数的返回值"
      
      func init() {
      	sim_board.RegisterDeck(Name, reflect.ValueOf(Create), GetHTML)
      }
      ```
2. 在`deck/all.go`中 import 自定义牌具的 package。
