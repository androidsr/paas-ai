<template>
  <div style="padding: 10px;">
    <a-row :gutter="15">
      <a-col :span="24">
        <a-form-item label="输出内容">
          <div class="h-input-group" v-for="(item, i) in formData.messages" :key="i" style="margin-bottom: 10px;">
            <a-input v-model:value="item.varName" placeholder="变量名" style="width: 150px; margin-right: 8px;" />
            <span class="h-input-addon">=</span>
            <!-- <a-auto-complete v-model:value="item.varValue" :options="parameters" placeholder="变量值" style="width: 200px; margin-right: 8px;" /> -->
            <a-button @click="delMessageLine(i)" type="link" danger>删除</a-button>
          </div>
          <a-button @click="addMessageLine" type="link">添加一行</a-button>
        </a-form-item>
      </a-col>

      <a-col :span="24">
        <a-button type="primary" @click="bindForm">设置</a-button>
        <a-button @click="setForm" style="margin-left: 10px;" type="default">重置</a-button>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import LogicFlow from '@logicflow/core';
import { message } from 'ant-design-vue';

export default {
  props: {
    lf: Object,
    currentNode: Object,
  },
  mounted() {
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      parameters: [],
      formData: {
        messages: []
      }
    };
  },
  methods: {
    setForm() {
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.lf.getProperties(id);
      if (properties.messages) {
        this.formData.messages = properties.messages;
      }
    },
    bindForm() {
      let id = this.$props.currentNode.data.id;
      this.$props.lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    addMessageLine() {
      this.formData.messages.push({ varName: "", varValue: "" });
    },
    delMessageLine(index) {
      this.formData.messages.splice(index, 1);
    },
    getParameters() {
      let flowData = this.$props.lf.getGraphRawData();
      let nodes = flowData["nodes"];
      let node = this.getStartNode(flowData);
      while (node && node?.id !== this.$props.currentNode?.data?.id && node?.type !== 'EndNode') {
        let properties = node?.properties;
        let parameterName = properties?.parameterName;

        if (parameterName) {
          if (properties?.inputHistory) {
            this.parameters.push({ label: "$" + parameterName + ".input", value: "$" + parameterName + ".input" });
          }
          if (properties?.resultHistory) {
            this.parameters.push({ label: "$" + parameterName + ".output", value: "$" + parameterName + ".output" });
          }
        }
        node = this.getNextNode(flowData, node.id);
      }
    },
    getStartNode(flowData) {
      return flowData["nodes"].find(node => node.type === 'StartNode');
    },
    getNode(flowData, id) {
      return flowData["nodes"].find(node => node.id === id);
    },
    getNextNode(flowData, id) {
      let edges = flowData["edges"];
      let edge = edges.find(edge => edge.sourceNodeId === id);
      return edge ? this.getNode(flowData, edge.targetNodeId) : null;
    },
  }
}
</script>

<style scoped>
.h-input-group {
  display: flex;
  align-items: center;
}
.h-input-addon {
  margin: 0 8px;
}
</style>
