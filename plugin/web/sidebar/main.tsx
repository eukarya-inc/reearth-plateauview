<<<<<<< HEAD
import ReactDOM from "react-dom";

import WdLayout from "../components/layout/WdLayout";
// import startMock from "../mock";

// startMock().then(async () => {
//   console.log(
// await fetch("https://example.com/user/aaa").then((r) => r.json());
//   );
// });
document.body.style.width = "370px";
document.body.style.height = "800px";
document.body.style.margin = "0";

ReactDOM.render(<WdLayout isInsideEditor={true} />, document.body);
=======
import ReactDOM from "react-dom/client";

import startMock from "../mock";

startMock().then(async () => {
  const element = document.getElementById("root");
  if (element) {
    const root = ReactDOM.createRoot(element);
    root.render(<h1>hello</h1>);
  }

  console.log(await fetch("https://example.com/user/aaa").then(r => r.json()));
});
>>>>>>> 98e02ea4a2e5d55e7c2254d513a717f3273d1bb9

export {};
