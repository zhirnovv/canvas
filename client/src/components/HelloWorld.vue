<template>
  <div class="Container">
    <div class="Navbar">
      <Typography :variant="'h3'">Drawing Board</Typography>
      <Typography :variant="'h5'">{{ authorizedUsername }}</Typography>
    </div>
    <DrawingBoard :isAutenticated="isAuthorized" />
    <Modal :isOpen="isUnauthorized">
      <template v-slot:header>
        <Typography :variant="'h3'">Please enter your name</Typography>
      </template>
      <template v-slot:body>
        <Input
          :placeholder="'alphanumeric username'"
          :disabled="isLoading"
          :error="usernameInputError"
          v-model="username"
          class="Auth-modal-input"
        />
      </template>
      <template v-slot:footer>
        <Button
          :disabled="isLoading"
          @click="register"
          class="Auth-modal-submit-button"
        >
          <Typography :color="'secondary'">Submit</Typography>
        </Button>
      </template>
    </Modal>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import DrawingBoard from "./Canvas.vue";
import Button from "../molecules/Button.vue";
import Modal from "../molecules/Modal.vue";
import Input from "../molecules/Input.vue";
import Typography from "../atoms/Typography.vue";

export interface ApiResponse {
  meta: {
    type?: string;
    instance: string;
    status: number;
  };
}

export interface ApiSuccessResponse<Data> extends ApiResponse {
  data: Data;
}

export interface ApiErrorResponse extends ApiResponse {
  error: {
    title: string;
    detail: string;
  };
}

enum AuthStatus {
  Initial = "initial",
  Authorized = "authorized",
  Unauthorized = "unauthorized",
}

@Component({
  components: {
    DrawingBoard,
    Modal,
    Typography,
    Button,
    Input,
  },
})
export default class HelloWorld extends Vue {
  public username = "";
  public usernameInputError = "";
  public isLoading = false;
  public authorizedUsername = "";

  private authStatus = AuthStatus.Initial;

  get isUnauthorized() {
    return this.authStatus === AuthStatus.Unauthorized;
  }

  get isAuthorized() {
    return this.authStatus === AuthStatus.Authorized;
  }

  mounted() {
    this.authorize();
  }

  public async register() {
    this.isLoading = true;
    console.log(this.username);

    const response = await fetch("http://dev.domain.com:8000/auth", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: this.username }),
    });

    if (response.ok) {
      this.authorize();
    } else {
      const result = (await response.json()) as ApiErrorResponse;
      this.usernameInputError = result.error.detail;
    }

    this.isLoading = false;
  }

  private async authorize() {
    this.isLoading = true;

    const response = await fetch("http://dev.domain.com:8000/auth", {
      method: "GET",
      credentials: "include",
    });

    if (response.ok) {
      const result = (await response.json()) as ApiSuccessResponse<{
        user: { name: string };
      }>;

      this.authorizedUsername = result.data.user.name;
      this.authStatus = AuthStatus.Authorized;
    } else {
      this.authStatus = AuthStatus.Unauthorized;
    }

    this.isLoading = false;
  }
}
</script>

<style scoped lang="scss">
@import "../assets/palette.scss";

.Container {
  height: 100vh;
  display: flex;
  flex-flow: column nowrap;
}

.Navbar {
  width: 100%;
  padding: 16px;
  display: flex;
  flex-flow: row nowrap;
  justify-content: space-between;
  align-items: center;

  background: white;
  border-radius: 5px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.25);
}
</style>
