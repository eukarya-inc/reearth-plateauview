import React from "react";

type Props = {
  onModalChange: () => void;
};

const TerrainLink: React.FC<Props> = onModalChange => {
  console.log(onModalChange);
  return <div></div>;
};
export default TerrainLink;
