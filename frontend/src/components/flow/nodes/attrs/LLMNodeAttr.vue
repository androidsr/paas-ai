<template>
  <div style="padding: 10px;">
    <a-row :gutter="[15, 15]">
      <a-col :span="24">
        <a-input v-model:value="formData.parameterName" placeholder="参数名称" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-input v-model:value="formData.name" placeholder="节点名称" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select v-model:value="formData.model" :options="models" placeholder="选择模型" style="width: 100%" />
      </a-col>
      <a-col :span="24">
        <a-button @click="addMessageLine" type="link">添加一行</a-button>
      </a-col>
      <a-col :span="24">
        <a-row v-for="(item, i) in formData.messages" :key="i" :gutter="[15, 15]" style="margin-bottom: 16px;">
          <a-col :span="24">
            <a-row style="margin-bottom: 8px;">
              <a-col :span="8">
                <a-select v-model:value="item.messageType" :options="messageType" style="width: 100%" />
              </a-col>
              <a-col :span="8">
                <a-select v-model:value="item.message" :options="parameters" style="width: 100%" />
              </a-col>
              <a-col :span="8">
                <a-button type="link" @click="delMessageLine(i)">删除</a-button>
              </a-col>
            </a-row>
            <a-row>
              <a-textarea v-model:value="item.message" placeholder="消息内容" :rows="5" style="width: 100%" />
            </a-row>
          </a-col>
        </a-row>
      </a-col>
      <a-col :span="24">
        <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
        <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
        <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
        <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
      </a-col>
      <a-col :span="24">
        <a-button @click="bindForm" type="primary">设置</a-button>
        <a-button @click="setForm" type="default" style="margin-left: 15px;">重置</a-button>
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
    this.getModel();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      models: [],
      messageType: [
        { label: "系统类型", value: "system" },
        { label: "用户类型", value: "human" },
        { label: "AI类型", value: "ai" }
      ],
      parameters: [],
      formData: {
        id: "",
        name: "",
        model: "",
        resultHistory: true,
        messages: [{ messageType: "system", message: "" }, { messageType: "human", message: "" }],
      }
    };
  },
  methods: {
    setForm() {
      this.formData.name = this.$props.currentNode.data.text.value;
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.lf.getProperties(id);
      this.formData.id = this.$props.currentNode.data.id;
      this.formData.model = properties.model;
      this.formData.printInput = properties.printInput;
      this.formData.printOutput = properties.printOutput;
      this.formData.parameterName = properties.parameterName;
      if (properties.resultHistory) {
        this.formData.resultHistory = properties.resultHistory;
      }
      this.formData.inputHistory = properties.inputHistory;
      if (properties.messages) {
        this.formData.messages = properties.messages;
      }
    },
    bindForm() {
      let id = this.$props.currentNode.data.id;
      this.$props.lf.updateText(id, this.formData.name);
      this.$props.lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    addMessageLine() {
      this.formData.messages.push({ messageType: "human", message: "" });
    },
    delMessageLine(index) {
      this.formData.messages.splice(index, 1);
    },
    getModel() {
      this.$get("/chat/getModel").then(res => {
        if (res.data.code === this.$success) {
          this.models = res.data.data;
        } else {
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
      for (let i = 0; i < nodes.length; i++) {
        if (nodes[i].type === 'StartNode') {
          return nodes[i];
        }
      }
    },
    getNode(flowData, id) {
      let nodes = flowData["nodes"];
      return nodes.find(node => node.id === id);
    },
    getNextNode(flowData, id) {
      let edges = flowData["edges"];
      let edge = edges.find(edge => edge.sourceNodeId === id);
      return edge ? this.getNode(flowData, edge.targetNodeId) : null;
    },
  }
};
</script>
