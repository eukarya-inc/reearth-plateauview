import { postMsg } from ".";

export const getOverriddenLayerByDataID = (dataID: string | undefined) =>
  new Promise<any>(resolve => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "getOverriddenLayerByDataID") {
        removeEventListener("message", eventListenerCallback);
        resolve(e.data.payload.overriddenLayer);
      }
    };
    addEventListener("message", eventListenerCallback);
    postMsg({
      action: "getOverriddenLayerByDataID",
      payload: {
        dataID,
      },
    });
  });
