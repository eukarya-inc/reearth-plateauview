import { useCallback, useEffect, useState } from "react";

import { Dataset } from "./components/content/Selection/DatasetCard";
import { useCurrentOverrides } from "./state";
import { ReearthApi } from "./types";
import { mergeProperty, postMsg } from "./utils";

export default () => {
  // ****************************************
  // Override Logic
  const [overrides, updateOverrides] = useCurrentOverrides();

  const handleOverridesUpdate = useCallback(
    (updatedProperties: Partial<ReearthApi>) => {
      updateOverrides([overrides, updatedProperties].reduce((p, v) => mergeProperty(p, v)));
    },
    [overrides],
  );

  useEffect(() => {
    postMsg({ action: "updateOverrides", payload: overrides });
  }, [overrides]);
  // ****************************************

  // ****************************************
  // Minimize Logic
  const [minimized, setMinimize] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      postMsg({ action: "minimize", payload: minimized });
    }, 250);
  }, [minimized]);
  // ****************************************

  // ****************************************
  // Dataset Logic
  const [selectedDatasets, updateDatasets] = useState<Dataset[]>([]);

  const handleDatasetAdd = useCallback((dataset: Dataset) => {
    updateDatasets(oldDatasets => [...oldDatasets, dataset]);
  }, []);

  const handleDatasetRemove = useCallback(
    (id: string) => updateDatasets(oldDatasets => oldDatasets.filter(d => d.id !== id)),
    [],
  );

  const handleDatasetRemoveAll = useCallback(() => updateDatasets([]), []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.type === "msgFromModal") {
        if (e.data.payload.dataset) {
          handleDatasetAdd(e.data.payload.dataset);
        }
      }
    };
    addEventListener("message", e => eventListenerCallback(e));
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []);
  // ****************************************

  const handleModalOpen = useCallback(() => {
    postMsg({ action: "modal-open" });
  }, []);

  return {
    selectedDatasets,
    overrides,
    minimized,
    setMinimize,
    handleDatasetRemove,
    handleDatasetRemoveAll,
    handleOverridesUpdate,
    handleModalOpen,
  };
};

addEventListener("message", e => {
  if (e.source !== parent) return;
  if (e.data.type) {
    if (e.data.type === "extended") {
      updateExtended(e.data.payload);
    } else if (e.data.type === "screenshot") {
      generatePrintView(e.data.payload);
    } else if (e.data.type === "screenshot-save") {
      const link = document.createElement("a");
      link.download = "screenshot.png";
      link.href = e.data.payload;
      link.click();
      link.remove();
    }
  }
});

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

function generatePrintView(payload?: string) {
  const doc = window.open()?.document;

  if (!doc || !payload) return;

  const css = `html,body{ margin: 0; }`;

  const styleTag = doc.createElement("style");
  styleTag.appendChild(document.createTextNode(css));
  styleTag.setAttribute("type", "text/css");
  doc.head.appendChild(styleTag);

  const iframe = doc.createElement("iframe");
  iframe.style.width = "100%";
  iframe.style.height = "100%";
  iframe.style.border = "none";

  doc.body.appendChild(iframe);

  const iframeDoc = iframe.contentWindow?.document;
  if (!iframeDoc) return;

  iframeDoc.open();

  const iframeHTML = `
  <div style="display: flex; flex-direction: column; max-width: 1200px; height: 100%; margin: 0 auto; padding: 20px;">
    <div style="display: flex; justify-content: right; align-items: center; gap: 8px; height: 60px;">
      <button style="padding: 8px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">Download map</button>
      <button style="padding: 9px; border: none; border-radius: 4px; background: #00BEBE; color: white; cursor: pointer;">Print</button>
    </div>
    <div style="display: flex; justify-content: center; width: 100%;">
      <img src="${payload}" style="max-width: 100%; object-fit: contain;" />
    </div>
    <div>
      <p>This map was created using https://plateauview.mlit.go.jp on ${new Date()}</p>
    </div>
  </div>
`;
  iframe.contentWindow?.document.write(iframeHTML);

  const iframeHtmlStyle = iframe.contentWindow?.document.createElement("style");
  if (iframeHtmlStyle) {
    iframeHtmlStyle.appendChild(document.createTextNode(css));
    iframeHtmlStyle.setAttribute("type", "text/css");
    iframe.contentWindow?.document.head.appendChild(iframeHtmlStyle);
  }

  iframe.contentWindow?.document.close();

  return iframe;
}
