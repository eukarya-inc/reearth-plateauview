import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import PageLayout from "@web/extensions/sidebar/modals/datacatalog/components/content/PageLayout";
import { useCallback, useState } from "react";

import { GroupBy } from "../../../api/api";
import { UserDataItem } from "../../../types";

import DatasetTree from "./DatasetTree";
import DatasetDetails, { Tag } from "./Details";

export type Props = {
  catalog?: DataCatalogItem[];
  addedDatasetDataIDs?: string[];
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem) => void;
};

const DatasetsPage: React.FC<Props> = ({ catalog, addedDatasetDataIDs, onDatasetAdd }) => {
  const [selectedDataset, setDataset] = useState<DataCatalogItem>();
  const [selectedTags, selectTags] = useState<Tag[]>([]);
  const [filter, setFilter] = useState<GroupBy>("city");

  const handleOpenDetails = useCallback((data?: DataCatalogItem) => {
    setDataset(data);
  }, []);

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

  return (
    <PageLayout
      left={
        <DatasetTree
          addedDatasetDataIDs={addedDatasetDataIDs}
          selectedDataset={selectedDataset}
          catalog={catalog}
          selectedTags={selectedTags}
          filter={filter}
          addDisabled={addDisabled}
          onFilter={handleFilter}
          onTagSelect={handleTagSelect}
          onOpenDetails={handleOpenDetails}
          onDatasetAdd={onDatasetAdd}
        />
      }
      right={
        <DatasetDetails
          dataset={selectedDataset}
          addDisabled={addDisabled}
          onTagSelect={handleTagSelect}
          onDatasetAdd={onDatasetAdd}
        />
      }
    />
  );
};

export default DatasetsPage;
