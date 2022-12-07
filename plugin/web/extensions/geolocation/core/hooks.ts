import { useCallback, useEffect, useState } from "react";

import { CurrentLocationInfo } from "./types";
import { postMsg } from "./utils";

const initialLocation: CurrentLocationInfo = {
  latitude: 35.70249,
  longitude: 139.7622,
  altitude: 5000,
};

export default () => {
  const [currentLocation, setCurrentLocation] = useState<CurrentLocationInfo>(initialLocation);

  const handleFlyToCurrentLocation = useCallback(() => {
    postMsg({ action: "flyTo", payload: { currentLocation } });
  }, []);

  useEffect(() => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
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
  }, []);

  return {
    handleFlyToCurrentLocation,
  };
};
