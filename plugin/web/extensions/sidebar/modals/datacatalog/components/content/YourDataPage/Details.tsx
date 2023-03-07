import DetailsComponent from "@web/extensions/sidebar/modals/datacatalog/components/content/DatasetDetails";
import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { styled } from "@web/theme";
import { useCallback } from "react";

export type Props = {
  isShareable?: boolean;
  requireLayerName?: boolean;
  dataset?: UserDataItem;
  onDatasetAdd: (dataset: UserDataItem) => void;
};

const DatasetDetails: React.FC<Props> = ({
  dataset,
  isShareable,
  requireLayerName,
  onDatasetAdd,
}) => {
  const handleDatasetAdd = useCallback(() => {
    if (!dataset) return;
    onDatasetAdd(dataset);
  }, [dataset, onDatasetAdd]);

  const ContentComponent: React.FC = () => (
    <Content
      dangerouslySetInnerHTML={{
        __html: dataset?.description as string,
      }}
    />
  );

  return dataset ? (
    <DetailsComponent
      dataset={dataset}
      isShareable={isShareable}
      requireLayerName={requireLayerName}
      addDisabled={false}
      onDatasetAdd={handleDatasetAdd}
      contentSection={ContentComponent}
    />
  ) : (
    <NoData>
      <NoDataMain>
        <StyledP>データセットがありません</StyledP>
      </NoDataMain>
    </NoData>
  );
};

export default DatasetDetails;

const NoData = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  color: rgba(0, 0, 0, 0.25);
  height: calc(100% - 24px);
  margin-bottom: 24px;
`;

const NoDataMain = styled.div`
  display: flex;
  justify-content: center;
  flex: 1;
  flex-direction: column;
`;

const StyledP = styled.p`
  margin: 0;
  text-align: center;
`;

const Content = styled.div`
  margin-top: 16px;
`;
