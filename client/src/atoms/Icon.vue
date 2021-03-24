<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :width="width"
    :height="height"
    :viewBox="viewBox"
    :aria-labelledby="iconName"
    :fill="'none'"
    role="presentation"
  >
    <title :id="iconName" lang="en">{{ iconName }}</title>
    <slot></slot>
  </svg>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

@Component
export default class Icon extends Vue {
  @Prop({
    required: true,
  })
  public iconName!: string;

  @Prop({
    required: true,
  })
  public viewBox!: string;

  @Prop({
    required: true,
  })
  public width!: number;

  @Prop({
    required: true,
  })
  public height!: number;
}
</script>
