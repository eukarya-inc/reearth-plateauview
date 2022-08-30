import React from "react";
import "./Wrapper.less";


type Props = {
  className?: string;
};
const Wrapper: React.FC<Props> = () => {
  return (
    <div className={"wrapper"}>
      <h2>Hello World</h2>
    </div>
  );
};
export default Wrapper;
