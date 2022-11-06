import { styled } from "@web/theme";
import React from "react";

import GoogleAnalystLink from "./components/googleAnalyticsLink";
import TerrainLink from "./components/terrainLink";
import useHook from "./hook";

const LocationWrapper: React.FC = () => {
  const { currentPoint, handlegoogleModalChange, handleTerrainModalChange } = useHook();
  return (
    <ContentWrapper>
      <Wrapper1>
        <LatWrapper>{currentPoint?.lat}</LatWrapper>
        <LngWrapper>{currentPoint?.lng}</LngWrapper>
        <DistanceWrapper>{currentPoint?.height}</DistanceWrapper>
      </Wrapper1>
      <Wrapper2>
        <GoogleAnalystLink onModalChange={handlegoogleModalChange} />
        <TerrainLink onModalChange={handleTerrainModalChange} />
      </Wrapper2>
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
const Wrapper1 = styled.div`
  display: flex;
`;
const Wrapper2 = styled.div`
  display: flex;
`;
