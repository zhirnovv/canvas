<template>
  <component :is="tag" :class="classObject">
    <slot></slot>
  </component>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";

// Extend this with more variants as the need arises
export enum TypographyVariant {
  Heading3 = "h3",
  Heading5 = "h5",
  Body1 = "body1",
  Subtitle1 = "subtitle1",
}

export enum TypographyColor {
  Primary = "primary",
  Secondary = "secondary",
}

@Component
export default class Typography extends Vue {
  @Prop({
    default: "p",
  })
  public tag!: keyof HTMLElementTagNameMap;

  @Prop({
    default: TypographyVariant.Body1,
  })
  public variant!: TypographyVariant;

  @Prop({
    default: TypographyColor.Primary,
  })
  public color!: TypographyColor;

  get classObject() {
    return {
      [`Typography__variant-${this.variant}`]: true,
      [`Typography__color-${this.color}`]: true,
    };
  }
}
</script>

<style scoped lang="scss">
@import "../assets/palette.scss";

.Typography {
  font-family: Roboto;

  &__color {
    &-primary {
      color: $typography-black;
    }

    &-secondary {
      color: $typography-white;
    }
  }

  /* 
    Variant determines the size and weight of the typography font.
    Format: Variant__<Variant enum value> 
  */
  &__variant {
    &-h3 {
      font-weight: bold;
      font-size: 32px;
      line-height: 44px;
    }

    &-h5 {
      font-weight: bold;
      font-size: 24px;
      line-height: 34px;
    }

    &-body1 {
      font-weight: normal;
      font-size: 16px;
      line-height: 22px;
    }

    &-subtitle1 {
      font-weight: normal;
      font-size: 12px;
      line-height: 16px;
    }
  }
}
</style>
