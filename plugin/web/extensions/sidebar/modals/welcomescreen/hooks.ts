import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

export default () => {
  const [showVideo, setShowVideo] = useState(false);
  const [dontShowAgain, setDontShowAgain] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    postMsg({ action: "initWelcome" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return null;
      if (e.data.type) {
        if (e.data.type === "msgToModal" && e.data.message) {
          setIsMobile(e.data.message);
        }
      }
    };
    (globalThis as any).addEventListener("message", (e: any) => eventListenerCallback(e));
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  });

  const handleDontShowAgain = useCallback(() => {
    setDontShowAgain(!dontShowAgain);
  }, [dontShowAgain]);

  const handleShowVideo = useCallback(() => {
    setShowVideo(true);
  }, []);

  const handleCloseVideo = useCallback(() => {
    setShowVideo(false);
  }, []);

  const handleClose = useCallback(() => {
    if (dontShowAgain)
      postMsg({
        action: "storageSave",
        payload: { key: "doNotShowWelcome", value: dontShowAgain },
      });
    postMsg({ action: "modalClose" });
  }, [dontShowAgain]);

  const handleOpenHelp = useCallback(() => {
    postMsg({ action: "modalClose" });
    postMsg({ action: "triggerHelpOpen" });
  }, []);

  const handleOpenCatalog = useCallback(() => {
    if (isMobile) {
      // postMsg({ action: "catalogModalOpen" });
    } else {
      postMsg({ action: "triggerCatalogOpen" });
    }
  }, [isMobile]);

  return {
    isMobile,
    showVideo,
    dontShowAgain,
    handleDontShowAgain,
    handleShowVideo,
    handleCloseVideo,
    handleClose,
    handleOpenHelp,
    handleOpenCatalog,
  };
};
