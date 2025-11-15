import axios from "axios"

const service = axios.create({
  baseURL: "http://localhost:8080/api",
  timeout: 3000,
  headers: {
    "Content-Type": "application/json;charset=UTF-8"
  }
})

service.interceptors.response.use(
  (response) => response,
  (error) => {
    if (!error.response) {
      return Promise.reject(new Error("网络错误or服务器无响应"))
    }

    const code = error.response.status
    let message = "未知错误"

    switch (code) {
      case 400:
        message = "请求错误"
        break
      case 401:
        message = "未授权"
        break
      case 403:
        message = "拒绝访问"
        break
      case 404:
        message = "资源未找到"
        break
      case 422:
        message = "请求参数验证错误"
        break
      case 500:
        message = "服务器异常"
        break
    }

    return Promise.reject(new Error(`${message} (${code})`))
  }
)

export default service
