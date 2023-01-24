import { CatalogItem } from "@web/extensions/sidebar/core/processCatalog";
import { useCallback, useState } from "react";

import PageLayout from "../PageLayout";

import Details from "./Details";
import FileSelectPane from "./FileSelect";

export type Props = {
  onDatasetAdd: (dataset: CatalogItem) => void;
};

const YourDataPage: React.FC<Props> = ({ onDatasetAdd }) => {
  const [selectedDataset, setDataset] = useState<CatalogItem>();

  const handleOpenDetails = useCallback((data?: CatalogItem) => {
    setDataset(data);
  }, []);

  return (
    <PageLayout
      left={<FileSelectPane onOpenDetails={handleOpenDetails} />}
      right={<Details dataset={selectedDataset} onDatasetAdd={onDatasetAdd} />}
    />
  );
};

export default YourDataPage;
