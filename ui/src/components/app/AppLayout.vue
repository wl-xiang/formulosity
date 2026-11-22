<template>
  <main class="mb-auto flex-grow">
    <div class="p-4">
      <div class="mb-4 flex justify-end">
        <PreferenceControls />
      </div>
      <slot />
    </div>
  </main>
</template>

<script setup lang="ts">
import PreferenceControls from '@/components/ui/PreferenceControls.vue'
</script>
