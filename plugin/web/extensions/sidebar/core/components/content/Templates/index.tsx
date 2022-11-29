import CommonPage from "@web/extensions/sidebar/core/components/content/CommonPage";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

type Template = {
  name: string;
};

const Templates: React.FC = () => {
  const [addedTemplates, changeTemplates] = useState<Template[]>([
    { name: "建物モデル" },
    { name: "ランドマーク" },
  ]);

  //   const handleTemplateAdd = useCallback(
  //     (template: Template) => {
  //       if (addedTemplates.includes(template)) return;
  //       changeTemplates([...addedTemplates, template]);
  //     },
  //     [addedTemplates],
  //   );

  const handleTemplateRemove = useCallback(
    (template: Template) => {
      if (!addedTemplates.includes(template)) return;
      changeTemplates(addedTemplates.filter(t => t !== template));
    },
    [addedTemplates],
  );

  return (
    <CommonPage>
      <Content>
        <Title>Template Editor</Title>
        <TemplateAddButton>
          <Icon icon="plus" size={16} /> New Template
        </TemplateAddButton>
        {addedTemplates.map((t, idx) => (
          <TemplateComponent key={idx}>
            {t.name} <StyledIcon icon="trash" size={16} onClick={() => handleTemplateRemove(t)} />
          </TemplateComponent>
        ))}
      </Content>
    </CommonPage>
  );
};

export default Templates;

const Content = styled.div`
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 8px;
`;

const Title = styled.p`
  margin: 0;
`;

const TemplateWrapper = styled.div`
  display: flex;
  align-items: center;
  height: 34px;
  width: 100%;
  background: #f5f5f5;
  cursor: pointer;
  transition: background 0.3s;

  :hover {
    background: #ffffff;
  }
`;

const TemplateAddButton = styled(TemplateWrapper)`
  justify-content: center;
  gap: 8px;
`;

const TemplateComponent = styled(TemplateWrapper)`
  justify-content: space-between;
  padding-left: 12px;
  padding-right: 10px;
`;

const StyledIcon = styled(Icon)`
  border-radius: 4px;
  padding: 2px;
  border-width: 0.5px;
  border-style: solid;
  border-color: transparent;

  :hover {
    background: #f5f5f5;
    border-color: black;
  }
`;
