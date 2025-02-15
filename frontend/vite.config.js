import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';
import { defineConfig } from 'vite';
import Pages from "vite-plugin-pages";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Pages({
      dirs: [{ dir: "src/pages", baseRoute: "" }],
      importMode: "async",
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'), // 设置 @ 别名为 src 目录
      '*': resolve('')
    },
    transpileDependencies: true,
    lintOnSave: false,
  },
  
})
