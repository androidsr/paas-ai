<template>
  <a-config-provider :locale="locale" :theme="{
    token: {
      borderRadius: 1,
      sizeStep: 4,
      sizeUnit: 4,
      wireframe: true,
      margin: '8px',
      fontSize: 14,
      colorPrimary: '#007bff'
    },
    components: {

    },
  }">
    <a-layout>
      <!-- <a-layout-sider :style="{ minHeight: '100vh' }">
        <a-menu mode="inline" @click="menuClick" theme="dark" :selectedKeys="[selectedKey]" v-show="$store.isAbout">
          <a-menu-item key="/aiChat">
            <HomeOutlined />&nbsp;&nbsp;AI对话
          </a-menu-item>
          <a-menu-item key="/prompt">
            <FieldTimeOutlined />&nbsp;&nbsp;提示词管理
          </a-menu-item>
          <a-menu-item key="/knowledge">
            <OrderedListOutlined />&nbsp;&nbsp;知识库管理
          </a-menu-item>
          <a-menu-item key="/function">
            <DatabaseOutlined />&nbsp;&nbsp;函数管理
          </a-menu-item>
          <a-menu-item key="/database">
            <DatabaseOutlined />&nbsp;&nbsp;数据库管理
          </a-menu-item>
          <a-menu-item key="/flow">
            <DatabaseOutlined />&nbsp;&nbsp;流程管理
          </a-menu-item>
          <a-menu-item key="/config">
            <SettingOutlined />&nbsp;&nbsp;设置
          </a-menu-item>
          <a-menu-item key="/about">
            <InfoCircleOutlined />&nbsp;&nbsp;关于
          </a-menu-item>
        </a-menu>
      </a-layout-sider> -->
      <a-layout>
        <a-layout-content style="background-color: white;padding: 10px;">
          <router-view></router-view>
        </a-layout-content>
      </a-layout>
    </a-layout>
  </a-config-provider>
</template>

<script>
import { EventsOn } from '../wailsjs/runtime';
import { GetConfig } from '../wailsjs/go/main/App';

export default {
  mounted() {
    // 监听来自 Go 后端的事件
    EventsOn('go-to-page', (page) => {
      this.goToPage(page);
    });
    this.$router.push(this.selectedKey);
    GetConfig().then((res) => {
      if (res.code == 200) {
        this.config = res.data;
      }
    })
  },
  watch: {
    '$route': function (newRoute) {
      if (newRoute.path.split('/').length - 1 == 1) {
        this.selectedKey = newRoute.path;
      }
    }
  },
  data() {
    return {
      config: {},
      selectedKey: "/aiChat",
      menus: [],
    };
  },
  methods: {
    goToPage(page) {
      this.$router.push(page);
    },
    menuClick(data) {
      if (!data.key) {
        return;
      }
      this.$store.setAction("list");
      this.$router.push(data.key)
    }
  },
};
</script>
<style scoped>
.iframeStyle {
  width: 100%;
  border: none;
  height: 89vh;
}
</style>