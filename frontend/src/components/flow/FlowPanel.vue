<template>
  <div class="flowCanvas">
    <NodePanel v-if="lf" :lf="lf"></NodePanel>
    <DataPanel v-if="lf" :lf="lf"></DataPanel>
    <div ref="container" style="display: flex;flex-grow: 1;border: 1px solid #1890FF;overflow: hidden;"></div>

    <div style="width: 300px;background-color: white;margin-left: 8px;border: 1px solid #1890FF;" v-show="drawer">
      <LLMNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'LLMNode'" @close="close"></LLMNodeAttr>
      <VectorNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'VectorNode'" @close="close"></VectorNodeAttr>
      <FunctionNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'FunctionNode'" @close="close"></FunctionNodeAttr>
      <LineAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'polyline'" @close="close"></LineAttr>
      <EndAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'EndNode'" @close="close"></EndAttr>
      <HttpNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'HttpNode'" @close="close"></HttpNodeAttr>
      <DatabaseNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'DatabaseNode'" @close="close"></DatabaseNodeAttr>
      <BugNodeAttr :lf="lf" :currentNode="currentNode" :key="this.currentNode?.data?.id"
        v-if="this.currentNode?.data?.type == 'BugNode'" @close="close"></BugNodeAttr>
    </div>
  </div>
</template>

<style>
.flowCanvas {
  position: relative;
  height: 100vh;
  margin: 0px;
  display: flex;
}
</style>

<script>
import LogicFlow from "@logicflow/core";

import "@logicflow/core/lib/style/index.css";
import theme from "@/components/flow/theme.js";
import {
  Menu,
  MiniMap,
  SelectionSelect,
  Snapshot,
  InsertNodeInPolyline,
  NodeResize
} from "@logicflow/extension";

import "@logicflow/extension/lib/style/index.css";
import NodePanel from "@/components/flow/NodePanel.vue";
import DataPanel from "@/components/flow/DataPanel.vue";
import LLMNodeAttr from "@/components/flow/nodes/attrs/LLMNodeAttr.vue";
import VectorNodeAttr from "@/components/flow/nodes/attrs/VectorNodeAttr.vue";
import FunctionNodeAttr from "@/components/flow/nodes/attrs/FunctionNodeAttr.vue";
import LineAttr from "@/components/flow/nodes/attrs/LineAttr.vue";
import EndAttr from "@/components/flow/nodes/attrs/EndAttr.vue";
import HttpNodeAttr from "@/components/flow/nodes/attrs/HttpNodeAttr.vue";
import DatabaseNodeAttr from "@/components/flow/nodes/attrs/DatabaseNodeAttr.vue";
import BugNodeAttr from "@/components/flow/nodes/attrs/BugNodeAttr.vue";

export default {
  data() {
    return {
      lf: null,
      currentNode: null,
      drawer: false,
    };
  },
  components: {
    NodePanel,
    DataPanel,
    LLMNodeAttr,
    VectorNodeAttr,
    FunctionNodeAttr,
    LineAttr,
    EndAttr,
    BugNodeAttr,
    HttpNodeAttr,
    DatabaseNodeAttr,
  },
  mounted() {
    LogicFlow.use(Menu);
    LogicFlow.use(SelectionSelect);
    LogicFlow.use(Snapshot);
    LogicFlow.use(InsertNodeInPolyline);
    LogicFlow.use(MiniMap);
    LogicFlow.use(InsertNodeInPolyline);
    //LogicFlow.use(NodeResize);

    //初始化
    this.lf = new LogicFlow({
      snapline: true,
      container: this.$refs.container,
      grid: {
        size: 10,
        visible: true,
        //type: "mesh",
        config: {
          color: "#ababab",
          thickness: 1,
        },
      },
      plugins: [NodeResize], // 启用节点大小调整插件
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
      /* eventCenter.on("node:resize", (oldNodeSize, newNodeSize) => {
        console.log(oldNodeSize, newNodeSize);
      }); */
      eventCenter.on("node:click", (args) => {
        this.drawer = false;
        if (args?.data?.type == "StartNode") {
          return;
        }
        this.drawer = true;
        this.currentNode = args;
      });
      eventCenter.on("edge:click", (args) => {
        this.drawer = true;
        this.currentNode = args;
      });
    },

    handleClose(done) {
      done();
    },
    close() {
      this.drawer = false;
    },
  },
};
</script>