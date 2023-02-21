import type { Field, Fields, Feature } from "@web/extensions/infobox/types";
import { Collapse, Button } from "@web/sharedComponents";
import { styled } from "@web/theme";
import update from "immutability-helper";
import { useCallback, useEffect, useState } from "react";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";

import type { FieldItem as FieldItemType } from "./FieldItem";
import PropertyItem from "./FieldItem";

type Props = {
  fields: Fields;
  feature?: Feature;
  saveFields: (fields: Fields) => void;
};

const Editor: React.FC<Props> = ({ feature, fields, saveFields, ...props }) => {
  const [fieldList, setFieldList] = useState<FieldItemType[]>([]);

  useEffect(() => {
    const fieldItems: FieldItemType[] = [];

    if (!fields?.fields || fields.fields?.length === 0) {
      feature?.properties.forEach(fp => {
        fieldItems.push({
          title: "",
          path: fp.key,
          value: fp.value,
          visible: true,
        });
      });
    } else {
      const processedFields: string[] = [];
      fields.fields.forEach(f => {
        fieldItems.push({
          ...f,
          value: feature?.properties.find(fp => fp.key === f.path)?.value ?? "",
        });
        processedFields.push(f.path);
      });
      feature?.properties
        .filter(fp => !processedFields.includes(fp.key))
        .forEach(fp => {
          fieldItems.push({
            title: "",
            path: fp.key,
            value: fp.value,
            visible: true,
          });
        });
    }

    setFieldList(fieldItems);
  }, [feature, fields]);

  const onCheckChange = useCallback((e: any) => {
    setFieldList(list => {
      const fieldItem = list.find(item => item.path === e.target["data-path"]);
      if (fieldItem) {
        fieldItem.visible = !!e.target.checked;
      }
      return [...list];
    });
  }, []);

  const onTitleChange = useCallback((e: any) => {
    setFieldList(list => {
      const fieldItem = list.find(item => item.path === e.target.dataset.path);
      if (fieldItem) {
        fieldItem.title = e.target.value;
      }
      return [...list];
    });
  }, []);

  const moveProperty = useCallback((dragIndex: number, hoverIndex: number) => {
    setFieldList((prevCards: FieldItemType[]) =>
      update(prevCards, {
        $splice: [
          [dragIndex, 1],
          [hoverIndex, 0, prevCards[dragIndex] as FieldItemType],
        ],
      }),
    );
  }, []);

  const onSave = useCallback(() => {
    const outputFields: Field[] = [];
    fieldList.forEach(f => {
      outputFields.push({
        path: f.path,
        title: f.title,
        visible: f.visible,
      });
    });
    saveFields({
      ...fields,
      fields: outputFields,
    });
  }, [fieldList, fields, saveFields]);

  return (
    <StyledPanel
      header={fields.name}
      key={fields.name}
      extra={
        <StyledButton size="small" onClick={onSave}>
          保存
        </StyledButton>
      }
      {...props}>
      <Wrapper>
        <PropertyHeader>
          <IconsWrapper />
          <ContentWrapper>
            <JsonPath>JSON Path</JsonPath>
            <Title>Title</Title>
            <Value>Value</Value>
          </ContentWrapper>
        </PropertyHeader>
        <DndProvider backend={HTML5Backend}>
          {fieldList.map((field, index) => (
            <PropertyItem
              id={field.path}
              index={index}
              key={field.path}
              field={field}
              onCheckChange={onCheckChange}
              onTitleChange={onTitleChange}
              moveProperty={moveProperty}
            />
          ))}
        </DndProvider>
      </Wrapper>
    </StyledPanel>
  );
};

const StyledPanel = styled(Collapse.Panel)`
  background: #f4f4f4;
  margin-bottom: 6px;
  box-shadow: 1px 2px 4px rgba(0, 0, 0, 0.25);
  border-radius: 4px !important;
  overflow: hidden;
`;

const StyledButton = styled(Button)`
  border-radius: 4px;
  margin-right: 10px;
`;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const PropertyHeader = styled.div`
  display: flex;
  align-items: flex-start;
  min-height: 32px;
  padding: 4px 0;
  gap: 12px;
  border-bottom: 1px solid #d9d9d9;
  font-size: 14px;
`;

const IconsWrapper = styled.div`
  width: 56px;
  flex-shrink: 0;
`;

const ContentWrapper = styled.div`
  width: 100%;
  display: flex;
  gap: 12px;
`;

const JsonPath = styled.div`
  width: 33%;
`;

const Title = styled.div`
  width: 33%;
`;

const Value = styled.div`
  width: 33%;
`;

export default Editor;
