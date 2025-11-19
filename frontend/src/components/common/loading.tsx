import { cn } from "@/lib/utils"
import { Item, ItemContent, ItemMedia, ItemTitle } from "../ui/item"
import { Spinner } from "../ui/spinner"

export function Loading({ content }: { content: string }) {
  return (
    <div
      className={cn("flex w-full max-w-xs", "flex-col gap-4 [--radius:1rem]")}
    >
      <Item variant="muted">
        <ItemMedia>
          <Spinner />
        </ItemMedia>
        <ItemContent>
          <ItemTitle className="line-clamp-1">{content}</ItemTitle>
        </ItemContent>
      </Item>
    </div>
  )
}
