import { Group } from "@web/extensions/sidebar/core/newTypes";
import { array_move } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

import { GroupItem, SwitchGroup } from "../../types";

export default ({
  value,
  fieldGroups,
  onUpdate,
}: {
  value: SwitchGroup;
  fieldGroups?: Group[];
  onUpdate: (property: SwitchGroup) => void;
}) => {
  const [groupItems, updateGroupItems] = useState<GroupItem[]>(value.groups);
  const [title, setTitle] = useState(value.title);
  const [selectedGroup, selectGroup] = useState(fieldGroups?.[0] ?? undefined);

  const handleTitleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setTitle(e.target.value);
      onUpdate({ ...value, title: e.target.value });
    },
    [value, onUpdate],
  );

  const handleGroupChoose = useCallback(
    (groupID: number) => {
      selectGroup(fieldGroups?.find(fg => fg.id === groupID));
    },
    [fieldGroups],
  );

  const handleItemAdd = useCallback(() => {
    if (!fieldGroups) return;
    const newItem = { title: "新グループ", groupID: fieldGroups[0].id };
    updateGroupItems(gi => [...gi, newItem]);
    onUpdate({ ...value, groups: [...value.groups, newItem] });
  }, [value, fieldGroups, onUpdate]);

  const handleItemRemove = useCallback(
    (idx: number) => {
      const newGroups = value.groups.filter((_, gidx) => gidx !== idx);
      onUpdate({ ...value, groups: newGroups });
    },
    [value, onUpdate],
  );

  const handleItemGroupChange = useCallback(
    (idx: number, groupID?: number) => {
      if (!groupID) return;
      const updatedGroups = value.groups;
      updatedGroups[idx].groupID = groupID;
      updateGroupItems(updatedGroups);
      onUpdate({ ...value, groups: updatedGroups });
    },
    [value, onUpdate],
  );

  const handleItemTitleChange = useCallback(
    (title: string, idx: number) => {
      const updatedGroups = value.groups;
      updatedGroups[idx].title = title;
      onUpdate({ ...value, groups: updatedGroups });
    },
    [value, onUpdate],
  );

  const handleItemMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      const newGroups = value.groups;
      array_move(newGroups, idx, idx - 1);
      onUpdate({ ...value, groups: newGroups });
    },
    [value, onUpdate],
  );

  const handleItemMoveDown = useCallback(
    (idx: number) => {
      if (idx >= value.groups.length - 1) return;
      const newGroups = value.groups;
      array_move(newGroups, idx, idx + 1);
      onUpdate({ ...value, groups: newGroups });
    },
    [value, onUpdate],
  );

  return {
    title,
    groupItems,
    selectedGroup,
    handleTitleChange,
    handleGroupChoose,
    handleItemGroupChange,
    handleItemTitleChange,
    handleItemAdd,
    handleItemRemove,
    handleItemMoveUp,
    handleItemMoveDown,
  };
};
