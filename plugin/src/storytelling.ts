import {
  PostMessageProps,
  Viewport,
  PluginMessage,
  PluginExtensionInstance,
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
  reearth.ui.postMessage({
    type: "pluginMessage",
    payload: pluginMessage.data,
  });
});

reearth.on("message", ({ type, payload }: PostMessageProps) => {
  switch (type) {
    case "resize":
      reearth.ui.resize(...payload);
      break;
    case "minimize":
      reearth.ui.resize(...payload);
      break;
    case "captureScene":
      reearth.ui.postMessage({
        type: "captureScene",
        payload: reearth.camera.position,
      });
      break;
    case "viewScene":
      reearth.camera.flyTo(payload, { duration: 2 });
      break;
    case "recaptureScene":
      reearth.ui.postMessage({
        type: "recaptureScene",
        payload: { camera: reearth.camera.position, id: payload },
      });
      break;
    case "editScene":
      reearth.modal.show(storyeditorHtml, { background: "transparent", width: 580, height: 320 });
      reearth.modal.postMessage({
        type: "editScene",
        payload,
      });
      break;
    case "closeSceneEditor":
      reearth.modal.close();
      break;
    case "saveScene":
      reearth.ui.postMessage({
        type: "saveScene",
        payload,
      });
      reearth.modal.close();
      break;
    case "getViewport":
      reearth.ui.postMessage({
        type: "viewport",
        payload: reearth.viewport,
      });
      break;
    case "shareStory":
      getSidebarId();
      console.log("share", sidebarId, {
        type: "shareStory",
        payload,
      });
      if (!sidebarId) return;
      reearth.plugins.postMessage(sidebarId, {
        type: "shareStory",
        payload,
      });
      break;
    default:
      break;
  }
});

reearth.on("resize", (viewport: Viewport) => {
  reearth.ui.postMessage({
    type: "viewport",
    payload: viewport,
  });
});
