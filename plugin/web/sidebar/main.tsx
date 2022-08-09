import ReactDOM from "react-dom/client";

import WdLayout from "../components/layout/WdLayout";

//import startMock from "../mock";

// startMock().then(async () => {
//   const element = document.getElementById("root");
//   if (element) {
//     const root = ReactDOM.createRoot(element);
//     root.render(<WdLayout isInsideEditor={false} />);
//   }

//   console.log(await fetch("https://example.com/user/aaa").then(r => r.json()));
// });
const element = document.getElementById("root");
if (element) {
  const root = ReactDOM.createRoot(element);

  root.render(<WdLayout isInsideEditor={false} />);
}

export {};
