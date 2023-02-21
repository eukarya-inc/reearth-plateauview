import type { Feature, Fields } from "@web/extensions/infobox/types";
import { Collapse } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useEffect, useState } from "react";

type Props = {
  key: number;
  feature: Feature;
  fields: Fields;
};

type FieldItemType = {
  path: string;
  title: string;
  value?: any;
};

const Viewer: React.FC<Props> = ({ feature, fields, key, ...props }) => {
  const [fieldList, setFieldList] = useState<FieldItemType[]>([]);

  useEffect(() => {
    const fieldItems: FieldItemType[] = [];

    if (!fields.fields || fields.fields?.length === 0) {
      feature.properties.forEach(p => {
        fieldItems.push({
          path: p.key,
          title: p.key,
          value: p.value,
        });
      });
    } else {
      const processedFields: string[] = [];
      fields.fields.forEach(f => {
        if (f.visible) {
          // field may not exist on feature
          const property = feature.properties.find(fp => fp.key === f.path);
          if (property) {
            fieldItems.push({
              path: f.path,
              title: f.title ?? f.path,
              value: property.value,
            });
          }
        }
        processedFields.push(f.path);
      });
      feature.properties
        .filter(fp => !processedFields.includes(fp.key))
        .forEach(fp => {
          fieldItems.push({
            path: fp.key,
            title: fp.key,
            value: fp.value,
          });
        });
    }
    setFieldList(fieldItems);
  }, [feature, fields]);

  return (
    <StyledPanel header={fields.name} key={key} {...props}>
      <Wrapper>
        {fieldList.map(field => (
          <PropertyItem key={field.path}>
            <Title>{field.title}</Title>
            <Value>{field.value}</Value>
          </PropertyItem>
        ))}
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

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const PropertyItem = styled.div`
  display: flex;
  align-items: flex-start;
  min-height: 32px;
  padding: 4px 0;
  gap: 12px;
  border-bottom: 1px solid #d9d9d9;
  font-size: 14px;
`;

const Title = styled.div`
  width: 50%;
`;

const Value = styled.div`
  width: 50%;
  word-break: break-all;
`;

export default Viewer;
