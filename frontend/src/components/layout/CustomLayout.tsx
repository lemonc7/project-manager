import { SidebarProvider } from "@/components/ui/sidebar"
import { Outlet } from "react-router"
import CustomAside from "./CustomAside"
import CustomHeader from "./CustomHeader"

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
    </SidebarProvider>
  )
}

export default CustomLayout
