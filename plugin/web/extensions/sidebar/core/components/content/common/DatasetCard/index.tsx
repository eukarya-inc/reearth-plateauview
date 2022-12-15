import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useEffect, useState } from "react";
import {
  Accordion,
  AccordionItem,
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemPanel,
  AccordionItemState,
} from "react-accessible-accordion";

import { Dataset as DatasetType, Field as FieldType } from "../types";

import Field from "./Field";

export type Dataset = DatasetType;

export type Field = FieldType;

export type Props = {
  dataset: Dataset;
  onRemove?: (id: string) => void;
};

const baseFields: Field[] = [
  { id: "zoom", title: "Ideal Zoom", icon: "mapPin", value: 1 },
  { id: "about", title: "About Data", icon: "about", value: "www.plateau.org/data-url" },
  { id: "remove", icon: "trash" },
];

const DatasetCard: React.FC<Props> = ({ dataset, onRemove }) => {
  const [visible, setVisibility] = useState(false);

  useEffect(() => {
    setVisibility(dataset.type !== "group" ? !!dataset.visible : false);
  }, [dataset]);

  return (
    <StyledAccordionComponent allowZeroExpanded allowMultipleExpanded>
      <AccordionItem>
        <Header>
          <HeaderContents>
            <LeftMain>
              <Icon
                icon={!visible ? "hidden" : "visible"}
                size={20}
                onClick={e => {
                  e?.stopPropagation();
                  setVisibility(!visible);
                }}
              />
              <Title>{dataset.name}</Title>
            </LeftMain>
            <AccordionItemState>
              {({ expanded }) => <ArrowIcon icon="arrowDown" size={16} expanded={expanded} />}
            </AccordionItemState>
          </HeaderContents>
        </Header>
        <BodyWrapper>
          <Content>
            {baseFields.map((field, idx) => (
              <BaseField key={idx} onClick={() => field.id === "remove" && onRemove?.(dataset.id)}>
                {field.icon && <Icon icon={field.icon} size={20} color="#00BEBE" />}
                {field.title && <FieldName>{field.title}</FieldName>}
              </BaseField>
            ))}
            {[
              { id: "camera", icon: undefined, title: "Camera" },
              { id: "legend", icon: undefined, title: "Legend" },
            ]?.map((field, idx) => (
              <Field key={idx} field={field} />
            ))}
          </Content>
        </BodyWrapper>
      </AccordionItem>
    </StyledAccordionComponent>
  );
};

export default DatasetCard;

const StyledAccordionComponent = styled(Accordion)`
  width: 100%;
  border-radius: 4px;
  box-shadow: 1px 2px 4px rgba(0, 0, 0, 0.25);
  margin: 8px 0;
  background: #ffffff;
`;

const Header = styled(AccordionItemHeading)`
  border-bottom-width: 1px;
  border-bottom-style: solid;
  border-bottom-color: transparent;
  border-bottom-color: #e0e0e0;
  display: flex;
  height: 46px;
`;

const HeaderContents = styled(AccordionItemButton)`
  display: flex;
  align-items: center;
  flex: 1;
  padding: 0 12px;
  outline: none;
  cursor: pointer;
`;

const BodyWrapper = styled(AccordionItemPanel)<{ noTransition?: boolean }>`
  width: 100%;
  border-radius: 0px 0px 4px 4px;
  background: #fafafa;
  padding: 12px;
`;

const LeftMain = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
`;

const Title = styled.p`
  margin: 0;
  font-size: 16px;
`;

const Content = styled.div`
  display: flex;
  align-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  // padding-right: 12px;
`;

const BaseField = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  flex: 1 0 auto;
  // min-width: 53px;
  padding: 8px;
  background: #ffffff;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  cursor: pointer;
`;

const ArrowIcon = styled(Icon)<{ expanded?: boolean }>`
  transition: transform 0.15s ease;
  transform: ${({ expanded }) => !expanded && "rotate(90deg)"};
`;

const FieldName = styled.p`
  margin: 0;
`;
