import WelcomePage from "@web/extensions/welcomePage/core";
import ReactDOM from "react-dom/client";

(async () => {
  const element = document.getElementById("root");
  if (element) {
    const root = ReactDOM.createRoot(element);
    root.render(<WelcomePage />);
  }
})();

export {};
