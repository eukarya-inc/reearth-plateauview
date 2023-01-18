import { Carousel, Icon, Pagination } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useRef, useState } from "react";
import { Remarkable } from "remarkable";

import type { Camera, Scene as SceneType } from "../../types";
import "./index.css";

type Props = {
  scenes: SceneType[];
  viewScene: (camera: Camera) => void;
  setPlayerHeight: (height: number) => void;
};

const minCarouselHeight = 131;
const maxCarouselHeight = 331;

const Player: React.FC<Props> = ({ scenes, viewScene, setPlayerHeight }) => {
  const sceneRefs = useRef<HTMLDivElement[]>([]);
  const setSceneRef = useCallback((dom: HTMLDivElement) => {
    sceneRefs.current.push(dom);
  }, []);

  const sceneTitleRefs = useRef<HTMLDivElement[]>([]);
  const setSceneTitleRef = useCallback((dom: HTMLDivElement) => {
    sceneTitleRefs.current.push(dom);
  }, []);

  const sceneContentRefs = useRef<HTMLDivElement[]>([]);
  const setSceneContentRef = useCallback((dom: HTMLDivElement) => {
    sceneContentRefs.current.push(dom);
  }, []);

  const updateHeight = useCallback(
    (index: number) => {
      if (scenes.length > 0 && sceneRefs.current[index]) {
        let carouselHeight =
          sceneTitleRefs.current[index].clientHeight +
          sceneContentRefs.current[index].clientHeight +
          8 + // gap
          24 + // description padding bottom
          12; // space

        carouselHeight =
          carouselHeight > maxCarouselHeight
            ? maxCarouselHeight
            : carouselHeight < minCarouselHeight
            ? minCarouselHeight
            : carouselHeight;
        sceneRefs.current[index].style.height = `${carouselHeight}px`;
        setPlayerHeight(carouselHeight + 24 + 40);
      } else {
        setPlayerHeight(minCarouselHeight + 24 + 40);
      }
    },
    [scenes, setPlayerHeight],
  );

  const carouselRef = useRef<any>(null);
  const [current, setCurrent] = useState<number>(0);
  const currentRef = useRef<number>(current);
  currentRef.current = current;

  const onSlideChange = useCallback(
    (oldSlide: number, currentSlide: number) => {
      if (currentSlide !== current) {
        const camera = scenes[currentSlide]?.camera;
        if (camera) {
          viewScene(camera);
        }
        setCurrent(currentSlide);
        updateHeight(currentSlide);
      }
    },
    [scenes, viewScene, current, setCurrent, updateHeight],
  );

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
      linkTarget: "_blank",
    }),
  );

  useEffect(() => {
    if (scenes.length === 0) {
      sceneRefs.current = [];
      sceneTitleRefs.current = [];
      sceneContentRefs.current = [];
      carouselRef.current = undefined;
      updateHeight(0);
    } else {
      if (currentRef.current !== 0) {
        carouselRef.current.goTo(0);
      } else {
        if (scenes[0]?.camera) {
          viewScene(scenes[0].camera);
        }
        updateHeight(0);
      }
    }
  }, [scenes, viewScene, updateHeight]);

  useEffect(() => {
    return () => {
      sceneRefs.current = [];
      sceneTitleRefs.current = [];
      sceneContentRefs.current = [];
      carouselRef.current = undefined;
    };
  }, []);

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
                draggable={true}
                speed={250}>
                {scenes.map((scene, index) => (
                  <div key={index}>
                    <StoryItem ref={setSceneRef}>
                      <Title ref={setSceneTitleRef}>{scene.title}</Title>
                      <Description>
                        <div
                          ref={setSceneContentRef}
                          dangerouslySetInnerHTML={{
                            __html: md.current.render(scene.description),
                          }}
                        />
                      </Description>
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
  pointer-events: ${({ disabled }) => (disabled ? "none" : "all")};
`;

const MainContent = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
`;

const CarouselContainer = styled.div`
  height: 100%;
  border: 1px solid rgba(0, 0, 0, 0.45);
  border-radius: 6px;
  overflow: hidden;
  width: 100%;
  box-sizing: content-box;
`;

const CarouselArea = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
`;

const PaginationContainer = styled.div`
  position: absolute;
  display: flex;
  flex-direction: row-reverse;
  right: 6px;
  bottom: 1px;
  background-color: #fff;
`;

const StoryItem = styled.div`
  width: 100%;
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
  padding: 0 12px 24px;
`;

export default Player;
