import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import {
  generateID,
  moveItemDown,
  moveItemUp,
  removeItem,
  postMsg,
} from "@web/extensions/sidebar/utils";
import { Icon, Dropdown, Menu, Radio } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState, useCallback, useEffect, useRef } from "react";

import { stringifyCondition } from "../../../utils";
import { BaseFieldProps, SwitchVisibility } from "../../types";

import ConditionItem from "./ConditionItem";

type UIStyles = "dropdown" | "radio";
export type ConditionItemType = SwitchVisibility["conditions"][0];

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
      const id = generateID();
      const nc: ConditionItemType = {
        id,
        condition: {
          key: id,
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

  const handleItemUpdate = (item: ConditionItemType, index: number) => {
    setConditions(c => {
      const nc = [...(c ?? [])];
      nc.splice(index, 1, item);
      return nc;
    });
  };

  //
  const [selectedVisibility, setSelectedVisibility] = useState<string | undefined>(
    value.userSettings?.selected ?? value.conditions?.[0]?.id,
  );

  const handleSelectVisibility = (id: string) => {
    setSelectedVisibility(id);
    postMsg({
      action: "unselect",
    });
  };

  const visibilityOptions = (
    <Menu
      items={conditions?.map((c, index) => {
        return {
          key: index,
          label: (
            <p style={{ margin: 0 }} onClick={() => handleSelectVisibility(c.id)}>
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

      const selectedCondition = conditions.find(c => c.id === selectedVisibility)?.condition;
      if (selectedCondition) {
        showConditions.unshift([stringifyCondition(selectedCondition), "true"]);
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
        },
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [conditions, selectedStyle, selectedVisibility]);

  useEffect(() => {
    if (
      conditions.length > 0 &&
      (!selectedVisibility || !conditions.find(c => c.id === selectedVisibility))
    ) {
      setSelectedVisibility(conditions[0].id);
    }
  }, [conditions, selectedVisibility]);

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
              onChange={e => handleSelectVisibility(e.target.value)}
              value={selectedVisibility}>
              {conditions?.map(c => (
                <StyledRadio key={c.id} value={c.id}>
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
                  <p style={{ margin: 0 }}>
                    {conditions.find(c => c.id === selectedVisibility)?.title}
                  </p>
                  <StyledIcon icon="arrowDownSimple" size={12} />
                </StyledDropdownButton>
              </Dropdown>
            </FieldValue>
          )
        ) : (
          <>項目なし</>
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
