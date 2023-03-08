// vite-configs/sidebar/popups/help/web.ts
import { defineConfig } from "file:///Users/lby/reearth/projects/reearth-plateauview/plugin/node_modules/vite/dist/node/index.js";

// vite.config.template.ts
import { resolve } from "path";
import react from "file:///Users/lby/reearth/projects/reearth-plateauview/plugin/node_modules/@vitejs/plugin-react/dist/index.mjs";
import importToCDN, { autoComplete } from "file:///Users/lby/reearth/projects/reearth-plateauview/plugin/node_modules/vite-plugin-cdn-import/dist/index.js";
import { viteSingleFile } from "file:///Users/lby/reearth/projects/reearth-plateauview/plugin/node_modules/vite-plugin-singlefile/dist/esm/index.js";
import svgr from "file:///Users/lby/reearth/projects/reearth-plateauview/plugin/node_modules/vite-plugin-svgr/dist/index.mjs";
var __vite_injected_original_dirname = "/Users/lby/reearth/projects/reearth-plateauview/plugin";
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
            path: "https://unpkg.com/styled-components@5.3.6/dist/styled-components.min.js"
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

// vite-configs/sidebar/popups/help/web.ts
var web_default = defineConfig(web({ name: "help", parent: "sidebar", type: "popup" }));
export {
  web_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS1jb25maWdzL3NpZGViYXIvcG9wdXBzL2hlbHAvd2ViLnRzIiwgInZpdGUuY29uZmlnLnRlbXBsYXRlLnRzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyJjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZGlybmFtZSA9IFwiL1VzZXJzL2xieS9yZWVhcnRoL3Byb2plY3RzL3JlZWFydGgtcGxhdGVhdXZpZXcvcGx1Z2luL3ZpdGUtY29uZmlncy9zaWRlYmFyL3BvcHVwcy9oZWxwXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCIvVXNlcnMvbGJ5L3JlZWFydGgvcHJvamVjdHMvcmVlYXJ0aC1wbGF0ZWF1dmlldy9wbHVnaW4vdml0ZS1jb25maWdzL3NpZGViYXIvcG9wdXBzL2hlbHAvd2ViLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9Vc2Vycy9sYnkvcmVlYXJ0aC9wcm9qZWN0cy9yZWVhcnRoLXBsYXRlYXV2aWV3L3BsdWdpbi92aXRlLWNvbmZpZ3Mvc2lkZWJhci9wb3B1cHMvaGVscC93ZWIudHNcIjsvLy8gPHJlZmVyZW5jZSB0eXBlcz1cInZpdGVzdFwiIC8+XG4vLy8gPHJlZmVyZW5jZSB0eXBlcz1cInZpdGUvY2xpZW50XCIgLz5cblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSBcInZpdGVcIjtcblxuaW1wb3J0IHsgd2ViIH0gZnJvbSBcIi4uLy4uLy4uLy4uL3ZpdGUuY29uZmlnLnRlbXBsYXRlXCI7XG5cbi8vIGh0dHBzOi8vdml0ZWpzLmRldi9jb25maWcvXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcod2ViKHsgbmFtZTogXCJoZWxwXCIsIHBhcmVudDogXCJzaWRlYmFyXCIsIHR5cGU6IFwicG9wdXBcIiB9KSk7XG4iLCAiY29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2Rpcm5hbWUgPSBcIi9Vc2Vycy9sYnkvcmVlYXJ0aC9wcm9qZWN0cy9yZWVhcnRoLXBsYXRlYXV2aWV3L3BsdWdpblwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiL1VzZXJzL2xieS9yZWVhcnRoL3Byb2plY3RzL3JlZWFydGgtcGxhdGVhdXZpZXcvcGx1Z2luL3ZpdGUuY29uZmlnLnRlbXBsYXRlLnRzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9Vc2Vycy9sYnkvcmVlYXJ0aC9wcm9qZWN0cy9yZWVhcnRoLXBsYXRlYXV2aWV3L3BsdWdpbi92aXRlLmNvbmZpZy50ZW1wbGF0ZS50c1wiOy8vLyA8cmVmZXJlbmNlIHR5cGVzPVwidml0ZXN0XCIgLz5cbi8vLyA8cmVmZXJlbmNlIHR5cGVzPVwidml0ZS9jbGllbnRcIiAvPlxuXG5pbXBvcnQgeyByZXNvbHZlIH0gZnJvbSBcInBhdGhcIjtcblxuaW1wb3J0IHJlYWN0IGZyb20gXCJAdml0ZWpzL3BsdWdpbi1yZWFjdFwiO1xuaW1wb3J0IHR5cGUgeyBVc2VyQ29uZmlnRXhwb3J0LCBQbHVnaW4gfSBmcm9tIFwidml0ZVwiO1xuaW1wb3J0IGltcG9ydFRvQ0ROLCB7IGF1dG9Db21wbGV0ZSB9IGZyb20gXCJ2aXRlLXBsdWdpbi1jZG4taW1wb3J0XCI7XG5pbXBvcnQgeyB2aXRlU2luZ2xlRmlsZSB9IGZyb20gXCJ2aXRlLXBsdWdpbi1zaW5nbGVmaWxlXCI7XG5pbXBvcnQgc3ZnciBmcm9tIFwidml0ZS1wbHVnaW4tc3ZnclwiO1xuXG5leHBvcnQgY29uc3QgcGx1Z2luID0gKG5hbWU6IHN0cmluZyk6IFVzZXJDb25maWdFeHBvcnQgPT4gKHtcbiAgYnVpbGQ6IHtcbiAgICBvdXREaXI6IFwiZGlzdC9wbHVnaW5cIixcbiAgICBlbXB0eU91dERpcjogZmFsc2UsXG4gICAgbGliOiB7XG4gICAgICBmb3JtYXRzOiBbXCJpaWZlXCJdLFxuICAgICAgLy8gaHR0cHM6Ly9naXRodWIuY29tL3ZpdGVqcy92aXRlL3B1bGwvNzA0N1xuICAgICAgZW50cnk6IGBzcmMvJHtuYW1lfS50c2AsXG4gICAgICBuYW1lOiBgUmVlYXJ0aFBsdWdpblBWXyR7bmFtZX1gLFxuICAgICAgZmlsZU5hbWU6ICgpID0+IGAke25hbWV9LmpzYCxcbiAgICB9LFxuICB9LFxufSk7XG5cbmV4cG9ydCBjb25zdCB3ZWIgPVxuICAoe1xuICAgIG5hbWUsXG4gICAgcGFyZW50LFxuICAgIHR5cGUgPSBcImNvcmVcIixcbiAgfToge1xuICAgIG5hbWU6IHN0cmluZztcbiAgICBwYXJlbnQ/OiBzdHJpbmc7XG4gICAgdHlwZT86IFwibW9kYWxcIiB8IFwiY29yZVwiIHwgXCJwb3B1cFwiO1xuICB9KTogVXNlckNvbmZpZ0V4cG9ydCA9PlxuICAoKSA9PiB7XG4gICAgY29uc3Qgcm9vdCA9XG4gICAgICBwYXJlbnQgJiYgdHlwZSAhPT0gXCJjb3JlXCJcbiAgICAgICAgPyBgLi93ZWIvZXh0ZW5zaW9ucy8ke3BhcmVudH0vJHt0eXBlfXMvJHtuYW1lfWBcbiAgICAgICAgOiBgLi93ZWIvZXh0ZW5zaW9ucy8ke25hbWV9LyR7dHlwZX1gO1xuICAgIGNvbnN0IG91dERpciA9XG4gICAgICBwYXJlbnQgJiYgdHlwZSAhPT0gXCJjb3JlXCJcbiAgICAgICAgPyBgLi4vLi4vLi4vLi4vLi4vZGlzdC93ZWIvJHtwYXJlbnR9LyR7dHlwZX1zLyR7bmFtZX1gXG4gICAgICAgIDogYC4uLy4uLy4uLy4uL2Rpc3Qvd2ViLyR7bmFtZX0vJHt0eXBlfWA7XG5cbiAgICByZXR1cm4ge1xuICAgICAgcGx1Z2luczogW1xuICAgICAgICByZWFjdCgpLFxuICAgICAgICBzZXJ2ZXJIZWFkZXJzKCksXG4gICAgICAgIHZpdGVTaW5nbGVGaWxlKCksXG4gICAgICAgIHN2Z3IoKSxcbiAgICAgICAgKGltcG9ydFRvQ0ROIC8qIHdvcmthcm91bmQgKi8gYXMgYW55IGFzIHsgZGVmYXVsdDogdHlwZW9mIGltcG9ydFRvQ0ROIH0pLmRlZmF1bHQoe1xuICAgICAgICAgIG1vZHVsZXM6IFtcbiAgICAgICAgICAgIGF1dG9Db21wbGV0ZShcInJlYWN0XCIpLFxuICAgICAgICAgICAgYXV0b0NvbXBsZXRlKFwicmVhY3QtZG9tXCIpLFxuICAgICAgICAgICAge1xuICAgICAgICAgICAgICBuYW1lOiBcInJlYWN0LWlzXCIsXG4gICAgICAgICAgICAgIHZhcjogXCJyZWFjdC1pc1wiLFxuICAgICAgICAgICAgICBwYXRoOiBcImh0dHBzOi8vdW5wa2cuY29tL3JlYWN0LWlzQDE4LjIuMC91bWQvcmVhY3QtaXMucHJvZHVjdGlvbi5taW4uanNcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICB7XG4gICAgICAgICAgICAgIG5hbWU6IFwiYW50ZFwiLFxuICAgICAgICAgICAgICB2YXI6IFwiYW50ZFwiLFxuICAgICAgICAgICAgICBwYXRoOiBcImh0dHBzOi8vY2RuanMuY2xvdWRmbGFyZS5jb20vYWpheC9saWJzL2FudGQvNC4yMi44L2FudGQubWluLmpzXCIsXG4gICAgICAgICAgICAgIGNzczogXCJodHRwczovL2NkbmpzLmNsb3VkZmxhcmUuY29tL2FqYXgvbGlicy9hbnRkLzQuMjIuOC9hbnRkLm1pbi5jc3NcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICB7XG4gICAgICAgICAgICAgIG5hbWU6IFwic3R5bGVkLWNvbXBvbmVudHNcIixcbiAgICAgICAgICAgICAgdmFyOiBcInN0eWxlZC1jb21wb25lbnRzXCIsXG4gICAgICAgICAgICAgIHBhdGg6IFwiaHR0cHM6Ly91bnBrZy5jb20vc3R5bGVkLWNvbXBvbmVudHNANS4zLjYvZGlzdC9zdHlsZWQtY29tcG9uZW50cy5taW4uanNcIixcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgXSxcbiAgICAgICAgfSksXG4gICAgICBdLFxuICAgICAgcHVibGljRGlyOiBmYWxzZSxcbiAgICAgIGVtcHR5T3V0RGlyOiBmYWxzZSxcbiAgICAgIHJvb3QsXG4gICAgICBidWlsZDoge1xuICAgICAgICBvdXREaXIsXG4gICAgICB9LFxuICAgICAgY3NzOiB7XG4gICAgICAgIHByZXByb2Nlc3Nvck9wdGlvbnM6IHtcbiAgICAgICAgICBsZXNzOiB7XG4gICAgICAgICAgICBqYXZhc2NyaXB0RW5hYmxlZDogdHJ1ZSxcbiAgICAgICAgICAgIG1vZGlmeVZhcnM6IHtcbiAgICAgICAgICAgICAgXCJwcmltYXJ5LWNvbG9yXCI6IFwiIzAwQkVCRVwiLFxuICAgICAgICAgICAgICBcImZvbnQtZmFtaWx5XCI6IFwiTm90byBTYW5zXCIsXG4gICAgICAgICAgICAgIFwidHlwb2dyYXBoeS10aXRsZS1mb250LXdlaWdodFwiOiBcIjUwMFwiLFxuICAgICAgICAgICAgICBcInR5cG9ncmFwaHktdGl0bGUtZm9udC1oZWlnaHRcIjogXCIyMS43OXB4XCIsXG4gICAgICAgICAgICB9LFxuICAgICAgICAgIH0sXG4gICAgICAgIH0sXG4gICAgICB9LFxuICAgICAgdGVzdDoge1xuICAgICAgICBnbG9iYWxzOiB0cnVlLFxuICAgICAgICBlbnZpcm9ubWVudDogXCJqc2RvbVwiLFxuICAgICAgICBzZXR1cEZpbGVzOiBcIi4vd2ViL3Rlc3Qvc2V0dXAudHNcIixcbiAgICAgIH0sXG4gICAgICByZXNvbHZlOiB7XG4gICAgICAgIGFsaWFzOiBbeyBmaW5kOiBcIkB3ZWJcIiwgcmVwbGFjZW1lbnQ6IHJlc29sdmUoX19kaXJuYW1lLCBcIndlYlwiKSB9XSxcbiAgICAgIH0sXG4gICAgfTtcbiAgfTtcblxuZnVuY3Rpb24gc2VydmVySGVhZGVycygpOiBQbHVnaW4ge1xuICByZXR1cm4ge1xuICAgIG5hbWU6IFwic2VydmVyLWhlYWRlcnNcIixcbiAgICBjb25maWd1cmVTZXJ2ZXIoc2VydmVyKSB7XG4gICAgICBzZXJ2ZXIubWlkZGxld2FyZXMudXNlKChfcmVxLCByZXMsIG5leHQpID0+IHtcbiAgICAgICAgcmVzLnNldEhlYWRlcihcIlNlcnZpY2UtV29ya2VyLUFsbG93ZWRcIiwgXCIvXCIpO1xuICAgICAgICBuZXh0KCk7XG4gICAgICB9KTtcbiAgICB9LFxuICB9O1xufVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUdBLFNBQVMsb0JBQW9COzs7QUNBN0IsU0FBUyxlQUFlO0FBRXhCLE9BQU8sV0FBVztBQUVsQixPQUFPLGVBQWUsb0JBQW9CO0FBQzFDLFNBQVMsc0JBQXNCO0FBQy9CLE9BQU8sVUFBVTtBQVRqQixJQUFNLG1DQUFtQztBQXlCbEMsSUFBTSxNQUNYLENBQUM7QUFBQSxFQUNDO0FBQUEsRUFDQTtBQUFBLEVBQ0EsT0FBTztBQUNULE1BS0EsTUFBTTtBQUNKLFFBQU0sT0FDSixVQUFVLFNBQVMsU0FDZixvQkFBb0IsVUFBVSxTQUFTLFNBQ3ZDLG9CQUFvQixRQUFRO0FBQ2xDLFFBQU0sU0FDSixVQUFVLFNBQVMsU0FDZiwyQkFBMkIsVUFBVSxTQUFTLFNBQzlDLHdCQUF3QixRQUFRO0FBRXRDLFNBQU87QUFBQSxJQUNMLFNBQVM7QUFBQSxNQUNQLE1BQU07QUFBQSxNQUNOLGNBQWM7QUFBQSxNQUNkLGVBQWU7QUFBQSxNQUNmLEtBQUs7QUFBQSxNQUNKLFlBQXdFLFFBQVE7QUFBQSxRQUMvRSxTQUFTO0FBQUEsVUFDUCxhQUFhLE9BQU87QUFBQSxVQUNwQixhQUFhLFdBQVc7QUFBQSxVQUN4QjtBQUFBLFlBQ0UsTUFBTTtBQUFBLFlBQ04sS0FBSztBQUFBLFlBQ0wsTUFBTTtBQUFBLFVBQ1I7QUFBQSxVQUNBO0FBQUEsWUFDRSxNQUFNO0FBQUEsWUFDTixLQUFLO0FBQUEsWUFDTCxNQUFNO0FBQUEsWUFDTixLQUFLO0FBQUEsVUFDUDtBQUFBLFVBQ0E7QUFBQSxZQUNFLE1BQU07QUFBQSxZQUNOLEtBQUs7QUFBQSxZQUNMLE1BQU07QUFBQSxVQUNSO0FBQUEsUUFDRjtBQUFBLE1BQ0YsQ0FBQztBQUFBLElBQ0g7QUFBQSxJQUNBLFdBQVc7QUFBQSxJQUNYLGFBQWE7QUFBQSxJQUNiO0FBQUEsSUFDQSxPQUFPO0FBQUEsTUFDTDtBQUFBLElBQ0Y7QUFBQSxJQUNBLEtBQUs7QUFBQSxNQUNILHFCQUFxQjtBQUFBLFFBQ25CLE1BQU07QUFBQSxVQUNKLG1CQUFtQjtBQUFBLFVBQ25CLFlBQVk7QUFBQSxZQUNWLGlCQUFpQjtBQUFBLFlBQ2pCLGVBQWU7QUFBQSxZQUNmLGdDQUFnQztBQUFBLFlBQ2hDLGdDQUFnQztBQUFBLFVBQ2xDO0FBQUEsUUFDRjtBQUFBLE1BQ0Y7QUFBQSxJQUNGO0FBQUEsSUFDQSxNQUFNO0FBQUEsTUFDSixTQUFTO0FBQUEsTUFDVCxhQUFhO0FBQUEsTUFDYixZQUFZO0FBQUEsSUFDZDtBQUFBLElBQ0EsU0FBUztBQUFBLE1BQ1AsT0FBTyxDQUFDLEVBQUUsTUFBTSxRQUFRLGFBQWEsUUFBUSxrQ0FBVyxLQUFLLEVBQUUsQ0FBQztBQUFBLElBQ2xFO0FBQUEsRUFDRjtBQUNGO0FBRUYsU0FBUyxnQkFBd0I7QUFDL0IsU0FBTztBQUFBLElBQ0wsTUFBTTtBQUFBLElBQ04sZ0JBQWdCLFFBQVE7QUFDdEIsYUFBTyxZQUFZLElBQUksQ0FBQyxNQUFNLEtBQUssU0FBUztBQUMxQyxZQUFJLFVBQVUsMEJBQTBCLEdBQUc7QUFDM0MsYUFBSztBQUFBLE1BQ1AsQ0FBQztBQUFBLElBQ0g7QUFBQSxFQUNGO0FBQ0Y7OztBRDFHQSxJQUFPLGNBQVEsYUFBYSxJQUFJLEVBQUUsTUFBTSxRQUFRLFFBQVEsV0FBVyxNQUFNLFFBQVEsQ0FBQyxDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
