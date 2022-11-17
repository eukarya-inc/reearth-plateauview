import { useCallback, useState } from "react";

import PageLayout from "../PageLayout";

import DatasetTree from "./DatasetTree";
import { Data } from "./DatasetTree/FileTree";
import { TEST_CATALOG_DATA } from "./DatasetTree/TEST_catalog_data";
import DatasetDetails, { Tag } from "./Details";

export type Props = {
  onDatasetAdd: (dataset: Data) => void;
};

const DatasetsPage: React.FC<Props> = ({ onDatasetAdd }) => {
  const catalog = TEST_CATALOG_DATA;
  const [selectedDataset, setDataset] = useState<Data>();
  const [selectedTags, selectTags] = useState<Tag[]>([]);

  const handleOpenDetails = useCallback((data?: Data) => {
    setDataset(data);
  }, []);

  const handleTagSelect = useCallback(
    (tag: Tag) =>
      selectTags(tags => (tags.includes(tag) ? [...tags.filter(t => t !== tag)] : [...tags, tag])),
    [],
  );

  return (
    <PageLayout
      left={
        <DatasetTree
          catalog={catalog}
          selectedTags={selectedTags}
          onTagSelect={handleTagSelect}
          onOpenDetails={handleOpenDetails}
        />
      }
      right={
        <DatasetDetails
          dataset={selectedDataset}
          onTagSelect={handleTagSelect}
          onDatasetAdd={onDatasetAdd}
        />
      }
    />
  );
};

export default DatasetsPage;
