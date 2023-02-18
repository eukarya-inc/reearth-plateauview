import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ComponentType, useCallback, useMemo, useState } from "react";

export type Props = {
  dataset: DataCatalogItem | UserDataItem;
  isShareable?: boolean;
  addDisabled: boolean;
  contentSection?: ComponentType;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem) => void;
};

const DatasetDetails: React.FC<Props> = ({
  dataset,
  isShareable,
  addDisabled,
  contentSection: ContentSection,
  onDatasetAdd,
}) => {
  const datasetType: "catalog" | "user" = useMemo(
    () => ("dataID" in dataset ? "catalog" : "user"),
    [dataset],
  );
  const [published, setPublished] = useState(
    datasetType === "user" ? undefined : (dataset as DataCatalogItem).public,
  );
  const showShareButton = false; // This code can be removed when decision about share button is made

  const handlePublish = useCallback(() => {
    if (datasetType === "user") return;

    setPublished(true);
  }, [datasetType, setPublished]);

  const handleDatasetAdd = useCallback(() => {
    if (!dataset || addDisabled) return;
    onDatasetAdd(dataset);
  }, [dataset, addDisabled, onDatasetAdd]);

  return (
    <>
      <TopWrapper>
        <HeaderWrapper>
          <Title>{dataset.name}</Title>
          {datasetType === "catalog" && (
            <PublishButton published={published} onClick={handlePublish}>
              <HoverText published={published}>公開</HoverText>
              <Text published={published}>{published ? "公開済み" : "未公開"}</Text>
            </PublishButton>
          )}
        </HeaderWrapper>
        <ButtonWrapper>
          <AddButton disabled={addDisabled} onClick={handleDatasetAdd}>
            {!addDisabled && <Icon icon="plusCircle" />}
            {addDisabled ? "シーンに追加済み" : "シーンに追加"}
          </AddButton>
          {showShareButton && (
            <ShareButton isShareable={isShareable}>
              <Icon icon="share" />
              シェア
            </ShareButton>
          )}
        </ButtonWrapper>
      </TopWrapper>
      {ContentSection && (
        <Wrapper>
          <ContentSection />
        </Wrapper>
      )}
    </>
  );
};

export default DatasetDetails;

const Wrapper = styled.div`
  padding: 16px 24px;
`;

const TopWrapper = styled(Wrapper)`
  border-bottom: 1px solid #c7c5c5;
`;

const HeaderWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
`;

const Title = styled.p`
  font-size: 16px;
  font-weight: 700;
  line-height: 22px;
  margin: 0;
`;

const ButtonWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 12px;
`;

const BaseButton = styled.button<{ disabled?: boolean }>`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  height: 40px;
  font-weight: 500;
  border-radius: 4px;
  border: 1px solid #e6e6e6;
`;

const AddButton = styled(BaseButton)<{ disabled?: boolean }>`
  color: ${({ disabled }) => (disabled ? "grey" : "#ffffff")};
  background-color: ${({ disabled }) => (disabled ? "#dcdcdc" : "#00bebe")};
  flex: 1;
  ${({ disabled }) => !disabled && "cursor: pointer;"}
`;

const ShareButton = styled(BaseButton)<{ isShareable?: boolean }>`
  display: ${({ isShareable }) => (isShareable !== false ? "flex" : "none")};
  color: #00bebe;
  background-color: #ffffff;
  flex: 1;
  cursor: pointer;
`;

const PublishButton = styled(BaseButton)<{ published?: boolean }>`
  display: flex;
  min-width: 120px;
  color: #ffffff;
  background-color: ${({ published }) => (published ? "#00bebe" : "#BFBFBF")};
  ${({ published }) => !published && "cursor: pointer;"};
  &:hover {
    background-color: #00bebe;
  }
`;

const HoverText = styled.p<{ published?: boolean }>`
  display: none;
  margin-bottom: 0;
  ${PublishButton}:hover & {
    display: ${({ published }) => (published ? "none" : "initial")};
  }
`;

const Text = styled.p<{ published?: boolean }>`
  display: initial;
  margin-bottom: 0;
  ${PublishButton}:hover & {
    display: ${({ published }) => (published ? "initial" : "none")};
  }
`;
