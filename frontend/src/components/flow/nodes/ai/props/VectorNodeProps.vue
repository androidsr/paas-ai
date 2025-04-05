<template>
  <a-form :label-col="{ style: { width: '80px' } }" label-align="left">
    <!-- 参数名称 -->
    <a-form-item label="参数名称">
      <a-input v-model:value="formData.parameterName" placeholder="参数名称" style="width: 100%;" />
    </a-form-item>

    <!-- 节点名称 -->
    <a-form-item label="节点名称">
      <a-input v-model:value="formData.name" placeholder="节点名称" style="width: 100%;" />
    </a-form-item>

    <a-form-item label="AI平台">
      <a-select v-model:value="formData.channelId" :options="channels" @select="getChannelItem" />
    </a-form-item>
    <a-form-item label="选择模型">
      <a-select v-model:value="formData.model" :options="models" />
    </a-form-item>


    <!-- 知识库集合 -->
    <a-form-item label="知识库集合">
      <a-select v-model:value="formData.collectionId" :options="collectionList" @change="collectionChange"
        placeholder="知识库集合" style="width: 100%;" />
    </a-form-item>

    <!-- 知识库文件 -->
    <a-form-item label="知识库文件">
      <a-select v-model:value="formData.filename" :options="filenameList" placeholder="知识库文件" style="width: 100%;" />
    </a-form-item>

    <!-- 检索数量 -->
    <a-form-item label="检索数量">
      <a-input-number v-model:value="formData.limit" :min="1" placeholder="检索数量" style="width: 100%;" />
    </a-form-item>

    <!-- 相似度评分 -->
    <a-form-item label="相似度评分">
      <a-slider v-model:value="formData.similarityScore" :min="0" :max="10" :step="0.1" :tooltip-visible="true"
        :marks="marks" />
    </a-form-item>

    <!-- 模型分析 -->
    <a-form-item label="模型分析">
      <a-checkbox v-model:checked="formData.analyse">模型分析</a-checkbox>
    </a-form-item>

    <!-- 系统提示词 -->
    <a-form-item label="系统提示词">
      <a-mentions v-model:value="formData.system" autofocus :options="parameters" placeholder="输入$补全" prefix="$"
        :rows="3" :disabled="!formData.analyse"></a-mentions>
    </a-form-item>

    <!-- 输入问题 -->
    <a-form-item label="输入问题">
      <a-mentions v-model:value="formData.message" autofocus :options="parameters" placeholder="输入$补全" prefix="$"
        :rows="3"></a-mentions>
    </a-form-item>
    <a-form-item label="输出设置">
      <a-checkbox v-model:checked="formData.printInput">打印输入</a-checkbox>
      <a-checkbox v-model:checked="formData.printOutput">打印输出</a-checkbox>
      <a-checkbox v-model:checked="formData.inputHistory">保存输入</a-checkbox>
      <a-checkbox v-model:checked="formData.resultHistory">保存输出</a-checkbox>
      <a-checkbox v-model:checked="formData.printComplete">完成输出</a-checkbox>
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
import { GetList as CollectionList } from '../../../../../../wailsjs/go/biz/CollectionBiz.js';
import { GetList as DocumentList } from '../../../../../../wailsjs/go/biz/DocumentBiz.js';
import { Get, GetList } from '../../../../../../wailsjs/go/biz/AiChannelBiz.js';

export default {
  props: {
    opend: false,
    lf: Object,
    currentNode: Object,
  },
  mounted() {
    this.getChannel();
    this.getCollectionList();
    this.getFilenameList();
    this.setForm();
    this.getParameters();
  },
  data() {
    return {
      channels: [],
      models: [],
      collectionList: [],
      filenameList: [],
      parameters: [],
      formData: {
        parameterName: "vector",
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
      this.formData.channelId = properties.channelId;
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
      this.formData.printComplete = properties.printComplete;
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
    getCollectionList() {
      CollectionList().then(res => {
        res.data.forEach(item => {
          this.collectionList.push({ label: item, value: item });
        });
      });
    },
    getFilenameList(data) {
      DocumentList(data).then(res => {
        if (res.code == 200) {
          this.filenameList = res.data;
        } else {
          this.filenameList = [];
          message.error(res.data.msg);
        }
      })
    },
    collectionChange(value) {
      this.formData.collectionId = value;
      this.getFilenameList();
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
  }
};
</script>
