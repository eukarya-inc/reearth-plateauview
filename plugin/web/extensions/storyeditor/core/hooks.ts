import { useCallback, useState, useRef, useEffect, useMemo } from "react";

import type { Camera } from "./types";
import { postMsg } from "./utils";

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

export type Story = {
  camera: Camera;
  title?: string;
  description?: string;
};

export default () => {
  const [minimized, setMinimized] = useState<boolean>(true);
  const minimizedRef = useRef<boolean>();
  minimizedRef.current = minimized;

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

  // Data
  const [stories, setStories] = useState<Story[]>([]);

  const addStory = useCallback((story: Story) => {
    setStories(stories => [...stories, story]);
  }, []);

  const loadStories = useCallback((stories: Story[]) => {
    setStories(stories);
  }, []);

  // mockup
  const mockStories: Story[] = [
    {
      camera: {
        lng: 139,
        lat: 40,
        height: 10,
        heading: 0,
        pitch: -90,
        roll: 0,
        fov: 0.5,
      },
      title: "1",
      description: "Some text here.",
    },
  ];

  useEffect(() => {
    loadStories(mockStories);
  }, []);

  const eventListenerCallback = (e: MessageEvent<any>) => {
    if (e.source !== parent) return;
    if (e.data.type === "initMsg") {
      console.log("initmsg", e.data.payload);
    }
  };

  // useEffect(() => {
  //   addEventListener("message", e => eventListenerCallback(e));
  //   console.log("plugin add message event listener");
  //   return () => {
  //     removeEventListener("message", eventListenerCallback);
  //   };
  // }, []);

  const [inited, setInited] = useState(false);

  if (!inited) {
    console.log("plugin add message event listener");
    addEventListener("message", e => eventListenerCallback(e));
    setInited(true);
  }

  return {
    minimized,
    handleMinimize,
    mode,
    setMode,
    contentTheme,
    stories,
    addStory,
  };
};
