<template>
  <a-form :model="formData" :label-col="{ style: { width: '80px' } }">
    <a-form-item label="显示名称">
      <a-input v-model:value="formData.name" placeholder="显示名称" />
    </a-form-item>
    <a-form-item label="条件表达式">
      <a-row :gutter="[8, 8]">
        <a-col :span="24">
          <a-mentions v-model:value="formData.condition" autofocus :options="parameters" placeholder="${aaa=='1'}$补全"
            prefix="$"></a-mentions>
        </a-col>
      </a-row>
    </a-form-item>

    <a-form-item label="遍历数据">
      <a-row :gutter="[8, 8]">
        <a-col :span="8">
          <a-select v-model:value="formData.loopType" :options="forType" @change="forTypeSelect" allowClear />
        </a-col>
        <a-col :span="12">
          <a-mentions v-model:value="formData.loopVariable" autofocus :options="parameters" placeholder="遍历变量$补全"
            prefix="$" :disabled="formData.loopType != 1"></a-mentions>
        </a-col>
        <a-col :span="4">
          <a-input-number v-model:value="formData.loopSize" placeholder="次数" style="width: 100%;"
            :disabled="formData.loopType != 1" />
        </a-col>
      </a-row>
    </a-form-item>
    <a-col :span="24">
      <a-space>
        <a-button type="primary" @click="bindForm">保存</a-button>
        <a-button @click="setForm" type="default">重置</a-button>
      </a-space>
    </a-col>
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
    this.getParameters();
    this.setForm();
  },
  data() {
    return {
      forType: [
        { label: "开始位置", value: '1' },
        { label: "结束位置", value: '2' },
        { label: "默认关闭", value: '0' },
      ],
      parameters: [],
      formData: {
        forType: "",
        condition: "",
        loopType: "0",
        loopVariable: "",
        loopSize: 0,
        name: ""
      }
    };
  },
  methods: {
    setForm() {
      let id = this.$props.currentNode.data.id;
      let properties = this.$props.lf.getProperties(id);
      this.formData.id = id;
      this.formData.name = properties.name;
      if (!!properties.forType) {
        this.formData.forType = properties.forType;
      }
      this.formData.condition = properties.condition;
      this.formData.loopType = properties.loopType;
      this.formData.loopVariable = properties.loopVariable;
      this.formData.loopSize = properties.loopSize;
    },
    bindForm() {
      let id = this.$props.currentNode.data.id;
      this.$props.lf.updateText(id, this.formData.name);
      this.$props.lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    forTypeSelect(data) {
      if (data == '1') {
        this.formData.name = "循环开始";
      } else if (data == '2') {
        this.formData.name = "循环结束";
      } else {
        this.formData.name = "";
      }
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
      let edges = flowData["edges"];
      for (var i = 0; i < edges.length; i++) {
        let edge = edges[i];
        if (edge.sourceNodeId == id) {
          return this.getNode(flowData, edge.targetNodeId);
        }
      }
    }
  }
}
</script>

<style scoped>
.h-input-group {
  display: flex;
  gap: 10px;
  width: 100%;
}
</style>
