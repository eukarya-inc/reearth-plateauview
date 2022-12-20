import { PostMessageProps } from "@web/extensions/sidebar/core/types";
import basicOperationHtml from "@web/extensions/sidebar/popups/BasicOperation/index.html?raw";
import clipFunctionHtml from "@web/extensions/sidebar/popups/ClipFunction/index.html?raw";
import shadowFunctionHtml from "@web/extensions/sidebar/popups/ShadowFunction/index.html?raw";
import tryMapInfoHtml from "@web/extensions/sidebar/popups/TryMapInfo/index.html?raw";

import html from "../dist/web/sidebar/core/index.html?raw";
import dataCatalogHtml from "../dist/web/sidebar/modals/datacatalog/index.html?raw";
import welcomeScreenHtml from "../dist/web/sidebar/modals/welcomescreen/index.html?raw";

const reearth = (globalThis as any).reearth;

let addedDatasets: string | undefined = undefined;

const doNotShowWelcome = true; // Make it `let doNotShowWelcome: boolean = false`, and then modify based on storage value when Storage API available

reearth.ui.show(html, { extended: true });

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  // Sidebar
  if (action === "init") {
    reearth.ui.postMessage({ type: "init", payload: { inEditor: reearth.scene.inEditor } });
    if (!doNotShowWelcome) {
      reearth.modal.show(welcomeScreenHtml, { background: "#000000bf" });
    }
  } else if (action === "updateOverrides") {
    reearth.visualizer.overrideProperty(payload);
  } else if (action === "addDatasetToScene") {
    // NEED TO HANDLE ADDING TO SCENE WHEN ABLE
  } else if (
    action === "screenshot" ||
    action === "screenshot-preview" ||
    action === "screenshot-save"
  ) {
    reearth.ui.postMessage({
      type: action,
      payload: reearth.scene.captureScreen(undefined, 0.01),
    });
  } else if (action === "msgFromModal") {
    reearth.ui.postMessage({ type: action, payload });
  } else if (action === "minimize") {
    if (payload) {
      reearth.ui.resize(undefined, undefined, false);
    } else {
      reearth.ui.resize(350, undefined, true);
    }
  } else if (action === "datacatalog-modal-open") {
    addedDatasets = payload;
    reearth.modal.show(dataCatalogHtml, { background: "transparent" });
    // Datacatalog modal
  } else if (action === "modal-close") {
    reearth.modal.close();
  } else if (action === "initDatasetCatalog") {
    reearth.modal.postMessage({ type: "msgFromSidebar", payload: addedDatasets });
  } else if (action === "welcome-modal-open") {
    reearth.modal.show(welcomeScreenHtml, { background: "transparent" });
  } else if (action === "show-popup") {
    if (payload === "Basic") {
      reearth.popup.show({ basicOperationHtml, position: "left" });
    } else if (payload === "map") {
      reearth.popup.show({ tryMapInfoHtml, position: "left" });
    } else if (payload === "shadow") {
      reearth.popup.show({ shadowFunctionHtml, position: "left" });
    } else if (payload === "clip") {
      reearth.popup.show({ clipFunctionHtml, position: "left" });
    }
  } else if (action === "close-popup") {
    reearth.popup.close();
  }
});

reearth.on("update", () => {
  reearth.ui.postMessage({
    type: "extended",
    payload: reearth.widget.extended,
  });
});
