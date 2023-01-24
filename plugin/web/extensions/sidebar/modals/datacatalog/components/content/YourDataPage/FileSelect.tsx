import { CatalogItem } from "@web/extensions/sidebar/core/processCatalog";
import { Tabs } from "@web/sharedComponents";
import { styled } from "@web/theme";

import LocalDataTab from "./LocalDataTab";
import WebDataTab from "./WebDataTab";

export type Props = {
  onOpenDetails?: (data?: CatalogItem) => void;
};

const FileSelectPane: React.FC<Props> = ({ onOpenDetails }) => {
  const options = [
    {
      value: "auto",
      label: "Auto-detect (Recommended)",
    },
    {
      value: "geojson",
      label: "GeoJSON",
    },
    {
      value: "kml",
      label: "KML or KMZ",
    },
    {
      value: "csv",
      label: "CSV",
    },
    {
      value: "czml",
      label: "CZML",
    },
    {
      value: "gpx",
      label: "GPX",
    },
    {
      value: "json",
      label: "JSON",
    },
    {
      value: "georss",
      label: "GeoRSS",
    },
    {
      value: "gltf",
      label: "GLTF",
    },
    {
      value: "shapefile",
      label: "ShapeFile (zip)",
    },
  ];

  return (
    <Wrapper>
      <Tabs defaultActiveKey="local" style={{ marginBottom: "12px" }}>
        <Tabs.TabPane tab="Add Local Data" key="local">
          <LocalDataTab options={options} onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
        <Tabs.TabPane tab="Add Web Data" key="web">
          <WebDataTab options={options} onOpenDetails={onOpenDetails} />
        </Tabs.TabPane>
      </Tabs>
    </Wrapper>
  );
};

export default FileSelectPane;

const Wrapper = styled.div`
  padding: 24px 12px;
`;
