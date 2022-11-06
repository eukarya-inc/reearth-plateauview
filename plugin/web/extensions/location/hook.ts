import { useCallback, useEffect, useState } from "react";

import { MouseEventData } from "./types";
import { postMsg } from "./utils";

export default () => {
  const [currentPoint, setCurrentPoint] = useState<MouseEventData>();
  const [showModal, setModal] = useState(false);

  const updateCurrentPoint = useCallback((mousedata: MouseEventData) => {
    setCurrentPoint(mousedata);
  }, []);

  useEffect(() => {
    (globalThis as any).addEventListener("message", (e: any) => {
      if (e.source !== parent) return;
      if (e.data.type) {
        if (e.data.type === "mousedata") {
          console.log(e.data);
          updateCurrentPoint(e.data);
        }
      }
    });
  });

  const handlegoogleModalChange = useCallback(() => {
    setModal(!showModal);
    postMsg({ action: !showModal ? "modal-google-open" : "modal-close" });
  }, [showModal]);
  const handleTerrainModalChange = useCallback(() => {
    setModal(!showModal);
    postMsg({ action: !showModal ? "modal-Terrain-open" : "modal-close" });
  }, [showModal]);

  return {
    currentPoint,
    handlegoogleModalChange,
    handleTerrainModalChange,
  };
};
