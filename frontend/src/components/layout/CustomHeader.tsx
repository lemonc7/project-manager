import { SidebarTrigger } from "@/components/ui/sidebar"

const CustomHeader = () => {
  return (
    <header className="sticky top-0 z-10 flex h-16 items-center justify-between bg-background">
      <SidebarTrigger />
    </header>
  )
}
export default CustomHeader
