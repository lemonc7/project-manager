import { SidebarProvider } from "@/components/ui/sidebar"
import { Outlet } from "react-router"
import CustomAside from "./custom-aside"
import CustomHeader from "./custom-header"
import { toast, Toaster } from "sonner"
import { useState } from "react"

function CustomLayout() {
  const [outletKey, setOutletKey] = useState(0)
  const handlerRefresh = () => {
    setOutletKey((prev) => prev + 1)
    toast.info("页面已刷新", { duration: 500 })
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
          <Outlet key={outletKey} />
        </div>
      </div>
      <Toaster richColors position="top-center" closeButton duration={2000} />
    </SidebarProvider>
  )
}

export default CustomLayout
