import { postMsg } from "@web/extensions/pedestrian/utils";
import { useCallback, useEffect, useState } from "react";

export default () => {
  const [mainButtonText, setMainButtonText] = useState<"Start" | "Exit">("Start");
  const [mode, setMode] = useState<"ready" | "picking" | "pedestrian">("ready");

  const [moveForwardOn, setMoveForwardOn] = useState(false);
  const [moveBackwardOn, setMoveBackwardOn] = useState(false);
  const [moveLeftOn, setMoveLeftOn] = useState(false);
  const [moveRightOn, setMoveRightOn] = useState(false);
  const [moveUpOn, setMoveUpOn] = useState(false);
  const [moveDownOn, setMoveDownOn] = useState(false);

  const handleMoveForwardClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveForwardOn : enable;
      setMoveForwardOn(on);
      postMsg("cameraMove", { moveType: "moveForward", on });
      if (on && moveBackwardOn) {
        setMoveBackwardOn(false);
      }
    },
    [moveForwardOn, moveBackwardOn],
  );

  const handleMoveBackwardClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveBackwardOn : enable;
      setMoveBackwardOn(on);
      postMsg("cameraMove", { moveType: "moveBackward", on });
      if (on && moveForwardOn) {
        setMoveForwardOn(false);
      }
    },
    [moveBackwardOn, moveForwardOn],
  );

  const handleMoveLeftClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveLeftOn : enable;
      setMoveLeftOn(on);
      postMsg("cameraMove", { moveType: "moveLeft", on });
      if (on && moveRightOn) {
        setMoveRightOn(false);
      }
    },
    [moveLeftOn, moveRightOn],
  );

  const handleMoveRightClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveRightOn : enable;
      setMoveRightOn(on);
      postMsg("cameraMove", { moveType: "moveRight", on });
      if (on && moveLeftOn) {
        setMoveLeftOn(false);
      }
    },
    [moveRightOn, moveLeftOn],
  );

  const handleMoveUpClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveUpOn : enable;
      setMoveUpOn(on);
      postMsg("cameraMove", { moveType: "moveUp", on });
      if (on && moveDownOn) {
        setMoveDownOn(false);
      }
    },
    [moveUpOn, moveDownOn],
  );

  const handleMoveDownClick = useCallback(
    (enable?: boolean) => {
      const on = enable === undefined ? !moveDownOn : enable;
      setMoveDownOn(on);
      postMsg("cameraMove", { moveType: "moveDown", on });
      if (on && moveUpOn) {
        setMoveUpOn(false);
      }
    },
    [moveDownOn, moveUpOn],
  );

  const onExit = useCallback(() => {
    setMode("ready");
    setMainButtonText("Start");
    postMsg("pedestrianExit");
  }, []);

  const onPicking = useCallback(() => {
    setMode("picking");
    setMainButtonText("Exit");
    postMsg("pickingStart");
  }, []);

  const handlePickingDone = useCallback(() => {
    setMode("pedestrian");
  }, []);

  const onMainButtonClick = useCallback(() => {
    if (mode === "ready") {
      onPicking();
    } else {
      onExit();
    }
  }, [mode, onPicking, onExit]);

  const onKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (mode !== "pedestrian") return;
      switch (e.code) {
        case "KeyW":
          handleMoveForwardClick(true);
          break;
        case "KeyA":
          handleMoveLeftClick(true);
          break;
        case "KeyS":
          handleMoveBackwardClick(true);
          break;
        case "KeyD":
          handleMoveRightClick(true);
          break;
        case "Space":
          handleMoveUpClick(true);
          break;
        case "ShiftLeft":
        case "ShiftRight":
          handleMoveDownClick(true);
          break;
        default:
          return undefined;
      }
    },
    [
      mode,
      handleMoveForwardClick,
      handleMoveBackwardClick,
      handleMoveLeftClick,
      handleMoveRightClick,
      handleMoveUpClick,
      handleMoveDownClick,
    ],
  );

  const onKeyUp = useCallback(
    (e: KeyboardEvent) => {
      if (mode !== "pedestrian") return;
      switch (e.code) {
        case "KeyW":
          handleMoveForwardClick(false);
          break;
        case "KeyA":
          handleMoveLeftClick(false);
          break;
        case "KeyS":
          handleMoveBackwardClick(false);
          break;
        case "KeyD":
          handleMoveRightClick(false);
          break;
        case "Space":
          handleMoveUpClick(false);
          break;
        case "ShiftLeft":
        case "ShiftRight":
          handleMoveDownClick(false);
          break;
        default:
          return undefined;
      }
    },
    [
      mode,
      handleMoveForwardClick,
      handleMoveBackwardClick,
      handleMoveLeftClick,
      handleMoveRightClick,
      handleMoveUpClick,
      handleMoveDownClick,
    ],
  );

  const onClose = useCallback(() => {
    if (mode !== "ready") {
      onExit();
    }
    postMsg("pedestrianClose");
  }, [mode, onExit]);

  const onMessage = useCallback(
    (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      switch (e.data.type) {
        case "pickingDone":
          handlePickingDone();
          break;
        default:
          break;
      }
    },
    [handlePickingDone],
  );

  useEffect(() => {
    document.documentElement.style.setProperty("--theme-color", "#00BEBE");
  }, []);

  useEffect(() => {
    (globalThis as any).parent.document.addEventListener("keydown", onKeyDown, false);
    (globalThis as any).parent.document.addEventListener("keyup", onKeyUp, false);

    return () => {
      (globalThis as any).parent.document.removeEventListener("keydown", onKeyDown);
      (globalThis as any).parent.document.removeEventListener("keyup", onKeyUp);
    };
  }, [onKeyDown, onKeyUp]);

  useEffect(() => {
    addEventListener("message", onMessage);
    return () => {
      removeEventListener("message", onMessage);
    };
  }, [onMessage]);

  return {
    mode,
    mainButtonText,
    moveForwardOn,
    moveBackwardOn,
    moveLeftOn,
    moveRightOn,
    moveUpOn,
    moveDownOn,
    handleMoveForwardClick,
    handleMoveBackwardClick,
    handleMoveLeftClick,
    handleMoveRightClick,
    handleMoveUpClick,
    handleMoveDownClick,
    onClose,
    onMainButtonClick,
  };
};
