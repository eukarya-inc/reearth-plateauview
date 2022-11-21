// vite-configs/location/web.ts
import { defineConfig } from "vite";

// vite.config.template.ts
import { resolve } from "path";
import react from "@vitejs/plugin-react";
import importToCDN, { autoComplete } from "vite-plugin-cdn-import";
import { viteSingleFile } from "vite-plugin-singlefile";
import svgr from "vite-plugin-svgr";
var __vite_injected_original_dirname = "C:\\plateau104\\reearth-plateauview\\plugin";
var web = ({
  name,
  parent,
  type = "core"
}) => () => {
  const root = parent && type !== "core" ? `./web/extensions/${parent}/${type}s/${name}` : `./web/extensions/${name}/${type}`;
  const outDir = parent && type !== "core" ? `../../../../../dist/web/${parent}/${type}s/${name}` : `../../../../dist/web/${name}/${type}`;
  return {
    plugins: [
      react(),
      serverHeaders(),
      viteSingleFile(),
      svgr(),
      importToCDN.default({
        modules: [
          autoComplete("react"),
          autoComplete("react-dom"),
          {
            name: "react-is",
            var: "react-is",
            path: "https://unpkg.com/react-is@18.2.0/umd/react-is.production.min.js"
          },
          {
            name: "antd",
            var: "antd",
            path: "https://cdnjs.cloudflare.com/ajax/libs/antd/4.22.8/antd.min.js",
            css: "https://cdnjs.cloudflare.com/ajax/libs/antd/4.22.8/antd.min.css"
          },
          {
            name: "styled-components",
            var: "styled-components",
            path: "https://unpkg.com/styled-components/dist/styled-components.min.js"
          }
        ]
      })
    ],
    publicDir: false,
    emptyOutDir: false,
    root,
    build: {
      outDir
    },
    css: {
      preprocessorOptions: {
        less: {
          javascriptEnabled: true,
          modifyVars: {
            "primary-color": "#00BEBE",
            "font-family": "Noto Sans",
            "typography-title-font-weight": "500",
            "typography-title-font-height": "21.79px"
          }
        }
      }
    },
    test: {
      globals: true,
      environment: "jsdom",
      setupFiles: "./web/test/setup.ts"
    },
    resolve: {
      alias: [{ find: "@web", replacement: resolve(__vite_injected_original_dirname, "web") }]
    }
  };
};
function serverHeaders() {
  return {
    name: "server-headers",
    configureServer(server) {
      server.middlewares.use((_req, res, next) => {
        res.setHeader("Service-Worker-Allowed", "/");
        next();
      });
    }
  };
}

