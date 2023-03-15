import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import PageLayout from "@web/extensions/sidebar/modals/datacatalog/components/content/PageLayout";
import { useCallback, useMemo, useState } from "react";

import { GroupBy } from "../../../api/api";
import { UserDataItem } from "../../../types";
import { TreeTab } from "../../hooks";

import DatasetTree from "./DatasetTree";
import DatasetDetails, { Tag } from "./Details";

export type Props = {
  catalog?: DataCatalogItem[];
  currentTreeTab: TreeTab;
  addedDatasetDataIDs?: string[];
  inEditor?: boolean;
  filter: GroupBy;
  onFilter: (filter: GroupBy) => void;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem) => void;
  onDatasetPublish: (dataID: string, publish: boolean) => void;
  onTreeTabChange: (tab: TreeTab) => void;
};

const DatasetsPage: React.FC<Props> = ({
  catalog,
  currentTreeTab,
  addedDatasetDataIDs,
  inEditor,
  filter,
  onFilter,
  onDatasetAdd,
  onDatasetPublish,
  onTreeTabChange,
}) => {
  const [selectedDatasetID, setDatasetID] = useState<string>();
  const [selectedTags, selectTags] = useState<Tag[]>([]);

  const handleOpenDetails = useCallback((data?: DataCatalogItem) => {
    setDatasetID(data?.dataID);
  }, []);

  const handleTagSelect = useCallback(
    (tag: Tag) =>
      selectTags(tags => {
        const selected = tags.find(selectedTag => selectedTag.name === tag.name)
          ? [...tags.filter(t => t.name !== tag.name)]
          : [...tags, tag];
        selected.length > 0 ? onFilter("tag") : onFilter("city");
        return selected;
      }),
    [onFilter],
  );

  const addDisabled = useCallback(
    (dataID: string) => !!addedDatasetDataIDs?.find(dataID2 => dataID2 === dataID),
    [addedDatasetDataIDs],
  );

  const selectedDataset = useMemo(
    () => catalog?.find(item => item.dataID === selectedDatasetID),
    [catalog, selectedDatasetID],
  );

  return (
    <PageLayout
      left={
        <DatasetTree
          addedDatasetDataIDs={addedDatasetDataIDs}
          selectedDataset={selectedDataset}
          catalog={catalog}
          currentTreeTab={currentTreeTab}
          selectedTags={selectedTags}
          filter={filter}
          addDisabled={addDisabled}
          onTagSelect={handleTagSelect}
          onTreeTabChange={onTreeTabChange}
          onOpenDetails={handleOpenDetails}
          onDatasetAdd={onDatasetAdd}
        />
      }
      right={
        <DatasetDetails
          dataset={selectedDataset}
          inEditor={inEditor}
          addDisabled={addDisabled}
          onTagSelect={handleTagSelect}
          onDatasetAdd={onDatasetAdd}
          onDatasetPublish={onDatasetPublish}
        />
      }
    />
  );
};

export default DatasetsPage;
