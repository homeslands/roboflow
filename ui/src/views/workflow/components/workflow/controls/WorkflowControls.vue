<script setup lang="ts">
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { ControlButton } from '@vue-flow/controls'
import { Panel, useVueFlow } from '@vue-flow/core'
import { FullscreenIcon, ScanIcon, ZoomInIcon, ZoomOutIcon } from 'lucide-vue-next'
import '@vue-flow/controls/dist/style.css'

const ANIMATION_DURATION = 300
const { fitView, zoomIn, zoomOut } = useVueFlow()

function handleCenterView() {
  fitView({ duration: ANIMATION_DURATION, padding: 0.2 })
}
function handleZoomOut() {
  zoomOut({ duration: ANIMATION_DURATION })
}
function handleZoomIn() {
  zoomIn({ duration: ANIMATION_DURATION })
}
function handleFixLayout() {}
</script>

<template>
  <Panel position="bottom-left" class="workflow-controls !shadow-md">
    <div class="workflow-controls-content">
      <ControlButton
        class="workflow-controls-button !bg-muted/50 hover:!bg-muted"
        @click="handleCenterView"
      >
        <Tooltip>
          <TooltipTrigger>
            <FullscreenIcon class="text-muted-foreground" />
          </TooltipTrigger>
          <TooltipContent>
            Center View
          </TooltipContent>
        </Tooltip>
      </ControlButton>

      <ControlButton
        class="workflow-controls-button !bg-muted/50 hover:!bg-muted"
        @click="handleZoomOut"
      >
        <Tooltip>
          <TooltipTrigger>
            <ZoomOutIcon class="text-muted-foreground" />
          </TooltipTrigger>
          <TooltipContent>
            Zoom Out
          </TooltipContent>
        </Tooltip>
      </ControlButton>

      <ControlButton
        class="workflow-controls-button !bg-muted/50 hover:!bg-muted"
        @click="handleZoomIn"
      >
        <Tooltip>
          <TooltipTrigger>
            <ZoomInIcon class="text-muted-foreground" />
          </TooltipTrigger>
          <TooltipContent>
            Zoom In
          </TooltipContent>
        </Tooltip>
      </ControlButton>

      <ControlButton
        class="workflow-controls-button !bg-muted/50 hover:!bg-muted"
        @click="handleFixLayout"
      >
        <Tooltip>
          <TooltipTrigger>
            <ScanIcon class="text-muted-foreground" />
          </TooltipTrigger>
          <TooltipContent>
            Fix Layout
          </TooltipContent>
        </Tooltip>
      </ControlButton>
    </div>
  </Panel>
</template>

<style lang="css" scoped>
.workflow-controls {
  border-radius: 4px;
  overflow: hidden;
}

.workflow-controls-content {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.workflow-controls-button {
  padding: 6px;
  height: 24px;
  width: 24px;
  border: none !important;
}

.workflow-controls-button svg {
  max-height: unset;
  max-width: unset;
  width: 100%;
  height: 100%;
}
</style>
