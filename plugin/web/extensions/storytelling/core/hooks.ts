import { useCallback, useState, useRef, useEffect, useMemo } from "react";

import type { Camera, Story } from "./types";
import { postMsg, generateId } from "./utils";

export const size = {
  mini: {
    width: 89,
    height: 40,
  },
  extend: {
    width: undefined,
    height: 178,
  },
};

export type Mode = "editor" | "play";

export default () => {
  const [minimized, setMinimized] = useState<boolean>(true);
  const minimizedRef = useRef<boolean>(minimized);

  useEffect(() => {
    if (minimized) {
      setTimeout(() => {
        if (minimizedRef.current) {
          postMsg("minimize", [size.mini.width, size.mini.height, false]);
        }
      }, 500);
    } else {
      postMsg("minimize", [size.extend.width, size.extend.height, true]);
    }
  }, [minimized]);

  const handleMinimize = useCallback(() => {
    setMinimized(minimized => !minimized);
  }, []);

  const [mode, setMode] = useState<Mode>("editor");

  const contentTheme = useMemo(() => (mode === "editor" ? "light" : "grey"), [mode]);

  // Stories
  const [stories, setStories] = useState<Story[]>([]);

  const addStory = useCallback((story: Story) => {
    setStories(stories => [...stories, story]);
  }, []);

  const captureScene = useCallback(() => {
    postMsg("captureScene");
  }, []);
  const handleCaptureScene = useCallback((camera: Camera) => {
    addStory({
      id: generateId(),
      title: "TiTle",
      description: "",
      camera,
    });
  }, []);

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

  useEffect(() => {
    // mock stories
    const stories = [];
    for (let i = 1; i < 1; i += 1) {
      stories.push({
        id: generateId(),
        title: `Title Title Title Title dotsThreeVertical ${i}`,
        description:
          "This is the first capture. Here you can see the new building of this city. This is the first capture. Here you can see the new building of this city.",
        camera: undefined,
      });
    }
    setStories(stories);
  }, []);

  const eventListenerCallback = useCallback((e: MessageEvent<any>) => {
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
      default:
        break;
    }
  }, []);

  useEffect(() => {
    addEventListener("message", eventListenerCallback);
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []);

  return {
    minimized,
    handleMinimize,
    mode,
    setMode,
    contentTheme,
    stories,
    captureScene,
    viewStory,
    recapture,
    deleteStory,
    editStory,
  };
};
