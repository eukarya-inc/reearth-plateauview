import { postMsg } from "@web/extensions/sidebar/core/utils";
import { useCallback } from "react";

export default () => {
  const handleClosePopup = useCallback(() => {
    postMsg({ action: "close-popup" });
  }, []);

  const handleShowMapModal = useCallback(() => {
    postMsg({ action: "show-map-modal" });
  }, []);

  const handleShowClipModal = useCallback(() => {
    postMsg({ action: "show-clip-modal" });
  }, []);
  return {
    handleClosePopup,
    handleShowMapModal,
    handleShowClipModal,
  };
};
