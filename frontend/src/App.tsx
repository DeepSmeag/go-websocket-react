import { ReactNode, useEffect, useRef, useState } from "react";
import "./App.css";

function App() {
  const [connection, setConnection] = useState<"NOT OPEN" | "OPEN">("NOT OPEN");
  const [numbers, setNumbers] = useState<Array<number>>(
    Array.from([1, 2, 3, 4])
  );
  const socketRef = useRef<WebSocket | null>(null);
  function initSocket() {
    if (socketRef.current != null) {
      console.log("Hey, it's still here");
    }
    const socket = new WebSocket("ws://localhost:8080/ws");
    socketRef.current = socket;
    socketRef.current.onopen = () => {
      setConnection("OPEN");
    };
    socketRef.current.onmessage = (e) => {
      const newNum: number = e.data;
      setNumbers((prev) => [...prev, newNum]);
    };
    socketRef.current.onclose = () => {
      console.log("Socket closed");
      setConnection("NOT OPEN");
    };
    socketRef.current.onerror = () => {
      console.error("ERROR WITH WEBSOCKET");
    };
  }

  useEffect(() => {
    initSocket();
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
    };
  }, []);
  return (
    <main>
      <h1 className="font-semibold mb-4">Numbers sent so far</h1>
      <p>Let's see here...</p>
      <h2>Connection: {connection}</h2>
      <p>{numbers.join(" ")}...</p>
      <div className="flex flex-row gap-4 justify-center my-10">
        <ButtonComponent
          text="Start/Resume"
          onClick={() => {
            if (socketRef.current) {
              socketRef.current.send("resume");
            } else {
              initSocket();
            }
          }}
        />
        <ButtonComponent
          text="Pause"
          onClick={() => {
            if (socketRef.current) {
              socketRef.current.send("pause");
            }
          }}
        />
        <ButtonComponent
          text="Close connection?"
          onClick={() => {
            console.log(
              "Not sure if I can temporarily disable and then reconnect, will check out later; for now just disconnecting"
            );
            if (socketRef.current) {
              socketRef.current.close(1000, "Cause I want to");
            }
          }}
        />
      </div>
    </main>
  );
}
function ButtonComponent({
  text,
  onClick,
}: {
  text: string;
  onClick: () => void;
}): ReactNode {
  return (
    <button
      className="p-4 h-10 w-18 flex items-center justify-center"
      onClick={onClick}
    >
      {text}
    </button>
  );
}

export default App;
