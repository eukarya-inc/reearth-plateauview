import { array_move } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { groupItem, SwitchGroup } from "../../types";

export default (value: SwitchGroup) => {
  const [switchGroupObj, updateGroups] = useState<SwitchGroup>(value);
  const [groupsTitle, setGroupsTitle] = useState(value.title);
  const [currentGroup, setCurrentGroup] = useState<groupItem>(value.groups[0]);
  const [modifiedGroups, updateModifiedGroups] = useState<SwitchGroup | undefined>();

  //initialize the helper array (modifiedGroups)
  useEffect(() => {
    updateModifiedGroups({ type: "switchGroup", title: "Temp", groups: [] });
  }, []);

  //add empty item each time we press on add item
  const handAddItem = useCallback(() => {
    updateModifiedGroups(l => {
      if (!l) return;
      l.groups.push({ id: 0, group: "", title: "" });
      return { ...l, groups: l.groups };
    });
  }, []);

  //modiy the group in the helper array
  const handleModifyGroup = (group: string, index: number) => {
    updateModifiedGroups(l => {
      let newArray: groupItem[] | undefined = undefined;
      if (!l || !l.groups) return;
      newArray = l.groups;
      newArray[index].group = group;
      return { ...l, groups: newArray };
    });
  };

  //modify the title in the helper array
  const handleModifyGroupTitle = (title: string, index: number) => {
    updateModifiedGroups(l => {
      let newArray: groupItem[] | undefined = undefined;
      if (!l || !l.groups) return;
      newArray = l.groups;
      newArray[index].title = title;
      return { ...l, groups: newArray };
    });
  };

  //modify the title of the switch group in config and reflect it on main switch group field component
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setGroupsTitle(e.target.value);
  }, []);

  //handle on change group in main field component
  const handleChooseGroup = useCallback((item: groupItem) => {
    setCurrentGroup(item);
  }, []);

  const handleMoveUp = useCallback((idx: number) => {
    if (idx === 0) return;
    updateModifiedGroups(l => {
      let newItems: groupItem[] | undefined = undefined;
      if (!l || !l.groups) return;
      newItems = l.groups;
      array_move(newItems, idx, idx - 1);
      return { ...l, groups: newItems };
    });
  }, []);

  const handleMoveDown = useCallback(
    (idx: number) => {
      if (modifiedGroups?.groups && idx >= modifiedGroups.groups.length - 1) return;
      updateModifiedGroups(l => {
        let newItems: groupItem[] | undefined = undefined;
        if (!l || !l.groups) return;
        newItems = l.groups;
        array_move(newItems, idx, idx + 1);
        return { ...l, groups: newItems };
      });
    },
    [modifiedGroups?.groups],
  );

  const handleRemove = useCallback((idx: number) => {
    updateModifiedGroups(l => {
      let newItems: groupItem[] | undefined = undefined;
      if (!l || !l.groups) return;
      newItems = l.groups.filter((_, idx2) => idx2 != idx);
      return { ...l, groups: newItems };
    });
  }, []);

  useEffect(() => {
    updateGroups(l => {
      return {
        ...l,
        groupsTitle,
      };
    });
    updateGroups(l => {
      l.groups?.forEach(item1 => {
        const itemFromArr2 = modifiedGroups?.groups.find(item2 => item2.group == item1.group);
        if (itemFromArr2) {
          item1.title = itemFromArr2.title;
        }
      });
      return { ...l };
    });
  }, [groupsTitle, modifiedGroups]);

  return {
    switchGroupObj,
    groupsTitle,
    currentGroup,
    modifiedGroups,
    handleModifyGroupTitle,
    handleModifyGroup,
    handAddItem,
    handleRemove,
    handleMoveDown,
    handleMoveUp,
    handleTitleChange,
    handleChooseGroup,
  };
};
