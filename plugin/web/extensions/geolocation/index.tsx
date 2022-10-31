import { styled } from "@web/theme";
import React from "react";

type Props = {};
const Geolocation: React.FC<Props> = () => {
  return <ContentWrapper></ContentWrapper>;
};
export default Geolocation;
const ContentWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: flex-start;
  padding: 5px 12px;
  gap: 2px;

  position: relative;
  width: 390px;
  height: 24px;
`;
