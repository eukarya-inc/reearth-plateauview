import { styled } from "@web/theme";
import React from "react";

type Props = {};
const LocationWrapper: React.FC<Props> = () => {
  return (
    <ContentWrapper>
      <LatWrapper></LatWrapper>
      <LngWrapper></LngWrapper>
      <DistanceWrapper></DistanceWrapper>
    </ContentWrapper>
  );
};
export default LocationWrapper;
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
const LatWrapper = styled.div``;
const LngWrapper = styled.div``;
const DistanceWrapper = styled.div``;
