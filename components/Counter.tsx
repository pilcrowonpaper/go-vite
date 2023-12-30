import { createSignal } from "solid-js";
import "./Counter.css"

export default function Counter () {
    const [count, setCount] = createSignal(0);
    return (
      <>
        <p>This is a client rendered Solid.js component!</p>
        <button onClick={() => setCount((curr) => curr + 1)} id="count-button">
          Clicks: {count()}
        </button>
      </>
    );
}