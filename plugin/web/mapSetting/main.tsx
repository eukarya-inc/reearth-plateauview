import ReactDOM from "react-dom";

import WdLayout from "../components/layout/WdLayout";
import startMock from "../mock";

startMock().then(async () => {
  console.log(
    await fetch("https://example.com/user/aaa").then((r) => r.json())
  );
});

ReactDOM.render(
  <WdLayout isInsideEditor={true} currentTab={"mapSetting"} />,
  document.body
);

export {};
