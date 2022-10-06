import { Checkbox, Row } from "@web/extensions/sharedComponents";
import mapBing from "@web/extensions/sidebar/assets/bgmap_bing.png";
import bgmap_darkmatter from "@web/extensions/sidebar/assets/bgmap_darkmatter.png";
import bgmap_gsi from "@web/extensions/sidebar/assets/bgmap_gsi.png";
import bgmap_tokyo from "@web/extensions/sidebar/assets/bgmap_tokyo.png";
import { styled } from "@web/theme";
import { memo, useCallback, useState } from "react";

import { API } from "../../types";

import CommonPage from "./CommonPage";

type TileSelection = "tokyo" | "bing" | "gsi" | "dark-matter";

type ViewSelection = "3d-terrain" | "3d-smooth" | "2d";

type BaseMapData = {
  key: TileSelection;
  url: string;
  title?: string;
  icon?: string;
};

type MapViewData = {
  key: ViewSelection;
  title: string;
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

const mapViewData: MapViewData[] = [
  {
    key: "3d-terrain",
    title: "3D地形",
  },
  { key: "3d-smooth", title: "3D地形なし" },
  { key: "2d", title: "2D" },
];

const baseMapData: BaseMapData[] = [
  {
    key: "tokyo",
    title: "National latest photo (seamless)",
    icon: bgmap_tokyo,
    url: "https://cyberjapandata.gsi.go.jp/xyz/pale/{z}/{x}/{y}.png",
  },
  {
    key: "bing",
    title: "Aerial photography (Bing)",
    icon: mapBing,
    url: "https://cyberjapandata.gsi.go.jp/xyz/std/{z}/{x}/{y}.png",
  },
  {
    key: "gsi",
    title: "GSI Maps (light color)",
    icon: bgmap_gsi,
    url: "https://cyberjapandata.gsi.go.jp/xyz/english/{z}/{x}/{y}.png",
  },
  {
    key: "dark-matter",
    title: "Dark Matter",
    icon: bgmap_darkmatter,
    url: "https://cyberjapandata.gsi.go.jp/xyz/lcm25k_2012/{z}/{x}/{y}.png",
  },
];

const MapSettings: React.FC = () => {
  const [currentTile, selectTile] = useState<BaseMapData>(baseMapData[0]);
  const [currentView, selectView] = useState<ViewSelection>("3d-terrain");
  const [currentChanges, updateCurrentChange] = useState<API>({
    default: {
      terrain: true,
      sceneMode: "3d",
    },
    tiles: [
      {
        id: baseMapData[0].key,
        tile_url: baseMapData[0].url,
      },
    ],
  });

  // useEffect(() => {
  //   addEventListener("message", (msg: any) => {
  //     if (msg.source !== parent) return;

  //     try {
  //       // const data = typeof msg.data === "string" ? JSON.parse(msg.data) : msg.data;
  //       // setMaps(data);
  //       // postMsg("setTile", data)
  //     } catch (error) {
  //       console.log("error: ", error);
  //     }
  //   });

  //   postMsg("getTiles");
  // }, []);

  const handleViewChange = useCallback((view: ViewSelection) => {
    selectView(view);
    postMsg("setView", view);
  }, []);

  const handleTileChange = useCallback((tile: BaseMapData) => {
    selectTile(tile);
    postMsg("setTile", tile.url);
  }, []);

  return (
    <CommonPage title="マップ設定">
      <>
        <SubTitle>マップビュー</SubTitle>
        <Section>
          <ViewWrapper>
            {mapViewData.map(({ key, title }) => (
              <MapViewButton
                key={key}
                value={key}
                selected={currentView === key}
                onClick={() => handleViewChange(key)}>
                <Text style={{ color: " #FFFFFF" }}>{title}</Text>
              </MapViewButton>
            ))}
          </ViewWrapper>
          <Checkbox>
            <Text>地下を隠す</Text>
          </Checkbox>
        </Section>
      </>
      <>
        <Title>ベースマップ</Title>
        <Section>
          <MapWrapper>
            {baseMapData.map(item => (
              <ImageButton
                key={item.key}
                selected={item.key === currentTile.key}
                onClick={() => handleTileChange(item)}
                style={{
                  backgroundImage: "url(" + item.icon + ")",
                  backgroundSize: "cover",
                  backgroundRepeat: "no-repeat",
                }}
              />
            ))}
          </MapWrapper>
        </Section>
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

const Section = styled(Row)`
  gap: 16px;
`;

const ViewWrapper = styled.div`
  display: flex;
  gap: 12px;
  width: 100%;
`;

const MapViewButton = styled.button<{ selected?: boolean }>`
  width: 91px;
  height: 29px;
  background: ${({ selected }) => (selected ? "#00bebe" : "#d1d1d1")};
  border-radius: 4px;
  border: none;
  padding: 4px 8px;
  cursor: pointer;
  transition: background 0.3s;

  :hover {
    background: #00bebe;
  }
`;

const MapWrapper = styled.div`
  display: flex;
  justify-content: start;
  gap: 8px;
  width: 100%;
`;

const ImageButton = styled.div<{ selected?: boolean }>`
  height: 64px;
  width: 64px;
  background: #d1d1d1;
  border: 2px solid ${({ selected }) => (selected ? "#00bebe" : "#d1d1d1")};
  border-radius: 2px;
  padding: 4px 8px;
  cursor: pointer;
`;
