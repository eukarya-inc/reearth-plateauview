import { useCallback, useState } from "react";

import PageLayout from "../PageLayout";

import DatasetTree from "./DatasetTree";
import { TEST_CATALOG_DATA } from "./DatasetTree/TEST_catalog_data";
import DatasetDetails, { Dataset as DatasetType } from "./Details";

export type Dataset = DatasetType;

export type Props = {
  onDatasetAdd: (dataset: Dataset) => void;
};

const DatasetsPage: React.FC<Props> = ({ onDatasetAdd }) => {
  const catalog = TEST_CATALOG_DATA;
  const [selectedDataset, setDataset] = useState<Dataset>();

  const handleOpenDetails = useCallback(() => {
    //HERE HANDLE SETTING DATASET ID
    setDataset(undefined); // Gotta do a lot to get this working
    /*
    Needs:
    title
    description
    tags
    etc.
    */
  }, []);

  return (
    <PageLayout
      left={<DatasetTree catalog={catalog} onOpenDetails={handleOpenDetails} />}
      right={<DatasetDetails dataset={selectedDataset} onDatasetAdd={onDatasetAdd} />}
    />
  );
};

export default DatasetsPage;
