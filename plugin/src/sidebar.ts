import html from "../dist/web/sidebar/index.html?raw";

(globalThis as any).reearth.ui.show(html, { width: 370, extended: true });

console.log(html);