// vite-configs/location/web.ts
var web_default = defineConfig(web({ name: "location" }));
export {
  web_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS1jb25maWdzL2xvY2F0aW9uL3dlYi50cyIsICJ2aXRlLmNvbmZpZy50ZW1wbGF0ZS50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIkM6XFxcXHBsYXRlYXUxMDRcXFxccmVlYXJ0aC1wbGF0ZWF1dmlld1xcXFxwbHVnaW5cXFxcdml0ZS1jb25maWdzXFxcXGxvY2F0aW9uXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJDOlxcXFxwbGF0ZWF1MTA0XFxcXHJlZWFydGgtcGxhdGVhdXZpZXdcXFxccGx1Z2luXFxcXHZpdGUtY29uZmlnc1xcXFxsb2NhdGlvblxcXFx3ZWIudHNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL0M6L3BsYXRlYXUxMDQvcmVlYXJ0aC1wbGF0ZWF1dmlldy9wbHVnaW4vdml0ZS1jb25maWdzL2xvY2F0aW9uL3dlYi50c1wiOy8vLyA8cmVmZXJlbmNlIHR5cGVzPVwidml0ZXN0XCIgLz5cbi8vLyA8cmVmZXJlbmNlIHR5cGVzPVwidml0ZS9jbGllbnRcIiAvPlxuXG5pbXBvcnQgeyBkZWZpbmVDb25maWcgfSBmcm9tIFwidml0ZVwiO1xuXG5pbXBvcnQgeyB3ZWIgfSBmcm9tIFwiLi4vLi4vdml0ZS5jb25maWcudGVtcGxhdGVcIjtcblxuLy8gaHR0cHM6Ly92aXRlanMuZGV2L2NvbmZpZy9cbmV4cG9ydCBkZWZhdWx0IGRlZmluZUNvbmZpZyh3ZWIoeyBuYW1lOiBcImxvY2F0aW9uXCIgfSkpO1xuIiwgImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJDOlxcXFxwbGF0ZWF1MTA0XFxcXHJlZWFydGgtcGxhdGVhdXZpZXdcXFxccGx1Z2luXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJDOlxcXFxwbGF0ZWF1MTA0XFxcXHJlZWFydGgtcGxhdGVhdXZpZXdcXFxccGx1Z2luXFxcXHZpdGUuY29uZmlnLnRlbXBsYXRlLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9DOi9wbGF0ZWF1MTA0L3JlZWFydGgtcGxhdGVhdXZpZXcvcGx1Z2luL3ZpdGUuY29uZmlnLnRlbXBsYXRlLnRzXCI7Ly8vIDxyZWZlcmVuY2UgdHlwZXM9XCJ2aXRlc3RcIiAvPlxuLy8vIDxyZWZlcmVuY2UgdHlwZXM9XCJ2aXRlL2NsaWVudFwiIC8+XG5cbmltcG9ydCB7IHJlc29sdmUgfSBmcm9tIFwicGF0aFwiO1xuXG5pbXBvcnQgcmVhY3QgZnJvbSBcIkB2aXRlanMvcGx1Z2luLXJlYWN0XCI7XG5pbXBvcnQgdHlwZSB7IFVzZXJDb25maWdFeHBvcnQsIFBsdWdpbiB9IGZyb20gXCJ2aXRlXCI7XG5pbXBvcnQgaW1wb3J0VG9DRE4sIHsgYXV0b0NvbXBsZXRlIH0gZnJvbSBcInZpdGUtcGx1Z2luLWNkbi1pbXBvcnRcIjtcbmltcG9ydCB7IHZpdGVTaW5nbGVGaWxlIH0gZnJvbSBcInZpdGUtcGx1Z2luLXNpbmdsZWZpbGVcIjtcbmltcG9ydCBzdmdyIGZyb20gXCJ2aXRlLXBsdWdpbi1zdmdyXCI7XG5cbmV4cG9ydCBjb25zdCBwbHVnaW4gPSAobmFtZTogc3RyaW5nKTogVXNlckNvbmZpZ0V4cG9ydCA9PiAoe1xuICBidWlsZDoge1xuICAgIG91dERpcjogXCJkaXN0L3BsdWdpblwiLFxuICAgIGVtcHR5T3V0RGlyOiBmYWxzZSxcbiAgICBsaWI6IHtcbiAgICAgIGZvcm1hdHM6IFtcImlpZmVcIl0sXG4gICAgICAvLyBodHRwczovL2dpdGh1Yi5jb20vdml0ZWpzL3ZpdGUvcHVsbC83MDQ3XG4gICAgICBlbnRyeTogYHNyYy8ke25hbWV9LnRzYCxcbiAgICAgIG5hbWU6IGBSZWVhcnRoUGx1Z2luUFZfJHtuYW1lfWAsXG4gICAgICBmaWxlTmFtZTogKCkgPT4gYCR7bmFtZX0uanNgLFxuICAgIH0sXG4gIH0sXG59KTtcblxuZXhwb3J0IGNvbnN0IHdlYiA9XG4gICh7XG4gICAgbmFtZSxcbiAgICBwYXJlbnQsXG4gICAgdHlwZSA9IFwiY29yZVwiLFxuICB9OiB7XG4gICAgbmFtZTogc3RyaW5nO1xuICAgIHBhcmVudD86IHN0cmluZztcbiAgICB0eXBlPzogXCJtb2RhbFwiIHwgXCJwbHVnaW5cIiB8IFwiY29yZVwiO1xuICB9KTogVXNlckNvbmZpZ0V4cG9ydCA9PlxuICAoKSA9PiB7XG4gICAgY29uc3Qgcm9vdCA9XG4gICAgICBwYXJlbnQgJiYgdHlwZSAhPT0gXCJjb3JlXCJcbiAgICAgICAgPyBgLi93ZWIvZXh0ZW5zaW9ucy8ke3BhcmVudH0vJHt0eXBlfXMvJHtuYW1lfWBcbiAgICAgICAgOiBgLi93ZWIvZXh0ZW5zaW9ucy8ke25hbWV9LyR7dHlwZX1gO1xuICAgIGNvbnN0IG91dERpciA9XG4gICAgICBwYXJlbnQgJiYgdHlwZSAhPT0gXCJjb3JlXCJcbiAgICAgICAgPyBgLi4vLi4vLi4vLi4vLi4vZGlzdC93ZWIvJHtwYXJlbnR9LyR7dHlwZX1zLyR7bmFtZX1gXG4gICAgICAgIDogYC4uLy4uLy4uLy4uL2Rpc3Qvd2ViLyR7bmFtZX0vJHt0eXBlfWA7XG5cbiAgICByZXR1cm4ge1xuICAgICAgcGx1Z2luczogW1xuICAgICAgICByZWFjdCgpLFxuICAgICAgICBzZXJ2ZXJIZWFkZXJzKCksXG4gICAgICAgIHZpdGVTaW5nbGVGaWxlKCksXG4gICAgICAgIHN2Z3IoKSxcbiAgICAgICAgKGltcG9ydFRvQ0ROIC8qIHdvcmthcm91bmQgKi8gYXMgYW55IGFzIHsgZGVmYXVsdDogdHlwZW9mIGltcG9ydFRvQ0ROIH0pLmRlZmF1bHQoe1xuICAgICAgICAgIG1vZHVsZXM6IFtcbiAgICAgICAgICAgIGF1dG9Db21wbGV0ZShcInJlYWN0XCIpLFxuICAgICAgICAgICAgYXV0b0NvbXBsZXRlKFwicmVhY3QtZG9tXCIpLFxuICAgICAgICAgICAge1xuICAgICAgICAgICAgICBuYW1lOiBcInJlYWN0LWlzXCIsXG4gICAgICAgICAgICAgIHZhcjogXCJyZWFjdC1pc1wiLFxuICAgICAgICAgICAgICBwYXRoOiBcImh0dHBzOi8vdW5wa2cuY29tL3JlYWN0LWlzQDE4LjIuMC91bWQvcmVhY3QtaXMucHJvZHVjdGlvbi5taW4uanNcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICB7XG4gICAgICAgICAgICAgIG5hbWU6IFwiYW50ZFwiLFxuICAgICAgICAgICAgICB2YXI6IFwiYW50ZFwiLFxuICAgICAgICAgICAgICBwYXRoOiBcImh0dHBzOi8vY2RuanMuY2xvdWRmbGFyZS5jb20vYWpheC9saWJzL2FudGQvNC4yMi44L2FudGQubWluLmpzXCIsXG4gICAgICAgICAgICAgIGNzczogXCJodHRwczovL2NkbmpzLmNsb3VkZmxhcmUuY29tL2FqYXgvbGlicy9hbnRkLzQuMjIuOC9hbnRkLm1pbi5jc3NcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICB7XG4gICAgICAgICAgICAgIG5hbWU6IFwic3R5bGVkLWNvbXBvbmVudHNcIixcbiAgICAgICAgICAgICAgdmFyOiBcInN0eWxlZC1jb21wb25lbnRzXCIsXG4gICAgICAgICAgICAgIHBhdGg6IFwiaHR0cHM6Ly91bnBrZy5jb20vc3R5bGVkLWNvbXBvbmVudHMvZGlzdC9zdHlsZWQtY29tcG9uZW50cy5taW4uanNcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgXSxcbiAgICAgICAgfSksXG4gICAgICBdLFxuICAgICAgcHVibGljRGlyOiBmYWxzZSxcbiAgICAgIGVtcHR5T3V0RGlyOiBmYWxzZSxcbiAgICAgIHJvb3QsXG4gICAgICBidWlsZDoge1xuICAgICAgICBvdXREaXIsXG4gICAgICB9LFxuICAgICAgY3NzOiB7XG4gICAgICAgIHByZXByb2Nlc3Nvck9wdGlvbnM6IHtcbiAgICAgICAgICBsZXNzOiB7XG4gICAgICAgICAgICBqYXZhc2NyaXB0RW5hYmxlZDogdHJ1ZSxcbiAgICAgICAgICAgIG1vZGlmeVZhcnM6IHtcbiAgICAgICAgICAgICAgXCJwcmltYXJ5LWNvbG9yXCI6IFwiIzAwQkVCRVwiLFxuICAgICAgICAgICAgICBcImZvbnQtZmFtaWx5XCI6IFwiTm90byBTYW5zXCIsXG4gICAgICAgICAgICAgIFwidHlwb2dyYXBoeS10aXRsZS1mb250LXdlaWdodFwiOiBcIjUwMFwiLFxuICAgICAgICAgICAgICBcInR5cG9ncmFwaHktdGl0bGUtZm9udC1oZWlnaHRcIjogXCIyMS43OXB4XCIsXG4gICAgICAgICAgICB9LFxuICAgICAgICAgIH0sXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgICAgdGVzdDoge1xuICAgICAgICBnbG9iYWxzOiB0cnVlLFxuICAgICAgICBlbnZpcm9ubWVudDogXCJqc2RvbVwiLFxuICAgICAgICBzZXR1cEZpbGVzOiBcIi4vd2ViL3Rlc3Qvc2V0dXAudHNcIixcbiAgICAgIH0sXG4gICAgICByZXNvbHZlOiB7XG4gICAgICAgIGFsaWFzOiBbeyBmaW5kOiBcIkB3ZWJcIiwgcmVwbGFjZW1lbnQ6IHJlc29sdmUoX19kaXJuYW1lLCBcIndlYlwiKSB9XSxcbiAgICAgIH0sXG4gICAgfTtcbiAgfTtcblxuZnVuY3Rpb24gc2VydmVySGVhZGVycygpOiBQbHVnaW4ge1xuICByZXR1cm4ge1xuICAgIG5hbWU6IFwic2VydmVyLWhlYWRlcnNcIixcbiAgICBjb25maWd1cmVTZXJ2ZXIoc2VydmVyKSB7XG4gICAgICBzZXJ2ZXIubWlkZGxld2FyZXMudXNlKChfcmVxLCByZXMsIG5leHQpID0+IHtcbiAgICAgICAgcmVzLnNldEhlYWRlcihcIlNlcnZpY2UtV29ya2VyLUFsbG93ZWRcIiwgXCIvXCIpO1xuICAgICAgICBuZXh0KCk7XG4gICAgICB9KTtcbiAgICB9LFxuICB9O1xufVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUdBLFNBQVMsb0JBQW9COzs7QUNBN0IsU0FBUyxlQUFlO0FBRXhCLE9BQU8sV0FBVztBQUVsQixPQUFPLGVBQWUsb0JBQW9CO0FBQzFDLFNBQVMsc0JBQXNCO0FBQy9CLE9BQU8sVUFBVTtBQVRqQixJQUFNLG1DQUFtQztBQXlCbEMsSUFBTSxNQUNYLENBQUM7QUFBQSxFQUNDO0FBQUEsRUFDQTtBQUFBLEVBQ0EsT0FBTztBQUNULE1BS0EsTUFBTTtBQUNKLFFBQU0sT0FDSixVQUFVLFNBQVMsU0FDZixvQkFBb0IsVUFBVSxTQUFTLFNBQ3ZDLG9CQUFvQixRQUFRO0FBQ2xDLFFBQU0sU0FDSixVQUFVLFNBQVMsU0FDZiwyQkFBMkIsVUFBVSxTQUFTLFNBQzlDLHdCQUF3QixRQUFRO0FBRXRDLFNBQU87QUFBQSxJQUNMLFNBQVM7QUFBQSxNQUNQLE1BQU07QUFBQSxNQUNOLGNBQWM7QUFBQSxNQUNkLGVBQWU7QUFBQSxNQUNmLEtBQUs7QUFBQSxNQUNKLFlBQXdFLFFBQVE7QUFBQSxRQUMvRSxTQUFTO0FBQUEsVUFDUCxhQUFhLE9BQU87QUFBQSxVQUNwQixhQUFhLFdBQVc7QUFBQSxVQUN4QjtBQUFBLFlBQ0UsTUFBTTtBQUFBLFlBQ04sS0FBSztBQUFBLFlBQ0wsTUFBTTtBQUFBLFVBQ1I7QUFBQSxVQUNBO0FBQUEsWUFDRSxNQUFNO0FBQUEsWUFDTixLQUFLO0FBQUEsWUFDTCxNQUFNO0FBQUEsWUFDTixLQUFLO0FBQUEsVUFDUDtBQUFBLFVBQ0E7QUFBQSxZQUNFLE1BQU07QUFBQSxZQUNOLEtBQUs7QUFBQSxZQUNMLE1BQU07QUFBQSxVQUNSO0FBQUEsUUFDRjtBQUFBLE1BQ0YsQ0FBQztBQUFBLElBQ0g7QUFBQSxJQUNBLFdBQVc7QUFBQSxJQUNYLGFBQWE7QUFBQSxJQUNiO0FBQUEsSUFDQSxPQUFPO0FBQUEsTUFDTDtBQUFBLElBQ0Y7QUFBQSxJQUNBLEtBQUs7QUFBQSxNQUNILHFCQUFxQjtBQUFBLFFBQ25CLE1BQU07QUFBQSxVQUNKLG1CQUFtQjtBQUFBLFVBQ25CLFlBQVk7QUFBQSxZQUNWLGlCQUFpQjtBQUFBLFlBQ2pCLGVBQWU7QUFBQSxZQUNmLGdDQUFnQztBQUFBLFlBQ2hDLGdDQUFnQztBQUFBLFVBQ2xDO0FBQUEsUUFDRjtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQSxNQUFNO0FBQUEsTUFDSixTQUFTO0FBQUEsTUFDVCxhQUFhO0FBQUEsTUFDYixZQUFZO0FBQUEsSUFDZDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsT0FBTyxDQUFDLEVBQUUsTUFBTSxRQUFRLGFBQWEsUUFBUSxrQ0FBVyxLQUFLLEVBQUUsQ0FBQztBQUFBLElBQ2xFO0FBQUEsRUFDRjtBQUNGO0FBRUYsU0FBUyxnQkFBd0I7QUFDL0IsU0FBTztBQUFBLElBQ0wsTUFBTTtBQUFBLElBQ04sZ0JBQWdCLFFBQVE7QUFDdEIsYUFBTyxZQUFZLElBQUksQ0FBQyxNQUFNLEtBQUssU0FBUztBQUMxQyxZQUFJLFVBQVUsMEJBQTBCLEdBQUc7QUFDM0MsYUFBSztBQUFBLE1BQ1AsQ0FBQztBQUFBLElBQ0g7QUFBQSxFQUNGO0FBQ0Y7OztBRDFHQSxJQUFPLGNBQVEsYUFBYSxJQUFJLEVBQUUsTUFBTSxXQUFXLENBQUMsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
