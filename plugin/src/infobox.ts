import {
  PostMessageProps,
  PluginExtensionInstance,
  SavePublicSetting,
  PluginMessage,
} from "@web/extensions/infobox/types";

import html from "../dist/web/infobox/core/index.html?raw";

const reearth = (globalThis as any).reearth;

let currentLayerId: string | undefined = reearth.layers.selected;

reearth.ui.show(html);

let sidebarId: string;
const getSidebarId = () => {
  if (sidebarId) return;
  sidebarId = reearth.plugins.instances.find(
    (instance: PluginExtensionInstance) => instance.extensionId === "sidebar",
  )?.id;
};
getSidebarId();

const infoboxFetchFields = () => {
  getSidebarId();
  if (!sidebarId) return;
  reearth.plugins.postMessage(sidebarId, {
    action: "infoboxFetchFields",
    payload: currentLayerId,
  });
};

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  if (action === "init") {
    reearth.ui.postMessage({
      action: "getInEditor",
      payload: reearth.scene.inEditor,
    });
    infoboxFetchFields();
  } else if (action === "savePublicSetting") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      action: "savePublicSetting",
      payload,
    } as SavePublicSetting);
  }
});

reearth.on("pluginMessage", (pluginMessage: PluginMessage) => {
  if (pluginMessage.data.action === "infoboxFields") {
    reearth.ui.postMessage({
      action: "fillData",
      payload: {
        primitive: reearth.layers.selectedFeature,
        fileds: pluginMessage.data,
      },
    });
  }
});

reearth.on("select", () => {
  if (reearth.layers.selected?.id !== currentLayerId) {
    currentLayerId = reearth.layers.selected.id;
    infoboxFetchFields();
    reearth.ui.postMessage({
      action: "setLoading",
    });
  } else {
    reearth.ui.postMessage({
      action: "fillData",
      payload: {
        primitive: reearth.layers.selectedFeature,
        fileds: undefined,
      },
    });
  }
});
