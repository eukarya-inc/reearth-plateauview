import { PostMessageProps } from "@web/extensions/storytelling/core/types";

import html from "../dist/web/storytelling/core/index.html?raw";
// import dataCatalogHtml from "../dist/web/storytelling/modals/datacatalog/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html, { width: 89, height: 40, extended: false });
reearth.ui.postMessage({
  type: "initMsg",
  payload: "ui inited.",
});

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  if (action === "minimize") {
    reearth.ui.resize(...payload);
  }
  // if (action === "updateOverrides") {
  //   reearth.visualizer.overrideProperty(payload);
  // } else if (action === "screenshot" || action === "screenshot-save") {
  //   reearth.ui.postMessage({
  //     type: action,
  //     payload: reearth.scene.captureScreen(),
  //   });
  // } else if (action === "msgFromModal") {
  //   reearth.ui.postMessage({ type: action, payload });
  // } else if (action === "modal-open") {
  //   reearth.modal.show(dataCatalogHtml, { background: "transparent" });
  // } else if (action === "modal-close") {
  //   reearth.modal.close();
  // }
});

// reearth.on("update", () => {
//   reearth.ui.postMessage({
//     type: "extended",
//     payload: reearth.widget.extended,
//   });
// });
