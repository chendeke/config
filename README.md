# config



## 配置使用规则



### 推荐使用yaml格式

```
logs:
  level: debug
  encode: json #json or console
  levelPort: 5053
  levelPattern: /handle/level
  initFields:
    company: yourcompany
    app: yourapp
```



- level 是日志等级
- encode 输出格式
- levelPort 日志等级服务端口, runtime调整level接口
- levelPattern 日志等级服务uri
- initFields 为日志with字段



### 配置文件目录

默认读取配置文件目录为: 运行目录/conf



`conf/config.yaml`为基本配置, 任何环境下都会读取此内容



**其他读取规则:**

可设置环境变量`runtime`值分别如下:

- dev 开发环境
- test 测试环境
- prod 正式环境



#### dev环境

读取`conf/dev.yaml`

如果内容有重复, 当前值会覆盖config.yaml基本配置



#### test环境

读取`conf/test.yaml`

如果内容有重复, 当前值会覆盖config.yaml基本配置



#### prod环境

读取`conf/prod.yaml`

如果内容有重复, 当前值会覆盖config.yaml基本配置



#### 无设置

会检查是否有dev, test, prod文件, 然后会依次读取, 后者值覆盖前者, 这样上线程序时, 只要提供prod的配置文件即可, 不需要设置环境变量, 尤其是当不方便时



### 配置使用说明

```
type WatchConfig struct {
	Level        string                 `json:"level" default:"debug"`
	Encode       string                 `json:"encode"`
	LevelPort    int                    `json:"level_port" default:"9090"`
	LevelPattern string                 `json:"level_pattern" default:"/handle/okokok"`
	InitFields   map[string]interface{} `json:"init_fields"`
}

func main() {
	wc := new(WatchConfig)
	err := config.Get("logs").Scan(wc)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(wc)
	}
}
```



支持**default tag**

