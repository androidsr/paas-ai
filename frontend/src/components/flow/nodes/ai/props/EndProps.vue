<template>
  <a-form :label-col="{ style: { width: '80px' } }">

    <a-form-item label="输出内容" v-for="(item, i) in formData.messages" :key="i">
      <a-row>
        <a-col :span="10">
          <a-input v-model:value="item.varName" placeholder="变量名" class="var-name" />
        </a-col>
        <a-col :span="1">
          =
        </a-col>
        <a-col :span="10">
          <a-mentions v-model:value="item.varValue" autofocus :options="parameters" placeholder="输入$补全"
            prefix="$"></a-mentions>
        </a-col>
        <a-col :span="3">
          <a-button type="link" danger @click="delMessageLine(i)">
            <DeleteOutlined />
          </a-button>
        </a-col>
      </a-row>
    </a-form-item>
    <a-button @click="addMessageLine" type="link">添加一行</a-button>
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
      parameters: [],
      loading: false, // 设置加载状态
      formData: {
        parameterName: "end",
        messages: []
      }
    };
  },
  methods: {
    setForm() {
      const { currentNode, lf } = this.$props;
      const id = currentNode.data.id;
      const properties = lf.getProperties(id);
      this.formData.messages = properties.messages || [];
    },

    bindForm() {
      this.loading = true; // 设置按钮加载状态
      const { currentNode, lf } = this.$props;
      const id = currentNode.data.id;
      lf.setProperties(id, this.formData);
      this.$emit('close', false);
      message.success("设置成功");
      this.loading = false; // 结束加载状态
    },

    addMessageLine() {
      this.formData.messages.push({ varName: "", varValue: "" });
    },

    delMessageLine(index) {
      this.formData.messages.splice(index, 1);
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
}
</script>
