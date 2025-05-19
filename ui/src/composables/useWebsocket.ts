import { ref } from "vue";

const socket = ref<WebSocket | null>(null);

export function useWebSocket() {
  const messages = ref<any[]>([]);

  const connect = () => {
    socket.value = new WebSocket("ws://localhost:8080/ws");

    socket.value.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.value.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      messages.value.push(msg);
      console.log("Received:", msg);
    };

    socket.value.onerror = (e) => {
      console.error("WebSocket error:", e);
    };

    socket.value.onclose = () => {
      console.warn("WebSocket closed");
    };
  };

  const send = (type: string, data: any) => {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      console.warn("WebSocket not connected");
      return;
    }

    socket.value.send(JSON.stringify({ type, data }));
  };

  return {
    connect,
    send,
    messages,
    socket,
  };
}
