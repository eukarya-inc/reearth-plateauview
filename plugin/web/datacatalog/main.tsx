import ReactDOM from "react-dom";

import Wrapper from "../components/Wrapper";

// import startMock from "../mock";

// startMock().then(async () => {
//   console.log(await fetch("https://example.com/user/aaa").then(r => r.json()));
// });
document.body.style.maxWidth = "370px";
document.body.style.width = "100%";
document.body.style.height = "100%";
document.body.style.margin = "0";

ReactDOM.render(<Wrapper />, document.body);

export {};
