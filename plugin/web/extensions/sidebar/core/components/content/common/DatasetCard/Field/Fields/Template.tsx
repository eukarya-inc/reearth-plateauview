import { Select } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useMemo } from "react";

import { BaseFieldProps } from "./types";

const Template: React.FC<BaseFieldProps<"template">> = ({
  value,
  editMode,
  templates,
  onUpdate,
}) => {
  console.log("TEMPLATES!", templates);

  const handleTemplateChange = useCallback(
    (id: string) => {
      console.log("CHANGINGG TEMPLATE", id);
      onUpdate({
        ...value,
        templateID: id,
      });
    },
    [value, onUpdate],
  );

  const templateOptions = useMemo(
    () =>
      templates?.map(t => {
        return {
          value: t.id,
          label: t.name,
        };
      }),
    [templates],
  );

  return editMode ? (
    <Wrapper>
      <Title>テンプレート</Title>
      <Select
        options={templateOptions}
        style={{ width: "100%" }}
        value={value.templateID ?? templates?.[0].id}
        onChange={handleTemplateChange}
        getPopupContainer={trigger => trigger.parentElement ?? document.body}
      />
    </Wrapper>
  ) : (
    <div>
      <p>CurrentTemplate: {templates?.find(t => t.id === value.templateID)?.name}</p>
    </div>
  );
};

export default Template;

const Wrapper = styled.div`
  display: flex;
`;

const Title = styled.p`
  margin: 0;
`;
