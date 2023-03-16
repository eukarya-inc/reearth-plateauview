import SelectionComponent from "@web/extensions/sidebar/core/components/content/Selection";
import { DataCatalogItem, BuildingSearch } from "@web/extensions/sidebar/core/types";
import { ReearthApi } from "@web/extensions/sidebar/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useEffect } from "react";

import PopupItem from "../sharedComponents/PopupItem";

type Props = {
  selectedDatasets: DataCatalogItem[];
  savingDataset: boolean;
  buildingSearch?: BuildingSearch;
  onDatasetSave: (datasetId: string) => void;
  onDatasetUpdate: (updatedDataset: DataCatalogItem) => void;
  onDatasetRemove: (id: string) => void;
  onDatasetRemoveAll: () => void;
  onBuildingSearch: (id: string) => void;
  onSceneUpdate: (updatedProperties: Partial<ReearthApi>) => void;
};

const Selection: React.FC<Props> = ({
  selectedDatasets,
  savingDataset,
  buildingSearch,
  onDatasetSave,
  onDatasetUpdate,
  onDatasetRemove,
  onDatasetRemoveAll,
  onBuildingSearch,
  onSceneUpdate,
}) => {
  useEffect(() => {
    postMsg({ action: "extendPopup" });
  }, []);

  return (
    <Wrapper>
      <PopupItem>
        <Title>Data Style Settings</Title>
      </PopupItem>
      <SelectionComponent
        selectedDatasets={selectedDatasets}
        savingDataset={savingDataset}
        buildingSearch={buildingSearch}
        onDatasetSave={onDatasetSave}
        onDatasetUpdate={onDatasetUpdate}
        onDatasetRemove={onDatasetRemove}
        onDatasetRemoveAll={onDatasetRemoveAll}
        onBuildingSearch={onBuildingSearch}
        onSceneUpdate={onSceneUpdate}
      />
    </Wrapper>
  );
};

export default Selection;

const Wrapper = styled.div`
  border-top: 1px solid #d9d9d9;
  height: calc(100% - 47px);
`;

const Title = styled.p`
  margin: 0;
`;
