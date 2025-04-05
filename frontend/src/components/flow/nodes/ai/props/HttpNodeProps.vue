<template>
  <a-form :model="formData" :label-col="{ style: { width: '80px' } }">
    <a-form-item label="参数名称">
      <a-input v-model:value="formData.parameterName" placeholder="参数名称" />
    </a-form-item>

    <a-form-item label="节点名称">
      <a-input v-model:value="formData.name" placeholder="节点名称" />
    </a-form-item>

    <a-form-item label="请求地址">
      <a-input v-model:value="formData.url" placeholder="请求地址" />
    </a-form-item>

    <a-form-item label="请求类型">
      <a-select v-model:value="formData.contentType" :options="contentTypes" placeholder="请求类型" />
    </a-form-item>

    <a-form-item label="请求方式">
      <a-select v-model:value="formData.method" :options="methods" placeholder="请求方式" />
    </a-form-item>

    <a-form-item label="Header">
      <a-input v-model:value="formData.headers" placeholder="Header" />
    </a-form-item>

    <a-form-item label="Cookie">
      <a-input v-model:value="formData.cookies" placeholder="Cookie" />
    </a-form-item>

    <a-form-item label="请求数据">
      <a-mentions v-model:value="formData.body" autofocus :options="parameters" placeholder="输入$补全" prefix="$"
        :rows="3"></a-mentions>
    </a-form-item>

    <a-form-item label="验证表达式">
      <a-input v-model:value="formData.expression" placeholder="${code == 200}" />
    </a-form-item>

    <a-form-item label="保留结果">
      <a-input v-model:value="formData.dataKey" placeholder="保留结果：data || result" />
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
        parameterName: "http",
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
