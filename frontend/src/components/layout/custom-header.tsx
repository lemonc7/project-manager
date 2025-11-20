import { SidebarTrigger } from "@/components/ui/sidebar"
import { useMatches, useNavigate } from "react-router"
import { Separator } from "../ui/separator"
import { ButtonGroup } from "../ui/button-group"
import { Button } from "../ui/button"
import { ArrowLeft, ArrowRight, HomeIcon, RefreshCw } from "lucide-react"
import { HEADER_ACTIONS_ID } from "../common/header-actions"

function CustomHeader({ onRefresh }: { onRefresh: () => void }) {
  const matches = useMatches()
  const handle = matches[matches.length - 1]?.handle as
    | { title?: string }
    | undefined
  const title = handle?.title ?? ""

  const navigate = useNavigate()

  return (
    <header className="sticky top-0 z-10 flex h-16 items-center justify-between bg-background">
      {/* 左侧aside折叠按钮 | 标题 | 主页 | 上一页 下一页 */}
      <div className="flex items-center space-x-2 h-4">
        <SidebarTrigger className="w-10 h-10" title="打开/关闭菜单" />
        <Separator orientation="vertical" className="bg-gray-400" />
        <h1 className="font-semibold">{title}</h1>
        <Separator orientation="vertical" className="bg-gray-400" />
        <Button
          title="主页"
          variant="outline"
          size="sm"
          className="ml-2"
          onClick={() => navigate("/home")}
        >
          <HomeIcon />
        </Button>
        <ButtonGroup className="hidden md:flex">
          <Button
            title="上一页"
            variant="outline"
            size="sm"
            className="w-12"
            onClick={() => navigate(-1)}
          >
            <ArrowLeft />
          </Button>
          <Button
            title="下一页"
            variant="outline"
            size="sm"
            className="w-12"
            onClick={() => navigate(1)}
          >
            <ArrowRight />
          </Button>
        </ButtonGroup>
      </div>
      {/* 右侧动态元素 | 刷新*/}
      <div className="flex items-center gap-2">
        <div id={HEADER_ACTIONS_ID} className="flex items-center gap-2" />
        <Button variant="outline" size="sm" title="刷新" onClick={onRefresh}>
          <RefreshCw />
        </Button>
      </div>
    </header>
  )
}
export default CustomHeader
