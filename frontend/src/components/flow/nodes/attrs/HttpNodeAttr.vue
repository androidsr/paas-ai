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
        <a-input v-model:value="formData.url" placeholder="请求地址" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select v-model:value="formData.contentType" :options="contentTypes" placeholder="请求类型"
          style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select v-model:value="formData.method" :options="methods" placeholder="请求方式" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-input v-model:value="formData.headers" placeholder="Header" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-input v-model:value="formData.cookies" placeholder="Cookie" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <!-- <a-auto-complete v-model:value="formData.body" :options="parameters" placeholder="Body" style="width: 100%;" /> -->
      </a-col>
      <a-col :span="24">
        <div class="h-input-group" style="display: flex; width: 100%;">
          <a-input v-model:value="formData.validateKey" style="width: 100px;" />
          <a-select v-model:value="formData.expression" :options="expressions"
            style="width: 100px; margin-left: 10px;" />
          <a-input v-model:value="formData.validateValue" style="width: 176px; margin-left: 10px;" />
        </div>
      </a-col>
      <a-col :span="24">
        <a-input v-model:value="formData.dataKey" placeholder="保留结果：data || result" style="width: 100%;" />
      </a-col>
      <a-col span="24">
        <a-space>
          <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
          <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
          <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
          <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
        </a-space>
      </a-col>

      <a-col span="24" style="margin-top: 15px;">
        <a-space>
          <a-button type="primary" @click="bindForm">设置</a-button>
          <a-button @click="setForm" type="default">重置</a-button>
        </a-space>
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
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      expressions: [
        { label: "等于", value: "=" },
        { label: "不等于", value: "!=" },
        { label: "大于", value: ">" },
        { label: "大于等于", value: ">=" },
        { label: "小于", value: "<" },
        { label: "小于等于", value: "<=" },
        { label: "包含", value: "include" },
        { label: "不包含", value: "exclusive" },
        { label: "以*开始", value: "startWith" },
        { label: "以*结束", value: "endWith" },
      ],
      contentTypes: [
        { label: "application/json", value: "application/json" },
        { label: "application/x-www-form-urlencoded", value: "application/x-www-form-urlencoded" },
        { label: "multipart/form-data", value: "multipart/form-data" },
        { label: "text/html", value: "text/html" },
      ],
      methods: [
        { label: "GET", value: "GET" },
        { label: "POST", value: "POST" },
        { label: "PUT", value: "PUT" },
        { label: "DELETE", value: "DELETE" },
      ],
      parameters: [],
      formData: {
        validateKey: "code",
        expression: "=",
        validateValue: "200",
        dataKey: "",
        method: "POST",
        contentType: "application/json",
        resultHistory: true,
      }
    };
  },
  methods: {
    setForm() {
      const { currentNode, lf } = this.$props;
      const properties = lf.getProperties(currentNode.data.id);
      this.formData = {
        ...this.formData,
        name: currentNode.data.text.value,
        url: properties.url,
        contentType: properties.contentType || this.formData.contentType,
        method: properties.method || this.formData.method,
        dataKey: properties.dataKey || this.formData.dataKey,
        validateKey: properties.validateKey || this.formData.validateKey,
        expression: properties.expression || this.formData.expression,
        validateValue: properties.validateValue || this.formData.validateValue,
        headers: properties.headers,
        cookies: properties.cookies,
        body: properties.body,
        printInput: properties.printInput,
        printOutput: properties.printOutput,
        inputHistory: properties.inputHistory,
        resultHistory: properties.resultHistory ?? true,
      };
    },
    bindForm() {
      const { currentNode, lf } = this.$props;
      lf.updateText(currentNode.data.id, this.formData.name);
      lf.setProperties(currentNode.data.id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    getParameters() {
      const flowData = this.$props.lf.getGraphRawData();
      let node = this.getStartNode(flowData);
      while (node && node.id !== this.$props.currentNode?.data?.id && node.type !== 'EndNode') {
        const { name, parameterName, inputHistory, resultHistory } = node.properties;
        if (parameterName) {
          if (inputHistory) {
            this.parameters.push({ label: `$${parameterName}.input`, value: `$${parameterName}.input` });
          }
          if (resultHistory) {
            this.parameters.push({ label: `$${parameterName}.output`, value: `$${parameterName}.output` });
          }
        }
        node = this.getNextNode(flowData, node.id);
      }
    },
    getStartNode(flowData) {
      return flowData.nodes.find(node => node.type === 'StartNode');
    },
    getNextNode(flowData, id) {
      const edge = flowData.edges.find(edge => edge.sourceNodeId === id);
      return edge ? this.getNode(flowData, edge.targetNodeId) : null;
    },
    getNode(flowData, id) {
      return flowData.nodes.find(node => node.id === id);
    },
  },
};
</script>
