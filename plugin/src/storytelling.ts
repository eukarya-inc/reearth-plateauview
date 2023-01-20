import {
  PostMessageProps,
  Viewport,
  PluginMessage,
  PluginExtensionInstance,
  ShareStory,
  SaveStoryData,
  CancelPlayStory,
} from "@web/extensions/storytelling/core/types";

import html from "../dist/web/storytelling/core/index.html?raw";
import storyeditorHtml from "../dist/web/storytelling/modals/sceneeditor/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html, { width: 89, height: 40, extended: false });

let sidebarId: string;
const getSidebarId = () => {
  if (sidebarId) return;
  sidebarId = reearth.plugins.instances.find(
    (instance: PluginExtensionInstance) => instance.extensionId === "sidebar",
  )?.id;
};
getSidebarId();

reearth.on("pluginmessage", (pluginMessage: PluginMessage) => {
  reearth.ui.postMessage(pluginMessage.data);
});

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  if (action === "resize") {
    reearth.ui.resize(...payload);
  } else if (action === "captureScene") {
    reearth.ui.postMessage({
      action: "captureScene",
      payload: reearth.camera.position,
    });
  } else if (action === "viewScene") {
    reearth.camera.flyTo(payload, { duration: 1.5 });
  } else if (action === "recaptureScene") {
    reearth.ui.postMessage({
      action: "recaptureScene",
      payload: { camera: reearth.camera.position, id: payload },
    });
  } else if (action === "editScene") {
    reearth.modal.show(storyeditorHtml, { background: "transparent", width: 580, height: 320 });
    reearth.modal.postMessage({
      action: "editScene",
      payload,
    });
  } else if (action === "closeSceneEditor") {
    reearth.modal.close();
  } else if (action === "saveScene") {
    reearth.ui.postMessage({
      action: "saveScene",
      payload,
    });
    reearth.modal.close();
  } else if (action === "getViewport") {
    reearth.ui.postMessage({
      action: "viewport",
      payload: reearth.viewport,
    });
  } else if (action === "shareStory") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      action: "shareStory",
      payload,
    } as ShareStory);
  } else if (action === "saveStoryData") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      action: "saveStoryData",
      payload,
    } as SaveStoryData);
  } else if (action === "cancelPlayStory") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      action: "cancelPlayStory",
      payload,
    } as CancelPlayStory);
  }
});

reearth.on("resize", (viewport: Viewport) => {
  reearth.ui.postMessage({
    action: "viewport",
    payload: viewport,
  });
});
