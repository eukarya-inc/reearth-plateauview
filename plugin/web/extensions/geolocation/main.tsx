import Geolocation from "@web/extensions/geolocation";
import ReactDOM from "react-dom/client";

(async () => {
  const element = document.getElementById("root");
  if (element) {
    const root = ReactDOM.createRoot(element);
    root.render(<Geolocation />);
  }
})();

export {};
