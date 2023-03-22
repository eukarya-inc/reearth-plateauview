import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import DetailsComponent from "@web/extensions/sidebar/modals/datacatalog/components/content/DatasetDetails";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import DOMPurify from "dompurify";
import { useMemo } from "react";

import { UserDataItem } from "../../../types";

import { Tag as TagType } from "./Tags";
// import Tags, {Tag as TagType} from "./Tags";

export type Tag = TagType;

export type Props = {
  dataset?: DataCatalogItem;
  isMobile?: boolean;
  inEditor?: boolean;
  addDisabled: (dataID: string) => boolean;
  onTagSelect?: (tag: TagType) => void;
  onDatasetAdd: (dataset: DataCatalogItem | UserDataItem, keepModalOpen?: boolean) => void;
  onDatasetPublish?: (dataID: string, publish: boolean) => void;
};

const DatasetDetails: React.FC<Props> = ({
  dataset,
  // isMobile,
  inEditor,
  addDisabled,
  // onTagSelect,
  onDatasetAdd,
  onDatasetPublish,
}) => {
  // const datasetTags = useMemo(
  //   () => (dataset?.type !== "group" ? dataset?.tags?.map(tag => tag) : undefined),
  //   [dataset],
  // );

  const ContentComponent: React.FC = () => {
    const descWithLink = useMemo(() => {
      if (!dataset?.desc) return "";
      const descHtml = dataset.desc.replace(
        /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/g,
        url => {
          return `<a href="${url}" target="_blank">${url}</a>`;
        },
      );
      return DOMPurify.sanitize(descHtml);
    }, []);

    return (
      <>
        {/* {!isMobile && <Tags tags={datasetTags} onTagSelect={onTagSelect} />} */}
        {dataset && dataset?.type !== "group" && (
          <Content dangerouslySetInnerHTML={{ __html: descWithLink }} />
        )}
      </>
    );
  };

  return dataset ? (
    <DetailsComponent
      dataset={dataset}
      addDisabled={addDisabled(dataset.dataID)}
      inEditor={inEditor}
      isPublishable={!!dataset.itemId}
      contentSection={ContentComponent}
      onDatasetAdd={onDatasetAdd}
      onDatasetPublish={onDatasetPublish}
    />
  ) : (
    <NoData>
      <NoDataMain>
        <Icon icon="empty" size={64} />
        <StyledP>データがない</StyledP>
        <br />
        <StyledP>データセットを選択してください(プレビューが表示されます)</StyledP>
      </NoDataMain>
      <NoDataFooter
        onClick={() =>
          window.open("https://www.geospatial.jp/ckan/dataset/plateau-tokyo23ku", "_blank")
        }>
        <Icon icon="newPage" size={16} />
        <StyledP> オープンデータ・ダウンロード(G空間情報センターへのリンク)</StyledP>
      </NoDataFooter>
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
  align-items: center;
  flex: 1;
  flex-direction: column;
`;

const NoDataFooter = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  margin: 0;
  color: #00bebe;
  cursor: pointer;
`;

const StyledP = styled.p`
  margin: 0;
  text-align: center;
`;

const Content = styled.div`
  margin-top: 16px;
  white-space: pre-wrap;

  a {
    color: #00bebe;
  }
`;
