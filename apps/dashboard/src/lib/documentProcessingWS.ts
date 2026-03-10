import useWebSocket, { ReadyState } from "react-use-websocket";
import { useEffect, useState } from "react";

import { ProcessStatus } from "./documentProcessing";

const API_KEY: string = process.env.NEXT_PUBLIC_API_AUTH_TOKEN || "";
const API_URL: string = process.env.NEXT_PUBLIC_API_URL_FRONT || "";
const WS_URL = `${API_URL.replace(/^http/, "ws")}`;

// WebSocket version - Hook for real-time process updates
export const useGetAllProcessWebSocket = () => {
  const [processes, setProcesses] = useState<ProcessStatus[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  const WS_URL_WITH_TOKEN = `${WS_URL}/ws/status?token=${API_KEY}`;

  const { lastMessage, readyState } = useWebSocket(
    WS_URL_WITH_TOKEN,
    {
      queryParams: {
        token: API_KEY,
      },
      shouldReconnect: () => true, // Auto-reconnect on disconnect
      reconnectAttempts: 10,
      reconnectInterval: 3000,
    },
    !!API_URL && !!API_KEY // Only connect if API_URL and API_KEY are available
  );

  useEffect(() => {
    setIsConnected(readyState === ReadyState.OPEN);
  }, [readyState]);

  useEffect(() => {
    if (lastMessage !== null) {
      try {
        const data = JSON.parse(lastMessage.data);
        if (Array.isArray(data)) {
          setProcesses(data);
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", error);
      }
    }
  }, [lastMessage]);

  return {
    processes,
    isConnected,
    readyState,
    connectionStatus: {
      [ReadyState.CONNECTING]: "Connecting",
      [ReadyState.OPEN]: "Connected",
      [ReadyState.CLOSING]: "Closing",
      [ReadyState.CLOSED]: "Closed",
      [ReadyState.UNINSTANTIATED]: "Uninstantiated",
    }[readyState],
  };
};
