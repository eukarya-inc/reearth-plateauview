import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback } from "react";
import {
  Accordion,
  AccordionItem,
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemPanel,
  AccordionItemState,
} from "react-accessible-accordion";

import fields, { BaseField as BaseFieldProps, FieldType } from "./Fields";

export type Field = BaseFieldProps<FieldType> & {
  value: any;
  icon?: string; // MAYBE NOT NEEDED
};

export type Props = {
  field: Field;
  editMode?: boolean;
  onGroupAdd?: () => void;
  onRemove?: () => void;
};

const FieldComponent: React.FC<Props> = ({ field, editMode, onGroupAdd, onRemove }) => {
  const FieldContent = fields[field.type];

  const handleGroupAdd = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      onGroupAdd?.();
    },
    [onGroupAdd],
  );

  const handleRemove = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      onRemove?.();
    },
    [onRemove],
  );

  return !editMode && field.type === "idealZoom" ? null : (
    <StyledAccordionComponent allowZeroExpanded>
      <AccordionItem>
        <AccordionItemState>
          {({ expanded }) => (
            <Header expanded={expanded}>
              {editMode ? (
                <HeaderContents>
                  <LeftContents>
                    <ArrowIcon icon="arrowDown" size={16} direction="right" expanded={expanded} />
                    <Title>{field.configTitle}</Title>
                  </LeftContents>
                  <RightContents>
                    <StyledIcon icon="group" size={16} onClick={handleGroupAdd} />
                    <StyledIcon icon="trash" size={16} onClick={handleRemove} />
                  </RightContents>
                </HeaderContents>
              ) : (
                <HeaderContents>
                  <Title>{field.title}</Title>
                  <ArrowIcon icon="arrowDown" size={16} direction="left" expanded={expanded} />
                </HeaderContents>
              )}
            </Header>
          )}
        </AccordionItemState>
        <BodyWrapper>{FieldContent && <FieldContent {...field} editMode={editMode} />}</BodyWrapper>
      </AccordionItem>
    </StyledAccordionComponent>
  );
};

export default FieldComponent;

const StyledAccordionComponent = styled(Accordion)`
  width: 100%;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  background: #ffffff;
`;

const Header = styled(AccordionItemHeading)<{ expanded?: boolean }>`
  border-bottom-width: 1px;
  border-bottom-style: solid;
  border-bottom-color: transparent;
  ${({ expanded }) => expanded && "border-bottom-color: #e0e0e0;"}
  display: flex;
  height: 30px;
`;

const HeaderContents = styled(AccordionItemButton)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
  padding: 0 12px;
  outline: none;
  cursor: pointer;
`;

const BodyWrapper = styled(AccordionItemPanel)`
  border-radius: 0px 0px 4px 4px;
  padding: 12px;
`;

const Title = styled.p`
  margin: 0;
`;

const StyledIcon = styled(Icon)`
  cursor: pointer;
`;

const LeftContents = styled.div`
  display: flex;
  align-items: center;
`;

const RightContents = styled.div`
  display: flex;
  gap: 4px;
`;

const ArrowIcon = styled(Icon)<{ direction: "left" | "right"; expanded?: boolean }>`
  transition: transform 0.15s ease;
  ${({ direction, expanded }) =>
    (direction === "right" && !expanded && "transform: rotate(-90deg);") ||
    (direction === "left" && !expanded && "transform: rotate(90deg);") ||
    null}
  ${({ direction }) => (direction === "left" ? "margin: 0 -4px 0 4px;" : "margin: 0 4px 0 -4px;")}
`;
