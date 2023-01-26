import { CatalogRawItem } from "@web/extensions/sidebar/core/processCatalog";
import { PostMessageProps } from "@web/extensions/sidebar/types";

import html from "../dist/web/sidebar/core/index.html?raw";
import clipVideoHtml from "../dist/web/sidebar/modals/clipVideo/index.html?raw";
import dataCatalogHtml from "../dist/web/sidebar/modals/datacatalog/index.html?raw";
import mapVideoHtml from "../dist/web/sidebar/modals/mapVideo/index.html?raw";
import welcomeScreenHtml from "../dist/web/sidebar/modals/welcomescreen/index.html?raw";
import helpPopupHtml from "../dist/web/sidebar/popups/help/index.html?raw";
import mobileDropdownHtml from "../dist/web/sidebar/popups/mobileDropdown/index.html?raw";

type PluginExtensionInstance = {
  id: string;
  pluginId: string;
  name: string;
  extensionId: string;
  extensionType: "widget" | "block";
  layout: any;
  extended: any;
  runTimes?: number;
};

const reearth = (globalThis as any).reearth;

// Defaults
let welcomePageIsOpen = false;
let mobileDropdownIsOpen = false;
const defaultLocation = { zone: "outer", section: "left", area: "middle" };
const mobileLocation = { zone: "outer", section: "center", area: "top" };

let rawCatalog: CatalogRawItem[] = [];

let addedDatasets: string | undefined = undefined;

const widgetInstance = reearth.plugins.instances.find(
  (i: PluginExtensionInstance) => i.id === reearth.widget.id,
);
// Defaults (end)

reearth.ui.show(html, { extended: true });

