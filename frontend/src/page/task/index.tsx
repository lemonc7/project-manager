import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

function TaskOverview() {
  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader>
        <CardTitle>概览</CardTitle>
      </CardHeader>
      <CardContent>内容</CardContent>
    </Card>
  )
}

function TaskTodo() {
  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader>
        <CardTitle>待办</CardTitle>
      </CardHeader>
      <CardContent>内容</CardContent>
    </Card>
  )
}

function TaskDoing() {
  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader>
        <CardTitle>进行中</CardTitle>
      </CardHeader>
      <CardContent>内容</CardContent>
    </Card>
  )
}

function TaskNext() {
  return (
    <Card className="hover:shadow-lg transition-shadow duration-300">
      <CardHeader>
        <CardTitle>下一步</CardTitle>
      </CardHeader>
      <CardContent>内容</CardContent>
    </Card>
  )
}

function TaskPage() {
  return (
    <main className="flex flex-col gap-4 p-4 h-full">
      <div className="w-full">
        <TaskOverview />
      </div>
      <div className="grid grid-cols-3 gap-4">
        <TaskDoing />
        <TaskNext />
        <TaskTodo />
      </div>
    </main>
  )
}

export default TaskPage
