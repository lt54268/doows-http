# dootask-workspace
## 接口使用说明
请求地址：http://服务器IP:5555
### 一、同步用户ID (GET)
 ```
http://服务器IP:5555/sync
 ```
### 二、设置创建工作区权限 (POST)
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
### 三、创建工作区 (POST)
 ```
http://服务器IP:5555/create
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
http://服务器IP:5555/delete-ws
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
http://服务器IP:5555/check
 ```

