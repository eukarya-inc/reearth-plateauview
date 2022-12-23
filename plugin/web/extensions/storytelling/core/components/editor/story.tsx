import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback } from "react";

import type { Camera, Story as StoryType } from "../../types";

type Props = StoryType & {
  key: number;
  viewStory: (camera: Camera) => void;
  recapture: (id: string) => void;
  deleteStory: (id: string) => void;
  editStory: (id: string) => void;
};

const Story: React.FC<Props> = ({
  id,
  title,
  description,
  camera,
  viewStory,
  recapture,
  deleteStory,
  editStory,
}) => {
  const hendleView = useCallback(() => {
    if (camera) {
      viewStory(camera);
    }
  }, [viewStory, camera]);

  const handleEdit = useCallback(() => {
    editStory(id);
  }, [editStory, id]);

  const handleRecapture = useCallback(() => {
    recapture(id);
  }, [recapture, id]);

  const handleDelete = useCallback(() => {
    deleteStory(id);
  }, [deleteStory, id]);

  const items = [
    { label: "View", key: "view", onClick: hendleView },
    { label: "Edit", key: "edit", onClick: handleEdit },
    { label: "Recapture", key: "recapture", onClick: handleRecapture },
    { label: "Delete", key: "delete", onClick: handleDelete },
  ];
  const menu = <Menu items={items} />;

  return (
    <StyledStory>
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

export default Story;
