# Data
与外部交互的基础层：包含 数据库、Redis、MQ 等

## 如何编写数据库的dao
1. 注入 data *Data
2. 获取 Db 对象时，禁止直接取属性，必须使用 GetDb 方法，此方法解析了上下文 Context 中是否开启了事务