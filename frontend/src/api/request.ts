import { getToken, removeToken } from "../utils/auth"
import axios from "axios"

const service = axios.create({
  baseURL: "/api",
  // baseURL: "http://localhost:5001/api",
  timeout: 3000,
  headers: {
    "Content-Type": "application/json;charset=UTF-8"
  }
})

// 请求拦截
service.interceptors.request.use(
  (config) => {
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截
service.interceptors.response.use(
  (response) => response,
  (error) => {
    if (!error.response) {
      return Promise.reject(new Error("网络错误or服务器无响应"))
    }

    const code = error.response.status
    if (code === 401) {
      removeToken()
      window.location.href = "/login"
    }

    return Promise.reject(error)
  }
)

export default service
