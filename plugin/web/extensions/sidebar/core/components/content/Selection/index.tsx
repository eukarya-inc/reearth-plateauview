import Footer from "@web/extensions/sidebar/core/components/Footer";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { DataCatalogItem, Template } from "../../../types";
import DatasetCard from "../common/DatasetCard";

export type Props = {
  className?: string;
  inEditor?: boolean;
  selectedDatasets?: DataCatalogItem[];
  templates?: Template[];
  savingDataset: boolean;
  onDatasetSave: (dataID: string) => void;
  onDatasetUpdate: (dataset: DataCatalogItem, cleanseOverride?: any) => void;
  onDatasetRemove: (dataID: string) => void;
  onDatasetRemoveAll: () => void;
  onModalOpen?: () => void;
  onThreeDTilesSearch: (id: string) => void;
};

const Selection: React.FC<Props> = ({
  className,
  inEditor,
  selectedDatasets,
  templates,
  savingDataset,
  onDatasetSave,
  onDatasetUpdate,
  onDatasetRemove,
  onDatasetRemoveAll,
  onModalOpen,
  onThreeDTilesSearch,
}) => {
  return (
    <Wrapper className={className}>
      <InnerWrapper>
        {onModalOpen && (
          <StyledButton onClick={onModalOpen}>
            <StyledIcon icon="plusCircle" size={20} />
            <ButtonText>カタログから検索する</ButtonText>
          </StyledButton>
        )}
        {selectedDatasets
          ?.map(d => (
            <DatasetCard
              key={d.id}
              dataset={d}
              templates={templates}
              savingDataset={savingDataset}
              inEditor={inEditor}
              onDatasetSave={onDatasetSave}
              onDatasetUpdate={onDatasetUpdate}
              onDatasetRemove={onDatasetRemove}
              onThreeDTilesSearch={onThreeDTilesSearch}
            />
          ))
          .reverse()}
      </InnerWrapper>
      <Footer datasetQuantity={selectedDatasets?.length} onRemoveAll={onDatasetRemoveAll} />
    </Wrapper>
  );
};

export default Selection;

const Wrapper = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

const InnerWrapper = styled.div`
  padding: 16px;
  flex: 1;
  overflow: auto;
`;

const StyledButton = styled.button`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  border: none;
  border-radius: 4px;
  background: #00bebe;
  color: #fff;
  padding: 10px;
  cursor: pointer;
`;

const ButtonText = styled.p`
  margin: 0;
  user-select: none;
`;

const StyledIcon = styled(Icon)`
  margin-right: 8px;
`;
