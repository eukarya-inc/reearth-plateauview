import { useCallback } from "react";

import { CurrentLocation } from "./types";
import { postMsg } from "./utils";

const initialLocation: CurrentLocation = {
  latitude: 35.68539,
  longitude: 139.72675,
  altitude: 219202.886,
};

export default () => {
  let currentLocation: CurrentLocation;

  const handleFlyToCurrentLocation = useCallback(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        function (position) {
          currentLocation = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            altitude: position.coords.altitude ?? 5000,
          };
          postMsg({ action: "flyTo", payload: { currentLocation } });
        },
        function (error) {
          console.error("Error Code = " + error.code + " - " + error.message);
          currentLocation = { ...initialLocation };
          postMsg({ action: "flyTo", payload: { currentLocation } });
        },
      );
    }
  }, []);

  return {
    handleFlyToCurrentLocation,
  };
};
