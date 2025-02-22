<template>
  <a-row :gutter="[16, 16]" style="padding: 10px;">
    <a-col :span="24">
      <a-input v-model:value="formData.parameterName" placeholder="参数名称" style="width: 100%;" />
    </a-col>
    <a-col :span="24">
      <a-input v-model:value="formData.name" placeholder="节点名称" style="width: 100%;" />
    </a-col>
    <a-col :span="24">
      <a-select v-model:value="formData.bugType" placeholder="爬取类型" style="width: 100%;">
        <a-select-option v-for="item in bugTypeData" :key="item.value" :value="item.value">
          {{ item.label }}
        </a-select-option>
      </a-select>
    </a-col>
    <a-col :span="24">
      <a-input v-model:value="formData.website" placeholder="爬取网点" style="width: 100%;" />
    </a-col>
    <a-col :span="24">
      <a-select v-model:value="formData.text" mode="tags" :options="parameters" placeholder="搜索文本" style="width: 100%;"/>
    </a-col>
    <a-col :span="24">
      <a-flex v-for="(item, i) in formData.selecters" :key="i" style="margin-bottom: 10px;">
        <a-input v-model:value="item.selecterKey" placeholder="保存名称" style="width: 45%;" />
        <a-input v-model:value="item.selecter" placeholder="标签选择器" style="width: 45%;" />
        <a-link class="h-input-addon link" @click="delSelecter(i)">删除</a-link>
      </a-flex>
    </a-col>
    <a-col :span="24">
      <a-link class="link" @click="addSelecter">添加一行</a-link>
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
      <a-space>
        <a-button type="primary" @click="bindForm">设置</a-button>
        <a-button type="default" @click="setForm">重置</a-button>
      </a-space>
    </a-col>
  </a-row>
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
    this.getWebsites();
    this.getParameters();
  },
  data() {
    return {
      websites: [],
      bugTypeData: [{ label: "爬取链接", value: "link" }, { label: "爬取文本", value: "text" }],
      parameters: [],
      formData: {
        bugType: "link",
        website: "",
        selecters: [],
        text: "",
      }
    }
  },
  methods: {
    setForm() {
      this.formData.name = this.$props.currentNode.data.text.value;
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.lf.getProperties(id);
      this.formData.parameterName = properties.parameterName;
      this.formData.website = properties.website;
      this.formData.text = properties.text;
      if (!!properties.bugType) {
        this.formData.bugType = properties.bugType;
      }
      if (!!properties.selecters) {
        this.formData.selecters = properties.selecters;
      }
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
      this.$props.lf.setProperties(id, Object.assign({ ...this.formData }, {
        similarityScore: Number((this.formData.similarityScore * 0.1).toFixed(2)),
      }));
      this.$emit('close', false);
      message.success("设置成功");
    },
    addSelecter() {
      this.formData.selecters.push({ selecterKey: "", selecter: "", });
    },
    delSelecter(index) {
      this.formData.selecters.splice(index, 1);
    },
    getWebsites() {
      this.$post("/crawlerselecter/list", { page: { current: 1, size: 1000 } }).then(res => {
        if (res.data.code == this.$success) {
          this.websites = res.data.data.rows;
        } else {
          this.websites = [];
          message.error(res.data.msg);
        }
      })
    },
    getParameters() {
      let flowData = (this.$props.lf).getGraphRawData();
      let nodes = flowData["nodes"];
      let node = this.getStartNode(flowData);
      while (!!node && node?.id != this.$props.currentNode?.data?.id && node?.type != 'EndNode') {
        var properties = node?.properties;
        var name = properties?.name;
        var parameterName = properties?.parameterName;

        if (!!parameterName) {
          if (properties?.inputHistory == true) {
            this.parameters.push({ label: "$" + parameterName + ".input", value: "$" + parameterName + ".input" });
          } if (properties?.resultHistory == true) {
            this.parameters.push({ label: "$" + parameterName + ".output", value: "$" + parameterName + ".output" });
          }
        }
        node = this.getNextNode(flowData, node.id);
      }
    },
    getStartNode(flowData) {
      let nodes = flowData["nodes"];
      for (var i = 0; i < nodes.length; i++) {
        let node = nodes[i];
        if (node.type == 'StartNode') {
          return node;
        }
      }
    },
    getNode(flowData, id) {
      let nodes = flowData["nodes"];
      for (var i = 0; i < nodes.length; i++) {
        let node = nodes[i];
        if (node.id == id) {
          return node;
        }
      }
    },
    getNextNode(flowData, id) {
      let edges = flowData["edges"]
      for (var i = 0; i < edges.length; i++) {
        let edge = edges[i];
        if (edge.sourceNodeId == id) {
          return this.getNode(flowData, edge.targetNodeId);
        }
      }
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
  margin: 0 5px;
}

.h-input-addon.link {
  color: #1890ff;
  cursor: pointer;
}
</style>
