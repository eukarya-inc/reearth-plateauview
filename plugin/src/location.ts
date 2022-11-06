import GoogleAnalyticstModal from "@web/extensions/location/modals/googleAnalyticsModal";
import TerrainModal from "@web/extensions/location/modals/terrainModal";
import { PostMessageProps, MouseEventData } from "@web/extensions/location/types";

import html from "../dist/web/location/core/index.html?raw";

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
reearth.on("message", ({ action }: PostMessageProps) => {
  if (action === "modal-google-open") {
    reearth.modal.show(GoogleAnalyticstModal);
  } else if (action === "modal-Terrain-open") {
    reearth.modal.show(TerrainModal);
  } else if (action === "modal-close") {
    reearth.modal.close();
  }
});
