import { Checkbox, Button, Radio, Row } from "@web/extensions/sharedComponents";
import mapBing from "@web/extensions/sidebar/assets/bgmap_bing.png";
import bgmap_darkmatter from "@web/extensions/sidebar/assets/bgmap_darkmatter.png";
import bgmap_gsi from "@web/extensions/sidebar/assets/bgmap_gsi.png";
import bgmap_tokyo from "@web/extensions/sidebar/assets/bgmap_tokyo.png";
import { styled } from "@web/theme";
import { memo, useState } from "react";

import CommonPage from "./CommonPage";

type TileSelection = "tokyo" | "bing" | "gsi" | "dark-matter";

type BaseMapData = {
  key: TileSelection;
  title: string;
  icon: string;
};

export function postMsg(act: string, payload?: any) {
  parent.postMessage(
    {
      act,
      payload,
    },
    "*",
  );
}

const MapSettings: React.FC = () => {
  const [currentTile, selectTile] = useState<TileSelection>("tokyo");

  // const [currentMaps, setMaps] = useState<[] | undefined>();

  const mapViewData = ["3D Terrain", "3D smooth", "2D"];

  const baseMapData: BaseMapData[] = [
    {
      key: "tokyo",
      title: "National latest photo (seamless)",
      icon: bgmap_tokyo,
    },
    {
      key: "bing",
      title: "Aerial photography (Bing)",
      icon: mapBing,
    },
    {
      key: "gsi",
      title: "GSI Maps (light color)",
      icon: bgmap_gsi,
    },
    {
      key: "dark-matter",
      title: "Dark Matter",
      icon: bgmap_darkmatter,
    },
  ];

  // useEffect(() => {
  //   addEventListener("message", (msg: any) => {
  //     if (msg.source !== parent) return;

  //     try {
  //       const data = typeof msg.data === "string" ? JSON.parse(msg.data) : msg.data;
  //       setMaps(data);
  //       // eslint-disable-next-line no-empty
  //     } catch (error) {}
  //   });
  //   postMsg("getTiles");
  // }, []);

  return (
    <CommonPage title="Map settings">
      <>
        <SubTitle>Map View</SubTitle>
        <MapViewSection>
          <Radio.Group defaultValue="3D Terrain" buttonStyle="solid">
            {mapViewData.map(item => (
              <MapViewButton key={item} value={item} type="primary">
                <Text style={{ color: " #FFFFFF" }}>{item}</Text>
              </MapViewButton>
            ))}
          </Radio.Group>
          <Checkbox>
            <Text>Terrain hides underground features</Text>
          </Checkbox>
        </MapViewSection>
      </>
      <>
        <Title>Base Map</Title>
        <BaseMapSection>
          <Radio.Group defaultValue={currentTile} onChange={e => selectTile(e.target.value)}>
            {baseMapData.map(item => (
              <ImageButton
                key={item.key}
                value={item.key}
                type="default"
                style={{
                  backgroundImage: "url(" + item.icon + ")",
                  backgroundSize: "cover",
                  backgroundRepeat: "no-repeat",
                }}
              />
            ))}
          </Radio.Group>
        </BaseMapSection>
      </>
    </CommonPage>
  );
};

export default memo(MapSettings);

const Title = styled.p`
  font-size: 16px;
`;

const SubTitle = styled.p`
  font-size: 14px;
`;

const Text = styled.p``;

const MapViewSection = styled(Row)`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  padding: 0px;
  gap: 16px;
  width: 296px;
  height: 103px;
`;

const MapViewButton = styled(Button)`
  width: 91px;
  height: 29px;
  border-radius: 4px;
  border-color: #d1d1d1;
  background: #d1d1d1;
  margin: 0px 0px 6px 6px;
  padding: 4px 8px;
`;

const BaseMapSection = styled(Row)`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  padding: 0px;
  gap: 8px;
  width: 326px;
  height: 86px;
`;

const ImageButton = styled(Radio.Button)`
  margin: 0px 0px 12px 12px;
  border-radius: 4px;
  background: #d1d1d1;
  padding: 4px 8px;
  width: 64px;
  height: 64px;
`;
