import { postMsg } from "@web/extensions/sidebar/core/utils";
import { useCallback } from "react";

export default () => {
  const handleClosePopup = useCallback(() => {
    postMsg({ action: "close-popup" });
  }, []);

  return {
    handleClosePopup,
  };
};
