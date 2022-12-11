import { useCallback } from "react";

import { CurrentLocationInfo } from "./types";
import { postMsg } from "./utils";

const initialLocation: CurrentLocationInfo = {
  latitude: 5.70249,
  longitude: 39.7622,
  altitude: 5000,
};

const goSuccess = (position: GeolocationPosition) => {
  return {
    latitude: position.coords.latitude,
    longitude: position.coords.longitude,
    altitude: position.coords.altitude ?? 5000,
  };
};

const goError = (err: GeolocationPositionError) => {
  console.error("Error Code = " + err.code + " - " + err.message);
  return initialLocation;
};

const getPosition = () =>
  new Promise((resolve, reject) =>
    navigator.geolocation.getCurrentPosition(
      pos => resolve(goSuccess(pos)),
      err => reject(goError(err)),
    ),
  );
export default () => {
  let currentLocation: CurrentLocationInfo;

  const handleFlyToCurrentLocation = useCallback(() => {
    getPosition()
      .then((position: CurrentLocationInfo | any) => {
        currentLocation = { ...position };
        postMsg({ action: "flyTo", payload: { currentLocation } });
      })
      .catch(error => console.error("error:", error.message));
  }, []);

  return {
    handleFlyToCurrentLocation,
  };
};
