import { Checkbox, Divider, Radio, Row, Slider, Space, Typography } from "antd";
import { memo } from "react";

import "../../styles/style.less";

import Bgmap_darkmatter from "../UI/Icon/Icons/bgmap_darkmatter.svg";
import Bgmap_gsi from "../UI/Icon/Icons/bgmap_gsi.svg";
import Bgmap_tokyo from "../UI/Icon/Icons/bgmap_tokyo.svg";
import MapBing from "../UI/Icon/Icons/mapBing.svg";

const mapViewData = ["3D Terrain", "3D smooth", "2D"];
const baseMapData = [
  {
    key: "1",
    title: "Aerial photography (Bing)",
    icon: MapBing,
  },
  {
    key: "2",
    title: "National latest photo (seamless)",
    icon: Bgmap_tokyo,
  },
  {
    key: "3",
    title: "GSI Maps (light color)",
    icon: Bgmap_gsi,
  },
  {
    key: "4",
    title: "Dark Matter",
    icon: Bgmap_darkmatter,
  },
];

const MapSettingTab: React.FC = () => {
  const { Text } = Typography;
  const marks = {
    0: {
      style: {
        color: "#f50",
      },
      label: <Text>quality</Text>,
    },
    100: {
      style: {
        color: "#f50",
      },
      label: <Text>performance</Text>,
    },
  };

  return (
    <Space direction="vertical">
      <Text className="sectionFonts">Map View</Text>
      <Row className="mapViewSection">
        <Radio.Group defaultValue="3D Terrain" buttonStyle="solid">
          {mapViewData.map(item => (
            <Radio.Button key={item} value={item} className={"radioButton"}>
              <Text className="buttonfonts">{item}</Text>
            </Radio.Button>
          ))}
        </Radio.Group>
        <Checkbox>
          <Text>Terrain hides underground features</Text>
        </Checkbox>
      </Row>
      <Divider />
      <Text className="sectionFonts">Base Map</Text>
      <Row className="baseMapSection">
        <Radio.Group defaultValue="a">
          {baseMapData.map(item => (
            <Radio.Button
              key={item.key}
              value={item.key}
              className={"imageButton"}
              type="primary"
              style={{
                // backgroundImage: `url(${item.icon}) `,
                backgroundImage: "url(" + item.icon + ")",
                backgroundSize: "cover",
                backgroundRepeat: "no-repeat",
              }}
            />
          ))}
        </Radio.Group>
      </Row>
      <Divider />
      <Text className="sectionFonts">Image optimization</Text>
      <Checkbox>
        <Text> use native device resolution</Text>
      </Checkbox>
      <Divider />
      <Text className="sectionFonts">Raster map quality</Text>
      <Slider marks={marks} className="rasterMapSection" />
    </Space>
  );
};
export default memo(MapSettingTab);
