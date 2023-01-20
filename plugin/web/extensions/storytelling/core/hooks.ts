import { ConfigProvider } from "@web/sharedComponents";
import update from "immutability-helper";
import { useCallback, useState, useRef, useEffect } from "react";

import type {
  Camera,
  Scene,
  Viewport,
  EditStory,
  SaveStory,
  DeleteStory,
  PlayStory,
  CancelPlayStory,
} from "./types";
import { postMsg, generateId } from "./utils";

export const sizes = {
  mini: {
    width: 89,
    height: 40,
  },
  editor: {
    width: undefined,
    height: 178,
  },
  player: {
    width: undefined,
    height: 195,
  },
};

export type Mode = "editor" | "player";
export type Size = { width: number | undefined; height: number };

export default () => {
  const [isMobile, setIsMobile] = useState<boolean>(window.innerWidth <= 768);
  const handleViewportResize = useCallback(
    (viewport: Viewport) => {
      if (viewport.isMobile !== isMobile) {
        setIsMobile(viewport.isMobile);
      }
    },
    [isMobile],
  );

  const [mode, setMode] = useState<Mode>("player");

  const [minimized, setMinimized] = useState<boolean>(true);
  const minimizedRef = useRef<boolean>(minimized);
  minimizedRef.current = minimized;

  const [size, setSize] = useState<Size>(sizes.mini);
  const sizeRef = useRef<Size>(size);
  sizeRef.current = size;
  const prevSizeRef = useRef<Size>(sizes.mini);

  const [playerHeight, setPlayerHeight] = useState<number>(sizes.player.height);
  const playerHeightRef = useRef<number>(playerHeight);
  playerHeightRef.current = playerHeight;

  const handleMinimize = useCallback(() => {
    setMinimized(minimized => !minimized);
  }, []);

  const handleSetMode = useCallback((mode: Mode) => {
    setMode(mode);
    if (mode === "editor") {
      setPlayerHeight(0);
      if (storyId.current) {
        postMsg("cancelPlayStory", {
          id: storyId.current,
        });
        storyId.current = undefined;
      }
    }
  }, []);

  useEffect(() => {
    prevSizeRef.current = sizeRef.current;

    setSize(
      minimized
        ? sizes.mini
        : mode === "editor"
        ? sizes.editor
        : { width: undefined, height: playerHeight },
    );
  }, [minimized, mode, playerHeight]);

  useEffect(() => {
    if (size.height > prevSizeRef.current.height) {
      postMsg("resize", [size.width, size.height, !minimizedRef.current]);
    } else if (size.height < prevSizeRef.current.height) {
      setTimeout(() => {
        if (sizeRef.current === size) {
          postMsg("resize", [size.width, size.height, !minimizedRef.current]);
        }
      }, 500);
    }
  }, [size]);

  // scenes
  const storyId = useRef<string>();
  const [scenes, setScenes] = useState<Scene[]>([]);

  const addScene = useCallback((scene: Scene) => {
    setScenes(scenes => [...scenes, scene]);
    postMsg("editScene", { id: scene.id, title: scene.title, description: scene.description });
  }, []);

  const captureScene = useCallback(() => {
    postMsg("captureScene");
  }, []);
  const handleCaptureScene = useCallback(
    (camera: Camera) => {
      addScene({
        id: generateId(),
        title: "",
        description: "",
        camera,
      });
    },
    [addScene],
  );

  const viewScene = useCallback((camera: Camera) => {
    postMsg("viewScene", camera);
  }, []);

  const recaptureScene = useCallback((id: string) => {
    postMsg("recaptureScene", id);
  }, []);
  const handleRecaptureScene = useCallback(({ camera, id }: { camera: Camera; id: string }) => {
    setScenes(scenes => {
      const scene = scenes.find(scene => scene.id === id);
      if (scene) {
        scene.camera = camera;
      }
      return [...scenes];
    });
  }, []);

  const deleteScene = useCallback((id: string) => {
    setScenes(scenes => {
      const index = scenes.findIndex(scene => scene.id === id);
      if (index !== -1) {
        scenes.splice(index, 1);
      }
      return [...scenes];
    });
  }, []);

  const editScene = useCallback(
    (id: string) => {
      const scene = scenes.find(scene => scene.id === id);
      if (scene) {
        postMsg("editScene", scene);
      }
    },
    [scenes],
  );

  const saveScene = useCallback((sceneInfo: Omit<Scene, "camera">) => {
    setScenes(scenes => {
      const scene = scenes.find(scene => scene.id === sceneInfo.id);
      if (scene) {
        scene.title = sceneInfo.title;
        scene.description = sceneInfo.description;
      }
      return [...scenes];
    });
  }, []);

  const moveScene = useCallback((dragIndex: number, hoverIndex: number) => {
    setScenes((prevScenes: Scene[]) =>
      update(prevScenes, {
        $splice: [
          [dragIndex, 1],
          [hoverIndex, 0, prevScenes[dragIndex] as Scene],
        ],
      }),
    );
  }, []);

  const clearStory = useCallback(() => {
    storyId.current = undefined;
    setScenes([]);
  }, []);

  const handleEditStory = useCallback(
    ({ id, scenes }: EditStory["payload"]) => {
      storyId.current = id;
      setScenes(scenes ? JSON.parse(scenes) : []);
      handleSetMode("editor");
      if (minimized) {
        handleMinimize();
      }
    },
    [handleSetMode, minimized, handleMinimize],
  );

  const handleSaveStory = useCallback(
    ({ id }: SaveStory["payload"]) => {
      postMsg("saveStoryData", {
        id,
        scenes: JSON.stringify(scenes),
      });
    },
    [scenes],
  );

  const handleDeleteStory = useCallback(({ id }: DeleteStory["payload"]) => {
    if (storyId.current === id) {
      storyId.current = undefined;
      setScenes([]);
    }
  }, []);

  const handlePlayStory = useCallback(
    ({ id, scenes }: PlayStory["payload"]) => {
      storyId.current = id;
      setScenes(JSON.parse(scenes ?? "[]"));
      handleSetMode("player");
      if (minimized) {
        handleMinimize();
      }
    },
    [handleSetMode, minimized, handleMinimize],
  );

  const handleCancelPlayStory = useCallback(
    ({ id }: CancelPlayStory["payload"]) => {
      if (storyId.current === id) {
        storyId.current = undefined;
        setScenes([]);
        if (!minimized) {
          handleMinimize();
        }
      }
    },
    [minimized, handleMinimize],
  );

  const shareStory = useCallback(() => {
    postMsg("shareStory", {
      scenes: JSON.stringify(scenes),
    });
  }, [scenes]);

  useEffect(() => {
    // theme
    const themeColor = "#00BEBE";
    document.documentElement.style.setProperty("--theme-color", themeColor);
    ConfigProvider.config({
      theme: {
        primaryColor: themeColor,
      },
    });

    // viewport
    postMsg("getViewport");
  }, []);

  const [contentWidth, setContentWidth] = useState<number>(document.body.clientWidth);
  useEffect(() => {
    const viewportResizeObserver = new ResizeObserver(entries => {
      const [entry] = entries;
      let width: number | undefined;

      if (entry.contentBoxSize) {
        const contentBoxSize = Array.isArray(entry.contentBoxSize)
          ? entry.contentBoxSize[0]
          : entry.contentBoxSize;
        width = contentBoxSize.inlineSize;
      } else if (entry.contentRect) {
        width = entry.contentRect.width;
      }

      setContentWidth(width ?? document.body.clientWidth);
    });

    viewportResizeObserver.observe(document.body);

    return () => {
      viewportResizeObserver.disconnect();
    };
  }, []);

  const onMessage = useCallback(
    (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      switch (e.data.type) {
        case "captureScene":
          handleCaptureScene(e.data.payload);
          break;
        case "recaptureScene":
          handleRecaptureScene(e.data.payload);
          break;
        case "saveScene":
          saveScene(e.data.payload);
          break;
        case "viewport":
          handleViewportResize(e.data.payload);
          break;
        case "editStory":
          handleEditStory(e.data.payload);
          break;
        case "saveStory":
          handleSaveStory(e.data.payload);
          break;
        case "deleteStory":
          handleDeleteStory(e.data.payload);
          break;
        case "playStory":
          handlePlayStory(e.data.payload);
          break;
        case "cancelPlayStory":
          handleCancelPlayStory(e.data.payload);
          break;
        default:
          break;
      }
    },
    [
      handleCaptureScene,
      handleRecaptureScene,
      saveScene,
      handleViewportResize,
      handleEditStory,
      handleSaveStory,
      handleDeleteStory,
      handlePlayStory,
      handleCancelPlayStory,
    ],
  );

  useEffect(() => {
    addEventListener("message", onMessage);
    return () => {
      removeEventListener("message", onMessage);
    };
  }, [onMessage]);

  return {
    size,
    mode,
    minimized,
    scenes,
    ConfigProvider,
    isMobile,
    playerHeight,
    contentWidth,
    setPlayerHeight,
    handleMinimize,
    handleSetMode,
    captureScene,
    viewScene,
    recaptureScene,
    deleteScene,
    editScene,
    moveScene,
    clearStory,
    shareStory,
  };
};
