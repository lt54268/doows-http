# dootask-workspace
## 接口使用说明
请求地址：http://服务器IP:5555
### 一、同步用户ID
 ```
http://服务器IP:5555/sync
 ```
### 二、设置创建工作区权限（POST）
 ```
http://服务器IP:5555/set
 ```
 ```
 headers
 "Content-Type: application/json" 
body
 {
   "user_id": 1,           // 用户ID
   "is_create": true       // true:允许创建工作区，false:不允许创建工作区
 }
 ```
### 三、创建工作区
 ```
{
  "action": "create",
  "data": {
    "user_id": 1       // 用户ID
  }
}
 ```
 ### 四、删除工作区
 ```
 {
  "action": "delete-ws",
  "data": {
    "user_id": 1       // 用户ID
  }
}
 ```
### 五、检查已创建的工作区数量
 ```
{
  "action": "check"
}
 ```
### 六、获取已创建工作区的用户ID
 ```
{
  "action": "get-users"
}
 ```
### 七、获取workspace_permission表中某个用户的所有字段
 ```
{
  "action": "get-workspace",
  "data": {
      "user_id": 1
  }
}
 ```
### 八、流式对话问答
 ```
{
  "action": "stream-chat",
  "data": {
    "message": "哈哈哈",
    "mode": "chat",
    "sessionId": 1,                // 对话ID，每个对话窗口唯一
    "slug": "workspace-for-user-1"
  }
}
 ```
### 九、常规对话问答
 ```
{
  "action": "chat",
  "data": {
    "message": "不哈哈",
    "mode": "chat",
    "sessionId": 2, 
    "slug": "workspace-for-user-1"
  }
}
 ```
### 十、存储最后一条聊天对话
 ```
{
  "action": "insert-message",
  "data": {
    "session_id": 1,
    "user_id": 3,                    // 用户ID
    "last_message": "嘻嘻嘻"         // 最后一条消息
  }
}
```
### 十一、获取history_aichat表中某个用户的所有字段
```
{
  "action": "get-history",
  "data": {
    "user_id": 1             // 用户ID
  }
}
```
