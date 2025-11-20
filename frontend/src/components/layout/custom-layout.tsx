import { SidebarProvider } from "@/components/ui/sidebar"
import { Outlet } from "react-router"
import CustomAside from "./custom-aside"
import CustomHeader from "./custom-header"
import { Toaster } from "sonner"
import { useQueryClient } from "@tanstack/react-query"

function CustomLayout() {
  const queryClient = useQueryClient()

  const handlerRefresh = async () => {
    await queryClient.invalidateQueries({ type: "active" })
  }

  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "12rem",
          "--sidebar-width-mobile": "16rem"
        } as React.CSSProperties
      }
    >
      <CustomAside />
      {/* 外侧容器 */}
      <div className="flex flex-1 flex-col bg-gray-50 p-4">
        {/* 内侧圆角容器 */}
        <div className="flex flex-1 flex-col rounded-xl bg-white shadow overflow-hidden px-4">
          <CustomHeader onRefresh={handlerRefresh} />
          <Outlet />
        </div>
      </div>
      <Toaster richColors position="top-center" closeButton duration={2000} />
    </SidebarProvider>
  )
}

export default CustomLayout
