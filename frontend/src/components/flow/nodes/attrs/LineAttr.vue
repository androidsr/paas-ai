<template>
  <a-row :gutter="[15, 15]" style="padding: 10px;">
    <a-col :span="24">
      <a-input v-model:value="formData.name" placeholder="节点名称" style="width: 100%;" />
    </a-col>

    <a-col>
      <div class="h-input-group" style="width: 100%">
        <a-input-group compact>
          <a-input-search v-model:value="formData.startCondition" :placeholder="'条件表达式'" style="width: 145px" />
          <a-select v-model:value="formData.expression" :options="expressions" style="width: 90px" bordered={false}
            allowClear>
          </a-select>
          <a-input-search v-model:value="formData.endCondition" :placeholder="'条件表达式'" style="width: 145px" />
        </a-input-group>
      </div>
    </a-col>

    <!-- Loop Settings -->
    <a-col :span="24">
      <div class="h-input-group" style="width: 100%">
        <a-input-group compact>
          <a-select v-model:value="formData.forType" :options="forType" @change="forTypeSelect" style="width: 100px"
            bordered={false} allowClear />
          <a-input-search v-model:value="formData.forData" :placeholder="'循环变量'" style="width: 200px" />
          <a-input-number v-model:value="formData.forSize" placeholder="次数" style="width: 80px" />
        </a-input-group>
      </div>
    </a-col>

    <a-col :span="24">
      <a-space>
        <a-button type="primary" @click="bindForm">设置</a-button>
        <a-button @click="setForm" type="default">重置</a-button>
      </a-space>
    </a-col>
  </a-row>
</template>

<script>
import { message} from 'ant-design-vue';

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
      forType: [
        { label: "开始位置", value: '1' },
        { label: "结束位置", value: '2' },
        { label: "默认关闭", value: '0' },
      ],
      expressions: [
        { label: "==", value: "=" },
        { label: "!=", value: "!=" },
        { label: ">", value: ">" },
        { label: ">=", value: ">=" },
        { label: "<", value: "<" },
        { label: "<=", value: "<=" },
        { label: "包含", value: "include" },
        { label: "不包含", value: "exclusive" },
        { label: "以*开始", value: "startWith" },
        { label: "以*结束", value: "endWith" },
        { label: "大小", value: "size" },
      ],
      parameters: [],
      formData: {
        forType: "",
        startCondition: "",
        endCondition: "",
        expression: "=",
        forData: "",
        forSize: null,
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
      this.formData.forData = properties.forData;
      this.formData.forSize = properties.forSize;
      this.formData.startCondition = properties.startCondition;
      this.formData.endCondition = properties.endCondition;
      if (!!properties.expression) {
        this.formData.expression = properties.expression;
      }
    },
    bindForm() {
      let id = this.$props.currentNode.data.id;
      this.$props.lf.updateText(id, this.formData.name);
      this.$props.lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
    },
    forTypeSelect(data) {
      if (data.value == '1') {
        this.formData.name = "循环开始";
      } else if (data.value == '2') {
        this.formData.name = "循环结束";
      } else if (data.value == '0') {
        this.formData.name = "";
      }
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
          }
          if (properties?.resultHistory == true) {
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
