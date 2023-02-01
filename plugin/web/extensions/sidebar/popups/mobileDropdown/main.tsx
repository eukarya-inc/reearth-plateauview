import ReactDOM from "react-dom/client";

import MobileDropdown from ".";

(async () => {
  const element = document.getElementById("root");
  const isMobile = true; // TODO: get isMobile value
  if (element) {
    const root = ReactDOM.createRoot(element);
    root.render(<MobileDropdown isMobile={isMobile} />);
  }
})();

export {};
