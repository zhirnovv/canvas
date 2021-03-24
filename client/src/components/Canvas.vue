<template>
  <div class="Canvas__container" ref="canvasBoundaries">
    <canvas
      class="Canvas"
      ref="canvas"
      @mousedown="startDrawing"
      @mouseup="stopDrawing"
      @mouseleave="stopDrawing"
      @mouseout="stopDrawing"
      @mousemove="draw"
    ></canvas>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import Input from "../molecules/Input.vue";
import Button from "../molecules/Button.vue";
import CloseIcon from "../molecules/icons/CloseIcon.vue";
import Modal from "../molecules/Modal.vue";
import Typography from "../atoms/Typography.vue";

enum MessageType {
  NewLine = "/canvas/add/line",
  All = "/canvas/all",
  Resize = "/canvas/resize",
}

interface CanvasMessage {
  type: MessageType;
  payload: { [key: string]: any };
  createdAt: Date;
}

interface Coordinates {
  x: number;
  y: number;
}

interface AddLineMessage extends CanvasMessage {
  type: MessageType.NewLine;
  payload: {
    lineData: {
      start: Coordinates;
      end: Coordinates;
    };
    stroke: number;
    color: string;
  };
}

interface ResizeCanvasMessage extends CanvasMessage {
  type: MessageType.Resize;
  payload: {};
}

interface AllMessages extends CanvasMessage {
  type: MessageType.All;
  payload: {
    messages: AddCanvasMessage[];
  };
}

type AddCanvasMessage = AddLineMessage;
type AvailableCanvasMessage = AllMessages | AddCanvasMessage;

@Component({
  components: {
    Button,
    Typography,
    Input,
    CloseIcon,
    Modal,
  },
})
export default class DrawingBoard extends Vue {
  @Prop() private isAutenticated!: boolean;

  private socket: WebSocket | null = null;
  private canvas: HTMLCanvasElement | null = null;
  private assignedColor: string = "";
  private drawHistory: AddCanvasMessage[] = [];
  private context: CanvasRenderingContext2D | null = null;
  private isDrawing = false;
  private isMoving = false;
  private currentMousePosition: Coordinates = {
    x: 0,
    y: 0,
  };
  private previousMousePosition: Coordinates | null = null;
  private intervalRef: number = 0;

  mounted() {
    const canvas = this.$refs.canvas as HTMLCanvasElement;

    if (canvas) {
      this.initializeCanvas(canvas);
      this.mountCanvasPoller();
    }
  }

  @Watch("isAutenticated")
  watchIsAuthenticated(value: boolean, oldValue: boolean) {
    if (oldValue === false && value === true) {
      const socket = new WebSocket("ws://dev.domain.com:8000/canvas/client");
      socket.addEventListener("open", (e) => console.log("hello"));
      socket.addEventListener("close", (e) => console.log("goodbye"));
      this.socket = socket;
    }
  }

  @Watch("socket")
  watchSocket() {
    if (this.socket) {
      let assignedColor = sessionStorage.getItem("assignedColor");
      if (!assignedColor) {
        assignedColor = `#${Math.floor(Math.random() * 16777215).toString(16)}`;
        sessionStorage.setItem("assignedColor", assignedColor);
      }
      this.assignedColor = assignedColor;
      this.socket.addEventListener("message", (event) => {
        const message = JSON.parse(event.data) as AllMessages;
        this.parseMessage(message);
      });
    }
  }

  startDrawing() {
    this.isDrawing = true;
  }

  stopDrawing() {
    this.isDrawing = false;
  }

  draw(e: MouseEvent) {
    const canvas = this.$refs.canvas as Element;

    if (canvas) {
      const { left, top, width, height } = canvas.getBoundingClientRect();

      this.currentMousePosition = {
        x: (e.clientX - left) / width,
        y: (e.clientY - top) / height,
      };
      this.isMoving = true;
    }
  }

  resizeCanvas(isInitial: boolean) {
    const canvas = this.$refs.canvas as HTMLCanvasElement;

    if (canvas && this.context) {
      const { height, width } = (this.$refs
        .canvasBoundaries as Element).getBoundingClientRect();
      console.log(height, width);

      if (width >= height) {
        canvas.height = height - 32;
        canvas.width = height - 32;
      } else {
        canvas.height = width - 32;
        canvas.width = width - 32;
      }

      if (!isInitial && this.socket) {
        this.socket.send(
          JSON.stringify({
            type: "/canvas/resize",
            payload: {},
          } as ResizeCanvasMessage)
        );
      }
    }
  }

  renderFromMessage(message: AddCanvasMessage) {
    const canvas = this.$refs.canvas as Element;

    if (this.context && canvas) {
      const { height, width } = canvas.getBoundingClientRect();
      this.drawHistory.push(message);
      switch (message.type) {
        case MessageType.NewLine: {
          if (this.context) {
            const { payload } = message;
            const { lineData, stroke, color } = payload;
            const { start, end } = lineData;

            this.context.beginPath();
            this.context.moveTo(start.x * width, start.y * height);
            this.context.lineTo(end.x * width, end.y * height);
            this.context.lineWidth = stroke;
            this.context.strokeStyle = color;
            this.context.stroke();
          }
          break;
        }
        default:
          break;
      }
    }
  }

  parseMessage(message: AvailableCanvasMessage) {
    switch (message.type) {
      case MessageType.All: {
        message.payload.messages.forEach((message) =>
          this.renderFromMessage(message)
        );
      }
      default:
        this.renderFromMessage(message as AddCanvasMessage);
        break;
    }
  }

  initializeCanvas(canvas: HTMLCanvasElement) {
    this.context = canvas.getContext("2d");

    this.resizeCanvas(true);
    window.addEventListener("resize", () => this.resizeCanvas(false));
  }

  mountCanvasPoller() {
    this.intervalRef = setInterval(() => {
      const {
        isAutenticated,
        socket,
        isDrawing,
        isMoving,
        currentMousePosition,
        previousMousePosition,
      } = this;

      if (
        isAutenticated &&
        socket &&
        isDrawing &&
        isMoving &&
        currentMousePosition &&
        previousMousePosition
      ) {
        socket.send(
          JSON.stringify({
            type: MessageType.NewLine,
            payload: {
              lineData: {
                start: previousMousePosition,
                end: currentMousePosition,
              },
              stroke: 2,
              color: this.assignedColor || "black",
            },
          } as Omit<AddLineMessage, "createdAt">)
        );
        this.isMoving = false;
      }
      this.previousMousePosition = { ...this.currentMousePosition };
    }, 30);
  }

}
</script>

<style scoped lang="scss">
.Canvas__container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-flow: row nowrap;
  justify-content: center;
  align-items: center;
}

.Canvas {
  background: white;
  border-radius: 5px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.25);

  &__disabled {
    background: rgba(0, 0, 0, 0.1);
  }
}

.Auth-modal-input {
  width: 100%;
}

.Auth-modal-submit-button {
  margin: 0 0 0 auto;
}
</style>
