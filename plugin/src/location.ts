import { PostMessageProps, MouseEventData } from "@web/extensions/location/core/types";

import html from "../dist/web/location/core/index.html?raw";
import GoogleAnalyticsHtml from "../dist/web/location/modals/googleAnalytics/index.html?raw";
import TerrainHtml from "../dist/web/location/modals/terrain/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html);

reearth.on("mousemove", (mousedata: MouseEventData) => {
  reearth.ui.postMessage(
    {
      type: "mousedata",
      payload: mousedata,
    },
    "*",
  );
});
reearth.on("cameramove", () => {
  reearth.ui.postMessage({
    type: "getLocations",
    payload: {
      point1: reearth.scene.getLocationFromScreenPosition(
        reearth.viewport.width / 2,
        reearth.viewport.height - 1,
      ),
      point2: reearth.scene.getLocationFromScreenPosition(
        reearth.viewport.width / 2 + 1,
        reearth.viewport.height - 1,
      ),
    },
  });
});
reearth.on("message", ({ action }: PostMessageProps) => {
  if (action === "modal-google-open") {
    reearth.modal.show(GoogleAnalyticsHtml, { background: "transparent" });
  } else if (action === "modal-terrain-open") {
    reearth.modal.show(TerrainHtml, { background: "transparent" });
  } else if (action === "modal-close") {
    reearth.modal.close();
  }
});
