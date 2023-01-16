import { Carousel, Icon, Pagination } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useRef, useState } from "react";
import { Remarkable } from "remarkable";

import type { Camera, Scene as SceneType } from "../../types";
import "./index.css";

type Props = {
  scenes: SceneType[];
  viewScene: (camera: Camera) => void;
};

const Player: React.FC<Props> = ({ scenes, viewScene }) => {
  const [current, setCurrent] = useState(0);
  const carouselRef = useRef<any>(null);
  const onSlideChange = (oldSlide: number, currentSlide: number) => {
    const camera = scenes[currentSlide]?.camera;
    if (camera) {
      viewScene(camera);
    }
    setCurrent(currentSlide);
  };

  const prev = useCallback(() => {
    if (carouselRef.current) {
      carouselRef.current.prev();
    }
  }, []);

  const next = useCallback(() => {
    if (carouselRef.current) {
      carouselRef.current.next();
    }
  }, []);

  const onPaginationChange = useCallback((current: number) => {
    if (carouselRef.current) {
      carouselRef.current.goTo(current - 1);
    }
  }, []);

  const md = useRef(
    new Remarkable({
      html: false,
      breaks: true,
      typographer: true,
      linkTarget: "__blank",
    }),
  );

  // auto view scene 1 if exist when active
  useEffect(() => {
    if (scenes[0]?.camera) {
      viewScene(scenes[0].camera);
    }
  }, [scenes, viewScene]);

  return (
    <Wrapper>
      <NavButton onClick={prev} disabled={current === 0}>
        <Icon icon="caretLeft" size={32} />
      </NavButton>
      <MainContent>
        <CarouselContainer>
          <CarouselArea>
            {scenes.length > 0 && (
              <Carousel
                beforeChange={onSlideChange}
                dots={false}
                ref={carouselRef}
                infinite={false}
                draggable={true}>
                {scenes.map((scene, index) => (
                  <div key={index}>
                    <StoryItem>
                      <Title>{scene.title}</Title>
                      <Description
                        dangerouslySetInnerHTML={{
                          __html: md.current.render(scene.description),
                        }}
                      />
                    </StoryItem>
                  </div>
                ))}
              </Carousel>
            )}
          </CarouselArea>
        </CarouselContainer>
        <PaginationContainer>
          {scenes.length > 0 && (
            <Pagination
              current={current + 1}
              size="small"
              total={scenes.length}
              pageSize={1}
              onChange={onPaginationChange}
            />
          )}
        </PaginationContainer>
      </MainContent>
      <NavButton onClick={next} disabled={current >= scenes.length - 1} className="next">
        <Icon icon="caretLeft" size={32} />
      </NavButton>
    </Wrapper>
  );
};

const Wrapper = styled.div`
  position: relative;
  display: flex;
  justify-content: space-between;
  height: 100%;
  padding: 12px;
  gap: 12px;
`;

const NavButton = styled.a<{ disabled: boolean }>`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 40px;
  &.next {
    transform: rotate(180deg);
  }
  color: ${({ disabled }) => (disabled ? "#ccc" : "--var(theme-color)")};
`;

const MainContent = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
`;

const CarouselContainer = styled.div`
  height: 107px;
  width: 100%;
`;

const CarouselArea = styled.div`
  position: absolute;
  width: 100%;
`;

const PaginationContainer = styled.div`
  position: relative;
  display: flex;
  flex-direction: row-reverse;
  top: -1px;
`;

const StoryItem = styled.div`
  width: 100%;
  height: 107px;
  border-bottom: 1px solid #c7c5c5;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Title = styled.div`
  font-size: 14px;
  font-weight: 700;
  line-height: 19px;
  flex-shrink: 0;
  padding: 12px 12px 0;
`;

const Description = styled.div`
  height: 100%;
  overflow: auto;
  font-size: 12px;
  line-height: 1.5;
  padding: 0 12px 12px;
`;

export default Player;
