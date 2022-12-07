import { PostMessageProps } from "@web/extensions/geolocation/core/types";

import html from "../dist/web/geolocation/core/index.html?raw";

const reearth = (globalThis as any).reearth;

reearth.ui.show(html, { width: 44, height: 44 });

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  if (action === "flyTo") {
    if (
      payload.currentLocation.latitude !== undefined &&
      payload.currentLocation.longitude !== undefined &&
      payload.currentLocation.altitude !== undefined
    ) {
      reearth.layers.add(
        {
          extensionId: "marker",
          isVisible: true,
          title: "myLocation",
          property: {
            default: {
              location: {
                lat: payload.currentLocation.latitude,
                lng: payload.currentLocation.longitude,
              },
              pointColor: "#12BDE2",
              style: "point",
            },
            customs: {
              id: "myLocation",
            },
          },
        },
        false,
      );

      const initCameraPos = reearth.camera.position;

      reearth.camera.flyTo(
        {
          lat: payload.currentLocation.latitude,
          lng: payload.currentLocation.longitude,
          height: payload.currentLocation.altitude,
          heading: initCameraPos?.heading ?? 0,
          pitch: -Math.PI / 2,
          roll: 0,
          fov: initCameraPos?.fov ?? 0.75,
        },
        {
          duration: 2,
        },
      );
    }
  }
});
