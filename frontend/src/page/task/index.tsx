import { createTeam } from "@/api/teams"

const TaskPage = () => {
  createTeam({
    name: "工艺",
    description: ""
  })
  return <main className="flex-1 overflow-auto"></main>
}

export default TaskPage
