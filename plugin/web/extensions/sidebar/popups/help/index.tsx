import { postMsg } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { Tab } from "../../core/components/content/Help/hooks";

import BasicOperation from "./BasicOperation";
import ClipFunction from "./ClipFunction";
import ShadowFunction from "./ShadowFunction";
import { PopupWrapper } from "./sharedComponent";
import TryMapInfo from "./TryMapInfo";

const Help: React.FC = () => {
  const [currentPopup, setCurrentPopup] = useState<Tab>();

  const handleClosePopup = useCallback(() => {
    postMsg({ action: "popupClose" });
  }, []);

  useEffect(() => {
    postMsg({ action: "initPopup" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return null;
      if (e.data.type) {
        if (e.data.type === "msgToPopup" && e.data.message) {
          setCurrentPopup(e.data.message);
        }
      }
    };
    (globalThis as any).addEventListener("message", (e: any) => eventListenerCallback(e));
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  });

  return (
    <PopupWrapper handleClose={handleClosePopup}>
      {currentPopup &&
        {
          basic: <BasicOperation />,
          map: <TryMapInfo />,
          shadow: <ShadowFunction />,
          clip: <ClipFunction />,
        }[currentPopup]}
    </PopupWrapper>
  );
};

export default Help;