if (widgetInstance.runTimes === 0) {
  if (reearth.viewport.isMobile) {
    reearth.clientStorage.setAsync("isMobile", true);
    reearth.widget.moveTo(mobileLocation);
  } else {
    reearth.clientStorage.setAsync("isMobile", false);
  }
} else {
  reearth.clientStorage.getAsync("isMobile").then((value: any) => {
    if (value && reearth.viewport.isMobile) {
      reearth.widget.moveTo(mobileLocation);
      reearth.clientStorage.setAsync("isMobile", false);
    } else if (!value && !reearth.viewport.isMobile) {
      reearth.widget.moveTo(defaultLocation);
      reearth.clientStorage.setAsync("isMobile", true);
    }
  });
}

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  // Mobile specific
  if (action === "mobileDropdownOpen") {
    reearth.popup.show(mobileDropdownHtml, {
      position: "bottom",
      width: reearth.viewport.width - 12,
    });
    mobileDropdownIsOpen = true;
  } else if (action === "msgToMobileDropdown") {
    reearth.popup.postMessage({ type: "msgToPopup", message: payload });
  }

  // Sidebar
  if (action === "initSidebar") {
    const sidebarInstance = reearth.plugins.instances.find(
      (i: PluginExtensionInstance) => i.id === reearth.widget.id,
    );
    if (sidebarInstance.runTimes === 0) {
      reearth.visualizer.overrideProperty(payload);
      reearth.clientStorage.setAsync("overrides", payload);
      reearth.ui.postMessage({
        type: action,
        payload: {
          isMobile: reearth.viewport.isMobile,
          projectID: reearth.viewport.query.projectID,
          inEditor: reearth.scene.inEditor,
          backendAccessToken: reearth.widget.property.default?.plateauAccessToken ?? "",
          backendURL: reearth.widget.property.default?.plateauURL ?? "",
          cmsURL: reearth.widget.property.default?.cmsURL ?? "",
          reearthURL: reearth.widget.property.default?.reearthURL ?? "",
        },
      });
    } else {
      reearth.clientStorage.getAsync("overrides").then((value: any) => {
        reearth.ui.postMessage({
          type: action,
          payload: {
            isMobile: reearth.viewport.isMobile,
            projectID: reearth.viewport.query.projectID,
            inEditor: reearth.scene.inEditor,
            backendAccessToken: reearth.widget.property.default?.plateauAccessToken ?? "",
            backendURL: reearth.widget.property.default?.plateauURL ?? "",
            cmsURL: reearth.widget.property.default?.cmsURL ?? "",
            reearthURL: reearth.widget.property.default?.reearthURL ?? "",
            storedOverrides: value,
          },
        });
      });
    }

    reearth.clientStorage.getAsync("doNotShowWelcome").then((value: any) => {
      if (!value && !reearth.scene.inEditor) {
        reearth.modal.show(welcomeScreenHtml, {
          width: reearth.viewport.width,
          height: reearth.viewport.height,
        });
        welcomePageIsOpen = true;
      }
    });
  } else if (action === "storageSave") {
    reearth.clientStorage.setAsync(payload.key, payload.value);
  } else if (action === "storageFetch") {
    reearth.clientStorage.getAsync(payload.key).then((value: any) => {
      reearth.ui.postMessage({
        type: "getAsync",
        payload: value,
      });
    });
  } else if (action === "storageKeys") {
    reearth.clientStorage.keysAsync().then((value: any) => {
      reearth.ui.postMessage({
        type: "keysAsync",
        payload: value,
      });
    });
  } else if (action === "storageDelete") {
    reearth.clientStorage.deleteAsync(payload.key);
  } else if (action === "updateOverrides") {
    reearth.visualizer.overrideProperty(payload);
    reearth.clientStorage.setAsync("overrides", payload);
  } else if (action === "addDatasetToScene") {
    // NEED TO HANDLE ADDING TO SCENE WHEN ABLE
  } else if (
    action === "screenshot" ||
    action === "screenshotPreview" ||
    action === "screenshotSave"
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
  } else if (action === "catalogModalOpen") {
    addedDatasets = payload.addedDatasets;
    rawCatalog = payload.rawCatalog;
    reearth.modal.show(dataCatalogHtml, { background: "transparent" });
  } else if (action === "modalClose") {
    reearth.modal.close();
    welcomePageIsOpen = false;
  } else if (action === "initDataCatalog") {
    reearth.modal.postMessage({
      type: action,
      payload: { rawCatalog, addedDatasets },
    });
  } else if (action === "helpPopupOpen") {
    reearth.popup.show(helpPopupHtml, { position: "right-start", offset: 4 });
  } else if (action === "initPopup") {
    reearth.ui.postMessage({ type: action });
  } else if (action === "initWelcome") {
    reearth.modal.postMessage({ type: "msgToModal", message: reearth.viewport.isMobile });
  } else if (action === "msgToPopup") {
    reearth.popup.postMessage({ type: "msgToPopup", message: payload });
  } else if (action === "msgFromPopup") {
    if (payload.height) {
      reearth.popup.update({ height: payload.height, width: reearth.viewport.width - 12 });
    }
  } else if (action === "popupClose") {
    reearth.popup.close();
    mobileDropdownIsOpen = false;
  } else if (action === "mapModalOpen") {
    reearth.modal.show(mapVideoHtml, { background: "transparent" });
  } else if (action === "clipModalOpen") {
    reearth.modal.show(clipVideoHtml, { background: "transparent" });
  }
});

reearth.on("update", () => {
  reearth.ui.postMessage({
    type: "extended",
    payload: reearth.widget.extended,
  });
});

reearth.on("resize", () => {
  // Modals
  if (welcomePageIsOpen) {
    reearth.modal.update({
      width: reearth.viewport.width,
      height: reearth.viewport.height,
    });
    reearth.modal.postMessage({ type: "msgToModal", payload: reearth.viewport.isMobile });
  }
  // Popups
  if (mobileDropdownIsOpen) {
    reearth.popup.update({
      width: reearth.viewport.width - 12,
    });
  }
});
