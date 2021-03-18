<template>
  <div class="hello">
    <h2>{{ uuid }}</h2>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { io } from "socket.io-client";

@Component
export default class HelloWorld extends Vue {
  @Prop() private msg!: string;

  private socket = new WebSocket("ws://localhost:8000/canvas/client");
  private uuid: string = "";

  public mounted() {
    this.register();
  }

  private async register() {
    const response = await fetch("http://localhost:8000/auth", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: "General Grievous" }),
    });

    if (response.ok) {
      const result = await response.json();

      console.log(result);
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
