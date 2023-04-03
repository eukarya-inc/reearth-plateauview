import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import PageLayout from "@web/extensions/sidebar/modals/datacatalog/components/content/PageLayout";
import { useCallback, useMemo, useState } from "react";

import { DataCatalogGroup, GroupBy } from "../../../api/api";
import { UserDataItem } from "../../../types";

import DatasetTree from "./DatasetTree";
import DatasetDetails, { Tag } from "./Details";

export type Props = {
  catalog?: DataCatalogItem[];
  addedDatasetDataIDs?: string[];
  inEditor?: boolean;
  selectedItem?: DataCatalogItem | DataCatalogGroup;
  expandedFolders?: { id?: string; name?: string }[];
  searchTerm: string;
  setExpandedFolders?: React.Dispatch<React.SetStateAction<{ id?: string; name?: string }[]>>;
  onSearch: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onSelect?: (item?: DataCatalogItem | DataCatalogGroup) => void;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem, keepModalOpen?: boolean) => void;
  onDatasetPublish: (dataID: string, publish: boolean) => void;
};

const DatasetsPage: React.FC<Props> = ({
  catalog,
  addedDatasetDataIDs,
  inEditor,
  selectedItem,
  expandedFolders,
  searchTerm,
  setExpandedFolders,
  onSearch,
  onSelect,
  onDatasetAdd,
  onDatasetPublish,
}) => {
  const [selectedTags, selectTags] = useState<Tag[]>([]);
  const [filter, setFilter] = useState<GroupBy>("city");

  const handleFilter = useCallback((filter: GroupBy) => {
    setFilter(filter);
  }, []);

  const handleTagSelect = useCallback(
    (tag: Tag) =>
      selectTags(tags => {
        const selected = tags.find(selectedTag => selectedTag.name === tag.name)
          ? [...tags.filter(t => t.name !== tag.name)]
          : [...tags, tag];
        selected.length > 0 ? handleFilter("tag") : handleFilter("city");
        return selected;
      }),
    [handleFilter],
  );

  const addDisabled = useCallback(
    (dataID: string) => !!addedDatasetDataIDs?.find(dataID2 => dataID2 === dataID),
    [addedDatasetDataIDs],
  );

  const selectedDataset = useMemo(() => {
    return selectedItem && "dataID" in selectedItem
      ? catalog?.find(item => item.dataID === selectedItem.dataID)
      : selectedItem;
  }, [catalog, selectedItem]);

  return (
    <PageLayout
      left={
        <DatasetTree
          addedDatasetDataIDs={addedDatasetDataIDs}
          catalog={catalog}
          selectedTags={selectedTags}
          filter={filter}
          selectedItem={selectedItem}
          expandedFolders={expandedFolders}
          searchTerm={searchTerm}
          setExpandedFolders={setExpandedFolders}
          onSearch={onSearch}
          onSelect={onSelect}
          addDisabled={addDisabled}
          onFilter={handleFilter}
          onTagSelect={handleTagSelect}
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
