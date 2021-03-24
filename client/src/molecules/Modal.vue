<template>
  <transition name="fade">
    <div class="Modal__backdrop" v-if="isOpen" @click.self="$emit('close')">
      <Box class="Modal">
        <div class="Modal__section Modal__header">
          <slot name="header">
            <Typography :variant="'h3'">bruh</Typography>
          </slot>
        </div>
        <div class="Modal__section">
          <slot name="body"> </slot>
        </div>
        <div class="Modal__section">
          <slot name="footer"></slot>
        </div>
      </Box>
    </div>
  </transition>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import CloseIcon from "../molecules/icons/CloseIcon.vue";
import Typography from "../atoms/Typography.vue";
import Box from "../atoms/Box.vue";

@Component({
  components: {
    Box,
    CloseIcon,
    Typography,
  },
})
export default class Modal extends Vue {
  @Prop() public isOpen!: boolean;
}
</script>

<style scoped lang="scss">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-leave-to,
.fade-enter {
  opacity: 0;
}

.fade-enter-to,
.fade-leave {
  opacity: 1;
}

.Modal__backdrop {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;

  display: flex;
  justify-content: center;
  align-items: center;

  background-color: rgba(0, 0, 0, 0.4);
}

.Modal {
  background-color: white;
  padding: 16px;

  display: flex;
  flex-flow: column nowrap;

  &__section {
    margin-bottom: 32px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  &__header {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
