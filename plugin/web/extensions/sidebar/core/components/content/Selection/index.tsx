import Footer from "@web/extensions/sidebar/core/components/Footer";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import DatasetCard, { Dataset } from "../common/DatasetCard";

export type Props = {
  inEditor?: boolean;
  selectedDatasets: Dataset[];
  // onDatasetUpdate?: (dataset: Dataset) => void;
  onDatasetRemove: (id: string) => void;
  onDatasetRemoveAll: () => void;
  onModalOpen?: () => void;
};

const exampleDatasets: Dataset[] = [
  { id: "haha", name: "渋谷", type: "3d-tiles", visible: true },
  { id: "haha2", name: "横浜", type: "3d-tiles", visible: false },
];

const Selection: React.FC<Props> = ({
  inEditor,
  selectedDatasets,
  // onDatasetUpdate,
  onDatasetRemove,
  onDatasetRemoveAll,
  onModalOpen,
}) => {
  return (
    <Wrapper>
      <InnerWrapper>
        <StyledButton onClick={onModalOpen}>
          <StyledIcon icon="plusCircle" size={20} />
          <ButtonText>カタログから検索する</ButtonText>
        </StyledButton>
        {(exampleDatasets ?? selectedDatasets)
          .map(d => (
            <DatasetCard key={d.id} dataset={d} inEditor={inEditor} onRemove={onDatasetRemove} />
          ))
          .reverse()}
      </InnerWrapper>
      <Footer datasetQuantity={selectedDatasets.length} onRemoveAll={onDatasetRemoveAll} />
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
`;

const StyledIcon = styled(Icon)`
  margin-right: 8px;
`;
