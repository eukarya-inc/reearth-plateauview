import { atom, useAtom } from "jotai";

const showGoogleModal = atom<boolean>(false);
export const useShowGoogleModal = () => useAtom(showGoogleModal);

const showTerrainModal = atom<boolean>(false);
export const useShowTerrainModal = () => useAtom(showTerrainModal);
