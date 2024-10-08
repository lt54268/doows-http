# dootask-workspace
## 接口使用说明
请求地址：http://127.0.0.1:5555
### 一、同步用户ID (GET)
 ```
http://127.0.0.1:5555/sync
 ```
### 二、设置创建工作区权限 (POST)
 ```
http://127.0.0.1:5555/set
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
### 三、创建工作区 (POST)
 ```
http://127.0.0.1:5555/create
 ```
 ```
headers
 "Content-Type: application/json" 
body
 {
   "user_id": 1           // 用户ID
 }
 ```
 ### 四、删除工作区 (DELETE)
 ```
http://127.0.0.1:5555/delete-ws
 ```
 ```
headers
 "Content-Type: application/json" 
body
 {
   "user_id": 1           // 用户ID
 }
 ```
### 五、检查已创建的工作区数量 (GET)
 ```
http://127.0.0.1:5555/check
 ```
### 六、新建对话窗口 (POST)
```
http://127.0.0.1:5555/new
```
```
headers
 "Content-Type: application/json" 
body
{
    "slug": "workspace-for-user-1",
    "model": "ChatGPT",
    "avatar": "sk123"
}
```
### 七、查询已获授权用户 (POST)
```
http://127.0.0.1:5555/get-user
```
```
headers
"Content-Type: application/json"
body
{
"user_id": 1
}
```
