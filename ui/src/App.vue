<script setup lang="ts">
import { useWebSocket } from "./composables/useWebsocket";
import { Button } from "./components/ui/button";
import { Input } from './components/ui/input';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {onMounted, ref} from "vue";
import MonacoEditor from './components/MonacoEditor.vue';

const url = ref<string>("https://rickandmortyapi.com/api/character");
const data = ref<unknown>();

const { connect, send, messages } = useWebSocket();

onMounted(() => {
  connect();
});

async function handleClick() {
  const response = await send<{ body: string, contentType: string, status: number }>("http.request", {
    method: "GET",
    url: url.value,
  });

  data.value = response.body;
}

async function handleRender() {
  const response = await send("render.template.request", {
    template: "hello {{ .name }}",
    variables: {
      name: "world",
    }
  })
}

const httpVerbs = [
  { value: "GET", label: "GET" },
  { value: "POST", label: "POST" },
  { value: "PUT", label: "PUT" },
  { value: "DELETE", label: "DELETE" },
  { value: "PATCH", label: "PATCH" },
]
</script>

<template>
  <div>
    <div class="flex w-full items-center max-w-xl gap-1.5">
      <Select>
        <SelectTrigger class="w-[180px]">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
            <SelectItem v-for="httpVerb in httpVerbs" :key="httpVerb.value" :value="httpVerb.value">
                {{ httpVerb.label }}
            </SelectItem>
        </SelectContent>
      </Select>
      <Input v-model="url" />
      <Button @click="handleClick">
        Send
      </Button>
    </div>
    <MonacoEditor v-model:model-value="data"/>
    <ul>
      <li>
        <div v-for="(message, index) in messages" :key="index">
          <pre>{{ message }}</pre>
        </div>
      </li>
    </ul>
  </div>
</template>
