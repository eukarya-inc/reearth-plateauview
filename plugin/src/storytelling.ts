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

reearth.on("message", ({ type, payload }: PostMessageProps) => {
  if (type === "resize") {
    reearth.ui.resize(...payload);
  } else if (type === "captureScene") {
    reearth.ui.postMessage({
      type: "captureScene",
      payload: reearth.camera.position,
    });
  } else if (type === "viewScene") {
    reearth.camera.flyTo(payload, { duration: 1.5 });
  } else if (type === "recaptureScene") {
    reearth.ui.postMessage({
      type: "recaptureScene",
      payload: { camera: reearth.camera.position, id: payload },
    });
  } else if (type === "editScene") {
    reearth.modal.show(storyeditorHtml, { background: "transparent", width: 580, height: 320 });
    reearth.modal.postMessage({
      type: "editScene",
      payload,
    });
  } else if (type === "closeSceneEditor") {
    reearth.modal.close();
  } else if (type === "saveScene") {
    reearth.ui.postMessage({
      type: "saveScene",
      payload,
    });
    reearth.modal.close();
  } else if (type === "getViewport") {
    reearth.ui.postMessage({
      type: "viewport",
      payload: reearth.viewport,
    });
  } else if (type === "shareStory") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      type: "shareStory",
      payload,
    } as ShareStory);
  } else if (type === "saveStoryData") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      type: "saveStoryData",
      payload,
    } as SaveStoryData);
  } else if (type === "cancelPlayStory") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      type: "cancelPlayStory",
      payload,
    } as CancelPlayStory);
  }
});

reearth.on("resize", (viewport: Viewport) => {
  reearth.ui.postMessage({
    type: "viewport",
    payload: viewport,
  });
});
