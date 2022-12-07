import { useCallback, useState } from "react";

import { CurrentLocationInfo } from "./types";
import { postMsg } from "./utils";

export default () => {
  const [currentLocation, setCurrentLocation] = useState<CurrentLocationInfo>();

  const handleFlyToCurrentLocation = useCallback(() => {
    if (navigator.geolocation) {
      navigator.geolocation.watchPosition(
        function (position) {
          setCurrentLocation({
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            altitude: position.coords.altitude ?? 5000,
          });
        },
        function (error) {
          console.error("Error Code = " + error.code + " - " + error.message);
        },
      );
    }
    postMsg({ action: "flyTo", payload: { currentLocation } });
  }, []);

  return {
    handleFlyToCurrentLocation,
  };
};
