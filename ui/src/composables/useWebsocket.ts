import { ref } from "vue";

export interface Deferred<T> {
  resolve: (value: T | PromiseLike<T>) => void;
  reject: (reason?: any) => void;
};

export interface WebSocketEnvelope<T> {
  type: string;
  correlationId: string | null;
  data: T;
}

const socket = ref<WebSocket | null>(null);

const requests: Record<string, Deferred<any>> = {}

function sendRequest<TResponse, TRequest = any>(type: string, request: TRequest): Promise<TResponse> {
  return new Promise<TResponse>((resolve, reject) => {
    const correlationId = crypto.randomUUID();
    requests[correlationId] = { resolve, reject };

    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      console.warn("WebSocket not connected");
      return;
    }

    socket.value.send(JSON.stringify({ type, correlationId, data: request }));
  })
}

export function useWebSocket() {
  const messages = ref<any[]>([]);

  const connect = () => {
    socket.value = new WebSocket("ws://localhost:8080/ws");

    socket.value.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.value.onmessage = (event) => {
      const envelope: WebSocketEnvelope<any> = JSON.parse(event.data);
      if (!envelope?.correlationId) {
        console.warn("Invalid message format:", event.data);
        return;
      }

      const request = requests[envelope.correlationId];
      if (!request) {
        console.warn("No matching request for correlationId:", envelope.correlationId);
        return;
      }

      delete requests[envelope.correlationId];
      if (envelope.type === "error") {
        request.reject(envelope.data);
        return;
      }

      request.resolve(envelope.data);
    };

    socket.value.onerror = (e) => {
      console.error("WebSocket error:", e);
    };

    socket.value.onclose = () => {
      console.warn("WebSocket closed");
    };
  };

  return {
    send: sendRequest,
    connect,
    messages,
    socket,
  };
}
