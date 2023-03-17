import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { generateID, moveItemDown, moveItemUp, removeItem } from "@web/extensions/sidebar/utils";
import { Icon, Dropdown, Menu, Radio } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState, useCallback, useEffect, useRef } from "react";

import { stringifyCondition } from "../../../utils";
import { BaseFieldProps, Cond, SwitchVisibility } from "../../types";

import ConditionItem from "./ConditionItem";

type UIStyles = "dropdown" | "radio";

const uiStyles: { [key: string]: string } = {
  dropdown: "ドロップダウン",
  radio: "ラジオ",
};

const SwitchVisibility: React.FC<BaseFieldProps<"switchVisibility">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  const [selectedStyle, selectStyle] = useState(value.uiStyle ?? "dropdown");

  const handleStyleChange = useCallback((style: UIStyles) => {
    selectStyle(style);
  }, []);

  const styleOptions = (
    <Menu
      items={Object.keys(uiStyles).map(key => {
        return {
          key: key,
          label: (
            <p style={{ margin: 0 }} onClick={() => handleStyleChange(key as UIStyles)}>
              {uiStyles[key]}
            </p>
          ),
        };
      })}
    />
  );

  const [conditions, setConditions] = useState<SwitchVisibility["conditions"]>(
    value.conditions ?? [],
  );

  const handleMoveUp = useCallback((idx: number) => {
    setConditions(c => moveItemUp(idx, c) ?? c);
  }, []);

  const handleMoveDown = useCallback((idx: number) => {
    setConditions(c => moveItemDown(idx, c) ?? c);
  }, []);

  const handleAdd = useCallback(() => {
    setConditions(c => {
      const nc: { condition: Cond<any>; title: string } = {
        condition: {
          key: generateID(),
          operator: "===",
          operand: true,
          value: true,
        },
        title: "",
      };
      return c ? [...c, nc] : [nc];
    });
  }, []);

  const handleRemove = useCallback((idx: number) => {
    setConditions(c => removeItem(idx, c) ?? c);
  }, []);

  const handleItemUpdate = (item: { condition: Cond<any>; title: string }, index: number) => {
    setConditions(c => {
      const nc = [...(c ?? [])];
      nc.splice(index, 1, item);
      return nc;
    });
  };

  //
  const [selectedVisibility, setSelectedVisibility] = useState(value.userSettings?.selected ?? 0);

  const visibilityOptions = (
    <Menu
      items={conditions?.map((c, index) => {
        return {
          key: index,
          label: (
            <p style={{ margin: 0 }} onClick={() => setSelectedVisibility(index)}>
              {c.title}
            </p>
          ),
        };
      })}
    />
  );

  const valueRef = useRef(value);
  const onUpdateRef = useRef<any>();
  valueRef.current = value;
  onUpdateRef.current = onUpdate;

  useEffect(() => {
    const timer = setTimeout(() => {
      const showConditions: [string, string][] = [["true", "false"]];

      if (conditions[selectedVisibility]) {
        showConditions.unshift([
          stringifyCondition(conditions[selectedVisibility].condition),
          "true",
        ]);
      }

      onUpdateRef.current({
        ...valueRef.current,
        uiStyle: selectedStyle,
        conditions,
        userSettings: {
          selected: selectedVisibility,
        },
        override: {
          marker: {
            show: {
              expression: {
                conditions: showConditions,
              },
            },
          },
          polyline: {
            show: {
              expression: {
                conditions: showConditions,
              },
            },
          },
          polygon: {
            show: {
              expression: {
                conditions: showConditions,
              },
            },
          },
          resource: {
            show: {
              expression: {
                conditions: showConditions,
              },
            },
          },
        },
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [conditions, selectedStyle, selectedVisibility]);

  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>UIスタイル</FieldTitle>
        <FieldValue>
          <Dropdown
            overlay={styleOptions}
            placement="bottom"
            trigger={["click"]}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}>
            <StyledDropdownButton>
              <p style={{ margin: 0 }}>{uiStyles[selectedStyle]}</p>
              <StyledIcon icon="arrowDownSimple" size={12} />
            </StyledDropdownButton>
          </Dropdown>
        </FieldValue>
      </Field>
      {conditions?.map((c, idx) => (
        <ConditionItem
          key={idx}
          index={idx}
          item={c}
          handleMoveDown={handleMoveDown}
          handleMoveUp={handleMoveUp}
          handleRemove={handleRemove}
          onItemUpdate={handleItemUpdate}
        />
      ))}
      <ButtonWrapper>
        <AddButton text="Add Condition" height={24} onClick={handleAdd} />
      </ButtonWrapper>
    </Wrapper>
  ) : (
    <Wrapper>
      <Field>
        {conditions.length > 0 ? (
          selectedStyle === "radio" ? (
            <Radio.Group
              onChange={e => setSelectedVisibility(e.target.value)}
              value={selectedVisibility}>
              {conditions?.map((c, index) => (
                <StyledRadio key={index} value={index}>
                  <Label>{c.title}</Label>
                </StyledRadio>
              ))}
            </Radio.Group>
          ) : (
            <FieldValue>
              <Dropdown
                overlay={visibilityOptions}
                placement="bottom"
                trigger={["click"]}
                getPopupContainer={trigger => trigger.parentElement ?? document.body}>
                <StyledDropdownButton>
                  <p style={{ margin: 0 }}>{conditions[selectedVisibility].title}</p>
                  <Icon icon="arrowDownSimple" size={12} />
                </StyledDropdownButton>
              </Dropdown>
            </FieldValue>
          )
        ) : (
          <>No Conditions</>
        )}
      </Field>
    </Wrapper>
  );
};

const StyledDropdownButton = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  height: 32px;
  padding: 0 16px;
  cursor: pointer;
`;

const StyledIcon = styled(Icon)`
  font-size: 0;
`;

const Field = styled.div<{ gap?: number }>`
  display: flex;
  align-items: center;
  ${({ gap }) => gap && `gap: ${gap}px;`}
`;

const Text = styled.p`
  margin: 0;
`;

const FieldTitle = styled(Text)`
  width: 82px;
`;

const FieldValue = styled.div`
  position: relative;
  display: flex;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  flex: 1;
  height: 100%;
  width: 100%;
`;

const StyledRadio = styled(Radio)`
  width: 100%;
  margin-top: 8px;
`;

const Label = styled.span`
  font-size: 14px;
`;

export default SwitchVisibility;
