//import ReactDOM from "react-dom";

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

document.body.style.maxWidth = "370px";
// document.body.style.maxHeight = "1078px";
document.body.style.width = "100%";
document.body.style.height = "100%";
document.body.style.margin = "0";
if (element) {
  const root = ReactDOM.createRoot(element);

  root.render(<WdLayout isInsideEditor={true} />);
}

export {};
