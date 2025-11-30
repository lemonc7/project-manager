<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <Button
        variant="outline"
        role="combobox"
        :aria-expanded="open"
        class="w-50 justify-between"
      >
        {{ selectedEmployee?.name || "选择员工..." }}
        <ChevronsUpDownIcon class="opacity-50"/>
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-40 p-0">
      <Command>
        <CommandInput class="h-9" placeholder="搜索员工..." />
        <CommandList>
          <CommandEmpty>未找到员工</CommandEmpty>
          <CommandGroup>
            <CommandItem
              v-for="emp in data"
              :key="emp.id"
              :value="emp.id"
              @select="
                (ev) => {
                  selectEmployee(Number(ev.detail.value))
                }
              "
            >
              {{ emp.name }}
              <CheckIcon
                :class="
                  cn(
                    'ml-auto',
                    props.employeeId === emp.id ? 'opacity-100' : 'opacity-0'
                  )
                "
              />
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>

<script setup lang="ts">
import useEmployees from "@/composables/use-employees"
import { computed, ref } from "vue"
import Popover from "../ui/popover/Popover.vue"
import PopoverTrigger from "../ui/popover/PopoverTrigger.vue"
import PopoverContent from "../ui/popover/PopoverContent.vue"
import Command from "../ui/command/Command.vue"
import CommandInput from "../ui/command/CommandInput.vue"
import CommandList from "../ui/command/CommandList.vue"
import CommandEmpty from "../ui/command/CommandEmpty.vue"
import CommandGroup from "../ui/command/CommandGroup.vue"
import CommandItem from "../ui/command/CommandItem.vue"
import { CheckIcon, ChevronsUpDownIcon } from "lucide-vue-next"
import { cn } from "@/lib/utils"

const props = defineProps<{
  employeeId?: number
}>()

const emit = defineEmits<{
  (e: "update:employeeId", value?: number): void
}>()

const { data } = useEmployees()

const open = ref(false)

const selectedEmployee = computed(() =>
  data.value?.find((e) => e.id === props.employeeId)
)

const selectEmployee = (selectedId: number) => {
  emit("update:employeeId", selectedId)
  open.value = false
}
</script>
