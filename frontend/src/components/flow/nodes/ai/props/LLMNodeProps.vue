<template>
  <a-form :model="formData" :label-col="{ style: { width: '80px' } }">
    <a-form-item label="参数名称">
      <a-input v-model:value="formData.parameterName" placeholder="参数名称" />
    </a-form-item>
    <a-form-item label="节点名称">
      <a-input v-model:value="formData.name" placeholder="节点名称" />
    </a-form-item>

    <a-form-item label="AI平台">
      <a-select v-model:value="formData.channelId" :options="channels" @select="getChannelItem" />
    </a-form-item>
    <a-form-item label="选择模型">
      <a-select v-model:value="formData.model" :options="models" />
    </a-form-item>

    <a-form-item label="消息内容" v-for="(item, i) in formData.messages" :key="i">
      <a-row>
        <a-col :span="20">
          <a-select v-model:value="item.messageType" :options="messageType" placeholder="选择消息类型" />
        </a-col>
        <a-col :span="4">
          <a-button type="link" danger @click="delMessageLine(i)">
            <DeleteOutlined />
          </a-button>
        </a-col>
      </a-row>
      <a-row :gutter="15" style="margin-top: 6px;">
        <a-col :span="24">
          <a-mentions v-model:value="item.message" autofocus :options="parameters" placeholder="输入$补全" prefix="$"
            :rows="3"></a-mentions>
        </a-col>
      </a-row>
    </a-form-item>

    <a-form-item>
      <a-button @click="addMessageLine" type="link">添加一行</a-button>
    </a-form-item>
    <a-form-item label="输出设置">
      <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
      <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
      <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
      <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
      <a-checkbox v-model:checked="formData.printComplete">完成输出</a-checkbox>
    </a-form-item>
    <br/>
    <a-form-item>
      <a-button @click="bindForm" type="primary">保存</a-button>
      <a-button @click="setForm" type="default" style="margin-left: 15px;">重置</a-button>
    </a-form-item>
  </a-form>

</template>

<script>
import { message } from 'ant-design-vue';
import { Get, GetList } from '../../../../../../wailsjs/go/biz/AiChannelBiz.js';

export default {
  props: {
    lf: Object,
    currentNode: Object,
  },
  mounted() {
    this.getChannel();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      channels: [],
      models: [],
      messageType: [
        { label: "系统类型", value: "system" },
        { label: "用户类型", value: "human" },
        { label: "AI类型", value: "ai" }
      ],
      parameters: [],
      formData: {
        parameterName: "llm",
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
      console.log(this.$props.currentNode.data)
      this.formData.name = this.$props.currentNode.data.text.value;
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.currentNode.data.properties;
      this.formData.id = this.$props.currentNode.data.id;
      this.formData.channelId = properties.channelId;
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
      this.formData.printComplete = properties.printComplete;
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
    getChannel() {
      GetList().then(res => {
        if (res.code == 200) {
          this.channels = res.data;
          if (!this.channelId) {
            this.getChannelItem(this.channels[0].value);
          } else {
            this.getChannelItem(this.channelId);
          }
        } else {
          message.error(res.data.msg);
        }
      });
    },
    getChannelItem(channelId) {
      if (!channelId) return;
      this.channelId = channelId;
      Get(this.channelId).then(res => {
        if (res.code == 200) {
          this.models = [];
          res.data.models.split(",").forEach(item => {
            this.models.push({ label: item, value: item });
          });
          this.model = this.models[0].value;
        } else {
          this.$message.error("获取设置失败！");
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
