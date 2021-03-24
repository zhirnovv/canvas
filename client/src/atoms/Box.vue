<template>
  <component :is="tag" class="Box" v-on="$listeners">
    <slot></slot>
  </component>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

@Component
export default class Box extends Vue {
  @Prop({
    default: "div",
  })
  public tag!: keyof HTMLElementTagNameMap;
}
</script>

<style scoped>
.Box {
  border-radius: 5px;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.25);
  background: #ffffff;
  appearance: none;
}
</style>
