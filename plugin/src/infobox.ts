import {
  PostMessageProps,
  PluginExtensionInstance,
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

const infoboxFieldsFetch = () => {
  getSidebarId();
  if (!sidebarId) return;
  reearth.plugins.postMessage(sidebarId, {
    action: "infoboxFieldsFetch",
    payload: currentLayerId,
  });
};

reearth.on("message", ({ action, payload }: PostMessageProps) => {
  if (action === "init") {
    reearth.ui.postMessage({
      action: "getInEditor",
      payload: reearth.scene.inEditor,
    });
    infoboxFieldsFetch();
  } else if (action === "saveFields") {
    getSidebarId();
    if (!sidebarId) return;
    reearth.plugins.postMessage(sidebarId, {
      action: "infoboxFieldsSave",
      payload,
    });
  }
});

reearth.on("pluginmessage", (pluginMessage: PluginMessage) => {
  reearth.ui.postMessage(pluginMessage.data);
  if (pluginMessage.data.action === "infoboxFields") {
    if (reearth.layers.selectedFeature) {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { attributes, ...rawProperties } = reearth.layers.selectedFeature.properties;
      const properties: { key: string; value?: any }[] = [];
      Object.keys(rawProperties).forEach(key => {
        properties.push({
          key,
          value: rawProperties[key],
        });
      });
      reearth.ui.postMessage({
        action: "fillData",
        payload: {
          feature: {
            properties,
          },
          fields: pluginMessage.data.payload,
        },
      });
    }
  }
});

// reearth.on("pluginMessage", (pluginMessage: PluginMessage) => {
//   console.log("get data from sidebar", pluginMessage);
//   if (pluginMessage.data.action === "infoboxFields") {
//     console.log(
//       "infoboxFields: get data from sidebar",
//       reearth.layers.selectedFeature,
//       pluginMessage.data,
//     );
//     reearth.ui.postMessage({
//       action: "fillData",
//       payload: {
//         feature: reearth.layers.selectedFeature,
//         fileds: pluginMessage.data,
//       },
//     });
//   }
// });

reearth.on("select", (layerId: string) => {
  currentLayerId = layerId;

  getSidebarId();
  if (!sidebarId) return;
  reearth.plugins.postMessage(sidebarId, {
    action: "infoboxFieldsFetch",
    payload: currentLayerId,
  });

  reearth.ui.postMessage({
    action: "setLoading",
  });

  // if (reearth.layers.selected?.id !== currentLayerId) {
  //   currentLayerId = reearth.layers.selected.id;
  //   infoboxFieldsFetch();
  //   reearth.ui.postMessage({
  //     action: "setLoading",
  //   });
  // } else {
  //   reearth.ui.postMessage({
  //     action: "fillData",
  //     payload: {
  //       feature: reearth.layers.selectedFeature,
  //     },
  //   });
  // }
});
