<template>
  <div style="padding: 10px;">
    <a-row :gutter="[16, 16]">
      <a-col :span="24">
        <a-input v-model:value="formData.parameterName" placeholder="参数名称" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-input v-model:value="formData.name" placeholder="节点名称" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select v-model:value="formData.sourceId" :options="sources" placeholder="数据库名称" style="width: 100%;">
        </a-select>
      </a-col>
      <a-col :span="24">
        <a-textarea v-model:value="formData.sql" :rows="5" placeholder="执行SQL" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select mode="tags" v-model:value="formData.wheres" :max-tags="1000" type="array" placeholder="WHERE条件"
          style="width: 100%;"></a-select>
      </a-col>
      <a-col :span="24">
        <!-- <a-auto-complete v-model:value="formData.args" :options="parameters" placeholder="输入参数" style="width: 100%;" /> -->
      </a-col>
      <a-col :span="24">
        <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
        <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
        <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
        <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
      </a-col>
      <a-col :span="24">
        <a-button type="primary" @click="bindForm">设置</a-button>
        <a-button @click="setForm" style="margin-left: 10px;">重置</a-button>
      </a-col>
    </a-row>
  </div>
</template>

<script>
import { message } from 'ant-design-vue';

export default {
  props: {
    lf: Object,
    currentNode: Object,
  },
  mounted() {
    this.getSources();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      sources: [],
      parameters: [],
      formData: {
        resultHistory: true,
      }
    }
  },
  methods: {
    setForm() {
      this.formData.name = this.$props.currentNode.data.text.value;
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.lf.getProperties(id);
      this.formData.sourceId = properties.sourceId;
      this.formData.args = properties.args;
      this.formData.wheres = properties.wheres || "";
      this.formData.sql = properties.sql;
      this.formData.parameterName = properties.parameterName;
      this.formData.printInput = properties.printInput;
      this.formData.printOutput = properties.printOutput;
      this.formData.inputHistory = properties.inputHistory;
      if (!!properties.resultHistory) {
        this.formData.resultHistory = properties.resultHistory;
      }
    },
    bindForm() {
      let id = this.$props.currentNode.data.id;
      this.$props.lf.updateText(id, this.formData.name);
      this.$props.lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    getSources() {
      this.$post("/dbdatasource/list", { page: { current: 1, size: 1000 } }).then(res => {
        if (res.data.code === this.$success) {
          this.sources = res.data.data.rows;
        } else {
          this.sources = [];
          message.error(res.data.msg);
        }
      });
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
      let nodes = flowData["nodes"];
      return nodes.find(node => node.type === 'StartNode');
    },
    getNextNode(flowData, id) {
      let edges = flowData["edges"];
      for (let edge of edges) {
        if (edge.sourceNodeId === id) {
          return this.getNode(flowData, edge.targetNodeId);
        }
      }
    },
    getNode(flowData, id) {
      let nodes = flowData["nodes"];
      return nodes.find(node => node.id === id);
    }
  }
}
</script>
