import useHooks from "@web/extensions/sidebar/core/components/hooks";
import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback } from "react";

import { CatalogItem, CatalogRawItem } from "../../core/processCatalog";

export default () => {
  const { handleProjectDatasetAdd } = useHooks();

  const handleClose = useCallback(() => {
    postMsg({ action: "popupClose" });
  }, []);

  const handleDatasetAdd = useCallback(
    (dataset: CatalogItem) => {
      handleProjectDatasetAdd(dataset as CatalogRawItem);
      handleClose();
    },
    [handleClose, handleProjectDatasetAdd],
  );

  return {
    handleDatasetAdd,
  };
};
