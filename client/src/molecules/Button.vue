<template>
  <Box :tag="'button'" :class="boxClassObject" v-on="$listeners">
    <slot></slot>
  </Box>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import Box from "../atoms/Box.vue";

export enum ButtonVariant {
  Filled = "filled",
}

export enum ButtonColor {
  Primary = "primary",
}

@Component({
  components: {
    Box,
  },
})
export default class Button extends Vue {
  @Prop({
    default: ButtonColor.Primary,
  })
  public color!: ButtonColor;

  @Prop({
    default: ButtonVariant.Filled,
  })
  public variant!: ButtonVariant;

  get boxClassObject() {
    return {
      ["Button-box"]: true,
      [`Button-box__${this.variant}-${this.color}`]: true,
    };
  }
}
</script>

<style scoped lang="scss">
@import "../assets/palette.scss";

.Button-box {
  padding: 8px 16px;
  outline: none;

  &:disabled {
    pointer-events: none;
    background-color: $royal-violet;
    opacity: 0.4;
  }

  &__filled {
    &-primary {
      border: none;
      background-color: $royal-violet;
      transition: opacity 0.2s ease;

      &:hover,
      &:focus {
        opacity: 0.8;
      }

      &:active {
        opacity: 0.9;
      }
    }
  }
}
</style>


