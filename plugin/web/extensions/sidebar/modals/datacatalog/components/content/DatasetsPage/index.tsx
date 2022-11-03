import PageLayout from "../PageLayout";

import DatasetTree from "./DatasetTree";
import DatasetDetails from "./Details";

export type Props = {
  onDatasetAdd: () => void;
};

const DatasetsPage: React.FC<Props> = ({ onDatasetAdd }) => {
  return (
    <PageLayout left={<DatasetTree />} right={<DatasetDetails onDatasetAdd={onDatasetAdd} />} />
  );
};

export default DatasetsPage;
