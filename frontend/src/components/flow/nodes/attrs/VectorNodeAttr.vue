<template>
  <div style="padding: 10px;">
    <a-row :gutter="[16,16]">
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
        <a-select v-model:value="formData.collectionId" :options="collectionList" @change="collectionChange" placeholder="知识库集合" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-select v-model:value="formData.filename" :options="filenameList" placeholder="知识库文件" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-input-number v-model:value="formData.limit" :min="1" placeholder="检索数量" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-slider v-model:value="formData.similarityScore" :min="0" :max="10" :step="0.1" :tooltip-visible="false" :marks="marks" />
      </a-col>
      <a-col :span="24">
        <a-checkbox v-model:checked="formData.analyse">模型分析</a-checkbox>
      </a-col>
      <a-col :span="24">
        <a-textarea v-model:value="formData.system" :rows="5" placeholder="系统提示词" :disabled="!formData.analyse" style="width: 100%;" />
      </a-col>
      <a-col :span="24">
        <a-textarea v-model:value="formData.message" :rows="5" placeholder="输入问题" style="width: 100%;" />
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
    this.getModel();
    this.getCollectionList();
    this.getFilenameList();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      models: [],
      collectionList: [],
      filenameList: [],
      parameters: [],
      formData: {
        collectionId: "",
        limit: 5,
        filename: "",
        collectionName: "",
        message: "",
        resultHistory: true,
        similarityScore: 7,
      },
      marks: {
        0: '0',
        10: '10'
      }
    };
  },
  methods: {
    setForm() {
      const properties = this.$props.lf.getProperties(this.$props.currentNode.data.id);
      this.formData.name = this.$props.currentNode.data.text.value;
      this.formData.model = properties.model;
      this.formData.printInput = properties.printInput;
      this.formData.printOutput = properties.printOutput;
      this.formData.limit = properties.limit || this.formData.limit;
      this.formData.collectionId = properties.collectionId;
      this.formData.collectionName = properties.collectionName;
      this.formData.system = properties.system;
      this.formData.message = properties.message;
      this.formData.parameterName = properties.parameterName;
      this.formData.filename = properties.filename;
      this.formData.inputHistory = properties.inputHistory;
      this.formData.resultHistory = properties.resultHistory;
      this.formData.analyse = properties.analyse;
      this.formData.similarityScore = properties.similarityScore > 0 ? properties.similarityScore * 10 : this.formData.similarityScore;
    },
    bindForm() {
      const id = this.$props.currentNode.data.id;
      this.$props.lf.updateText(id, this.formData.name);
      this.$props.lf.setProperties(id, {
        ...this.formData,
        similarityScore: (this.formData.similarityScore * 0.1).toFixed(2),
      });
      this.$emit('close', false);
      message.success("设置成功");
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
    getCollectionList() {
      this.$post("/ollamacollection/list", { page: { current: 1, size: 1000 } }).then(res => {
        if (res.data.code === this.$success) {
          this.collectionList = res.data.data.rows;
        } else {
          message.error(res.data.msg);
        }
      });
    },
    getFilenameList() {
      if (this.formData.collectionId) {
        this.$post("/ollamaembedding/list", { vars: { "a.collection_id": this.formData.collectionId }, page: { current: 1, size: 1000 } }).then(res => {
          if (res.data.code === this.$success) {
            this.filenameList = res.data.data.rows;
          } else {
            message.error(res.data.msg);
          }
        });
      }
    },
    collectionChange(value) {
      this.formData.collectionId = value;
      this.getFilenameList();
    },
    getParameters() {
      const flowData = this.$props.lf.getGraphRawData();
      let node = this.getStartNode(flowData);
      while (node && node.id !== this.$props.currentNode.data.id && node.type !== 'EndNode') {
        const properties = node.properties;
        if (properties.parameterName) {
          if (properties.inputHistory) {
            this.parameters.push({ label: `$${properties.parameterName}.input`, value: `$${properties.parameterName}.input` });
          }
          if (properties.resultHistory) {
            this.parameters.push({ label: `$${properties.parameterName}.output`, value: `$${properties.parameterName}.output` });
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
  }
};
</script>
