import { SidebarProvider } from "@/components/ui/sidebar"
import { Outlet } from "react-router"
import CustomAside from "./CustomAside"
import CustomHeader from "./CustomHeader"
import { Toaster } from "sonner"

const CustomLayout = () => {
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
          <CustomHeader />
          <Outlet />
        </div>
      </div>
      <Toaster richColors position="top-right" closeButton duration={2000}/>
    </SidebarProvider>
  )
}

export default CustomLayout
