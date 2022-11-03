import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { postMsg } from "../../../core/utils";

import DatasetsPage from "./content/DatasetsPage";
import YourDataPage from "./content/YourData";

export type Tab = "dataset" | "your-data";

const DataCatalog: React.FC = () => {
  const [currentTab, changeTabs] = useState<Tab>("dataset");

  const handleDatasetAdd = useCallback(() => {
    const id = `dataset-${Math.floor(Math.random() * 100)}`;
    postMsg({
      action: "msgFromModal",
      payload: {
        dataset: {
          id,
          name: `建物の${id}`,
          hidden: false,
          idealZoom: { lat: 20, lon: 30, height: 100 },
          dataUrl: "www.example.com",
          fields: [],
        },
      },
    });
  }, []);

  const handleClose = useCallback(() => {
    postMsg({ action: "modal-close" });
  }, []);

  return (
    <Wrapper>
      <Header>
        <Title>Data Catalogue</Title>
        <TabsWrapper>
          <Tab selected={currentTab === "dataset"} onClick={() => changeTabs("dataset")}>
            <TabName>PLATEAU Dataset</TabName>
          </Tab>
          <Tab selected={currentTab === "your-data"} onClick={() => changeTabs("your-data")}>
            <TabName>Your Data</TabName>
          </Tab>
        </TabsWrapper>
        <MinimizeButton>
          <Icon size={32} icon="close" onClick={handleClose} />
        </MinimizeButton>
      </Header>
      {currentTab === "your-data" ? (
        <YourDataPage />
      ) : (
        <DatasetsPage onDatasetAdd={handleDatasetAdd} />
      )}
    </Wrapper>
  );
};

export default DataCatalog;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  width: 1155px;
  height: 753px;
  background: #f4f4f4;
`;

const Header = styled.div`
  display: flex;
  background: #dcdcdc;
  height: 48px;
  position: relative;
`;

const Title = styled.p`
  align-self: center;
  font-size: 14px;
  font-weight: 700;
  margin: 0 12px;
  color: #4a4a4a;
`;

const TabsWrapper = styled.div`
  display: flex;
  margin: 8px 10px 0 10px;
  gap: 10px;
`;

const Tab = styled.div<{ selected?: boolean }>`
  border-width: 1px 1px 0px 1px;
  border-style: solid;
  border-color: ${({ selected }) => (selected ? "#f4f4f4" : "#c8c8c8")};
  border-radius: 2px 2px 0px 0px;
  background: ${({ selected }) => (selected ? "#f4f4f4" : "#c8c8c8")};
  color: ${({ selected }) => (selected ? "#00BEBE" : "#898989")};
  padding: 8px 12px;
  cursor: pointer;
`;

const TabName = styled.p`
  margin: 0;
  user-select: none;
`;

const MinimizeButton = styled.button`
  position: absolute;
  right: 0;
  border: none;
  height: 48px;
  width: 48px;
  background: #00bebe;
  cursor: pointer;
  transition: background 0.3s;
`;
