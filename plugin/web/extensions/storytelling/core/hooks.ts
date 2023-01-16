import { ConfigProvider } from "@web/sharedComponents";
import update from "immutability-helper";
import { useCallback, useState, useRef, useEffect } from "react";

import type { Camera, Story, Viewport } from "./types";
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

  const [mode, setMode] = useState<Mode>(isMobile ? "player" : "editor");
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

  // Stories
  const [stories, setStories] = useState<Story[]>([]);

  const addStory = useCallback((story: Story) => {
    setStories(stories => [...stories, story]);
    postMsg("editStory", { id: story.id, title: story.title, description: story.description });
  }, []);

  const captureScene = useCallback(() => {
    postMsg("captureScene");
  }, []);
  const handleCaptureScene = useCallback(
    (camera: Camera) => {
      addStory({
        id: generateId(),
        title: "",
        description: "",
        camera,
      });
    },
    [addStory],
  );

  const viewStory = useCallback((camera: Camera) => {
    postMsg("viewStory", camera);
  }, []);

  const recapture = useCallback((id: string) => {
    postMsg("recapture", id);
  }, []);
  const handleRecapture = useCallback(({ camera, id }: { camera: Camera; id: string }) => {
    setStories(stories => {
      const story = stories.find(story => story.id === id);
      if (story) {
        story.camera = camera;
      }
      return [...stories];
    });
  }, []);

  const deleteStory = useCallback((id: string) => {
    setStories(stories => {
      const index = stories.findIndex(story => story.id === id);
      if (index !== -1) {
        stories.splice(index, 1);
      }
      return [...stories];
    });
  }, []);

  const editStory = useCallback(
    (id: string) => {
      const story = stories.find(story => story.id === id);
      if (story) {
        postMsg("editStory", story);
      }
    },
    [stories],
  );

  const saveStory = useCallback((storyInfo: Omit<Story, "camera">) => {
    setStories(stories => {
      const story = stories.find(story => story.id === storyInfo.id);
      if (story) {
        story.title = storyInfo.title;
        story.description = storyInfo.description;
      }
      return [...stories];
    });
  }, []);

  const moveStory = useCallback((dragIndex: number, hoverIndex: number) => {
    setStories((prevStories: Story[]) =>
      update(prevStories, {
        $splice: [
          [dragIndex, 1],
          [hoverIndex, 0, prevStories[dragIndex] as Story],
        ],
      }),
    );
  }, []);

  const share = useCallback(() => {
    postMsg("shareStoryTelling", {
      storyTellingId: "",
      stories,
    });
  }, [stories]);

  useEffect(() => {
    // mock stories
    const stories = [];
    for (let i = 1; i < 1; i += 1) {
      stories.push({
        id: generateId(),
        title: `Title ${i}`,
        description: `# Header 1
## Header 2
### Header 3`,
        camera: undefined,
      });
    }
    setStories(stories);

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
        case "recapture":
          handleRecapture(e.data.payload);
          break;
        case "saveStory":
          saveStory(e.data.payload);
          break;
        case "viewport":
          handleViewportResize(e.data.payload);
          break;
        default:
          break;
      }
    },
    [handleCaptureScene, handleRecapture, saveStory, handleViewportResize],
  );

  useEffect(() => {
    addEventListener("message", onMessage);
    return () => {
      removeEventListener("message", onMessage);
    };
  }, [onMessage]);

  return {
    size,
    handleMinimize,
    mode,
    handleSetMode,
    stories,
    captureScene,
    viewStory,
    recapture,
    deleteStory,
    editStory,
    moveStory,
    share,
    ConfigProvider,
    isMobile,
  };
};
