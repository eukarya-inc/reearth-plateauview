import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";
import type { Identifier, XYCoord } from "dnd-core";
import { useCallback, useRef } from "react";
import { useDrag, useDrop } from "react-dnd";

import type { Camera, Scene as SceneType } from "../../types";

type Props = SceneType & {
  index: number;
  viewScene: (camera: Camera) => void;
  recaptureScene: (id: string) => void;
  deleteScene: (id: string) => void;
  editScene: (id: string) => void;
  moveScene: (dragIndex: number, hoverIndex: number) => void;
};

const Scene: React.FC<Props> = ({
  id,
  title,
  description,
  camera,
  index,
  viewScene,
  recaptureScene,
  deleteScene,
  editScene,
  moveScene,
}) => {
  const hendleView = useCallback(() => {
    if (camera) {
      viewScene(camera);
    }
  }, [viewScene, camera]);

  const handleEdit = useCallback(() => {
    editScene(id);
  }, [editScene, id]);

  const handleRecapture = useCallback(() => {
    recaptureScene(id);
  }, [recaptureScene, id]);

  const handleDelete = useCallback(() => {
    deleteScene(id);
  }, [deleteScene, id]);

  const items = [
    { label: "View", key: "view", onClick: hendleView },
    { label: "Edit", key: "edit", onClick: handleEdit },
    { label: "Recapture", key: "recapture", onClick: handleRecapture },
    { label: "Delete", key: "delete", onClick: handleDelete },
  ];
  const menu = <Menu items={items} />;

  interface DragItem {
    index: number;
    id: string;
    type: string;
  }

  const ref = useRef<HTMLDivElement>(null);

  const [{ handlerId }, drop] = useDrop<DragItem, void, { handlerId: Identifier | null }>({
    accept: "scene",
    collect(monitor) {
      return {
        handlerId: monitor.getHandlerId(),
      };
    },
    hover(item: DragItem, monitor) {
      if (!ref.current) {
        return;
      }
      const dragIndex = item.index;
      const hoverIndex = index;

      if (dragIndex === hoverIndex) {
        return;
      }

      const hoverBoundingRect = ref.current?.getBoundingClientRect();
      const hoverMiddleX = (hoverBoundingRect.right - hoverBoundingRect.left) / 2;
      const clientOffset = monitor.getClientOffset();
      const hoverClientX = (clientOffset as XYCoord).x - hoverBoundingRect.left;

      if (dragIndex < hoverIndex && hoverClientX < hoverMiddleX) {
        return;
      }

      if (dragIndex > hoverIndex && hoverClientX > hoverMiddleX) {
        return;
      }

      moveScene(dragIndex, hoverIndex);

      item.index = hoverIndex;
    },
  });

  const [{ isDragging }, drag] = useDrag({
    type: "scene",
    item: () => {
      return { id, index };
    },
    collect: (monitor: any) => ({
      isDragging: monitor.isDragging(),
    }),
  });

  const opacity = isDragging ? 0 : 1;
  drag(drop(ref));

  return (
    <StyledStory ref={ref} style={{ opacity }} data-handler-id={handlerId}>
      <Header>
        <Title>{title}</Title>
        <ActionsBtn>
          <Dropdown trigger={["click"]} overlay={menu} placement={"topRight"}>
            <Icon icon="dotsThreeVertical" size={24} />
          </Dropdown>
        </ActionsBtn>
      </Header>
      <Description>{description}</Description>
    </StyledStory>
  );
};

const StyledStory = styled.div`
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  width: 170px;
  height: 114px;
  flex-shrink: 0;
  background: #f8f8f8;
  border-radius: 8px;
  border: 1px solid #c7c5c5;
  padding: 6px 0;
`;

const Header = styled.div`
  display: flex;
  width: 100%;
  height: 36px;
  padding: 6px 12px;
  align-items: center;
  justify-content: space-between;
`;

const Title = styled.div`
  width: 100%;
  font-weight: 700;
  font-size: 14px;
  line-height: 19px;
  color: #000;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

const ActionsBtn = styled.a`
  display: flex;
  flex-shrink: 0;
`;

const Description = styled.div`
  height: 100%;
  padding: 6px 12px;
  color: #3a3a3a;
  font-size: 12px;
  line-height: 18px;
  font-weight: 500;
  overflow: hidden;
  display: -webkit-box !important;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  white-space: normal;
`;

export default Scene;
