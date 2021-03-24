<template>
  <Icon
    :viewBox="'0 0 34 34'"
    :height="height"
    :width="width"
    :iconName="'close'"
  >
    <path
      d="M29.75 17c0-7.04-5.71-12.75-12.75-12.75S4.25 9.96 4.25 17 9.96 29.75 17 29.75 29.75 24.04 29.75 17z"
      stroke="#000"
      stroke-width="2"
      stroke-miterlimit="10"
    />
    <path
      d="M12.75 21.25l8.5-8.5m0 8.5l-8.5-8.5 8.5 8.5z"
      stroke="#000"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </Icon>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Icon from "../../atoms/Icon.vue";

@Component({
  components: {
    Icon,
  },
})
export default class CloseIcon extends Vue {
  @Prop({
    default: 34,
  })
  public width!: number;

  @Prop({
    default: 34,
  })
  public height!: number;
}
</script>
