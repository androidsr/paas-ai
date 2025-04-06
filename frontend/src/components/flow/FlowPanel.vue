<template>
  <div class="flowCanvas">
    <NodePanel v-if="lf" :lf="lf"></NodePanel>
    <DataPanel v-if="lf" :lf="lf"></DataPanel>
    <div ref="container" @click.stop="resetPanel"
      style="display: flex;flex-grow: 1;border: 1px solid #1890FF;overflow: hidden;"></div>

    <a-drawer v-model:open="visible" title="节点属性" :mask="false" :width="500"
      :maskStyle="{ background: 'rgba(0, 0, 0, 0)' }">
      <LLMNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'LLMNode'" @close="close"></LLMNodeProps>
      <VectorNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'VectorNode'" @close="close"></VectorNodeProps>
      <FunctionNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'FunctionNode'" @close="close"></FunctionNodeProps>
      <LineProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'polyline' || this.currentNode?.data?.type == 'bezier' || this.currentNode?.data?.type == 'line'"
        @close="close"></LineProps>
      <EndProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'EndNode'" @close="close"></EndProps>
      <HttpNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'HttpNode'" @close="close"></HttpNodeProps>
      <DatabaseNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'DatabaseNode'" @close="close"></DatabaseNodeProps>
      <ScriptNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'ScriptNode'" @close="close"></ScriptNodeProps>
      <McpNodeProps :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'McpNode'" @close="close"></McpNodeProps>
    </a-drawer>
  </div>
</template>

<style>
.flowCanvas {
  position: relative;
  height: 96vh;
  margin: 0px;
  display: flex;
}
</style>

<script>
import LogicFlow from "@logicflow/core";
import '@logicflow/core/es/index.css'

import "@logicflow/core/lib/style/index.css";
import theme from "./theme.js";
import {
  Menu,
  SelectionSelect,
  Snapshot,
  InsertNodeInPolyline,
  NodeResize
} from "@logicflow/extension";

import "@logicflow/extension/lib/style/index.css";
import NodePanel from "./NodePanel.vue";
import DataPanel from "./DataPanel.vue";
import LLMNodeProps from "./nodes/ai/props/LLMNodeProps.vue";
import VectorNodeProps from "./nodes/ai/props/VectorNodeProps.vue";
import FunctionNodeProps from "./nodes/ai/props/FunctionNodeProps.vue";
import LineProps from "./nodes/ai/props/LineProps.vue";
import EndProps from "./nodes/ai/props/EndProps.vue";
import HttpNodeProps from "./nodes/ai/props/HttpNodeProps.vue";
import DatabaseNodeProps from "./nodes/ai/props/DatabaseNodeProps.vue";
import ScriptNodeProps from "./nodes/ai/props/ScriptNodeProps.vue";
import McpNodeProps from "./nodes/ai/props/McpNodeProps.vue";

export default {
  data() {
    return {
      lf: null,
      currentNode: null,
      drawer: false,
      visible: false,
    };
  },
  components: {
    NodePanel,
    DataPanel,
    LLMNodeProps,
    VectorNodeProps,
    FunctionNodeProps,
    LineProps,
    EndProps,
    HttpNodeProps,
    DatabaseNodeProps,
    ScriptNodeProps,
    McpNodeProps,
  },
  mounted() {
    LogicFlow.use(Menu);
    LogicFlow.use(SelectionSelect);
    LogicFlow.use(Snapshot);
    LogicFlow.use(InsertNodeInPolyline);
    LogicFlow.use(InsertNodeInPolyline);
    //初始化
    this.lf = new LogicFlow({
      snapline: true,
      edgeType: this.flowType == 1 ? 'bezier' : 'polyline',
      container: this.$refs.container,
      grid: {
        size: 25,
        //visible: true,
        //type: "mesh",
        config: {
          color: "#ababab",
          thickness: 1,
        },
      },
      plugins: [NodeResize], // 启用节点大小调整插件
      /* pluginsOptions: {
        proximityConnect: {
          enable: true, // 插件是否启用
          distance: 50, // 渐进连线阈值
          reverseDirection: false, // 连线方向
        },
      }, */
      keyboard: {
        enabled: true,
        shortcuts: [
          {
            keys: ["Del"],
            callback: () => {
              const elements = this.lf.getSelectElements(true);
              this.lf.clearSelectElements();
              elements.edges.forEach((edge) => this.lf.deleteEdge(edge.id));
              elements.nodes.forEach((node) => this.lf.deleteNode(node.id));
            },
          },
        ],
      },
      style: theme,
    });
    //绑定事件
    const { eventCenter } = this.lf.graphModel;
    this.bindEvent(eventCenter);
    this.lf.render();
  },

  methods: {
    bindEvent(eventCenter) {
      eventCenter.on("node:click", (args) => {
        console.log(args)
        this.drawer = false;
        if (args?.data?.type == "StartNode") {
          this.visible = false;
          return;
        }
        this.currentNode = args;
        this.drawer = true;
        this.visible = true;
      });
      eventCenter.on("edge:click", (args) => {
        this.drawer = true;
        this.currentNode = args;
        this.visible = true;
      });
    },
    handleClose(done) {
      done();
    },
    close() {
      this.visible = false;
      this.drawer = false;
      this.currentNode = null;
    },
  },
};
</script>