import { Select } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { BaseFieldProps } from "../../types";

import useHooks from "./hooks";

const BuildingShadow: React.FC<BaseFieldProps<"buildingShadow">> = ({
  value,
  dataID,
  editMode,
  onUpdate,
}) => {
  const { options, handleUpdateSelect } = useHooks({
    value,
    dataID,
    onUpdate,
  });

  return editMode ? (
    <StyledSelect
      defaultValue="disabled"
      style={{ width: "100%" }}
      value={options.shadow}
      onChange={handleUpdateSelect("shadow")}
      options={[
        {
          value: "disabled",
          label: "なし",
        },
        {
          value: "enabled",
          label: "投影と受光",
        },
        {
          value: "cast_only",
          label: "投影のみ",
        },
        {
          value: "receive_only",
          label: "受光のみ",
        },
      ]}
    />
  ) : null;
};

export default BuildingShadow;

const StyledSelect = styled(Select)`
  width: 100%;
`;
