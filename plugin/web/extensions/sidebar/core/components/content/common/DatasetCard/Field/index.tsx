import { Group, Template } from "@web/extensions/sidebar/core/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useMemo, useState } from "react";
import {
  Accordion,
  AccordionItem,
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemPanel,
  AccordionItemState,
} from "react-accessible-accordion";

import fields from "./Fields";
import {
  ConfigData,
  FieldComponent as FieldComponentType,
  fieldName,
  generalFieldName,
  pointFieldName,
  polygonFieldName,
  polylineFieldName,
  threeDFieldName,
} from "./Fields/types";

export type Props = {
  index: number;
  field: FieldComponentType;
  dataID?: string;
  isActive: boolean;
  editMode?: boolean;
  templates?: Template[];
  selectGroups?: Group[];
  configData?: ConfigData[];
  onUpdate?: (id: string) => (property: any) => void;
  onRemove?: (id: string) => void;
  onMoveUp?: (index: number) => void;
  onMoveDown?: (index: number) => void;
  onGroupsUpdate?: (groups: Group[], selectedGroup?: string) => void;
  onCurrentGroupUpdate?: (fieldGroupID: string) => void;
};

const getFieldGroup = (field: string) => {
  if (field in generalFieldName) {
    return "一般";
  } else if (field in pointFieldName) {
    return "ポイント";
  } else if (field in polygonFieldName) {
    return "ポリゴン";
  } else if (field in threeDFieldName) {
    return "3Dタイル";
  } else if (field in polylineFieldName) {
    return "ポリライン";
  }
};

const FieldComponent: React.FC<Props> = ({
  index,
  field,
  dataID,
  isActive,
  editMode,
  templates,
  selectGroups,
  configData,
  onUpdate,
  onRemove,
  onMoveUp,
  onMoveDown,
  onGroupsUpdate,
  onCurrentGroupUpdate,
}) => {
  const Field = fields[field.type];
  const [groupPopupOpen, setGroupPopup] = useState(false);

  const handleGroupSelectOpen = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      postMsg({
        action: "groupSelectOpen",
        payload: { groups: selectGroups, selected: field.group },
      });
      setGroupPopup(true);
    },
    [field.group, selectGroups],
  );

  const handleRemove = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      onRemove?.(field.id);
    },
    [field, onRemove],
  );

  const handleUpClick = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      onMoveUp?.(index);
    },
    [index, onMoveUp],
  );

  const handleDownClick = useCallback(
    (e: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      e?.stopPropagation();
      onMoveDown?.(index);
    },
    [index, onMoveDown],
  );

  useEffect(() => {
    const eventListenerCallback = (e: any) => {
      if (e.source !== parent) return;
      if (groupPopupOpen) {
        if (e.data.action === "saveGroups") {
          onGroupsUpdate?.(e.data.payload.groups, e.data.payload.selected);
          setGroupPopup(false);
        } else if (e.data.action === "popupClose") {
          setGroupPopup(false);
        }
      }
    };
    (globalThis as any).addEventListener("message", eventListenerCallback);
    return () => {
      (globalThis as any).removeEventListener("message", eventListenerCallback);
    };
  }, [groupPopupOpen, onGroupsUpdate]);

  const title = useMemo(() => `${fieldName[field.type]}(${getFieldGroup(field.type)})`, [field]);

  return !editMode && !isActive ? null : (
    <StyledAccordionComponent
      allowZeroExpanded
      preExpanded={[field.id]}
      hide={!editMode && !Field?.hasUI}>
      <AccordionItem uuid={field.id}>
        <AccordionItemState>
          {({ expanded }) => (
            <Header showBorder={expanded}>
              {editMode ? (
                <HeaderContents>
                  <LeftContents>
                    {Field && (
                      <ArrowIcon icon="arrowDown" size={16} direction="right" expanded={expanded} />
                    )}
                    <Title>{title}</Title>
                  </LeftContents>
                  <RightContents>
                    <StyledIcon icon="arrowUpThin" size={16} onClick={handleUpClick} />
                    <StyledIcon icon="arrowDownThin" size={16} onClick={handleDownClick} />
                    <StyledIcon
                      icon="group"
                      color={field.group ? "#00BEBE" : "inherit"}
                      size={16}
                      onClick={handleGroupSelectOpen}
                    />
                    <StyledIcon icon="trash" size={16} onClick={handleRemove} />
                  </RightContents>
                </HeaderContents>
              ) : (
                <HeaderContents>
                  <Title>{title}</Title>
                  <ArrowIcon icon="arrowDown" size={16} direction="left" expanded={expanded} />
                </HeaderContents>
              )}
            </Header>
          )}
        </AccordionItemState>
        {Field?.Component && (
          <BodyWrapper>
            <Field.Component
              value={{ ...field }}
              editMode={editMode}
              isActive={isActive}
              templates={templates}
              fieldGroups={selectGroups}
              configData={configData}
              dataID={dataID}
              onUpdate={onUpdate?.(field.id)}
              onCurrentGroupUpdate={onCurrentGroupUpdate}
            />
          </BodyWrapper>
        )}
      </AccordionItem>
    </StyledAccordionComponent>
  );
};

export default FieldComponent;

const StyledAccordionComponent = styled(Accordion)<{ hide: boolean }>`
  ${({ hide }) => hide && "display: none;"}
  width: 100%;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  background: #ffffff;
`;

const Header = styled(AccordionItemHeading)<{ showBorder?: boolean }>`
  border-bottom-width: 1px;
  border-bottom-style: solid;
  border-bottom-color: transparent;
  ${({ showBorder }) => showBorder && "border-bottom-color: #e0e0e0;"}
  display: flex;
  height: auto;
`;

const HeaderContents = styled(AccordionItemButton)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
  padding: 12px;
  outline: none;
  cursor: pointer;
`;

const BodyWrapper = styled(AccordionItemPanel)`
  position: relative;
  border-radius: 0px 0px 4px 4px;
  padding: 12px;
`;

const Title = styled.p`
  margin: 0;
  user-select: none;
  width: 160px;
  overflow-wrap: break-word;
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
