模拟登录支付宝获取Cookie
==================================

## 实现方式

1. 安装docker-selenium

    ```ssh
        docker run --rm  -p 4433:4433 -v /dev/shm:/dev/shm selenium/standalone-chrome
    ```
1. 使用Selenium登录支付宝,截取验证码

1. 调用打码平台接口，解析验证码

1. 登录获取Cookie

## 接口

地址：
```html
 /api/alicookies
```

方法

POST

调用参数（JSON）

```json

{
	"username": "支付宝用户名",
	"password": "支付宝密码"
}
```

返回

```json

{
    "ret": 0, // 0:成功, 1:失败
    "msg": "失败原因", //只有ret == 1时才返回
    "data": {} //Cookie数据，只有ret==0时才返回
}