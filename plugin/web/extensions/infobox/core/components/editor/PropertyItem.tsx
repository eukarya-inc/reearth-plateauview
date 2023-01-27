import type { PublicProperty } from "@web/extensions/infobox/types";
import { Icon, Checkbox } from "@web/sharedComponents";
import { styled } from "@web/theme";
import type { Identifier, XYCoord } from "dnd-core";
import { useRef } from "react";
import { useDrag, useDrop } from "react-dnd";

export type PropertyItem = PublicProperty & {
  value?: any;
};

type Props = {
  id: string;
  index: number;
  property: PropertyItem;
  onCheckChange: (e: any) => void;
  moveProperty: (dragIndex: number, hoverIndex: number) => void;
};

interface DragItem {
  index: number;
  id: string;
  type: string;
}

const PropertyItem: React.FC<Props> = ({ id, index, property, onCheckChange, moveProperty }) => {
  const dragRef = useRef<HTMLDivElement>(null);
  const previewRef = useRef<HTMLDivElement>(null);

  const [{ isDragging }, drag, preview] = useDrag({
    type: "propertyItem",
    item: () => {
      return { id, index };
    },
    collect: (monitor: any) => ({
      isDragging: monitor.isDragging(),
    }),
  });

  const [{ handlerId }, drop] = useDrop<DragItem, void, { handlerId: Identifier | null }>({
    accept: "propertyItem",
    collect(monitor) {
      return {
        handlerId: monitor.getHandlerId(),
      };
    },
    hover(item: DragItem, monitor) {
      if (!previewRef.current) {
        return;
      }
      const dragIndex = item.index;
      const hoverIndex = index;

      if (dragIndex === hoverIndex) {
        return;
      }

      const hoverBoundingRect = previewRef.current?.getBoundingClientRect();
      const hoverMiddleY = (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;
      const clientOffset = monitor.getClientOffset();
      const hoverClientY = (clientOffset as XYCoord).y - hoverBoundingRect.top;

      if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
        return;
      }

      if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
        return;
      }

      moveProperty(dragIndex, hoverIndex);

      item.index = hoverIndex;
    },
  });

  const opacity = isDragging ? 0.2 : 1;

  drag(dragRef);
  drop(preview(previewRef));

  return (
    <StyledPropertyItem
      disabled={property.hidden}
      ref={previewRef}
      data-handler-id={handlerId}
      style={{ opacity }}>
      <IconsWrapper>
        <DragHandle ref={dragRef}>
          <Icon icon="dotsSixVertical" size={16} />
        </DragHandle>
        <StyledCheckbox
          onChange={onCheckChange}
          data-key={property.key}
          checked={!property.hidden}
        />
      </IconsWrapper>
      <ContentWrapper>
        <JsonPath>{property.key}</JsonPath>
        <Title>
          <TitleInput disabled={!!property.hidden} />
        </Title>
        <Value>{property.value}</Value>
      </ContentWrapper>
    </StyledPropertyItem>
  );
};

const StyledPropertyItem = styled.div<{ disabled?: boolean }>`
  display: flex;
  align-items: flex-start;
  min-height: 32px;
  padding: 4px 0;
  gap: 12px;
  border-bottom: 1px solid #d9d9d9;
  font-size: 14px;
  color: ${({ disabled }) => (disabled ? "rgba(0, 0, 0, 0.25)" : "#000")};
`;

const IconsWrapper = styled.div`
  width: 56px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-around;
`;

const ContentWrapper = styled.div`
  width: 100%;
  display: flex;
  gap: 12px;
`;

const DragHandle = styled.div`
  cursor: pointer;
`;

const StyledCheckbox = styled(Checkbox)`
  .ant-checkbox-inner {
    border-color: transparent !important;
  }
`;

const TitleInput = styled.input`
  width: 100%;
  height: 24px;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  padding: 0 8px;
  font-size: 14px;
`;

const JsonPath = styled.div`
  width: 33%;
`;

const Title = styled.div`
  width: 33%;
`;

const Value = styled.div`
  width: 33%;
  word-break: break-all;
`;

export default PropertyItem;
