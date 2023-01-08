import { useState } from "react";

import { Tab } from "../../core/components/content/Help/hooks";
import BasicOperation from "../BasicOperation";
import ClipFunction from "../ClipFunction";
import ShadowFunction from "../ShadowFunction";
import TryMapInfo from "../TryMapInfo";

const Popup: React.FC = () => {
  const [currentPopup, setCurrentPopup] = useState<Tab>("basic");

  addEventListener("message", e => {
    if (e.source !== parent) return null;
    if (e.data.type) {
      if (e.data.type === "msgFromHelp" && e.data.message) {
        setCurrentPopup(e.data.message);
      }
    }
  });

  return {
    basic: <BasicOperation />,
    map: <TryMapInfo />,
    shadow: <ShadowFunction />,
    clip: <ClipFunction />,
  }[currentPopup];
};

export default Popup;
