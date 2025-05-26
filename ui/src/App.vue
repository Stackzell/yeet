<script setup lang="ts">
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
</script>

<template>
  <div class="h-screen w-screen">
    <ResizablePanelGroup direction="horizontal">
      <ResizablePanel :default-size="15">
        <div>Sidebar</div>
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel>
        <div class="h-full">
          <RouterView />
        </div>
      </ResizablePanel>
    </ResizablePanelGroup>
  </div>
</template>
