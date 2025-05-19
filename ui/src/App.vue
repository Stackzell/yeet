<script setup lang="ts">
import { useWebSocket } from "./composables/useWebsocket";
import { Button } from "./components/ui/button";
import {onMounted} from "vue";

const { connect, send, messages } = useWebSocket();

onMounted(() => {
  connect();
});

function handleClick() {
  send("http.request", {
    method: "GET",
    url: "{{ .baseUrl }}/api/character",
    variables: {
      baseUrl: "https://blah.com",
    }
  });
}
</script>

<template>
  <div>
    <Button>Testing Shad</Button>
    <Button @click="handleClick" class="p-2"> Send WebSocket Message </Button>
    <ul>
      <li>
        <div v-for="(message, index) in messages" :key="index">
          <pre>{{ message }}</pre>
        </div>
      </li>
    </ul>
  </div>
</template>
