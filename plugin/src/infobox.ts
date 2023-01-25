import { PostMessageProps } from "@web/extensions/infobox/types";

import html from "../dist/web/infobox/core/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html);

reearth.on("message", ({ action }: PostMessageProps) => {
  if (action === "getInEditor") {
    reearth.ui.postMessage({
      action: "getInEditor",
      payload: reearth.scene.inEditor,
    });
  }
});
