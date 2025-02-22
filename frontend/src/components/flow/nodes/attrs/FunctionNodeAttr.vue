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
        <a-select v-model:value="formData.model" :options="models" placeholder="选择模型" style="width: 100%;" />
      </a-col>

      <a-col :span="24">
        <a-select v-model:value="formData.funcCall" :options="funcs" placeholder="函数调用" style="width: 100%;" />
      </a-col>

      <a-col :span="24">
        <a-button type="link" @click="addMessageLine">添加一行</a-button>
      </a-col>

      <a-col :span="24">
        <div v-for="(item, i) in formData.messages" :key="i" style="margin-bottom: 16px;" >
          <a-row :gutter="[16, 16]">
            <a-col :span="8" style="margin-bottom: 8px;">
              <a-select v-model:value="item.messageType" :options="messageType" placeholder="选择消息类型"
                style="width: 100%;" />
            </a-col>
            <a-col :span="8">
              <a-select v-model:value="item.message" :options="parameters" placeholder="选择消息" style="width: 100%;" />
            </a-col>
            <a-col :span="8">
              <a-button type="link" danger @click="delMessageLine(i)">删除</a-button>
            </a-col>
          </a-row>

          <a-row :gutter="15">
            <a-col :span="24">
              <a-textarea v-model:value="item.message" placeholder="消息内容" :rows="5" style="width: 100%;" />
            </a-col>
          </a-row>
        </div>
      </a-col>

      <a-col :span="24">
        <a-space>
          <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
          <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
          <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
          <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
        </a-space>
      </a-col>

      <a-col :span="24">
        <a-button type="primary" @click="bindForm">设置</a-button>
        <a-button type="default" @click="setForm" style="margin-left: 10px;">重置</a-button>
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
    this.getFunctionsDatas();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      models: [],
      messageType: [{ label: "系统类型", value: "system" }, { label: "用户类型", value: "human" }, { label: "AI类型", value: "ai" }],
      funcs: [],
      parameters: [],
      formData: {
        id: "",
        name: "",
        model: "",
        funcCall: "",
        resultHistory: true,
        messages: [{ messageType: "system", message: "" }, { messageType: "human", message: "" }],
        printInput: false,
        printOutput: false,
        inputHistory: false,
      },
    };
  },
  methods: {
    setForm() {
      const properties = this.$props.lf.getProperties(this.$props.currentNode.data.id);
      this.formData.name = this.$props.currentNode.data.text.value;
      this.formData.id = this.$props.currentNode.data.id;
      this.formData.model = properties.model;
      this.formData.parameterName = properties.parameterName;
      this.formData.printInput = properties.printInput;
      this.formData.printOutput = properties.printOutput;
      this.formData.funcCall = properties.funcCall;
      this.formData.inputHistory = properties.inputHistory;
      if (properties.resultHistory) {
        this.formData.resultHistory = properties.resultHistory;
      }
      if (properties.messages) {
        this.formData.messages = properties.messages;
      }
    },
    bindForm() {
      const id = this.$props.currentNode.data.id;
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
    getFunctionsDatas() {
      this.$post("/ollamafunctions/list", { page: { current: 1, size: 1000 } }).then(res => {
        if (res.data.code === this.$success) {
          this.funcs = res.data.data.rows;
        }
      });
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
      const flowData = this.$props.lf.getGraphRawData();
      let node = this.getStartNode(flowData);
      while (node && node.id !== this.$props.currentNode.data.id && node.type !== 'EndNode') {
        const properties = node.properties;
        const parameterName = properties?.parameterName;

        if (parameterName) {
          if (properties.inputHistory) {
            this.parameters.push({ label: "$" + parameterName + ".input", value: "$" + parameterName + ".input" });
          }
          if (properties.resultHistory) {
            this.parameters.push({ label: "$" + parameterName + ".output", value: "$" + parameterName + ".output" });
          }
        }
        node = this.getNextNode(flowData, node.id);
      }
    },
    getStartNode(flowData) {
      return flowData.nodes.find(node => node.type === 'StartNode');
    },
    getNode(flowData, id) {
      return flowData.nodes.find(node => node.id === id);
    },
    getNextNode(flowData, id) {
      const edge = flowData.edges.find(edge => edge.sourceNodeId === id);
      return edge ? this.getNode(flowData, edge.targetNodeId) : null;
    },
  }
};
</script>

<style scoped>
.h-panel {
  display: flex;
  align-items: center;
}

.h-panel-bar {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.h-split {
  width: 1px;
  background-color: #ccc;
  margin: 0 10px;
}

.h-panel-right {
  cursor: pointer;
}
</style>
