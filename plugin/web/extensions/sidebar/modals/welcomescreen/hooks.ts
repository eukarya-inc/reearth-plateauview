import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

export default () => {
  const [ShowVideo, setShowVideo] = useState(false);
  const [dontShowAgain, setDontShowAgain] = useState(false);

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
    postMsg({ action: "modalClose" });
    postMsg({
      action: "storageSaveWelcomeScreen",
      payload: { key: "doNotShowWelcome", value: dontShowAgain },
    });
  }, [dontShowAgain]);

  return {
    handleDontShowAgain,
    ShowVideo,
    dontShowAgain,
    handleShowVideo,
    handleCloseVideo,
    handleClose,
  };
};
