import { ref } from "vue"

const BASE_STATS_LENGTH = 3
export function useTasks() {
  // 状态 (State): 持久化的数组
  const lengths = ref<number[]>(Array(BASE_STATS_LENGTH).fill(0))

  // 动作 (Actions): 唯一的修改状态的方法
  function setLength(index: number, length: number) {
    if (index >= 0 && index < BASE_STATS_LENGTH) {
      lengths.value[index] = length
    }
  }

  return { lengths, setLength }
}
