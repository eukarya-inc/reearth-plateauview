import { Checkbox, Divider, Radio, Slider, Space, Typography } from "antd";
import { memo } from "react";

import "../../styles/style.less";
import icons from "../UI/Icon/icons";

const mapViewData = ["3D Terrain", "3D smooth", "2D"];
const baseMapData = [
  {
    key: "1",
    title: "Aerial photography (Bing)",
    icon: icons.mapBing,
  },
  {
    key: "2",
    title: "National latest photo (seamless)",
    icon: icons.bgmap_tokyo,
  },
  {
    key: "3",
    title: "GSI Maps (light color)",
    icon: icons.bgmap_gsi,
  },
  {
    key: "4",
    title: "Dark Matter",
    icon: icons.bgmap_darkmatter,
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
      <Radio.Group defaultValue="3D Terrain" buttonStyle="solid">
        {mapViewData.map((item) => (
          <Radio.Button key={item} value={item} className={"radioButton"}>
            <Text className="buttonfonts">{item}</Text>
          </Radio.Button>
        ))}
      </Radio.Group>
      <Checkbox>
        <Text>Terrain hides underground features</Text>
      </Checkbox>
      <Divider />
      <Text className="sectionFonts">Base Map</Text>
      <Radio.Group defaultValue="a">
        {baseMapData.map((item) => (
          <Radio.Button
            key={item.key}
            value={item.key}
            className={"imageButton"}
            type="primary"
            style={{
              backgroundImage: `url(${icons.mapBing}) `,
            }}
          />
        ))}
      </Radio.Group>
      <Divider />
      <Text className="sectionFonts">Image optimization</Text>
      <Checkbox>
        <Text> use native device resolution</Text>
      </Checkbox>
      <Divider />
      <Text className="sectionFonts">Raster map quality</Text>
      <Slider marks={marks} />
    </Space>
  );
};
export default memo(MapSettingTab);
