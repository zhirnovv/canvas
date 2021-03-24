<template>
  <div>
    <input
      :class="inputClassObject"
      :value="value"
      :placeholder="placeholder"
      @input="$emit('input', $event.target.value)"
    />
    <Typography variant="subtitle1">{{ error }}</Typography>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Box from "../atoms/Box.vue";
import Typography from "../atoms/Typography.vue";

@Component({
  components: {
    Box,
    Typography,
  },
})
export default class Input extends Vue {
  @Prop({
    default: "",
  })
  public error!: string;

  @Prop() public value!: string;
  @Prop() public placeholder!: string;

  get inputClassObject() {
    return ["Input", { ["Input_error"]: this.error !== "" }];
  }
}
</script>

<style scoped lang="scss">
@import "../assets/palette.scss";

.Input {
  padding: 6px 16px;

  font-family: Roboto;
  background: white;
  border: 1px solid white;
  border-radius: 5px;
  box-shadow: 0 0 4px rgba(0, 0, 0, 0.25);

  transition: border 0.2s ease, opacity 0.2s ease, color 0.2s ease;

  outline: none;

  &::placeholder {
    color: rgba(0, 0, 0, 0.4);
  }

  &:focus {
    border: 1px solid black;
  }

  &_error {
    color: $dangerous-red;
    border: 1px solid $dangerous-red;

    &:focus {
      border: 1px solid $dangerous-red;
    }
  }
}
</style>
