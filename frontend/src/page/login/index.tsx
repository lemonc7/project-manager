import { login } from "@/api/auth"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { isLoggedIn } from "@/util/auth"
import { useState } from "react"
import { Navigate } from "react-router"

export default function Login() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState("")

  if (isLoggedIn()) {
    return <Navigate to="/home" replace />
  }

  const handlerSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    setLoading(true)
    setErr("")

    try {
      await login({ username, password })
      window.location.href = "/home"
    } catch (err) {
      if (err instanceof Error) {
        setErr(err.message)
      } else {
        setErr(String(err) || "登录失败")
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>登录</CardTitle>
        </CardHeader>
        <form onSubmit={handlerSubmit}>
          <CardContent>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2">
                <Label htmlFor="username">用户名</Label>
                <Input
                  id="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  placeholder="请输入用户名"
                  required
                />
              </div>
              <div className="grid gap-2">
                <Label htmlFor="password">密码</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="请输入密码"
                  required
                />
              </div>
              {err && <p className="text-red-500">{err}</p>}
            </div>
          </CardContent>
          <CardFooter className="flex-col gap-2 mt-6">
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "登录中..." : "登录"}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
