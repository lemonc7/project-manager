import { Button } from "@/components/ui/button"
import React from "react"
import { useNavigate } from "react-router"

const NotFoundPage: React.FC = () => {
  const navigate = useNavigate()

  return (
    <div className="flex flex-1 flex-col items-center justify-center text-center px-4">
      <h1 className="text-6xl font-bold text-gray-800 mb-4 select-none">404</h1>
      <p className="text-lg text-gray-600 mb-6">抱歉，你访问的页面不存在。</p>
      <Button variant="default" onClick={() => navigate("/home")}>
        返回首页
      </Button>
    </div>
  )
}

export default NotFoundPage
