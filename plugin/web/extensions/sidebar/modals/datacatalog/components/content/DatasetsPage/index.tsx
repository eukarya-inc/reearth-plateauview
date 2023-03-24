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
  selectedDatasetID?: string;
  selectedItem?: DataCatalogItem;
  expandedFolders?: { id?: string; name?: string }[];
  searchTerm: string;
  setExpandedFolders?: React.Dispatch<React.SetStateAction<{ id?: string; name?: string }[]>>;
  filter: GroupBy;
  onSearch: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSelect?: (item?: DataCatalogItem) => void;
  onOpenDetails?: (data?: DataCatalogItem) => void;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem, keepModalOpen?: boolean) => void;
  onFilter: (filter: GroupBy) => void;
  onDatasetPublish: (dataID: string, publish: boolean) => void;
  onTreeTabChange: (tab: TreeTab) => void;
};

const DatasetsPage: React.FC<Props> = ({
  catalog,
  currentTreeTab,
  addedDatasetDataIDs,
  inEditor,
  selectedDatasetID,
  selectedItem,
  expandedFolders,
  searchTerm,
  filter,
  setExpandedFolders,
  onSearch,
  onSelect,
  onOpenDetails,
  onFilter,
  onDatasetAdd,
  onDatasetPublish,
  onTreeTabChange,
}) => {
  const [selectedTags, selectTags] = useState<Tag[]>([]);
  const [filter, setFilter] = useState<GroupBy>("city");

  const handleFilter = useCallback((filter: GroupBy) => {
    setFilter(filter);
  }, []);

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
          selectedItem={selectedItem}
          expandedFolders={expandedFolders}
          searchTerm={searchTerm}
          setExpandedFolders={setExpandedFolders}
          onSearch={onSearch}
          onSelect={onSelect}
          addDisabled={addDisabled}
          onTagSelect={handleTagSelect}
          onOpenDetails={handleOpenDetails}
          onTreeTabChange={onTreeTabChange}
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
