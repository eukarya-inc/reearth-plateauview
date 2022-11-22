import { styled } from "@web/theme";

import useHooks from "./hooks";

const LocationWrapper: React.FC = () => {
  const { currentPoint, currentDistance, handleGoogleModalOpen, handleTerrainModalOpen } =
    useHooks();

  return (
    <ContentWrapper>
      <LocationsWrapper>
        <Text>Lat {currentPoint?.lat?.toFixed(5)} ° N</Text>
        <Text>Lon {currentPoint?.lng?.toFixed(5)} ° E</Text>
        <DistanceLegend>
          <DistanceLegendLabel>{currentDistance.label}</DistanceLegendLabel>
          <Line uniteLine={currentDistance.uniteLine} />
        </DistanceLegend>
      </LocationsWrapper>
      <ModalsWrapper>
        <GoogleAnalyticsLink onClick={handleGoogleModalOpen}>
          Google Analyticsの利用について
        </GoogleAnalyticsLink>
        <TerrainLink onClick={handleTerrainModalOpen}>地形データ</TerrainLink>
      </ModalsWrapper>
    </ContentWrapper>
  );
};

export default LocationWrapper;

const ContentWrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  padding: 4px 12px;
  background: #dcdcdc;
  height: 100%;
  width: 100%;
`;

const LocationsWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  padding: 0px;
  gap: 35px;
`;

const ModalsWrapper = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
`;

const DistanceLegend = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 0;
`;

const Text = styled.p`
  font-size: 10px;
  width: 84px;
  margin: 0;
  color: #262626;
`;

const DistanceLegendLabel = styled.p`
  font-size: 10px;
  margin: 0;
  color: #262626;
`;

const Line = styled.div<{ uniteLine?: number }>`
  height: 1px;
  background: #000;
  color: #262626;
  width: ${({ uniteLine }) => uniteLine + "px"};
`;

const GoogleAnalyticsLink = styled.a`
  font-size: 10px;
  color: #434343;
  text-decoration-line: underline;
`;

const TerrainLink = styled.a`
  font-size: 10px;
  color: #434343;
  text-decoration-line: underline;
`;
