import { defineStore } from "pinia"
import { ref } from "vue"
import type { Process } from "../models/processes"

export const useProcessStore = defineStore("process", () => {
  const maxSequence = ref(1)
  function setMaxSequence(processes?: Process[]) {
    if (processes?.length) {
      const maxVal = processes.reduce((max, p) => Math.max(max, p.sequence), 0)
      maxSequence.value = maxVal + 1
    } else {
      maxSequence.value = 1
    }
  }

  const sequenceToScroll = ref<number | null>(null)
  function setSequenceToScroll(seq: number) {
    sequenceToScroll.value = seq
  }
  function clearSequenceToScroll() {
    sequenceToScroll.value = null
  }

  return {
    maxSequence,
    setMaxSequence,
    sequenceToScroll,
    setSequenceToScroll,
    clearSequenceToScroll
  }
})
