import Footer from "@web/extensions/sidebar/core/components/Footer";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import DatasetWrapper, { Dataset } from "./DatasetCard";

export type Props = {
  selectedDatasets: Dataset[];
  onDatasetRemove: (id: string) => void;
  onDatasetRemoveAll: () => void;
  onModalOpen?: () => void;
};

const Selection: React.FC<Props> = ({
  selectedDatasets,
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
        {selectedDatasets
          .map(d => <DatasetWrapper key={d.id} dataset={d} onRemove={onDatasetRemove} />)
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
  color: #4a4a4a;
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
  transition: background 0.2s;

  :hover {
    background: #4cc2c2;
  }
`;

const ButtonText = styled.p`
  margin: 0;
`;

const StyledIcon = styled(Icon)`
  margin-right: 8px;
`;
