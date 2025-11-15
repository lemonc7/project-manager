import { useParams } from "react-router"

const ProcessPage = () => {
  const params = useParams<{ id: string }>()
  return <div>当前项目 ID: {params.id}</div>
}
export default ProcessPage
