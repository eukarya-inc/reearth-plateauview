import { mergeProperty, postMsg } from "@web/utils";
import { useCallback, useEffect, useState } from "react";

import { useCurrentOverrides } from "./state";
import { ReearthApi } from "./types";

export default () => {
  const [overrides, updateOverrides] = useCurrentOverrides();
  const [showModal, setModal] = useState(false);

  const handleOverridesUpdate = useCallback(
    (updatedProperties: Partial<ReearthApi>) => {
      updateOverrides([overrides, updatedProperties].reduce((p, v) => mergeProperty(p, v)));
    },
    [overrides],
  );

  const handleModalChange = useCallback(() => {
    setModal(!showModal);
    postMsg({ action: !showModal ? "modal-open" : "modal-close" });
  }, [showModal]);

  useEffect(() => {
    postMsg({ action: "updateOverrides", payload: overrides });
  }, [overrides]);

  return {
    overrides,
    handleOverridesUpdate,
    handleModalChange,
  };
};

function updateExtended(e: { vertically: boolean }) {
  const html = document.querySelector("html");
  const body = document.querySelector("body");
  const root = document.getElementById("root");

  if (e?.vertically) {
    html?.classList.add("extended");
    body?.classList.add("extended");
    root?.classList.add("extended");
  } else {
    html?.classList.remove("extended");
    body?.classList.remove("extended");
    root?.classList.remove("extended");
  }
}

addEventListener("message", e => {
  if (e.source !== parent) return;
  if (e.data.type) {
    if (e.data.type === "extended") {
      updateExtended(e.data.payload);
    }
  }
});
