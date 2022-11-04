import React from "react";

type Props = {
  onModalChange: () => void;
};

const GoogleAnalyticstLink: React.FC<Props> = onModalChange => {
  console.log(onModalChange);
  return <div></div>;
};
export default GoogleAnalyticstLink;
