<template>
  <a-card size="small" :title="$store.menus.title">
    <paas-page-table :columns="columns" module="/fwconfig"></paas-page-table>
  </a-card>
</template>

<script>
import { message, Modal } from 'ant-design-vue';

export default {
  mounted() {
  },
  created() {
  },
  data() {
    return {
      columns: [
        { title: "流程名称", dataIndex: "name" },
        {
          title: "流程类型", dataIndex: "flowType",
          customRender: (text) => this.$dictReader("flowType", text.value, this.$colors),
        }, 
        {
          title: "流程状态", dataIndex: "flowState",
          customRender: (text) => this.$dictReader("flowState", text.value, this.$colors),
        },
        { title: "流程描述", dataIndex: "remark" },
        { title: '操作', key: 'action', fixed: 'right', width: 180 }
      ],
    };
  },
  provide() {
    return {
      parentMethods: {
        configClick: this.configClick
      }
    };
  },
  methods: {
    configClick(data) {
      this.$router.push({
        path: "/flow/design",
        query: {
          id: data.id,
        },
      });
    },
  },
};
</script>