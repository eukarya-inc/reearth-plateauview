import { useCallback, useEffect, useState } from "react";

import { postMsg } from "../core/utils";
import { Primitive, PublicSetting } from "../types";

import { TEST_PRIMITIVES, TEST_PUBLIC_SETTINGS } from "./TEST_DATA";

type Mode = "edit" | "view" | "pending";

export default () => {
  const [mode, setMode] = useState<Mode>("pending");
  const [primitives, setPrimitives] = useState<Primitive[]>([]);
  const [publicSettings, setPublicSettings] = useState<PublicSetting[]>([]);

  useEffect(() => {
    setMode("edit"); // DEV ONLY
    setPrimitives(TEST_PRIMITIVES); // DEV ONLY
    setPublicSettings(TEST_PUBLIC_SETTINGS); // DEV ONLY
  }, []);

  const handleInEditor = useCallback((inEditor: boolean) => {
    setMode(inEditor ? "edit" : "view");
  }, []);

  useEffect(() => {
    postMsg("getInEditor");
  }, []);

  const onMessage = useCallback(
    (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      switch (e.data.action) {
        case "getInEditor":
          handleInEditor(e.data.payload);
          break;
        default:
          break;
      }
    },
    [handleInEditor],
  );

  useEffect(() => {
    addEventListener("message", onMessage);
    return () => {
      removeEventListener("message", onMessage);
    };
  }, [onMessage]);

  return {
    mode,
    primitives,
    publicSettings,
  };
};
