import { postMsg } from "@web/extensions/storytelling/core/utils";
import { useCallback, useEffect, useRef } from "react";

export default () => {
  const storyId = useRef<string>();
  const titleRef = useRef<HTMLInputElement>(null);
  const descriptionRef = useRef<HTMLTextAreaElement>(null);

  const onCancel = useCallback(() => {
    postMsg("closeStoryEditor");
  }, []);

  const onSave = useCallback(() => {
    postMsg("saveStory", {
      id: storyId.current,
      title: titleRef.current?.value,
      description: descriptionRef.current?.value,
    });
  }, []);

  useEffect(() => {
    if ((window as any).editStory && titleRef.current && descriptionRef.current) {
      storyId.current = (window as any).editStory.id;
      titleRef.current.value = (window as any).editStory.title;
      descriptionRef.current.value = (window as any).editStory.description;
    }
  }, []);

  return {
    titleRef,
    descriptionRef,
    onCancel,
    onSave,
  };
};
