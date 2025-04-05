<template>
  <a-form :label-col="{ style: { width: '80px' } }" label-align="left">
    <!-- 参数名称 -->
    <a-form-item label="参数名称">
      <a-input v-model:value="formData.parameterName" placeholder="参数名称" />
    </a-form-item>

    <!-- 节点名称 -->
    <a-form-item label="节点名称">
      <a-input v-model:value="formData.name" placeholder="节点名称" />
    </a-form-item>

    <!-- 代码块内容 -->
    <a-form-item label="代码块内容">
      <a-mentions v-model:value="formData.content" autofocus :options="parameters" placeholder="输入$补全" prefix="$"
        :rows="10"></a-mentions>
    </a-form-item>

    <a-form-item label="输出设置">
      <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
      <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
      <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
      <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
    </a-form-item>
    <br />
    <a-form-item>
      <a-button @click="bindForm" type="primary">保存</a-button>
      <a-button @click="setForm" type="default" style="margin-left: 15px;">重置</a-button>
    </a-form-item>
  </a-form>
</template>

<script>
import { message } from 'ant-design-vue';
import { GetList } from '../../../../../../wailsjs/go/biz/DbconfigBiz.js';

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
        parameterName: "database",
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
      this.formData.wheres = properties.wheres || [];
      this.formData.sql = properties.sql;
      this.formData.parameterName = properties.parameterName;
      this.formData.printInput = properties.printInput;
      this.formData.printOutput = properties.printOutput;
      this.formData.inputHistory = properties.inputHistory;
      this.formData.content = properties.content;
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
      GetList().then(res => {
        if (res.code === 200) {
          this.sources = res.data;
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
            this.parameters.push({ label: properties?.name + "-输入", value: "input." + parameterName });
          }
          if (properties?.resultHistory) {
            this.parameters.push({ label: properties?.name + "-输出", value: "output." + parameterName });
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
