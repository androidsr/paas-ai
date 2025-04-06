// vite.config.js
import vue from "file:///D:/dev/go/work/paas-ai/frontend/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { resolve } from "path";
import { defineConfig } from "file:///D:/dev/go/work/paas-ai/frontend/node_modules/vite/dist/node/index.js";
import Pages from "file:///D:/dev/go/work/paas-ai/frontend/node_modules/vite-plugin-pages/dist/index.js";
var __vite_injected_original_dirname = "D:\\dev\\go\\work\\paas-ai\\frontend";
var vite_config_default = defineConfig({
  plugins: [
    vue(),
    Pages({
      dirs: [{ dir: "src/pages", baseRoute: "" }],
      importMode: "async"
    })
  ],
  resolve: {
    alias: {
      "@": resolve(__vite_injected_original_dirname, "src"),
      "*": resolve("")
    },
    transpileDependencies: true,
    lintOnSave: false
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJEOlxcXFxkZXZcXFxcZ29cXFxcd29ya1xcXFxwYWFzLWFpXFxcXGZyb250ZW5kXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ZpbGVuYW1lID0gXCJEOlxcXFxkZXZcXFxcZ29cXFxcd29ya1xcXFxwYWFzLWFpXFxcXGZyb250ZW5kXFxcXHZpdGUuY29uZmlnLmpzXCI7Y29uc3QgX192aXRlX2luamVjdGVkX29yaWdpbmFsX2ltcG9ydF9tZXRhX3VybCA9IFwiZmlsZTovLy9EOi9kZXYvZ28vd29yay9wYWFzLWFpL2Zyb250ZW5kL3ZpdGUuY29uZmlnLmpzXCI7aW1wb3J0IHZ1ZSBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUnO1xuaW1wb3J0IHsgcmVzb2x2ZSB9IGZyb20gJ3BhdGgnO1xuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSc7XG5pbXBvcnQgUGFnZXMgZnJvbSBcInZpdGUtcGx1Z2luLXBhZ2VzXCI7XG5cbi8vIGh0dHBzOi8vdml0ZWpzLmRldi9jb25maWcvXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICBwbHVnaW5zOiBbXG4gICAgdnVlKCksXG4gICAgUGFnZXMoe1xuICAgICAgZGlyczogW3sgZGlyOiBcInNyYy9wYWdlc1wiLCBiYXNlUm91dGU6IFwiXCIgfV0sXG4gICAgICBpbXBvcnRNb2RlOiBcImFzeW5jXCIsXG4gICAgfSlcbiAgXSxcbiAgcmVzb2x2ZToge1xuICAgIGFsaWFzOiB7XG4gICAgICAnQCc6IHJlc29sdmUoX19kaXJuYW1lLCAnc3JjJyksIC8vIFx1OEJCRVx1N0Y2RSBAIFx1NTIyQlx1NTQwRFx1NEUzQSBzcmMgXHU3NkVFXHU1RjU1XG4gICAgICAnKic6IHJlc29sdmUoJycpXG4gICAgfSxcbiAgICB0cmFuc3BpbGVEZXBlbmRlbmNpZXM6IHRydWUsXG4gICAgbGludE9uU2F2ZTogZmFsc2UsXG4gIH0sXG4gIFxufSlcbiJdLAogICJtYXBwaW5ncyI6ICI7QUFBMlIsT0FBTyxTQUFTO0FBQzNTLFNBQVMsZUFBZTtBQUN4QixTQUFTLG9CQUFvQjtBQUM3QixPQUFPLFdBQVc7QUFIbEIsSUFBTSxtQ0FBbUM7QUFNekMsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsU0FBUztBQUFBLElBQ1AsSUFBSTtBQUFBLElBQ0osTUFBTTtBQUFBLE1BQ0osTUFBTSxDQUFDLEVBQUUsS0FBSyxhQUFhLFdBQVcsR0FBRyxDQUFDO0FBQUEsTUFDMUMsWUFBWTtBQUFBLElBQ2QsQ0FBQztBQUFBLEVBQ0g7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEtBQUssUUFBUSxrQ0FBVyxLQUFLO0FBQUEsTUFDN0IsS0FBSyxRQUFRLEVBQUU7QUFBQSxJQUNqQjtBQUFBLElBQ0EsdUJBQXVCO0FBQUEsSUFDdkIsWUFBWTtBQUFBLEVBQ2Q7QUFFRixDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
