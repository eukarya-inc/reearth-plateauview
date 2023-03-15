import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import PageLayout from "@web/extensions/sidebar/modals/datacatalog/components/content/PageLayout";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useMemo, useState } from "react";

import { GroupBy } from "../../../api/api";
import { UserDataItem } from "../../../types";

import DatasetTree from "./DatasetTree";
import DatasetDetails, { Tag } from "./Details";

export type Props = {
  catalog?: DataCatalogItem[];
  addedDatasetDataIDs?: string[];
  inEditor?: boolean;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem) => void;
  onDatasetPublish: (dataID: string, publish: boolean) => void;
};

const DatasetsPage: React.FC<Props> = ({
  catalog,
  addedDatasetDataIDs,
  inEditor,
  onDatasetAdd,
  onDatasetPublish,
}) => {
  const [selectedDatasetID, setDatasetID] = useState<string>();
  const [selectedTags, selectTags] = useState<Tag[]>([]);
  const [filter, setFilter] = useState<GroupBy>("city");

  const handleOpenDetails = useCallback((data?: DataCatalogItem) => {
    setDatasetID(data?.dataID);
  }, []);

  const handleFilter = useCallback((filter: GroupBy) => {
    setFilter(filter);
    postMsg({ action: "storageSave", payload: { key: "filter", value: filter } });
  }, []);

  useEffect(() => {
    postMsg({ action: "catalogModalOpen" }); // Needed to trigger sending client storage data from Sidebar
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "catalogModalOpen") {
        if (e.data.payload.filter) {
          setFilter(e.data.payload.filter);
        }
      }
    };
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, [handleFilter]);

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
