import html from "../dist/web/sidebar/index.html?raw";

(globalThis as any).reearth.ui.show(html, {
  width: 370,
  height: 880,
  margin: 0,
});
console.log(html);
