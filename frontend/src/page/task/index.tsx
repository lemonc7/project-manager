import { Loading } from "@/components/common/loading"

function TaskPage() {
  return (
    <main className="flex items-center justify-center h-full">
        <Loading content="任务列表加载中..." />
    </main>
  )
}

export default TaskPage
