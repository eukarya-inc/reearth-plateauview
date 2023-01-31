import { array_move } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

export type groupItem = {
  title: string;
  group: string;
  id: number;
};

export type SwitchGroupObj = {
  title: string;
  groups: groupItem[];
};

export default (value: SwitchGroupObj) => {
  const [groups, updateGroups] = useState<SwitchGroupObj>(value);
  const [groupsTitle, setGroupsTitle] = useState(value.title);

  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setGroupsTitle(e.target.value);
  }, []);

  const handleMoveUp = useCallback((idx: number) => {
    if (idx === 0) return;
    updateGroups(l => {
      let newItems: groupItem[] | undefined = undefined;
      if (l.groups) {
        newItems = l.groups;
        array_move(newItems, idx, idx - 1);
      }
      return { ...l, items: newItems };
    });
  }, []);
  const handleMoveDown = useCallback(
    (idx: number) => {
      if (groups.groups && idx >= groups.groups.length - 1) return;
      updateGroups(l => {
        let newItems: groupItem[] | undefined = undefined;
        if (l.groups) {
          newItems = l.groups;
          array_move(newItems, idx, idx + 1);
        }
        return { ...l, items: newItems };
      });
    },
    [groups.groups],
  );

  const handleAdd = useCallback(() => {
    alert("ADD ITEM");
  }, []);

  const handleRemove = useCallback((idx: number) => {
    updateGroups(l => {
      let newItems: groupItem[] | undefined = undefined;
      if (l.groups) {
        newItems = l.groups.filter((_, idx2) => idx2 != idx);
      }
      return { ...l, items: newItems };
    });
  }, []);

  useEffect(() => {
    updateGroups(l => {
      return {
        ...l,
        groupsTitle,
      };
    });
  }, [groupsTitle]);

  return {
    groups,
    groupsTitle,
    handleTitleChange,
    handleMoveUp,
    handleMoveDown,
    handleAdd,
    handleRemove,
  };
};
