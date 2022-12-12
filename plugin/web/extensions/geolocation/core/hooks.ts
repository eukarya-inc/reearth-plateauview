import { useCallback } from "react";

import { CurrentLocation } from "./types";
import { postMsg } from "./utils";

const initialLocation: CurrentLocation = {
  latitude: 5.70249,
  longitude: 39.7622,
  altitude: 5000,
};

export default () => {
  let currentLocation: CurrentLocation;

  const handleFlyToCurrentLocation = useCallback(() => {
    if (navigator.geolocation) {
      navigator.geolocation.watchPosition(
        function (position) {
          currentLocation = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            altitude: position.coords.altitude ?? 5000,
          };
        },
        function (error) {
          console.error("Error Code = " + error.code + " - " + error.message);
          currentLocation = { ...initialLocation };
        },
      );
    }
    postMsg({ action: "flyTo", payload: { currentLocation } });
  }, []);

  return {
    handleFlyToCurrentLocation,
  };
};
