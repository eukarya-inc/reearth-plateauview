import { Switch } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const Search: React.FC<BaseFieldProps<"search">> = ({ value, editMode, onUpdate }) => {
  const [enabled, setEnabled] = useState(!!value.enabled);

  const handleEnabledChange = useCallback(() => {
    setEnabled(!enabled);
    onUpdate({
      type: "search",
      enabled: !enabled,
    });
  }, [enabled, onUpdate]);

  return editMode ? (
    <Wrapper>
      <Title>Enable</Title>
      <Switch checked={enabled} size="small" onChange={handleEnabledChange} />
    </Wrapper>
  ) : null;
};

const Wrapper = styled.div`
  width: 100%;
  height: 32px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 4px;
`;

const Title = styled.div`
  width: 87px;
  font-size: 14px;
  display: flex;
  align-items: center;
`;

export default Search;
