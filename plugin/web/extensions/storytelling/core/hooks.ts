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

  const [mode, setMode] = useState<Mode>("editor");
  const [size, setSize] = useState<Mode | "mini">("mini");
  const sizeRef = useRef<Mode | "mini">();
  sizeRef.current = size;
  const prevSizeRef = useRef<Mode | "mini">("mini");

  const handleMinimize = useCallback(() => {
    prevSizeRef.current = size;
    setSize(size => (size === "mini" ? mode : "mini"));
  }, [mode, size]);

  const handleSetMode = useCallback(
    (mode: Mode) => {
      prevSizeRef.current = size;
      setMode(mode);
      setSize(mode);
      // notice sidebar cancel story play when switch to edit mode
      if (mode === "editor" && storyId.current) {
        postMsg("cancelPlayStory", {
          id: storyId.current,
        });
        storyId.current = undefined;
      }
      if (mode === "editor") {
        setPlayerHeight(sizes.player.height);
      }
    },
    [size],
  );

  useEffect(() => {
    if (size === "mini") {
      setTimeout(() => {
        if (sizeRef.current === "mini") {
          postMsg("resize", [sizes.mini.width, sizes.mini.height, false]);
        }
      }, 500);
    } else if (size === "editor") {
      if (prevSizeRef.current === "player") {
        setTimeout(() => {
          if (sizeRef.current === "editor") {
            postMsg("resize", [sizes.editor.width, sizes.editor.height, true]);
          }
        }, 500);
      } else {
        postMsg("resize", [sizes.editor.width, sizes.editor.height, true]);
      }
    } else {
      postMsg("resize", [sizes.player.width, sizes.player.height, true]);
    }
  }, [size]);

  const [playerHeight, setPlayerHeight] = useState<number>(sizes.player.height);
  const playerHeightRef = useRef<number>();
  playerHeightRef.current = playerHeight;
  const playerPrevHeightRef = useRef<number>(sizes.player.height);
  const handlePlayerHeight = useCallback(
    (height: number) => {
      playerPrevHeightRef.current = playerHeight;
      setPlayerHeight(height);
    },
    [playerHeight],
  );

  useEffect(() => {
    if (mode === "player") {
      if (playerHeight > playerPrevHeightRef.current) {
        postMsg("resize", [undefined, playerHeight, true]);
      } else {
        setTimeout(() => {
          if (playerHeightRef.current === playerHeight) {
            postMsg("resize", [undefined, playerHeight, true]);
          }
        }, 500);
      }
    }
  }, [playerHeight, mode]);

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
    },
    [handleSetMode],
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
      setScenes(JSON.parse(scenes));
      handleSetMode("player");
    },
    [handleSetMode],
  );

  const handleCancelPlayStory = useCallback(
    ({ id }: CancelPlayStory["payload"]) => {
      if (storyId.current === id) {
        storyId.current = undefined;
        setScenes([]);
        if (size !== "mini") {
          handleMinimize();
        }
      }
    },
    [size, handleMinimize],
  );

  const shareStory = useCallback(() => {
    postMsg("shareStory", {
      scenes: JSON.stringify(scenes),
    });
  }, [scenes]);

  useEffect(() => {
    // mock scenes
    //     const scenes = [];
    //     for (let i = 1; i < 3; i += 1) {
    //       scenes.push({
    //         id: generateId(),
    //         title: `Title ${i}`,
    //         description: `# Header 1
    // ## Header 2
    // ### Header 3
    // ### Header 4
    // ### Header 5`,
    //         camera: undefined,
    //       });
    //       scenes.push({
    //         id: generateId(),
    //         title: `Title ${i}`,
    //         description: `# Header 1
    // ## Header 2
    // ### Header 3
    // ### Header 4
    // ### Header 5
    // ### Header 6
    // ### Header 7`,
    //         camera: undefined,
    //       });
    //     }
    //     setScenes(scenes);

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
    scenes,
    ConfigProvider,
    isMobile,
    playerHeight,
    handlePlayerHeight,
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
