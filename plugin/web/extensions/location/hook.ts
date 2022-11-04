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

  const handleModalChange = useCallback(() => {
    setModal(!showModal);
    postMsg({ action: !showModal ? "modal-open" : "modal-close" });
  }, [showModal]);

  return {
    currentPoint,
    handleModalChange,
  };
};
