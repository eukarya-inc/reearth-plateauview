import { PostMessageProps } from "@web/extensions/storytelling/core/types";

import html from "../dist/web/storytelling/core/index.html?raw";
import storyeditorHtml from "../dist/web/storytelling/modals/storyeditor/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html, { width: 89, height: 40, extended: false });

reearth.on("message", ({ type, payload }: PostMessageProps) => {
  switch (type) {
    case "minimize":
      reearth.ui.resize(...payload);
      break;
    case "captureScene":
      reearth.ui.postMessage({
        type: "captureScene",
        payload: reearth.camera.position,
      });
      break;
    case "viewStory":
      reearth.camera.flyTo(payload, { duration: 2 });
      break;
    case "recapture":
      reearth.ui.postMessage({
        type: "recapture",
        payload: { camera: reearth.camera.position, id: payload },
      });
      break;
    case "editStory":
      reearth.modal.show(storyeditorHtml, { background: "transparent", width: 580, height: 320 });
      reearth.modal.postMessage({
        type: "editStory",
        payload,
      });
      break;
    case "closeStoryEditor":
      reearth.modal.close();
      break;
    case "saveStory":
      reearth.ui.postMessage({
        type: "saveStory",
        payload,
      });
      reearth.modal.close();
      break;
    default:
      break;
  }
});
