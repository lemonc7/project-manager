import { Button } from "@/components/ui/button"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from "@/components/ui/sidebar"
import { Users, ListTodo, Sheet, User, LogOut } from "lucide-react"
import { Link, useNavigate } from "react-router"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from "../ui/dropdown-menu"
import { useQuery } from "@tanstack/react-query"
import { getEmployeeMe } from "@/api/employees"
import { removeToken } from "@/util/auth"
import type { EmployeeMe } from "@/model/employees"

const items = [
  {
    title: "任务",
    url: "/home",
    icon: ListTodo
  },
  {
    title: "项目总览",
    url: "/projects",
    icon: Sheet
  },
  {
    title: "成员",
    url: "/employees",
    icon: User
  }
]

export default function CustomAside() {
  const { data } = useQuery<EmployeeMe, Error>({
    queryKey: ["user"],
    queryFn: getEmployeeMe
  })

  const navigate = useNavigate()
  const logout = () => {
    removeToken()
    navigate("/login", { replace: true })
  }

  return (
    <div className="[&_div[data-slot=sidebar-inner]]:bg-gray-50">
      <Sidebar className="border-none!">
        <SidebarHeader className="mt-2">
          <Button variant="ghost" onClick={() => navigate("/home")}>
            <Users />
            项目管理平台
          </Button>
        </SidebarHeader>
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>项目管理</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <Link to={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter className="flex justify-start py-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost">
                <User />
                <span>{data?.name || "未知"}</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={logout}>
                <LogOut />
                登出
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarFooter>
      </Sidebar>
    </div>
  )
}
