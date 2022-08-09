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

export {};
