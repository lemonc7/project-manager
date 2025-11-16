import { Button } from "@/components/ui/button"
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from "@/components/ui/sidebar"
import { Users, ListTodo, Sheet, User } from "lucide-react"
import { useNavigate } from "react-router"

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

const CustomAside = () => {
  const navigate = useNavigate()
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
                      <a href={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        {/* <SidebarFooter>设置</SidebarFooter> */}
      </Sidebar>
    </div>
  )
}

export default CustomAside
