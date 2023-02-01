import { CatalogRawItem } from "@web/extensions/sidebar/core/processCatalog";
import DatasetTree, {
  Tag,
} from "@web/extensions/sidebar/modals/datacatalog/components/content/DatasetsPage/DatasetTree";
import { postMsg } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useCallback, useEffect, useState } from "react";

import PopupItem from "../sharedComponents/PopupItem";

type Props = {
  isMobile?: boolean;
  rawCatalog?: CatalogRawItem[];
};

const Catalog: React.FC<Props> = ({ isMobile, rawCatalog }) => {
  useEffect(() => {
    postMsg({ action: "extendPopup" });
  }, []);

  const [selectedTags, selectTags] = useState<Tag[]>([]);

  // const handleOpenDetails = useCallback((data?: CatalogItem) => {
  //   setDataset(data);
  // }, []);

  const handleTagSelect = useCallback(
    (tag: Tag) =>
      selectTags(tags => (tags.includes(tag) ? [...tags.filter(t => t !== tag)] : [...tags, tag])),
    [],
  );

  return (
    <Wrapper>
      <PopupItem>
        <Title>データカタログ</Title>
      </PopupItem>
      <DatasetTree
        isMobile={isMobile}
        rawCatalog={rawCatalog}
        selectedTags={selectedTags}
        onTagSelect={handleTagSelect}
      />
    </Wrapper>
  );
};

export default Catalog;

const Wrapper = styled.div`
  border-top: 1px solid #d9d9d9;
`;

const Title = styled.p`
  margin: 0;
`;
