import { styled } from "@web/theme";

import useHooks from "./hooks";

const LocationWrapper: React.FC = () => {
  const { currentPoint, currentDistance, handlegoogleModalChange, handleTerrainModalChange } =
    useHooks();

  return (
    <ContentWrapper>
      <LocationsWrapper>
        <Text>Lat {Math.round(currentPoint?.lat ?? 0 * 100000) / 100000} ° N</Text>
        <Text>Lon {Math.round(currentPoint?.lng ?? 0 * 100000) / 100000} ° E</Text>
        <DistanceLegend>
          <Text>{currentDistance.label}</Text>
          <UnderLinedText uniteLine={currentDistance.uniteLine} />
        </DistanceLegend>
      </LocationsWrapper>
      <ModalsWrapper>
        <GoogleAnalyticsLink onClick={handlegoogleModalChange}>
          Google Analyticsの利用について
        </GoogleAnalyticsLink>
        <TerrainLink onClick={handleTerrainModalChange}>地形データ</TerrainLink>
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
  gap: 4px;
  background: #dcdcdc;
  height: 100%;
`;

const DistanceLegend = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 0;
`;

const Text = styled.p`
  font-size: 10px;
  margin: 0;
  color: #262626;
`;

const UnderLinedText = styled.div<{ uniteLine?: number }>`
  height: 2px;
  background: #000;
  color: #262626;
  width: ${({ uniteLine }) => uniteLine + "px"};
`;

const LocationsWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 0px;
  gap: 12px;
  height: 50%;
  flex: none;
  order: 0;
  align-self: stretch;
  flex-grow: 0;
`;

const ModalsWrapper = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 0px;
  gap: 12px;
  width: 100%
  height: 50%;
  flex: none;
  order: 1;
  flex-grow: 0;
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
